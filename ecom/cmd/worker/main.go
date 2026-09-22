package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image/jpeg"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/abhinayjangde/ecom/internal/config"
	"github.com/abhinayjangde/ecom/internal/db"
	"github.com/abhinayjangde/ecom/internal/lib"
	"github.com/abhinayjangde/ecom/internal/queue"
	"github.com/abhinayjangde/ecom/internal/utils"
	"github.com/disintegration/imaging"
)

const (
	maxDimension = 1200
	quality      = 80
	pollInterval = 2 * time.Second
)

func main() {
	cfg := config.MustLoad()

	logger, closer, err := utils.NewLogger(cfg.LogFile)
	if err != nil {
		slog.Error("logger initialization failed", "err", err)
		os.Exit(1)
	}
	defer closer.Close()
	slog.SetDefault(logger)

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer database.Close()

	s3c, err := lib.NewS3Client(cfg.AWSRegion, cfg.S3Bucket)
	if err != nil {
		logger.Error("s3 client initialization failed", "err", err)
		os.Exit(1)
	}

	q := queue.NewPostgres(database)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// recover jobs orphaned by a previous worker crash
	if err := q.RequeueStale(ctx); err != nil {
		logger.Error("failed to requeue stale jobs", "err", err)
		os.Exit(1)
	}

	logger.Info("worker started", "poll_interval", pollInterval.String())

	for {
		select {
		case <-ctx.Done():
			logger.Info("worker shutting down")
			return
		default:
		}

		job, err := q.ClaimNext(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				time.Sleep(pollInterval)
				continue
			}
			logger.Error("claim failed", "err", err)
			time.Sleep(pollInterval)
			continue
		}

		logger.Info("processing job",
			"job_id", job.ID, "image_id", job.ImageID, "attempt", job.Attempts)

		if err := processImage(ctx, s3c, database, job, cfg); err != nil {
			logger.Error("image processing failed", "job_id", job.ID, "err", err)
			if qerr := q.MarkFailed(ctx, job.ID, err.Error()); qerr != nil {
				logger.Error("failed to mark job failed", "job_id", job.ID, "err", qerr)
			}
			continue
		}

		if qerr := q.MarkDone(ctx, job.ID); qerr != nil {
			logger.Error("failed to mark job done", "job_id", job.ID, "err", qerr)
		}
		logger.Info("image processed", "job_id", job.ID, "image_id", job.ImageID)
	}
}

func processImage(ctx context.Context, s3c *lib.S3Client, database *sql.DB, job *queue.Job, cfg config.Config) error {
	data, err := s3c.GetObject(ctx, job.ObjectKey)
	if err != nil {
		return fmt.Errorf("get original: %w", err)
	}

	// AutoOrientation applies the JPEG EXIF orientation tag before we re-encode.
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	// Only downscale (the guard prevents upscaling small uploads).
	if img.Bounds().Dx() > maxDimension || img.Bounds().Dy() > maxDimension {
		img = imaging.Fit(img, maxDimension, maxDimension, imaging.Lanczos)
	}

	processedKey := strings.TrimSuffix(job.ObjectKey, filepath.Ext(job.ObjectKey)) + ".processed.jpg"

	// Re-encoding to JPEG drops all EXIF/metadata automatically.
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return fmt.Errorf("encode jpeg: %w", err)
	}

	if err := s3c.PutObject(ctx, processedKey, buf.Bytes(), "image/jpeg"); err != nil {
		return fmt.Errorf("put processed: %w", err)
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	if _, err := database.ExecContext(ctx, `
		UPDATE images
		SET status = 'completed',
		    processed_key = $2,
		    width = $3,
		    height = $4,
		    mime = 'image/jpeg',
		    url = $5,
		    error = NULL
		WHERE id = $1`,
		job.ImageID, processedKey, width, height, publicURL(cfg, processedKey),
	); err != nil {
		return fmt.Errorf("update image row: %w", err)
	}

	// Original no longer needed; best-effort to avoid failing a completed job.
	if err := s3c.DeleteObject(ctx, job.ObjectKey); err != nil {
		slog.Warn("failed to delete original object", "key", job.ObjectKey, "err", err)
	}

	return nil
}

func publicURL(cfg config.Config, key string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		cfg.S3Bucket, cfg.AWSRegion, strings.TrimPrefix(key, "/"),
	)
}

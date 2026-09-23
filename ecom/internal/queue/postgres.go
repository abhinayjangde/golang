package queue

import (
	"context"
	"database/sql"
	"errors"
)

const (
	MaxRetries = 3
)

type Job struct {
	ID        string
	ImageID   string
	ObjectKey string
	Mime      string
	Attempts  int
}

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// RequeueStale resets 'processing' jobs back to 'pending' (worker crash recovery).
// Call once on startup before the poll loop.
func (q *Postgres) RequeueStale(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, `
		UPDATE image_jobs
		SET status = 'pending', updated_at = NOW()
		WHERE status = 'processing'`)
	return err
}

// ClaimNext atomically grabs the oldest pending job and marks it processing.
// Returns sql.ErrNoRows when the queue is empty (worker sleeps on that).
func (q *Postgres) ClaimNext(ctx context.Context) (*Job, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var j Job
	err = tx.QueryRowContext(ctx, `
		SELECT j.id, j.image_id, i.object_key, i.mime, j.attempts
		FROM image_jobs j
		JOIN images i ON i.id = j.image_id
		WHERE j.status = 'pending'
		ORDER BY j.created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED`,
	).Scan(&j.ID, &j.ImageID, &j.ObjectKey, &j.Mime, &j.Attempts)
	if err != nil {
		return nil, err // sql.ErrNoRows when queue is empty
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE image_jobs
		SET status = 'processing', attempts = attempts + 1, updated_at = NOW()
		WHERE id = $1`, j.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &j, nil
}

// MarkDone completes a job successfully.
func (q *Postgres) MarkDone(ctx context.Context, jobID string) error {
	_, err := q.db.ExecContext(ctx, `
		UPDATE image_jobs
		SET status = 'done', last_error = NULL, updated_at = NOW()
		WHERE id = $1`, jobID)
	return err
}

// MarkFailed records a failure. Retries while attempts < MaxRetries,
// otherwise the job becomes 'failed' (requeue via RequeueStale won't touch it).
func (q *Postgres) MarkFailed(ctx context.Context, jobID string, cause string) error {
	_, err := q.db.ExecContext(ctx, `
		UPDATE image_jobs
		SET status = CASE WHEN attempts >= $2 THEN 'failed' ELSE 'pending' END,
		    last_error = $3,
		    updated_at = NOW()
		WHERE id = $1`,
		jobID, MaxRetries, cause)
	return err
}

// IsEmpty reports whether there are pending jobs (used by tests/logging).
func (q *Postgres) IsEmpty(ctx context.Context) (bool, error) {
	var one int
	err := q.db.QueryRowContext(ctx,
		`SELECT 1 FROM image_jobs WHERE status = 'pending' LIMIT 1`).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

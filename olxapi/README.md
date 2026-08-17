

This guide teaches you how to containerize Go applications using Docker.
https://docs.docker.com/guides/golang/

## automatically rebuild/reload docker-compose file
```yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    env_file:
      - .env
      
    develop:
      watch:
        # Watch the go.mod and go.sum files
        - path: ./go.mod
          action: rebuild
        - path: ./go.sum
          action: rebuild
        
        # Watch all Go source files in your project
        - path: .
          target: /app
          action: rebuild
          ignore:
            - .git/
            - .dockerignore
            - Dockerfile
```


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bin ./cmd
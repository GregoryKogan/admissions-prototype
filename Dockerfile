# Stage 1: Build the Vue.js frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app

COPY ui/package.json ui/yarn.lock ./
RUN yarn install

COPY ui/ ./
RUN yarn build

# Stage 2: Build the Go backend
FROM golang:1.23.1-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o admissions ./cmd/admissions

# Stage 3: Create the final image
FROM alpine:latest

RUN adduser -D -u 1001 app

WORKDIR /root/

COPY --from=backend-builder --chown=app:app /app/admissions ./admissions
COPY --from=frontend-builder --chown=app:app /app/dist ./ui/dist

EXPOSE 8888
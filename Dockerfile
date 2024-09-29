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

RUN go build -o server .

# Stage 3: Create the final image
FROM alpine:latest

WORKDIR /root/

COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/dist ./ui/dist

EXPOSE 8888

CMD ["./server"]
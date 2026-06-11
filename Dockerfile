FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.24-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
COPY server/*.go server/
COPY --from=frontend /app/frontend/dist ./server/dist
RUN go build -o bin/ghhealth ./cmd/ghhealth

FROM alpine:latest
WORKDIR /app
COPY --from=backend /app/bin/ghhealth ./
EXPOSE 8080
CMD ["./ghhealth", "server"]

FROM node:24-alpine AS web
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.26-alpine AS api
WORKDIR /src
COPY go.mod ./
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /server ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=api /server /app/server
COPY --from=web /src/frontend/dist /app/frontend/dist
ENV PORT=8080 STATIC_DIR=/app/frontend/dist
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/server"]

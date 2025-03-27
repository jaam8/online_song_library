FROM golang:latest AS build
WORKDIR /app
COPY . .

ENV CGO_ENABLED=0
ENV GOPROXY=proxy.golang.org
RUN cp .env.example .env
RUN go mod download
RUN go build -o ./online_song_library ./cmd/main.go

FROM alpine:latest AS run
RUN apk add --no-cache bash
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=build /app/online_song_library /online_song_library
COPY --from=build /app/db app/db
COPY .env /app/.env

WORKDIR /app
EXPOSE ${REST_PORT}
CMD ["/wait-for-it.sh", "postgres_container:5432", "--", "/online_song_library"]

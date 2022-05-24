FROM golang:1.17-alpine AS base

RUN apk add build-base

# This is for development
FROM base as dev

WORKDIR /app

RUN go get github.com/cespare/reflex
COPY reflex.conf /
ENTRYPOINT ["reflex", "-c", "/reflex.conf"]

# This is for production
FROM base as builder

ENV CGO_ENABLED=0

WORKDIR /bre4ker

COPY . .
RUN go mod download

RUN go build -o ./serve

FROM alpine

WORKDIR /src

COPY --from=builder /bre4ker/serve .
COPY --from=builder /bre4ker/migrations/ migrations/

CMD ["./serve"]
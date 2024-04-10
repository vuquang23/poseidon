# Vendor stage
FROM golang:1.21.3 as dep
WORKDIR /build
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . .
RUN go mod vendor

# Build binary stage
FROM golang:1.21.3 as build
WORKDIR /build
COPY --from=dep /build .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o app -tags nethttpomithttp2 ./cmd/app

# Minimal image
FROM alpine:latest
WORKDIR /app
COPY migrations migrations
COPY --from=build /build/app app
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates
RUN apk --no-cache add tzdata
CMD ["./app"]

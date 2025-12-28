FROM --platform=$TARGETPLATFORM golang:alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=$CGO_ENABLED go build -ldflags "-extldflags='-static'" -o ehome

FROM alpine:3.22.2
WORKDIR /app
COPY --from=builder /build/ehome /app/ehome
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/static /app/static
ENTRYPOINT ["/app/ehome"]

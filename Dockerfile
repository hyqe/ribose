FROM golang as builder
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ribose

FROM debian
WORKDIR /
COPY --from=builder /ribose .
COPY ./migrations/* /migrations/
EXPOSE 80
ENTRYPOINT ["/ribose"]
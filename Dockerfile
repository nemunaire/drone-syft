FROM golang:1-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/plugin .

FROM anchore/syft

COPY --from=builder /bin/plugin /bin/plugin
ENTRYPOINT ["/bin/plugin"]

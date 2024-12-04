FROM golang:1.23.4-alpine as build_deps

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -ldflags '-w -s' -a -installsuffix cgo -o webhook main.go

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /workspace/webhook webhook

ENTRYPOINT ["/webhook"]

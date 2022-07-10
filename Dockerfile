FROM golang:1.17.1 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /portAPI cmd/portAPI.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /portDomainService cmd/portDomainService.go

FROM alpine:3.13.6 as portAPI
WORKDIR /root/
COPY --from=builder /portAPI ./
EXPOSE 5000
ENTRYPOINT ["./portAPI", "-log-level", "debug"]

FROM alpine:3.13.6 as portDomainService
WORKDIR /root/
COPY --from=builder /portDomainService ./
EXPOSE 5000
ENTRYPOINT ["./portDomainService",  "-log-level", "debug"]

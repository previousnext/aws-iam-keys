FROM golang:1.8
ADD workspace /go
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/sshd-iam-user_linux_amd64 /usr/local/bin/sshd-iam-user
CMD ["sshd-iam-user"]

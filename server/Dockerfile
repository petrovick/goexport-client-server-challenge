#Step 1
FROM golang:1.20-alpine as BUILDER

RUN set -ex && \
    apk add --no-cache gcc musl-dev

RUN set -ex && \
    rm -f /usr/libexec/gcc/x86_64-alpine-linux-musl/6.4.0/cc1obj && \
    rm -f /usr/libexec/gcc/x86_64-alpine-linux-musl/6.4.0/lto1 && \
    rm -f /usr/libexec/gcc/x86_64-alpine-linux-musl/6.4.0/lto-wrapper && \
    rm -f /usr/bin/x86_64-alpine-linux-musl-gcj


WORKDIR /usr/src/app

RUN export CGO_ENABLED=1

COPY . .

RUN go mod tidy

RUN go env -w CGO_ENABLED=1

CMD cd /usr/src/app && \
    CGO_ENABLED=1 go run main.go

# # Step 2
# FROM golang:1.20-alpine as final

# WORKDIR /usr/src/app

# COPY --from=BUILDER /usr/src/app/goserver ./goserver

# CMD ["./goserver"]
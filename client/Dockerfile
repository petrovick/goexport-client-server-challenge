#Step 1
FROM golang:1.20-alpine as BUILDER

WORKDIR /usr/src/app

COPY go.mod main.go ./

RUN go build

CMD ["tail", "-f", "/dev/null"]

# Step 2
FROM scratch as final

WORKDIR /usr/src/app

COPY --from=BUILDER /usr/src/app/client ./client

CMD ["./client"]
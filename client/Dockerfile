#Step 1
FROM golang:1.20-alpine as BUILDER

WORKDIR /usr/src/app

COPY . .

CMD cd /usr/src/app && \
    go run main.go

# # Step 2
# FROM golang:1.20-alpine as final

# WORKDIR /usr/src/app

# COPY --from=BUILDER /usr/src/app/goserver ./goserver

# CMD ["./goserver"]
FROM golang:1.16-alpine
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .

RUN addgroup -S -g 1001 radix-non-root-group

# Add a new user "radix-non-root-user" with user id 1001 and include in group
RUN adduser -S -u 1001 -G radix-non-root-group radix-non-root-user

RUN ls -al /app

USER 1001
EXPOSE 8080
CMD ["/app/main"]

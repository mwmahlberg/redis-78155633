FROM --platform=linux/amd64 golang:1.22-alpine3.19 AS build
WORKDIR /app
COPY . .
RUN go build -o /app/myapp

FROM --platform=linux/amd64 alpine:3.19
COPY --from=build /app/myapp /usr/local/bin/myapp
ENTRYPOINT ["/usr/local/bin/myapp"]

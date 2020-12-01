FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY . .
RUN go build -o nd-lobby .

EXPOSE 8080:8080/tcp

CMD ["/app/nd-lobby"]
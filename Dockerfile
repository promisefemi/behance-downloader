FROM golang:1.18-alpine
WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/app ./web
EXPOSE 8000
CMD [ "./bin/app" ]
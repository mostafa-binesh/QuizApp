# Dockerfile
FROM golang:latest

# WORKDIR /go/src/app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

# RUN go get -d -v ./...
# RUN go install -v ./...
RUN go mod download
# COPY . . // copy all files
COPY . .
RUN go build -o /docker
EXPOSE 8070
CMD [ "/docker" ]


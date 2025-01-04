FROM golang:1.22
WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN go build -o go-auth-service .
EXPOSE 8080
CMD [ "./go-auth-service" ]

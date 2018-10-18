FROM golang

RUN mkdir /app
ADD . /app/
WORKDIR /app

# get dependencies
RUN go get -d -v 

#build the binary
RUN go build -o /app_binary

# Make 8080 available (docker run -p 8080:8080)
EXPOSE 8080

ENTRYPOINT ["/app_binary"]


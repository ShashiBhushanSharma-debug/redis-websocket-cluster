FROM golang:1.26.1-alpine

WORKDIR /app

# initialize a temporary go module and fetch the Gorilla Dependency
RUN go mod init chat && \
    go get github.com/gorilla/websocket && \
    go get github.com/redis/go-redis/v9

# Download everything listed in the files
RUN go mod download

# Copy the chat server and the html interface
COPY *.go ./
COPY index ./   

# Build the GO binary 
RUN go build -o chat-server .

# Expose the default port the chat app uses 
EXPOSE 8080

# Run the compiled binary 
CMD ["./chat-server"]
#FROM golang:alpine
FROM golang:alpine
# Set necessary environmet variables needed for our image

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN pip install unoconv
RUN apk update
#RUN apk add --upgrade mupdf-tools
RUN apk add --upgrade libreoffice
RUN apk add --upgrade libreoffice-lang-cs
RUN apk add --upgrade libreoffice-writer



# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Export necessary port
#EXPOSE 8080
#EXPOSE 9000

# Command to run when starting the container
CMD ["/dist/main"]
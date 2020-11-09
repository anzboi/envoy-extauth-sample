FROM golang:1.13.7-buster
WORKDIR /usr/src/app
COPY auth.go auth.go
CMD [ "go", "run", "auth.go" ]
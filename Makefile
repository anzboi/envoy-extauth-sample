.PHONY: up
up: bin/service
	docker-compose up

bin/service:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/service ./service.go

clean:
	rm -rf bin
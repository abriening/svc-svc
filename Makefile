build: clean format
	go build -o bin/svc-svc

format:
	go fmt ./...

clean:
	rm -f bin/*

all: build format clean
	GOOS=windows GOARCH=amd64 go build -o bin/svc-svc.exe
	GOOS=linux GOARCH=amd64 go build -o bin/svc-svc-linux
	GOOS=darwin GOARCH=amd64 go build -o bin/svc-svc-mac

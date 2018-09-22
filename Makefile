.PHONY: clean build

clean:
	rm -rf ./dist

build:
	GOOS=linux GOARCH=amd64 go build -o ./dist/main main.go
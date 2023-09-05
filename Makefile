build:
	@go build -C cmd -o ../bin/authJwt.exe

run: build
	@./bin/authJwt

test:
	@go test -v ./...
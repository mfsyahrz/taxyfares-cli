build:
	go build -o taxyfare main.go
setup:
	go mod download
test:
	go test -v ./...
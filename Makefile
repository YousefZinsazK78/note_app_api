build:
	@go build -o bin/noteapp ./main.go

run: build
	@./bin/noteapp
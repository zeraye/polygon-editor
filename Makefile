build:
	go build -o bin/polygon-editor cmd/polygon-editor/main.go 

run: build
	./bin/polygon-editor

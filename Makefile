LD_LIBRARY_PATH=lib
export LD_LIBRARY_PATH

build:
	go build -o bin/polygon-editor

run: build
	./bin/polygon-editor

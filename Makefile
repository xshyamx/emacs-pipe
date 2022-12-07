.phony: clean prepare build linux windows test

clean:
	rm -fr build

prepare:
	mkdir -p build

build: linux windows

linux: prepare
	go build -o build/ep main.go

windows: prepare
	GOOS=windows go build -o build/ep.exe main.go

test: build
	./simple.sh | ./build/ep

install: build
	cp -fv build/ep $(HOME)/.local/bin

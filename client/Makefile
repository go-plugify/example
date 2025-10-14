NOW_TIME := $(shell date +%Y%m%d%H%M%S)

run:
	rm -rf build
	mkdir -p build
	cp main.go build/
	touch build/go.mod
	echo "module plugin${NOW_TIME}" > build/go.mod
	echo "\ngo 1.23.10" >> build/go.mod
	cd build && GO111MODULE=on CGO_ENABLED=1 go build -ldflags="-s -w" -buildmode=plugin -o patch.so .
	curl --insecure -s -X POST -F "file=@build/patch.so" "http://localhost:8080/api/v1/plugin/init"
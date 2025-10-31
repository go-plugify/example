NOW_TIME := $(shell date +%Y%m%d%H%M%S)
NEW_MODULE_NAME := plugin${NOW_TIME}

run:
	rm -rf build
	mkdir -p build
	cp -r src/* build/
	rm -rf build/go.mod
	touch build/go.mod
	echo "module plugin${NOW_TIME}" > build/go.mod
	echo "\ngo 1.23.10" >> build/go.mod
	find ./build -type f \( -name "*.go" -o -name "go.mod" \) | while read file; do \
		sed -i "s/plugin[0-9]\{14\}/${NEW_MODULE_NAME}/g" $$file; \
	done
	cd build && GO111MODULE=on CGO_ENABLED=1 go build -ldflags="-s -w" -buildmode=plugin -o patch.so .
	curl --insecure -s -X POST -F "file=@build/patch.so" -F "meta={\"id\": \"example\", \"name\": \"Example Plugin\", \"description\": \"An example plugin\", \"author\": \"Your Name\", \"version\": \"v0.0.1\", \"loader\": \"native_plugin_http\"}"  "http://localhost:8080/api/v1/plugin/init" | jq
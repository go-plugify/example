
run:
	curl --insecure -s -X POST -F "file=@main.go" -F "meta={\"id\": \"example\", \"name\": \"Example Plugin\", \"description\": \"An example plugin\", \"author\": \"Your Name\", \"version\": \"v0.0.1\", \"loader\": \"yaegi_http\"}"  "http://localhost:8080/api/v1/plugin/init" | jq
route69:
	go build -o dist/route69 ./cmd/route69

run: route69
	sudo ./dist/route69
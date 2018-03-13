build:
	go build -buildmode=plugin -o dist/csv_loader.so

build-darwin:
	env GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -o dist/darwin/csv_loader.so

build-linux:
	env GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o dist/linux/csv_loader.so


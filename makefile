all: run

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"  -o paperlink .
	zip -r paperlink.zip ./tmpl ./paperlink ./makefile
	rm paperlink

run: npm
	@go run .

npm:
	@npm install

test:
	go test .

clean:
	rm paperlink*


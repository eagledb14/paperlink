all: run

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"  -o paperlink .
	zip -r paperlink.zip ./resources ./paperlink
	rm paperlink


run: 
	@go run .

clean:
	rm paperlink*


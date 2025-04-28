run:
	go run cmd/main.go

build-install:
	go build cmd/main.go
	sudo mv main /usr/local/bin/gakun

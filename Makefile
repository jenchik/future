
dependency:
	go get -u github.com/jenchik/thread
	go get -u github.com/stretchr/testify/assert

test:
	go test -v ./...

BINARY_NAME=clipp2p

build:
	go build -o ${BINARY_NAME} ./cmd/clipp2p/main.go

run:
	go build -o ${BINARY_NAME} ./cmd/clipp2p/main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test -v 

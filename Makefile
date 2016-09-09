

all:
	go build ./
	go build ./support/natsSupport

ex:
	mkdir -p bin
	go build ./examples/microblag
	go build -o bin/user-dds ./examples/cmd/user-dds
	go build -o bin/search-dds ./examples/cmd/search-dds
	go build -o bin/entry-dds ./examples/cmd/entry-dds
	go build -o bin/user-producer ./examples/cmd/user-producer
	go build -o bin/user-subscriber ./examples/cmd/user-subscriber


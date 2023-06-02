all: bin Port-Gopper Port-Gopper-Client

Port-Gopper:
	go build -o ./bin/Port-Gopper ./src/server.go

Port-Gopper-Client:
	go build -o ./bin/Port-Gopper-Client ./src/client.go

bin:
	mkdir -p ./bin

clean:
	rm -rf ./bin/*
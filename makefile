all: Port-Gopper

Port-Gopper:
	go build -o ./bin/Port-Gopper ./src/main.go

clean:
	rm -rf ./bin/*

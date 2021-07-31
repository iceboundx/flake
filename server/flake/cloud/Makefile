.PHONY: all docker clean

all:
	go build -o cloud-controller main.go

docker:
	docker build . -t flake/cloud-controller:v1.0.0

clean:
	rm -rf ./controller
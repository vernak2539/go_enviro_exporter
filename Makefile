build:
	docker run --rm -v $(shell pwd):/go/src/github.com/vernak2539/go_enviro_exporter -w /go/src/github.com/vernak2539/go_enviro_exporter docker.elastic.co/beats-dev/golang-crossbuild:1.16.6-armel-debian10 --build-cmd "go mod init;go mod tidy;go build -o enviro_exporter.linux-armv6" -p "linux/armv6"
	mv enviro_exporter.linux-armv6 $(shell pwd)/builds/
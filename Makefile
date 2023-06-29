build: build_arm6 build_arm7 build_amd64 build_arm64

build_arm6:
	env GOOS=linux GOARM=6 GOARCH=arm go build -o ./build/enviro_exporter-$(VERSION).linux-armv6 ./cmd/main.go

build_arm7:
	env GOOS=linux GOARM=7 GOARCH=arm go build -o ./build/enviro_exporter-$(VERSION).linux-armv7 ./cmd/main.go

build_amd64:
	env GOOS=linux GOARCH=amd64 go build -o ./build/enviro_exporter-$(VERSION).linux-amd64 ./cmd/main.go

build_arm64:
	env GOOS=linux GOARCH=arm64 go build -o ./build/enviro_exporter-$(VERSION).linux-arm64 ./cmd/main.go

build:
	env GOOS=linux GOARM=6 GOARCH=arm go build -o ./build/enviro_exporter-$(VERSION).linux-armv6 ./main.go
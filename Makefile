build:
	env GOOS=linux GOARM=6 GOARCH=arm go build -o ./builds/enviro_exporter.linux-armv6 ./main.go
run:
	go run cmd/main.go

build-ios:
	fyne package --os ios -appID com.gabrielross.timer-go
set shell := ["powershell.exe", "-Command"]

# build main
build:
    go build -o elvui-updater.exe main.go

# run main
run:
    go run main.go

# run test
test:
    go test -v ./...
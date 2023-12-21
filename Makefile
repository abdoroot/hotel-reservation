run:
	go run main.go
test:
	go test -v ./... --count=1
seed:
	go run scripts/seed.go
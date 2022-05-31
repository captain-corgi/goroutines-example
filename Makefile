playground:
	go run cmd/playground/main.go

zero:
	go run cmd/zero-concurency/main.go > cmd/zero-concurency/results.txt

ranging:
	go run cmd/ranging-over-channels/main.go
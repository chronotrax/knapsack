.PHONY: build
build:
	go build -o ./build/knapsack
	GOOS=linux go build -o ./build/knapsack-linux .
	GOOS=windows go build -o ./build/knapsack-windows .
	GOOS=darwin go build -o ./build/knapsack-darwin .

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./...

.PHONY: demo
demo: build
	./build/knapsack 672 13 "Bat" 3 5 9 18 38 75 155 310
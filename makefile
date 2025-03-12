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
demo:
	#./build/knapsack 672 13 "Bat" 3,5,9,18,38,75,155,310
	./build/knapsack 491 41 "Hello World!" 2,3,7,14,30,57,120,251 3
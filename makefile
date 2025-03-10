.PHONY: build
build:
	GOOS=linux go build -o ./build/knapsack-linux .
	GOOS=windows go build -o ./build/knapsack-windows .
	GOOS=darwin go build -o ./build/knapsack-darwin .

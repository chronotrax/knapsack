package main

import (
	"fmt"
	"github.com/chronotrax/knapsack/knapsack"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: generate-set [blockSize]")
		return
	}

	size, err := strconv.Atoi(os.Args[1])
	if err != nil || size <= 0 {
		fmt.Println("argument is not a positive integer:", os.Args[1])
		return
	}

	s, err := knapsack.RandomSet(size)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Superincreasing?:", s.IsSuperincreasing())
	fmt.Println("size:", len(s))

	fmt.Println(knapsack.BigIntsToStr(s))
}

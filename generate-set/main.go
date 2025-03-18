package main

import (
	"fmt"
	"github.com/chronotrax/knapsack/knapsack"
	"log"
	"os"
	"strconv"
)

func main() {
	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	s, err := knapsack.RandomSet(size)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Superincreasing?:", s.IsSuperincreasing())
	fmt.Println("size:", len(s))

	fmt.Println(knapsack.BigIntsToStr(s))
}

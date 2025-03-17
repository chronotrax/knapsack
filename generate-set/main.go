package main

import (
	"fmt"
	"github.com/chronotrax/knapsack/types"
	"log"
	"os"
	"strconv"
)

func main() {
	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	s := types.NewS(size)

	fmt.Println("Superincreasing?:", s.IsSuperincreasing())
	fmt.Println("size:", len(s))

	for i, x := range s {
		if i == len(s)-1 {
			fmt.Printf("%v\n", x)
			continue
		}
		fmt.Printf("%v,", x)
	}
}

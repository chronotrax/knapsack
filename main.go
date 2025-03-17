package main

import (
	"fmt"
	"github.com/chronotrax/knapsack/crypt"
	"github.com/chronotrax/knapsack/types"
	"os"
	"strconv"
	"strings"
)

func main() {
	var private types.PrivateKey
	var public types.PublicKey
	var data []byte
	maxK := uint64(5)

	//for i, arg := range os.Args {
	//	fmt.Printf("%d: %v\n", i, arg)
	//}

	l := len(os.Args)

	if l == 1 {
		// no args given, use random cryptosystem
		knapsack := types.RandomKnapsack(8)
		private = knapsack.Private

		public = types.NewPublicKey(knapsack.S, private)

		data = []byte("Hello World!")
	} else if l >= 5 {
		// use given the crypto system
		u, err := strconv.ParseUint(os.Args[1], 10, 64)
		if err != nil {
			fmt.Println("u is not a uint: ", err)
			return
		}

		v, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("v is not a uint: ", err)
			return
		}

		private, err = types.NewPrivateKey(u, v)
		if err != nil {
			fmt.Println("invalid u & v: ", err)
			return
		}

		maxKeys, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println("maxKeys is not a uint: ", err)
			return
		}
		maxK = maxKeys

		data = []byte(os.Args[4])

		s := make([]uint64, 0)
		for _, a := range strings.Split(os.Args[5], ",") {
			ua, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				fmt.Println("invalid s value: ", err)
				return
			}
			s = append(s, ua)
		}

		public = types.NewPublicKey(s, private)
	} else {
		// print help
		fmt.Println("usage: ./knapsack [u] [v] [maxKeys to brute force] [data to encrypt] [s1,s1,...,s8]")
		return
	}

	fmt.Printf("private key: u = %d, v = %d\n", private.U, private.V)
	fmt.Println("public key: ", public)
	fmt.Println("data: ", data, string(data))

	fmt.Println("\nencrypting...")
	cipher := crypt.Encrypt(data, public)
	fmt.Println("ciphertext: ", cipher)

	fmt.Println("\ndecrypting...")
	plain := crypt.Decrypt(cipher, private, public)

	fmt.Println("plaintext: ", plain, string(plain))
	fmt.Println("original data: ", data, string(data))

	fmt.Println("\n\nbrute forcing decryption...")
	crypt.BruteForce(cipher, public, data, maxK)
}

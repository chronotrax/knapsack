package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var private PrivateKey
	var public PublicKey
	var data []byte
	maxKeys := uint64(5)

	//for i, arg := range os.Args {
	//	fmt.Printf("%d: %v\n", i, arg)
	//}

	l := len(os.Args)

	if l == 1 {
		// no args given, use random cryptosystem
		knapsack := RandomKnapsack(8)
		private = knapsack.Private

		public = NewPublicKey(knapsack.S, private)

		data = []byte("Hello World!")
	} else if l >= 5 {
		// use given crypto system
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

		private, err = NewPrivateKey(u, v)
		if err != nil {
			fmt.Println("invalid u & v: ", err)
			return
		}

		data = []byte(os.Args[3])

		s := make([]uint64, 0)
		for _, a := range strings.Split(os.Args[4], ",") {
			ua, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				fmt.Println("invalid s value: ", err)
				return
			}
			s = append(s, ua)
		}

		public = NewPublicKey(s, private)

		m, err := strconv.ParseUint(os.Args[5], 10, 64)
		if err != nil {
			fmt.Println("m is not a uint: ", err)
			return
		}
		maxKeys = m
	} else {
		// print help
		fmt.Println("usage: ./knapsack [u] [v] [data to encrypt] [s1,s1,...,s8] [# of keys to brute force]")
		return
	}

	fmt.Printf("private key: u = %d, v = %d\n", private.U, private.V)
	fmt.Println("public key: ", public)
	fmt.Println("data: ", data, string(data))

	fmt.Println("\nencrypting...")
	cipher := Encrypt(data, public)
	fmt.Println("ciphertext: ", cipher)

	fmt.Println("\ndecrypting...")
	plain := Decrypt(cipher, private, public)

	fmt.Println("plaintext: ", plain, string(plain))
	fmt.Println("original data: ", data, string(data))

	fmt.Println("\n\nbrute forcing decryption...")
	BruteForce(cipher, public, data, maxKeys)
}

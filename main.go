package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var private PrivateKey
	var public PublicKey
	var data []byte
	maxVal := ^uint(0)

	//for i, arg := range os.Args {
	//	fmt.Printf("%d: %v\n", i, arg)
	//}

	l := len(os.Args)

	if l == 1 {
		// no args given, use random cryptosystem
		knap := RandomKnapsack()
		private = knap.Private

		public = NewPublicKey(knap.S, private)

		data = []byte("Hello World!")
	} else if l == 12 {
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

		private, err = NewPrivateKey(uint(u), uint(v))
		if err != nil {
			fmt.Println("invalid u & v: ", err)
			return
		}

		data = []byte(os.Args[3])

		s := S{}
		for i, a := range os.Args[4:] {
			ua, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				fmt.Println("invalid s value: ", err)
				return
			}
			s[i] = uint(ua)
		}

		public = NewPublicKey(s, private)
	} else {
		// print help
		fmt.Println("usage: ./knapsack [u] [v] [data to encrypt] [s1] [s2] ... [s8]")
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
	BruteForce(cipher, public, data, maxVal)
}

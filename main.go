package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var private PrivateKey
	var public PublicKey
	var data []byte
	var maxVal uint

	for i, arg := range os.Args {
		fmt.Printf("%d: %v\n", i, arg)
	}

	var err error
	if len(os.Args) == 1 {
		// use defaults
		private, err = NewPrivateKey(672, 13)
		if err != nil {
			log.Fatalln(err)
		}

		public, err = NewPublicKey(S([]uint{3, 5, 9, 18, 38, 75, 155, 310}), private)
		if err != nil {
			log.Fatalln(err)
		}

		data = []byte{byte('B'), byte('a'), byte('t')}

		maxVal = 1000
	}

	fmt.Printf("private key: u = %d, v = %d\n", private.U, private.V)
	fmt.Println("public key: ", private)
	fmt.Println("data: ", data)
	fmt.Println("encrypting...")

	cipher := Encrypt(data, public)

	fmt.Println("ciphertext: ", cipher)
	fmt.Println("decrypting...")

	plain := Decrypt(cipher, private, public)

	fmt.Println("plaintext: ", plain)
	fmt.Println("original data: ", data)

	fmt.Println("brute forcing decryption...")
	BruteForce(cipher, public, data, maxVal)
}

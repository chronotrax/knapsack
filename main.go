package main

import (
	"fmt"
	"github.com/chronotrax/knapsack/knapsack"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func main() {
	var k *knapsack.Knapsack
	var data []byte
	maxKeys := uint64(5)

	l := len(os.Args)

	if l == 1 {
		// no args given, use random cryptosystem
		fmt.Println("using a random cryptosystem")
		var err error
		k, err = knapsack.NewKnapsack(1)
		if err != nil {
			fmt.Println(err)
			return
		}

		data = []byte("Hello World!")
	} else if l >= 5 {
		// use given the crypto system
		fmt.Println("using the given cryptosystem")
		// read v argument
		v, success := new(big.Int).SetString(os.Args[1], 10)
		if !success {
			fmt.Println("v is not an integer")
			return
		}

		// read u argument
		u, success := new(big.Int).SetString(os.Args[2], 10)
		if !success {
			fmt.Println("u is not an integer")
			return
		}

		private := &knapsack.PrivateKey{
			U: u,
			V: v,
		}

		// read maxKeys argument
		maxK, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println("maxKeys is not a uint:", err)
			return
		}
		maxKeys = maxK

		// read data argument
		for _, dStr := range strings.Split(os.Args[4], ",") {
			s := strings.TrimSpace(dStr)
			if len(s) == 0 {
				continue
			}

			if len(s) != 2 {
				fmt.Println("not a hex string (0x00)", s)
				return
			}

			d, err := strconv.ParseUint(s, 16, 8)
			if err != nil {
				fmt.Println("invalid hex", s)
				return
			}
			data = append(data, byte(d))
		}

		// read set argument
		s := make([]*big.Int, 0)
		for _, sStr := range strings.Split(os.Args[5], ",") {
			si, success := new(big.Int).SetString(strings.TrimSpace(sStr), 10)
			if !success {
				fmt.Println("invalid s value:", err)
				return
			}
			s = append(s, si)
		}

		k, err = knapsack.NewKnapsackCustom(len(s)/8, private, s)
	} else {
		// print help
		fmt.Println("usage: ./knapsack [v] [u] [max # of keys to brute force] [hex string to encrypt] [s1,s1,...,s8]")
		return
	}

	if k == nil {
		fmt.Println("k == nil, something went wrong")
		return
	}

	fmt.Println("block size (in bytes): ", k.BlockSize)
	fmt.Println("block size (in bits): ", k.BlockSize*8)
	fmt.Printf("private key: v=%d, u=%d\n", k.Private.V, k.Private.U)
	fmt.Println("public key: ", k.Public)
	fmt.Println("data: ", data, string(data))

	fmt.Println("\nencrypting...")
	plain := k.NewPlaintext(data)
	cipher := k.Encrypt(plain)
	fmt.Println("ciphertext: ", cipher)

	fmt.Println("\ndecrypting...")
	newPlain, err := k.Decrypt(cipher)
	if err != nil {
		fmt.Println(err)
		return
	}
	newData := k.FromPlaintext(newPlain)

	fmt.Println("decrypted data: ", newData, string(newData))
	fmt.Println("original data: ", data, string(data))

	fmt.Println("\n\nStarting Shamir attack...")
	knapsack.Attack(k.BlockSize, cipher, k.Public, data)

	fmt.Println("\n\nbrute forcing decryption...")
	knapsack.BruteForce(k.BlockSize, cipher, k.Public, data, maxKeys)
}

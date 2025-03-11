package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"runtime"
	"slices"
	"time"
)

// Encrypt encrypts data using a PublicKey.
func Encrypt(data []byte, public PublicKey) []uint {
	result := make([]uint, 0)
	for _, d := range data {
		var sum uint = 0
		var bit byte = 128

		// loop through public key (left->right) and sum up values where bits match
		for i := 0; i < 8; i++ {
			if d&bit != 0 {
				sum += public[i]
			}
			bit = bit >> 1
		}

		result = append(result, sum)
	}
	return result
}

// Decrypt decrypts cipher using PrivateKey & PublicKey.
func Decrypt(cipher []uint, private PrivateKey, public PublicKey) []byte {
	inverse := EEA(private.V, private.U)

	// re-create superincreasing set S
	s := S{}
	for i := 0; i < 8; i++ {
		s[i] = (public[i] * inverse) % private.U
	}

	result := make([]byte, 0)
	for _, c := range cipher {
		t := (c * inverse) % private.U
		sum := byte(0)
		bit := byte(1)

		// loop through S (right->left) and convert t to binary
		for i := 7; i >= 0; i-- {
			if t >= s[i] {
				t -= s[i]
				sum += bit
			}
			bit = bit << 1
		}

		result = append(result, sum)
	}

	return result
}

// BruteForce finds private keys given cipher & PublicKey.
// expected is the original data to compare against.
// maxVal is the highest PrivateKey values to brute force until.
func BruteForce(cipher []uint, public PublicKey, expected []byte, maxVal uint) {
	u, v := uint(1), uint(1)
	t := time.Now()
	keys := make(chan PrivateKey)
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	ctx := context.Background()
	keysFound := 0
loop:
	for {
		select {
		case p := <-keys:
			fmt.Printf("found private key! u = %d, v = %d\n", p.U, p.V)
			fmt.Println("time taken: ", time.Now().Sub(t))
			keysFound++
		default:
			err := sem.Acquire(ctx, 1)
			if err != nil {
				fmt.Println(err)
			}

			//fmt.Printf("trying u = %d, v = %d\n", u, v)
			go bruteHelper(cipher, public, expected, u, v, keys, sem)

			if v < u {
				v++
			} else if u == maxVal {
				fmt.Println("max value reached")
				break loop
			} else {
				u++
				v = 1
			}
		}
	}
	fmt.Println("# of keys found: ", keysFound)
}

func bruteHelper(data []uint, public PublicKey, expected []byte, u, v uint, keys chan<- PrivateKey, sem *semaphore.Weighted) {
	test := Decrypt(data, PrivateKey{
		U: u,
		V: v,
	}, public)

	//time.Sleep(time.Duration(rand.Int()%8+2) * time.Second)

	if slices.Equal(test, expected) {
		//fmt.Println(string(test))
		keys <- PrivateKey{
			U: u,
			V: v,
		}
	}
	sem.Release(1)
	return
}

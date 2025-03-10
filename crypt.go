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

// Decrypt decrypts data using PrivateKey & PublicKey.
func Decrypt(data []uint, private PrivateKey, public PublicKey) []byte {
	inverse := EEA(private.V, private.U)

	// re-create S
	s := [8]uint{}
	for i := 0; i < 8; i++ {
		s[i] = (public[i] * inverse) % private.U
	}

	res := make([]byte, 0)
	for _, d := range data {
		t := (d * inverse) % private.U
		sum := byte(0)
		bit := byte(1)

		for j := 7; j >= 0; j-- {
			if t >= s[j] {
				t -= s[j]
				sum += bit
			}
			bit = bit << 1
		}

		res = append(res, sum)
	}

	return res
}

func BruteForce(data []uint, public PublicKey, expected []byte, maxVal uint) {
	u, v := uint(1), uint(1)
	t := time.Now()
	done := make(chan PrivateKey)
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	ctx := context.Background()
	keysFound := 0
loop:
	for {
		select {
		case p := <-done:
			fmt.Printf("found private key! u = %d, v = %d\n", p.U, p.V)
			fmt.Println("time taken: ", time.Now().Sub(t))
			keysFound++
		default:
			err := sem.Acquire(ctx, 1)
			if err != nil {
				fmt.Println(err)
			}

			//fmt.Printf("trying u = %d, v = %d\n", u, v)
			go bruteHelper(data, public, expected, u, v, done, sem)

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

func bruteHelper(data []uint, public PublicKey, expected []byte, u, v uint, done chan<- PrivateKey, sem *semaphore.Weighted) {
	test := Decrypt(data, PrivateKey{
		U: u,
		V: v,
	}, public)

	if slices.Equal(test, expected) {
		done <- PrivateKey{
			U: u,
			V: v,
		}
	}
	sem.Release(1)
	return
}

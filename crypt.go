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
func Encrypt(data []byte, public PublicKey) []uint64 {
	result := make([]uint64, 0)
	for _, d := range data {
		var sum uint64 = 0
		var bit byte = 1 << (len(public) - 1)

		// loop through public key (left->right) and sum up values where bits match
		for i := 0; i < len(public); i++ {
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
func Decrypt(cipher []uint64, private PrivateKey, public PublicKey) []byte {
	inverse := EEA(private.V, private.U)

	// re-create superincreasing set S
	s := make([]uint64, len(public))
	for i := 0; i < len(public); i++ {
		s[i] = (public[i] * inverse) % private.U
	}

	result := make([]byte, 0)
	for _, c := range cipher {
		t := (c * inverse) % private.U
		sum := byte(0)
		bit := byte(1)

		// loop through S (right->left) and convert t to binary
		for i := len(public) - 1; i >= 0; i-- {
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
// maxKeys is the max # of keys to brute force before stopping
func BruteForce(cipher []uint64, public PublicKey, expected []byte, maxKeys uint64) {
	maxUInt := ^uint64(0)
	u, v := uint64(1), uint64(1)
	keys := make(chan PrivateKey)
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	ctx := context.Background()
	keysFound := uint64(0)
	t := time.Now()
loop:
	for {
		select {
		case p := <-keys:
			fmt.Printf("found private key! u = %d, v = %d\n", p.U, p.V)
			fmt.Println("time taken: ", time.Now().Sub(t))
			keysFound++
			if keysFound >= maxKeys {
				break loop
			}
		default:
			err := sem.Acquire(ctx, 1)
			if err != nil {
				fmt.Println(err)
			}

			//fmt.Printf("trying u = %d, v = %d\n", u, v)
			go bruteHelper(cipher, public, expected, u, v, keys, sem)

			if v < u {
				v++
			} else if u == maxUInt {
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

func bruteHelper(data []uint64, public PublicKey, expected []byte, u, v uint64, keys chan<- PrivateKey, sem *semaphore.Weighted) {
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

//func Shamir(cipher []uint, public PublicKey, expected []byte, maxKeys uint) {
//
//}
//
//func Norm() {}
//
//func GS(b [8]uint) (X [8]uint, Y [8][8]uint) {
//	X = [8]uint{}
//	Y = [8][8]uint{}
//	for i := 0; i < 8; i++ {
//		Y[i] = [8]uint{}
//	}
//
//	x0 := b[0]
//	for j := 1; j < len(X); j++ {
//		xj := b[j]
//		for i := 0; i < j-1; i++ {
//
//		}
//	}
//
//	return X, Y
//}

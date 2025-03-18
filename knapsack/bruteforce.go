package knapsack

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"slices"
	"time"
)

// BruteForce finds private keys given a Ciphertext & PublicKey.
// expected is the original data to compare against.
// maxKeys is the max # of keys to brute force before stopping
func BruteForce(blockSize int, cipher Ciphertext, public PublicKey, expected []byte, maxKeys uint64) {
	maxVal := big.NewInt(math.MaxInt)
	//maxVal := big.NewInt(4)
	u := big.NewInt(1)

	validKeys := make(chan *PrivateKey)
	keysFound := uint64(0)

	threads := int64(runtime.NumCPU())
	//threads := int64(1)
	workers := make(chan *big.Int, threads)
	fmt.Printf("using %d thread(s)\n", threads)

	ctx, cancel := context.WithCancel(context.Background())

	t := time.Now()

	// ticker goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		tracker := big.NewInt(0)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t2 := time.Now().Sub(t).Round(time.Second)
				change := new(big.Int).Sub(u, tracker)
				speed := float64(change.Int64() / 10)
				fmt.Printf("time elapsed: %v, speed: (u per second) %.02f/s, currently on: u=%d\n", t2, speed, u)
				tracker = new(big.Int).Set(u)
			}
		}
	}()

	// validKeys goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-validKeys:
				fmt.Printf("time taken: %v, found private key! v=%d u=%d\n", time.Now().Sub(t), p.V, p.U)

				keysFound++
				if keysFound >= maxKeys {
					fmt.Println("max valid validKeys found")
					cancel()
				}
			}
		}
	}()

	// worker goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tryU := <-workers: // acquire a thread and a `u` value to try
				//fmt.Printf("sending v=%d, u=%d\n", u)
				k := &Knapsack{
					BlockSize: blockSize,
					Private: &PrivateKey{
						U: tryU,
					},
					Public: public,
				}

				go func() {
					worker(ctx, k, cipher, expected, validKeys) // blocks until worker completes

					<-workers // release a thread
					//fmt.Println("worker goroutine completed")
				}()
			}
		}
	}()

	working := true
	for working {
		select {
		case <-ctx.Done():
			working = false
		default:
			// send a `u` value to worker channel
			workers <- new(big.Int).Set(u)

			// if u == maxVal
			if u.Cmp(maxVal) == 0 {
				fmt.Println("max value reached")
				cancel()
				working = false
			} else {
				u.Add(u, big.NewInt(1))
			}
		}
	}
	fmt.Println("# of valid validKeys found: ", keysFound)
}

func worker(ctx context.Context, k *Knapsack, cipher Ciphertext, expected []byte, keys chan<- *PrivateKey) {
	// for v < u
	for v := big.NewInt(1); v.Cmp(k.Private.U) == -1; v.Add(v, big.NewInt(1)) {
		select {
		case <-ctx.Done():
			return
		default:
			k.Private.V = v

			//time.Sleep(time.Duration(mathRand.Int()%10+1) * time.Second)
			//k.Private.V = big.NewInt(70000)
			//k.Private.U = big.NewInt(70001)

			plain, err := k.Decrypt(cipher)
			if err != nil {
				//fmt.Printf("FAIL: error: v=%d u=%d\n", k.Private.V, k.Private.U)
				continue
			}

			data := k.FromPlaintext(plain)

			if slices.Equal(data, expected) {
				keys <- k.Private
			} // else {
			//fmt.Printf("FAIL: not equal: v=%d u=%d\n", k.Private.V, k.Private.U)
			//}
		}
	}
}

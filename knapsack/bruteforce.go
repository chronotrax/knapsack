package knapsack

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"slices"
	"sync"
	"time"
)

// BruteForce finds private keys given a Ciphertext & PublicKey.
// expected is the original data to compare against.
// maxKeys is the max # of keys to brute force before stopping
func BruteForce(blockSize int, cipher Ciphertext, public PublicKey, expected []byte, maxKeys uint64) {
	maxVal := big.NewInt(math.MaxInt)
	//maxVal := big.NewInt(4)
	u := big.NewInt(1)
	v := big.NewInt(1)
	keys := make(chan *PrivateKey)
	keysFound := uint64(0)

	threads := int64(runtime.NumCPU()) * 2
	//threads := int64(2)
	workers := make(chan struct{}, threads)
	fmt.Printf("using %d threads\n", threads)

	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	ctx, cancel := context.WithCancel(context.Background())

	ticker := time.NewTicker(3 * time.Second)
	t := time.Now()

	// ticker goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t2 := time.Now().Sub(t).Round(time.Second)
				fmt.Printf("currently on: (v=%d u=%d), time elapsed: %v\n", v, u, t2)
			}
		}
	}()

	// keys goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-keys:
				fmt.Printf("found private key! v=%d u=%d\n", p.V, p.U)
				fmt.Println("time taken: ", time.Now().Sub(t))
				keysFound++
				if keysFound >= maxKeys {
					cancel()
				}
			}
		}
	}()

	working := true
	for working {
		select {
		case <-ctx.Done():
			working = false
		case workers <- struct{}{}: // acquire a thread
			// set up a worker
			wg.Add(1)
			m.Lock() // prevent the main thread from changing v, u prematurely
			go func() {
				defer wg.Done()
				//fmt.Printf("sending v=%d, u=%d\n", v, u)
				k := &Knapsack{
					BlockSize: blockSize,
					Private: &PrivateKey{
						U: new(big.Int).Set(u),
						V: new(big.Int).Set(v),
					},
					Public: public,
				}

				m.Unlock() // done working with v, u

				worker(k, cipher, expected, keys) // blocks until worker completes

				<-workers // release a thread
			}()

			m.Lock() // change v, u
			// if v < u
			if v.Cmp(u) == -1 {
				v.Add(v, big.NewInt(1))
			} else if u.Cmp(maxVal) == 0 {
				fmt.Println("max value reached")
				cancel()
				working = false
			} else {
				u.Add(u, big.NewInt(1))
				v = big.NewInt(1)
			}
			m.Unlock() // done changing v, u
		}
	}
	wg.Wait()
	fmt.Println("# of keys found: ", keysFound)
}

func worker(k *Knapsack, cipher Ciphertext, expected []byte, keys chan<- *PrivateKey) {
	//time.Sleep(time.Duration(mathRand.Int()%10+1) * time.Second)

	plain, err := k.Decrypt(cipher)
	if err != nil {
		//fmt.Printf("FAIL: error: v=%d u=%d\n", k.Private.V, k.Private.U)
		return
	}

	data := k.FromPlaintext(plain)

	if slices.Equal(data, expected) {
		keys <- k.Private
	} else {
		//fmt.Printf("FAIL: not equal: v=%d u=%d\n", k.Private.V, k.Private.U)
	}
	return
}

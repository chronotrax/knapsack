package knapsack

import (
	"fmt"
	"math"
	"math/big"
	"slices"
)

// gs is the Gram–Schmidt algorithm.
func gs(b, mu, B []float64, n int) {
	bStar := make([]float64, n*n)

	for i := 0; i < n; i++ {
		// b*_i = b_i
		for j := 0; j < n; j++ {
			bStar[j*n+i] = b[j*n+i]
		}

		// Initialize a term to keep track of the summation term
		// in the gram schmidt algorithm
		summand := make([]float64, n)

		// Calculate mu_ij and summand value
		for j := 0; j < i; j++ {
			// mu_ij = <b_i, b*_j>/<b*_j, b*_j>
			var numeratorMu float64
			for k := 0; k < n; k++ {
				numeratorMu = numeratorMu + b[k*n+i]*bStar[k*n+j]
			}
			mu[i*n+j] = numeratorMu / B[j]

			for k := 0; k < n; k++ {
				summand[k] = summand[k] + mu[i*n+j]*bStar[k*n+j]
			}
		}

		// Calculate new b_star values and B values
		for j := 0; j < n; j++ {
			bStar[j*n+i] = bStar[j*n+i] - summand[j]
		}

		var dot float64
		for j := 0; j < n; j++ {
			dot = dot + bStar[j*n+i]*bStar[j*n+i]
		}
		B[i] = dot
	}
}

// lll is the Lenstra–Lenstra–Lovász lattice basis reduction algorithm
func lll(b []float64, n int, lc, uc float64) {
	mu := make([]float64, n*n)
	B := make([]float64, n)

	gs(b, mu, B, n)

	k := 1
	for k < n {
		// Step 1 of LLL, achieve m_(k,k-1) <= lc
		// This is condition (1.18) from LLL
		if mu[k*n+k-1] > lc {
			r := math.Round(mu[k*n+k-1])

			// replace b_k by b_k - r*b_(k-1)
			for i := 0; i < n; i++ {
				b[k+i*n] = b[k+i*n] - r*b[k-1+i*n]
			}

			// Change the mu's accordingly
			for i := 0; i < k-1; i++ {
				mu[k*n+i] = mu[k*n+i] - r*mu[(k-1)*n+i]
			}
			mu[k*n+k-1] = mu[k*n+k-1] - r
		}

		// This completes step one of the algorithm
		// Step 2
		// Case 1: B_k + (mu(k,k-1)^2)B_(k-1) < uc*B_(k-1)
		left := B[k] + mu[k*n+k-1]*mu[k*n+k-1]*B[k-1]
		right := uc * B[k-1]

		if left < right && k > 0 {
			// First swap b_k and b_(k-1)
			for i := 0; i < n; i++ {
				temp := b[i*n+k]
				b[i*n+k] = b[i*n+k-1]
				b[i*n+k-1] = temp
			}

			// We need to save three values for the rest of Step 2 case 1, these
			// are B_k, B_k + mu_(k,k-1)*B[k-1], and mu_(k,k-1)
			BTemp := B[k]
			C := B[k] + mu[k*n+k-1]*mu[k*n+k-1]*B[k-1]
			muTemp := mu[k*n+k-1]

			// Now we continue with the algorithm, first adjust B
			mu[k*n+k-1] = muTemp * B[k-1] / C
			B[k] = B[k-1] * B[k] / C
			B[k-1] = C

			// All other B values stay the same
			// Next we adjust mu
			for i := k + 1; i < n; i++ {
				temp := mu[i*n+k-1]
				mu[i*n+k-1] = mu[i*n+k-1]*mu[k*n+k-1] + mu[i*n+k]*BTemp/C
				mu[i*n+k] = temp - mu[i*n+k]*muTemp
			}

			for i := 0; i < k-1; i++ {
				temp := mu[(k-1)*n+i]
				mu[(k-1)*n+i] = mu[k*n+i]
				mu[k*n+i] = temp
			}

			// All other mu values stay the same.
			// Decrement k
			k = k - 1
			// This concludes step 2 case 1
		} else { // Case 2: B_k + (mu(k,k-1)^2)B_(k-1) >= uc*B_(k-1) or k == 0
			l := k
			for l > 0 {
				l = l - 1
				if math.Abs(mu[k*n+l]) > lc {
					r := math.Round(mu[k*n+l])

					// b_k = b_k - r*b_l
					for i := 0; i < n; i++ {
						b[i*n+k] = b[i*n+k] - r*b[i*n+l]
					}

					for i := 0; i < l; i++ {
						mu[k*n+i] = mu[k*n+i] - r*mu[l*n+i]
					}

					mu[k*n+l] = mu[k*n+l] - r
					// The other mu are unchanged
					l = k
				}
			}
			k = k + 1
		}
	}
}

// checkColumn checks if the c column of Matrix m is in the correct form:
// 0 to n-2 are all 0's or 1's, and n-1 is a 0.
func checkColumn(c []float64, col, n int) bool {
	// 0 to n-2 must be a 0 or 1
	for i := 0; i < n-1; i++ {
		t := c[col*n+i]
		// if t != 0 and t != 1
		if t != 0 && t != 1 {
			return false
		}
	}

	// the last value must be == 0
	if c[n*(n-1)+col] != 0 {
		return false
	}

	return true
}

func Attack(blockSize int, cipher Ciphertext, public PublicKey, expected []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Attack had an panic:", r)
		}
	}()

	k := &Knapsack{
		BlockSize: blockSize,
		Public:    public,
	}

	// n is 1 larger than the original block n
	n := len(public) + 1
	m := make([]float64, n*n)

	// make 0 to n-1 an identity matrix
	for i := 0; i < n-1; i++ {
		m[i*n+i] = 1
	}

	// make bottom row (0 to n-1) the public key
	for i := 0; i < n-1; i++ {
		m[n*(n-1)+i] = float64(public[i].Int64())
	}

	// make bottom right corner -1*cipher
	m[n*n-1] = float64(cipher[0].Int64() * -1)

	fmt.Println("initial matrix:")
	printM(m, n)

	lll(m, n, 0.5, 0.75)

	fmt.Println("reduced matrix")
	printM(m, n)

	found := false
	for i := 0; i < n; i++ {
		if checkColumn(m, i, n) {
			found = true
			fmt.Printf("suspected plaintext found at column %d\n", i)

			plain := make(Plaintext, len(public))
			for j := 0; j < len(public); j++ {
				plain[j] = new(big.Int).SetInt64(int64(m[i+j*n]))
			}

			data := k.FromPlaintext(plain)
			if slices.Equal(data, expected) {
				fmt.Println("suspected plaintext matches original! :D")
			} else {
				fmt.Println("but it does NOT match the original plaintext :(")
			}
		}
	}

	if !found {
		fmt.Println("no suspected plaintext found in reduced matrix :(")
	}
}

func printM(m []float64, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%v ", m[i*n+j])
		}
		fmt.Println("")
	}
}

package knapsack

import (
	"fmt"
	"github.com/chronotrax/knapsack/matrix"
	"math/big"
)

// copySlice copies s to a new []*big.Int to prevent modifying the original slice.
func copySlice(s []*big.Int) []*big.Int {
	m := make([]*big.Int, len(s))
	for i := 0; i < len(s); i++ {
		m[i] = new(big.Int).Set(s[i])
	}
	return m
}

// GS is the Gram–Schmidt algorithm.
func GS(b matrix.Matrix) (matrix.Matrix, error) {
	size := len(b)
	x := matrix.NewMatrixEmpty(size, size)

	// x0 = b0
	x[0] = copySlice(b[0])

	// for j = 1 to n
	for j := 1; j < size; j++ {
		// xj = bj
		x[j] = copySlice(b[j])

		// for i = 0 to j - 1
		for i := 0; i < j-1; i++ {
			// xi * bj
			prod, err := matrix.DotProduct(x[i], b[j])
			if err != nil {
				return nil, err
			}
			y := new(big.Rat).SetInt(prod)

			// ||xi||^2
			prod, err = matrix.DotProduct(x[i], x[i])
			if err != nil {
				return nil, err
			}

			// yij = (xi * bj) / ||xi||^2
			y.Quo(y, new(big.Rat).SetInt(prod))

			// xj = xj - yij * xi
			sub, err := matrix.SliceSubtract(x[j], matrix.MultiplyRatOnSlice(x[i], y))
			x[j] = sub
		}
	}

	return x, nil
}

func LLL(b matrix.Matrix, delta *big.Rat, maxIterations int) (matrix.Matrix, error) {
	j := 1
	n := len(b)
	iter := 0

	x, err := GS(b)
	if err != nil {
		return nil, err
	}

	// for j = 1 to n
	for j < n && iter < maxIterations {
		iter++

		// for i = j - 1 to 0
		for i := j - 1; i >= 0; i-- {
			prod, err := matrix.DotProduct(x[i], b[j])
			if err != nil {
				return nil, err
			}
			q := new(big.Rat).SetInt(prod)

			prod, err = matrix.DotProduct(x[i], x[i])
			if err != nil {
				return nil, err
			}
			q.Quo(q, new(big.Rat).SetInt(prod))
			roundQ := new(big.Int).Set(q.Num())
			roundQ.Div(roundQ, q.Denom())
			sub, err := matrix.SliceSubtract(b[j], matrix.MultiplyIntOnSlice(b[i], roundQ))
			b[j] = sub

			x, err = GS(b)
			if err != nil {
				return nil, err
			}
		}

		prod, err := matrix.DotProduct(x[j-1], x[j-1])
		if err != nil {
			return nil, err
		}
		d := new(big.Rat).Mul(delta, new(big.Rat).SetInt(prod))

		prod, err = matrix.DotProduct(x[j], x[j])
		if err != nil {
			return nil, err
		}
		v := new(big.Rat).SetInt(prod)

		prod, err = matrix.DotProduct(x[j-1], x[j])
		if err != nil {
			return nil, err
		}
		v.Add(v, new(big.Rat).SetInt(prod))

		if d.Cmp(v) < 0 {
			b[j], b[j-1] = b[j-1], b[j]
			j = max(j-1, 1)
			x, err = GS(b)
			if err != nil {
				return nil, err
			}
		} else {
			j++
		}
	}

	return b, nil
}

func Attack(blockSize int, cipher Ciphertext, public PublicKey, expected []byte) error {
	size := len(public) + 1
	m := matrix.NewMatrixEmpty(size, size)

	for i := 0; i < size-1; i++ {
		m[i][i] = big.NewInt(1)
	}

	for i := 0; i < size-1; i++ {
		m[size-1][i] = public[i]
	}

	m[size-1][size-1] = new(big.Int).Mul(cipher[0], big.NewInt(-1))

	fmt.Printf("initial matrix:\n%s\n\n", m)

	reduced, err := LLL(m, big.NewRat(3, 4), 1000)
	if err != nil {
		return err
	}

	fmt.Printf("reduced matrix:\n%s\n\n", reduced)

	s := make([]*big.Int, size)
	for i, v := range reduced {
		s[i] = v[i]
	}

	fmt.Println("suspected private set:", s)
	return nil
}

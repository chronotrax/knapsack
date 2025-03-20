package knapsack

import (
	"fmt"
	"github.com/chronotrax/knapsack/matrix"
	"math/big"
	"slices"
)

// gs is the Gram–Schmidt algorithm.
func gs(b matrix.Matrix) (x, y matrix.Matrix) {
	n := b.Height()
	x = matrix.NewMatrixEmpty(n, n)
	y = matrix.NewMatrixEmpty(n, n)

	// x0 = b0 (b0 is M's first column, not row)
	x.SetCol(0, b.Col(0))

	// for j = 1 to n
	for j := 1; j < n; j++ {
		// xj = bj
		x.SetCol(j, b.Col(j))

		// for i = 0 to j - 1 (inclusive)
		for i := 0; i <= j-1; i++ {
			// xi * bj
			prod := new(big.Rat).Set(matrix.DotProduct(x.Col(i), b.Col(j)))

			// yij = (xi * bj) / ||xi||^2
			co := new(big.Rat).Quo(prod, matrix.DotProduct(x.Col(i), x.Col(i)))
			// we must keep track of coefficients for LLL
			y[j][i] = co

			// xj = xj - yij * xi
			x.SetCol(j, matrix.SubtractVec(x.Col(j), matrix.MulRatOnVec(co, x.Col(i))))
		}
	}

	return x, y
}

// lll is the Lenstra–Lenstra–Lovász lattice basis reduction algorithm
func lll(b matrix.Matrix, delta *big.Rat, maxIterations int) matrix.Matrix {
	// matrix is n by n square
	n := b.Height()

	// (X,Y) = GS(M)
	x, y := gs(b)
	//fmt.Printf("x after first GS:\n%s\n\n", x)
	//fmt.Printf("y after first GS:\n%s\n\n", y)

	// since we cannot run forever, go until maxIterations
	for iter := 1; iter <= maxIterations; iter++ {
		//fmt.Println("ITERATION", iter)

		// for j = 1 to n
		for j := 1; j < n; j++ {

			// for i = j - 1 to 0
			for i := j - 1; i >= 0; i-- {
				// |yij|
				abs := new(big.Rat).Abs(y[j][i])

				// if |yij| > 1/2
				if abs.Cmp(big.NewRat(1, 2)) > 0 {
					// yij + 1/2
					sum := new(big.Rat).Add(y[j][i], big.NewRat(1, 2))

					// floor(yij + 1/2)
					flo := new(big.Int).Quo(sum.Num(), sum.Denom())

					// bj - floor(yij + 1/2) * bi
					dif := matrix.SubtractVec(b.Col(j), matrix.MulRatOnVec(new(big.Rat).SetInt(flo), b.Col(i)))

					// bj = bj - floor(yij + 1/2) * bi
					b.SetCol(j, dif)
				}
			}
		}

		//fmt.Printf("x after first half of LLL:\n%s\n\n", x)
		//fmt.Printf("y after first half of LLL:\n%s\n\n", y)

		// (X,Y) = GS(M)
		x, y = gs(b)
		//fmt.Printf("x after second GS:\n%s\n\n", x)
		//fmt.Printf("y after second GS:\n%s\n\n", y)

		// for j = 0 to n - 1
		for j := 0; j < n-1; j++ {
			// 3/4 * ||xj||^2
			right := new(big.Rat).Mul(delta, matrix.DotProduct(x.Col(j), x.Col(j)))

			// yj,j+1 * xj
			prod := matrix.MulRatOnVec(y[j+1][j], x.Col(j))

			// xj+1 + yj,j+1 * xj
			sum := make(matrix.Vector, n)
			for i := 0; i < len(prod); i++ {
				sum[i] = new(big.Rat).Add(prod[i], x[i][j+1])
			}

			// ||xj+1 + yj,j+1 * xj||^2
			left := matrix.DotProduct(sum, sum)

			// if ||xj+1 + yj,j+1 * xj||^2 < 3/4 * ||xj||^2
			if left.Cmp(right) < 0 {
				// swap(bj, bj+1)
				temp := b.Col(j)
				b.SetCol(j, b.Col(j+1))
				b.SetCol(j+1, temp)
				break
			}
		}

		//fmt.Printf("x after second half of LLL:\n%s\n\n", x)
		//fmt.Printf("y after second half of LLL:\n%s\n\n", y)
	}

	return b
}

// checkColumn checks if the c column of Matrix m is in the correct form:
// 0 to n-2 are all 0's or 1's, and n-1 is a 0.
func checkColumn(c matrix.Vector) bool {
	// 0 to n-2 must be a 0 or 1
	for i := 0; i < len(c)-1; i++ {
		t := new(big.Int).Set(c[i].Num())
		// if t != 0 and t != 1
		if t.Cmp(big.NewInt(0)) != 0 && t.Cmp(big.NewInt(1)) != 0 {
			return false
		}
	}

	// the last value must be == 0
	if c[len(c)-1].Num().Cmp(big.NewInt(0)) != 0 {
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

	// size is 1 larger than the original block size
	size := len(public) + 1
	m := matrix.NewMatrixEmpty(size, size)

	// make 0 to n-1 an identity matrix
	for i := 0; i < size-1; i++ {
		m[i][i] = big.NewRat(1, 1)
	}

	// make bottom row (0 to n-1) the public key
	for i := 0; i < size-1; i++ {
		m[size-1][i] = new(big.Rat).SetInt(public[i])
	}

	// make bottom right corner -1*cipher
	m[size-1][size-1] = new(big.Rat).Mul(new(big.Rat).SetInt(cipher[0]), big.NewRat(-1, 1))

	fmt.Printf("initial matrix:\n%s\n\n", m)

	reduced := lll(m, big.NewRat(3, 4), 1000)

	fmt.Printf("reduced matrix:\n%s\n\n", reduced)

	found := false
	for i := 0; i < size; i++ {
		col := reduced.Col(i)
		if checkColumn(col) {
			found = true
			fmt.Printf("suspected plaintext found at column %d: %v\n", i, col)

			plain := make(Plaintext, len(public))
			for j := 0; j < len(public); j++ {
				plain[j] = new(big.Int).Set(col[j].Num())
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

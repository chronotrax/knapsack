package matrix

import (
	"fmt"
	"math/big"
	"strings"
)

// Vector is a slice of *big.Int.
// Represents a row or column of a Matrix.
type Vector []*big.Rat

func (v Vector) String() string {
	s := strings.Builder{}
	for i := 0; i < len(v); i++ {
		if v[i].IsInt() {
			s.WriteString(v[i].Num().String())
		} else {
			s.WriteString(v[i].String())
		}

		if i != len(v)-1 {
			s.WriteString(" ")
		}
	}

	return s.String()
}

// SubtractVec subtracts b's values from a (a-b) into a new Vector.
func SubtractVec(a, b Vector) Vector {
	if len(a) != len(b) {
		panic(fmt.Sprintf("slice lengths don't match. a=%d, b=%d", len(a), len(b)))
	}

	res := make(Vector, len(a))

	for i := range a {
		res[i] = new(big.Rat).Sub(a[i], b[i])
	}

	return res
}

// MulRatOnVec multiplies r onto v in a new Vector.
func MulRatOnVec(r *big.Rat, v Vector) Vector {
	res := make(Vector, len(v))

	for i, val := range v {
		res[i] = new(big.Rat).Mul(val, r)
	}

	return res
}

// DotProduct takes the dot product of a and b into a new *big.Rat.
func DotProduct(a, b Vector) *big.Rat {
	if len(a) != len(b) {
		panic(fmt.Sprintf("slice lengths don't match. a=%d, b=%d", len(a), len(b)))
	}

	sum := new(big.Rat)

	for i := range a {
		product := new(big.Rat).Mul(a[i], b[i])
		sum.Add(sum, product)
	}

	return sum
}

// Matrix of the form [y][x], where [0][0] is the top left corner, expanding down and right.
// Height is along the y value, Width is along the x value.
type Matrix []Vector

func NewMatrixEmpty(height, width int) Matrix {
	m := make([]Vector, height)

	for c := 0; c < height; c++ {
		m[c] = make(Vector, width)
		for r := 0; r < width; r++ {
			m[c][r] = new(big.Rat)
		}
	}

	return m
}

// NewMatrixFull fills the Matrix with copies of v's values.
func NewMatrixFull(height, width int, v Vector) Matrix {
	if len(v) != height*width {
		panic(fmt.Sprintf("array length != height * width (%d != %d)", len(v), height*width))
	}

	i := 0
	m := make([]Vector, height)

	for c := 0; c < height; c++ {
		m[c] = make(Vector, width)
		for r := 0; r < width; r++ {
			m[c][r] = new(big.Rat).Set(v[i])
			i++
		}
	}

	return m
}

func (m Matrix) Height() int {
	return len(m)
}

func (m Matrix) Width() int {
	return len(m[0])
}

// Row returns a copy of the row r.
func (m Matrix) Row(r int) Vector {
	v := make(Vector, m.Width())
	for i := 0; i < m.Width(); i++ {
		v[i] = new(big.Rat).Set(m[r][i])
	}
	return v
}

// Col returns a copy of the column c.
func (m Matrix) Col(c int) Vector {
	v := make(Vector, m.Height())
	for i := 0; i < m.Height(); i++ {
		v[i] = new(big.Rat).Set(m[i][c])
	}
	return v
}

// SetRow sets values at row r to a copy of v's values.
func (m Matrix) SetRow(r int, v Vector) {
	for i := 0; i < m.Width(); i++ {
		m[r][i] = new(big.Rat).Set(v[i])
	}
}

// SetCol sets values at column c to a copy of v's values.
func (m Matrix) SetCol(c int, v Vector) {
	for i := 0; i < m.Height(); i++ {
		m[i][c] = new(big.Rat).Set(v[i])
	}
}

func (m Matrix) String() string {
	s := strings.Builder{}
	for c := 0; c < m.Height(); c++ {
		for r := 0; r < m.Width(); r++ {
			if m[c][r].IsInt() {
				s.WriteString(m[c][r].Num().String())
			} else {
				s.WriteString(m[c][r].String())
			}

			if r != m.Width()-1 {
				s.WriteString("\t")
			}
		}
		if c != m.Height()-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

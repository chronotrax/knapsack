package matrix

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type Matrix [][]*big.Int

func NewMatrixEmpty(col, row int) Matrix {
	m := make([][]*big.Int, col)

	for c := 0; c < col; c++ {
		m[c] = make([]*big.Int, row)
		for r := 0; r < row; r++ {
			m[c][r] = new(big.Int)
		}
	}

	return m
}

func NewMatrixFull(col, row int, arr []*big.Int) (Matrix, error) {
	if len(arr) != col*row {
		return nil, fmt.Errorf("array length != col * row (%d != %d)", len(arr), col*row)
	}

	i := 0
	m := make([][]*big.Int, col)

	for c := 0; c < col; c++ {
		m[c] = make([]*big.Int, row)
		for r := 0; r < row; r++ {
			m[c][r] = new(big.Int).Set(arr[i])
			i++
		}
	}

	return m, nil
}

func (m Matrix) ColLen() int {
	return len(m)
}

func (m Matrix) RowLen() int {
	return len(m[0])
}

func (m Matrix) String() string {
	s := strings.Builder{}
	for c := 0; c < m.ColLen(); c++ {
		for r := 0; r < m.RowLen(); r++ {
			s.WriteString(m[c][r].String())
			if r != m.RowLen()-1 {
				s.WriteString("\t")
			}
		}
		if c != m.ColLen()-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

func Add(a, b Matrix) (Matrix, error) {
	if a.ColLen() != b.ColLen() || a.RowLen() != b.RowLen() {
		return nil, errors.New("matrix dimensions do not match")
	}

	m := NewMatrixEmpty(a.ColLen(), a.RowLen())

	for c := 0; c < a.ColLen(); c++ {
		for r := 0; r < a.RowLen(); r++ {
			m[c][r] = m[c][r].Add(a[c][r], b[c][r])
		}
	}

	return m, nil
}

func Mul(a, b Matrix) (Matrix, error) {
	if a.ColLen() != b.RowLen() || a.RowLen() != b.ColLen() {
		return nil, errors.New("matrix dimensions not compatible")
	}

	shortCol := min(a.ColLen(), b.ColLen())
	shortRow := min(a.RowLen(), b.RowLen())

	m := NewMatrixEmpty(shortCol, shortRow)

	mc := 0
	mr := 0
	for c := 0; c < shortCol; c++ {
		for r := 0; r < shortRow; r++ {
			for i := 0; i < a.RowLen(); i++ {
				x := new(big.Int).Mul(a[c][i], b[i][r])
				m[mc][mr] = x.Add(m[mc][mr], x)
			}
			if mr == shortRow-1 {
				mr = 0
				mc++
			} else {
				mr++
			}
		}
	}

	return m, nil
}

func DotProduct(a, b []*big.Int) (*big.Int, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("slice lengths don't match. a=%d, b=%d", len(a), len(b))
	}

	r := new(big.Int)

	for i := range a {
		x := new(big.Int).Mul(a[i], b[i])
		r.Add(r, x)
	}

	return r, nil
}

func SliceSubtract(a, b []*big.Int) ([]*big.Int, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("slice lengths don't match. a=%d, b=%d", len(a), len(b))
	}

	r := make([]*big.Int, len(a))

	for i := range a {
		r[i] = new(big.Int).Sub(a[i], b[i])
	}

	return r, nil
}

func MultiplyIntOnSlice(s []*big.Int, i *big.Int) []*big.Int {
	r := make([]*big.Int, len(s))

	for j, val := range s {
		r[j] = new(big.Int).Mul(val, i)
	}

	return r
}

func MultiplyRatOnSlice(s []*big.Int, rat *big.Rat) []*big.Int {
	r := make([]*big.Int, len(s))

	for i, val := range s {
		r[i] = new(big.Rat).Mul(new(big.Rat).SetInt(val), rat).Num()
	}

	return r
}

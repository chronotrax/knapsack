package types

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/chronotrax/knapsack/help"
)

type Plaintext []byte

type Ciphertext []uint64

// S is a superincreasing set.
type S []uint64

func (s S) IsSuperincreasing() bool {
	sum := s[0]
	for _, si := range s[1:] {
		if si <= sum {
			return false
		}
		sum += si
	}
	return true
}

func NewS(size int) S {
	a := make([]uint64, size)
	a[0] = 1
	sum := uint64(0)
	for i := 1; i < size; i++ {
		sum += a[i-1]
		a[i] = sum + 1
	}

	return a
}

func RandomS(size int) S {
	result := S(make([]uint64, size))
	r := make([]byte, size)
	for {
		_, _ = rand.Reader.Read(r)
		result[0] = binary.BigEndian.Uint64(r)
		for i := 1; i < size; i++ {
			_, _ = rand.Reader.Read(r)
			result[i] = binary.BigEndian.Uint64(r) + result[i-1]
		}

		// s should be superincreasing & largest # is < max uint / 2 (leave room for PrivateKey.U)
		if result.IsSuperincreasing() && result[size-1] < (^uint64(0)/2) {
			break
		}
	}
	return result
}

// PublicKey requires that:
//
// PublicKey[i] = (PrivateKey.V * S[i]) % PrivateKey.U
type PublicKey []uint64

func NewPublicKey(s S, private PrivateKey) PublicKey {
	result := make([]uint64, len(s))
	for i, si := range s {
		result[i] = (si * private.V) % private.U
	}

	return result
}

// PrivateKey requires that:
//
// 1. GCD(U, V) = 1
//
// 2. U > 2*Sn (where Sn is the largest value in the PublicKey)
type PrivateKey struct {
	U uint64
	V uint64
}

func NewPrivateKey(u, v uint64) (PrivateKey, error) {
	if help.GCD(v, u) != 1 {
		return PrivateKey{}, fmt.Errorf("GCD(%d, %d) != 1", u, v)
	}

	return PrivateKey{
		U: u,
		V: v,
	}, nil
}

func RandomPrivateKey(s S) PrivateKey {
	r := make([]byte, len(s))
	for {
		// generate random u > 2*Sn
		var u uint64
		for {
			_, _ = rand.Reader.Read(r)
			u = binary.BigEndian.Uint64(r)
			if u > 2*s[len(s)-1] {
				break
			}
		}
		_, _ = rand.Reader.Read(r)
		v := binary.BigEndian.Uint64(r)
		private, err := NewPrivateKey(u, v)
		if err != nil {
			continue
		}

		return private
	}
}

type Knapsack struct {
	S       S
	Private PrivateKey
}

func RandomKnapsack(size int) Knapsack {
	// generate random superincreasing set S
	s := RandomS(size)

	// generate random large private key
	p := RandomPrivateKey(s)

	return Knapsack{
		S:       s,
		Private: p,
	}
}

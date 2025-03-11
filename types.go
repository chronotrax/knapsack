package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// S is a superincreasing set.
type S [8]uint

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

func RandomS() S {
	result := S{}
	r := make([]byte, 8)
	for {
		rand.Reader.Read(r)
		result[0] = uint(binary.BigEndian.Uint64(r))
		for i := 1; i < 8; i++ {
			rand.Reader.Read(r)
			result[i] = uint(binary.BigEndian.Uint64(r)) + result[i-1]
		}
		// s should be superincreasing & largest # is < max uint / 2 (leave room for PrivateKey.U)
		if result.IsSuperincreasing() && result[7] < (^uint(0)/2) {
			break
		}
	}
	return result
}

// PublicKey requires that:
//
// PublicKey[i] = (PrivateKey.V * S[i]) % PrivateKey.U
type PublicKey [8]uint

func NewPublicKey(s S, private PrivateKey) PublicKey {
	result := [8]uint{}
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
	U uint
	V uint
}

func NewPrivateKey(u, v uint) (PrivateKey, error) {
	if GCD(v, u) != 1 {
		return PrivateKey{}, fmt.Errorf("GCD(%d, %d) != 1", u, v)
	}

	return PrivateKey{
		U: u,
		V: v,
	}, nil
}

func RandomPrivateKey(s S) PrivateKey {
	r := make([]byte, 8)
	for {
		// generate random u > 2*Sn
		var u uint
		for {
			rand.Reader.Read(r)
			u = uint(binary.BigEndian.Uint64(r))
			if u > 2*s[7] {
				break
			}
		}
		rand.Reader.Read(r)
		v := uint(binary.BigEndian.Uint64(r))
		priv, err := NewPrivateKey(u, v)
		if err != nil {
			continue
		}

		return priv
	}
}

type Knapsack struct {
	S       S
	Private PrivateKey
}

func RandomKnapsack() Knapsack {
	// generate random superincreasing set S
	s := RandomS()

	// generate random large private key
	p := RandomPrivateKey(s)

	return Knapsack{
		S:       s,
		Private: p,
	}
}

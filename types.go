package main

import "fmt"

// S is a superincreasing set.
type S [8]uint

func (s S) IsSuperincreasing() bool {
	sum := s[0]
	for i := 1; i < 8; i++ {
		if s[i] <= sum {
			return false
		}
		sum += s[i]
	}
	return true
}

// PublicKey requires that:
//
// PublicKey[i] = (PrivateKey.V * S[i]) % PrivateKey.U
type PublicKey [8]uint

func NewPublicKey(s S, private PrivateKey) (PublicKey, error) {
	result := [8]uint{}
	for i, val := range s {
		result[i] = (val * private.V) % private.U
	}

	if private.U <= 2*result[7] {
		return PublicKey{}, fmt.Errorf("U (%d) <= 2*Sn (%d)", private.U, 2*result[7])
	}

	return result, nil
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

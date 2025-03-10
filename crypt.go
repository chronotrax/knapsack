package main

// Encrypt encrypts data using a PublicKey.
func Encrypt(data []byte, public PublicKey) []uint {
	result := make([]uint, 0)
	for _, d := range data {
		var sum uint = 0
		var bit byte = 128

		// loop through public key (left->right) and sum up values where bits match
		for i := 0; i < 8; i++ {
			if d&bit != 0 {
				sum += public[i]
			}
			bit = bit >> 1
		}

		result = append(result, sum)
	}
	return result
}

// Decrypt decrypts data using PrivateKey & PublicKey.
func Decrypt(data []uint, private PrivateKey, public PublicKey) []byte {
	inverse := EEA(private.V, private.U)

	// re-create S
	s := [8]uint{}
	for i := 0; i < 8; i++ {
		s[i] = (public[i] * inverse) % private.U
	}

	res := make([]byte, 0)
	for _, d := range data {
		t := (d * inverse) % private.U
		sum := byte(0)
		bit := byte(1)

		for j := 7; j >= 0; j-- {
			if t >= s[j] {
				t -= s[j]
				sum += bit
			}
			bit = bit << 1
		}

		res = append(res, sum)
	}

	return res
}

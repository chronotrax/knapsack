package main

type pubKey [8]uint

type privKey struct {
	u uint
	v uint
}

func encrypt(data []byte, pub pubKey) []uint {
	res := make([]uint, 0)
	for _, d := range data {
		var sum uint = 0
		var bit byte = 128
		for i := 0; i < 8; i++ {
			if d&bit != 0 {
				sum += pub[i]
			}
			bit /= 2
		}
		res = append(res, sum)
	}
	return res
}

func decrypt(data []uint, priv privKey) []byte {
	return []byte{}
}

// eea returns inverse of a mod m.
func eea(a int, m int) uint {
	m2 := m
	var tn int
	t := []int{0, 1}
	n := 2
	for m%a != 0 {
		q := m / a
		r := m % a
		tn = t[n-2] - q*t[n-1]
		t = append(t, tn)
		m = a
		a = r
		n += 1
	}
	for tn < 0 {
		tn += m2
	}
	return uint(tn)
}

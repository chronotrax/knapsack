package main

func GCD(a, b uint64) uint64 {
	for b%a != 0 {
		r := b % a
		b = a
		a = r
	}
	return a
}

// EEA returns inverse of a mod m.
func EEA(a, m uint64) uint64 {
	m2 := int(m)
	var tn int
	t := []int{0, 1}
	n := 2
	for m%a != 0 {
		q := int(m / a)
		r := m % a
		tn = t[n-2] - q*t[n-1]
		t = append(t, tn)
		m = a
		a = r
		n += 1
	}
	// make sure the result is positive
	for tn < 0 {
		tn += m2
	}
	return uint64(tn)
}

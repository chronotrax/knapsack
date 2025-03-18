package knapsack

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// sMax is the largest big.Int to use when generating a Set's random number.
var sMax = big.NewInt(8)

type PrivateKey struct {
	V *big.Int
	U *big.Int
}

func validGCD(v, u *big.Int) bool {
	gcd := new(big.Int).GCD(nil, nil, v, u)
	return gcd.Cmp(big.NewInt(1)) == 0
}

func RandomPrivateKey(s Set) (*PrivateKey, error) {
	pMax := new(big.Int).Set(s[len(s)-1])
	pMax = pMax.Mul(pMax, big.NewInt(10))

	var err error
	for {
		// generate random u > 2*Sn
		var u *big.Int
		for {
			u, err = rand.Int(rand.Reader, pMax)
			if err != nil {
				return nil, err
			}

			// if u > 2 * Sn, stop generating new u values
			if u.Cmp(new(big.Int).Mul(big.NewInt(2), s[len(s)-1])) == 1 {
				break
			}
		}
		var v *big.Int
		for {
			v, err = rand.Int(rand.Reader, pMax)
			if err != nil {
				return nil, err
			}

			// if GCD(v,u) == 1, stop generating new v values
			if validGCD(v, u) {
				break
			}
		}

		return &PrivateKey{
			U: u,
			V: v,
		}, nil
	}
}

type Set []*big.Int

func (s Set) IsSuperincreasing() bool {
	sum := new(big.Int).Set(s[0])

	for _, si := range s[1:] {
		// if s1 <= sum
		if si.Cmp(sum) != 1 {
			return false
		}
		sum.Add(sum, si)
	}

	return true
}

// randomSetHelper generates a big.Int between [2, sMax+2)
func randomSetHelper() (*big.Int, error) {
	i, err := rand.Int(rand.Reader, sMax)
	if err != nil {
		return nil, err
	}
	i = i.Add(i, big.NewInt(2))
	return i, nil
}

// RandomSet returns a new superincreasing Set of size 8 * blockSize.
// Starting with a value < sMax and increasing
func RandomSet(blockSize int) (Set, error) {
	size := 8 * blockSize
	a := make([]*big.Int, size)

	// start with a random int
	var err error
	a[0], err = randomSetHelper()
	if err != nil {
		return nil, err
	}

	sum := big.NewInt(0)

	// each step, increase by sum + (r < sMax)
	for i := 1; i < size; i++ {
		sum.Add(sum, a[i-1])

		r, err := randomSetHelper()
		if err != nil {
			return nil, err
		}

		a[i] = new(big.Int).Add(sum, r)
	}

	return a, nil
}

type PublicKey []*big.Int

// NewPublicKey creates a new PublicKey with PrivateKey.V * Set[i] % PrivateKey.U.
func NewPublicKey(private *PrivateKey, s Set) PublicKey {
	public := make([]*big.Int, len(s))

	for i, si := range s {
		ai := new(big.Int).Mul(private.V, si)
		public[i] = ai.Mod(ai, private.U)
	}

	return public
}

type Plaintext []*big.Int

func (k *Knapsack) NewPlaintext(data []byte) Plaintext {
	if len(data) == 0 {
		return nil
	}

	plain := make([]*big.Int, 0)

	for i := 0; i < len(data); i += k.BlockSize {
		end := i + k.BlockSize
		if end > len(data) {
			end = len(data)
		}

		block := make([]byte, k.BlockSize)
		copy(block, data[i:end])

		bigInt := new(big.Int).SetBytes(block)
		plain = append(plain, bigInt)
	}

	return plain
}

func (k *Knapsack) FromPlaintext(plain Plaintext) []byte {
	if len(plain) == 0 {
		return nil
	}

	data := make([]byte, 0)

	for _, b := range plain {
		data = append(data, b.Bytes()...)
	}

	// remove padded 0's
	clean := make([]byte, 0)
	for _, d := range data {
		if d != 0 {
			clean = append(clean, d)
		}
	}

	return clean
}

type Ciphertext []*big.Int

type Knapsack struct {
	BlockSize int
	Private   *PrivateKey
	Set       Set
	Public    PublicKey
}

func validateBlockSize(blockSize int) error {
	if blockSize < 1 || 8 < blockSize {
		return fmt.Errorf("blockSize must be between [1,8]")
	}

	return nil
}

func newKnapsack(blockSize int, private *PrivateKey, s Set) (*Knapsack, error) {
	k := new(Knapsack)
	k.BlockSize = blockSize
	k.Private = private
	k.Set = s

	k.Public = NewPublicKey(k.Private, k.Set)

	return k, nil
}

func NewKnapsack(blockSize int) (*Knapsack, error) {
	err := validateBlockSize(blockSize)
	if err != nil {
		return nil, err
	}

	s, err := RandomSet(blockSize)
	if err != nil {
		return nil, err
	}

	private, err := RandomPrivateKey(s)
	if err != nil {
		return nil, err
	}

	return newKnapsack(blockSize, private, s)
}

func NewKnapsackCustom(blockSize int, private *PrivateKey, s Set) (*Knapsack, error) {
	err := validateBlockSize(blockSize)
	if err != nil {
		return nil, err
	}

	if !validGCD(private.V, private.U) {
		return nil, fmt.Errorf("GCD(%v, %v) != 1", private.V, private.U)
	}

	return newKnapsack(blockSize, private, s)
}

func (k *Knapsack) Encrypt(plain Plaintext) Ciphertext {
	cipher := make([]*big.Int, 0)

	// loop through every block
	for _, block := range plain {
		sum := big.NewInt(0)

		// loop through the bytes in each block (left to right)
		for i, b := range block.Bytes() {
			bit := byte(1)

			// loop through plaintext byte (right to left) and increment sum by the matching PublicKey index
			for j := 7; j >= 0; j-- {
				if b&bit != 0 {
					sum.Add(sum, k.Public[(i*8)+j])
				}
				bit = bit << 1 // bitshift left by 1
			}
		}

		cipher = append(cipher, sum)
	}

	return cipher
}

func (k *Knapsack) Decrypt(cipher Ciphertext) (Plaintext, error) {
	v := k.Private.V
	u := k.Private.U
	// if u == 0
	if u.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("u == 0")
	}

	inverse := new(big.Int).ModInverse(v, u)

	if inverse == nil || big.NewInt(0).Cmp(inverse) == 0 {
		return nil, fmt.Errorf("inverse is nil or 0")
	}

	// recreate Set
	size := 8 * k.BlockSize
	s := make([]*big.Int, size)
	for i := 0; i < size; i++ {
		si := new(big.Int)
		si = si.Mul(k.Public[i], inverse)
		s[i] = si.Mod(si, u)
	}

	plain := make([]*big.Int, 0)

	// loop through every block
	for _, block := range cipher {
		t := new(big.Int).Set(block)
		t = t.Mul(t, inverse)
		t = t.Mod(t, u)

		sum := big.NewInt(0)

		bit := 0

		// loop through t, subtract si if possible, create binary
		for i := size - 1; i >= 0; i-- {
			// if t >= s[i]
			if t.Cmp(s[i]) >= 0 {
				t.Sub(t, s[i])
				sum.SetBit(sum, bit, 1)
			}
			bit++
		}

		plain = append(plain, sum)
	}

	return plain, nil
}

func BigIntsToStr(ints []*big.Int) string {
	s := strings.Builder{}

	for i, b := range ints {
		s.WriteString(fmt.Sprintf("%v", b.String()))
		if i < len(ints)-1 {
			s.WriteString(", ")
		}
	}

	return s.String()
}

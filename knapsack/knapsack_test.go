package knapsack

import (
	"crypto/rand"
	"math/big"
	mathRand "math/rand/v2"
	"reflect"
	"strconv"
	"testing"
)

func TestKnapsack(t *testing.T) {
	type args struct {
		blockSize int
		private   *PrivateKey
		set       []int64
		data      []byte
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				blockSize: 1,
				private: &PrivateKey{
					U: big.NewInt(672),
					V: big.NewInt(13),
				},
				set:  []int64{3, 5, 9, 18, 38, 75, 155, 310},
				data: []byte("Bat"),
			},
		},
		{
			name: "2",
			args: args{
				blockSize: 2,
				private: &PrivateKey{
					U: big.NewInt(476729), // 476729 = 2 * 238364 + 1
					V: big.NewInt(476728),
				},
				set:  []int64{8, 17, 29, 56, 118, 234, 464, 931, 1862, 3724, 7448, 14900, 29794, 59591, 119183, 238364},
				data: []byte("Bat"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// convert to set
			s := Set{}
			for _, v := range tt.args.set {
				s = append(s, big.NewInt(v))
			}

			if !s.IsSuperincreasing() {
				t.Logf("WARNING: Set is not superincreasing: %#v\n", tt.args.set)
			}

			if len(s) != 8*tt.args.blockSize {
				t.Fatalf("set (len=%d) does not match blocksize (%d)", len(s), 8*tt.args.blockSize)
			}

			k, err := NewKnapsackCustom(tt.args.blockSize, tt.args.private, s)
			if err != nil {
				t.Fatal(err)
			}

			plain := k.NewPlaintext(tt.args.data)

			cipher := k.Encrypt(plain)

			newPlain, err := k.Decrypt(cipher)
			if err != nil {
				t.FailNow()
			}

			got := k.FromPlaintext(newPlain)

			if !reflect.DeepEqual(got, tt.args.data) {
				t.Errorf("got %#v, want %#v", got, tt.args.data)
			}
		})
	}
}

func TestRandomKnapsack(t *testing.T) {
	for i := range 100 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			//r := (rand.Int() % 8) + 1
			r := (mathRand.Int() % 8) + 1

			t.Logf("creating knapsack with blocksize %d", r)

			k, err := NewKnapsack(r)
			if err != nil {
				t.Fatal(err)
			}

			data := []byte(rand.Text())

			plain := k.NewPlaintext(data)

			cipher := k.Encrypt(plain)

			newPlain, err := k.Decrypt(cipher)
			if err != nil {
				t.FailNow()
			}

			got := k.FromPlaintext(newPlain)

			if !reflect.DeepEqual(got, data) {
				t.Errorf("got %#v, want %#v", got, data)
			}
		})
	}
}

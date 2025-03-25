package knapsack

import (
	"math/big"
	"reflect"
	"testing"
)

func Test_gs(t *testing.T) {
	type args struct {
		n int
		b []*big.Float
	}
	tests := []struct {
		name string
		args args
		want []*big.Float
	}{
		{
			name: "1",
			args: args{
				n: 3,
				b: []*big.Float{big.NewFloat(1), big.NewFloat(1), big.NewFloat(1),
					big.NewFloat(-1), big.NewFloat(0), big.NewFloat(1),
					big.NewFloat(1), big.NewFloat(1), big.NewFloat(2)},
			},
			want: []*big.Float{big.NewFloat(1), big.NewFloat(1 / 3), big.NewFloat(-1 / 2),
				big.NewFloat(-1), big.NewFloat(2 / 3), big.NewFloat(0),
				big.NewFloat(1), big.NewFloat(1 / 3), big.NewFloat(1 / 2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := emptySlice(tt.args.n * tt.args.n)
			B := emptySlice(tt.args.n)
			if gs(tt.args.b, mu, B, tt.args.n); !reflect.DeepEqual(tt.args.b, tt.want) {
				t.Errorf("gs() = %v, want %v", tt.args.b, tt.want)
			}
		})
	}
}

func Test_lll(t *testing.T) {
	type args struct {
		n     int
		b     []*big.Float
		delta *big.Float
	}
	tests := []struct {
		name string
		args args
		want []*big.Float
	}{
		{
			name: "1",
			args: args{
				n: 2,
				b: []*big.Float{big.NewFloat(47), big.NewFloat(95),
					big.NewFloat(215), big.NewFloat(460)},
				delta: big.NewFloat(3 / 4),
			},
			want: []*big.Float{big.NewFloat(1), big.NewFloat(40),
				big.NewFloat(30), big.NewFloat(5)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if lll(tt.args.b, tt.args.n, big.NewFloat(1/2), tt.args.delta); !reflect.DeepEqual(tt.args.b, tt.want) {
				t.Errorf("lll() = %v, want %v", tt.args.b, tt.want)
			}
		})
	}
}

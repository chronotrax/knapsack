package knapsack

import (
	"github.com/chronotrax/knapsack/matrix"
	"math/big"
	"reflect"
	"testing"
)

func Test_gs(t *testing.T) {
	type args struct {
		b matrix.Matrix
	}
	tests := []struct {
		name string
		args args
		want matrix.Matrix
	}{
		{
			name: "1",
			args: args{
				b: matrix.NewMatrixFull(3, 3,
					matrix.Vector{new(big.Rat).SetInt64(1), new(big.Rat).SetInt64(1), new(big.Rat).SetInt64(1),
						new(big.Rat).SetInt64(-1), new(big.Rat).SetInt64(0), new(big.Rat).SetInt64(1),
						new(big.Rat).SetInt64(1), new(big.Rat).SetInt64(1), new(big.Rat).SetInt64(2)}),
			},
			want: matrix.NewMatrixFull(3, 3,
				matrix.Vector{new(big.Rat).SetInt64(1), new(big.Rat).SetFrac64(1, 3), new(big.Rat).SetFrac64(-1, 2),
					new(big.Rat).SetInt64(-1), new(big.Rat).SetFrac64(2, 3), new(big.Rat).SetInt64(0),
					new(big.Rat).SetInt64(1), new(big.Rat).SetFrac64(1, 3), new(big.Rat).SetFrac64(1, 2)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := gs(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lll(t *testing.T) {
	type args struct {
		b             matrix.Matrix
		delta         *big.Rat
		maxIterations int
	}
	tests := []struct {
		name string
		args args
		want matrix.Matrix
	}{
		{
			name: "1",
			args: args{
				b: matrix.NewMatrixFull(2, 2,
					matrix.Vector{new(big.Rat).SetInt64(47), new(big.Rat).SetInt64(95),
						new(big.Rat).SetInt64(215), new(big.Rat).SetInt64(460)}),
				delta:         new(big.Rat).SetFrac64(3, 4),
				maxIterations: 1000,
			},
			want: matrix.NewMatrixFull(2, 2,
				matrix.Vector{new(big.Rat).SetInt64(1), new(big.Rat).SetInt64(40),
					new(big.Rat).SetInt64(30), new(big.Rat).SetInt64(5)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lll(tt.args.b, tt.args.delta, tt.args.maxIterations); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lll() = %v, want %v", got, tt.want)
			}
		})
	}
}

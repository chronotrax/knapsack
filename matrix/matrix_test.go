package matrix

import (
	"math/big"
	"reflect"
	"testing"
)

func TestSubtractVec(t *testing.T) {
	type args struct {
		a Vector
		b Vector
	}
	tests := []struct {
		name string
		args args
		want Vector
	}{
		{
			name: "1",
			args: args{
				a: Vector{new(big.Rat).SetInt64(1), new(big.Rat).SetInt64(-2)},
				b: Vector{new(big.Rat).SetInt64(2), new(big.Rat).SetInt64(-3)},
			},
			want: Vector{new(big.Rat).SetInt64(-1), new(big.Rat).SetInt64(1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubtractVec(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubtractVec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMulRatOnVec(t *testing.T) {
	type args struct {
		r *big.Rat
		v Vector
	}
	tests := []struct {
		name string
		args args
		want Vector
	}{
		{
			name: "1",
			args: args{
				r: new(big.Rat).SetInt64(2),
				v: Vector{new(big.Rat).SetInt64(4), new(big.Rat).SetInt64(-3)},
			},
			want: Vector{new(big.Rat).SetInt64(8), new(big.Rat).SetInt64(-6)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MulRatOnVec(tt.args.r, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MulRatOnVec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDotProduct(t *testing.T) {
	type args struct {
		a Vector
		b Vector
	}
	tests := []struct {
		name string
		args args
		want *big.Rat
	}{
		{
			name: "1",
			args: args{
				a: Vector{new(big.Rat).SetInt64(-6), new(big.Rat).SetInt64(8)},
				b: Vector{new(big.Rat).SetInt64(5), new(big.Rat).SetInt64(12)},
			},
			want: new(big.Rat).SetInt64(66),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DotProduct(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DotProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

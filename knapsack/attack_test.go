package knapsack

import (
	"reflect"
	"testing"
)

func Test_gs(t *testing.T) {
	type args struct {
		n int
		b []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "1",
			args: args{
				n: 3,
				b: []float64{1, 1, 1,
					-1, 0, 1,
					1, 1, 2},
			},
			want: []float64{1, 1 / 3, -1 / 2,
				-1, 2 / 3, 0,
				1, 1 / 3, 1 / 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := make([]float64, tt.args.n*tt.args.n)
			B := make([]float64, tt.args.n)
			if gs(tt.args.b, mu, B, tt.args.n); !reflect.DeepEqual(tt.args.b, tt.want) {
				t.Errorf("gs() = %v, want %v", tt.args.b, tt.want)
			}
		})
	}
}

func Test_lll(t *testing.T) {
	type args struct {
		n     int
		b     []float64
		delta float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "1",
			args: args{
				n: 2,
				b: []float64{47, 94,
					215, 460},
				delta: 3 / 4,
			},
			want: []float64{1, 40,
				30, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if lll(tt.args.b, tt.args.n, 1/2, tt.args.delta); !reflect.DeepEqual(tt.args.b, tt.want) {
				t.Errorf("lll() = %v, want %v", tt.args.b, tt.want)
			}
		})
	}
}

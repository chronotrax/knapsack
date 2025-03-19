package matrix

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNewMatrixFull(t *testing.T) {
	type args struct {
		col int
		row int
		arr []*big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    Matrix
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				col: 3,
				row: 2,
				arr: []*big.Int{big.NewInt(1), big.NewInt(2),
					big.NewInt(10), big.NewInt(20),
					big.NewInt(100), big.NewInt(200)},
			},
			want: [][]*big.Int{{big.NewInt(1), big.NewInt(2)},
				{big.NewInt(10), big.NewInt(20)},
				{big.NewInt(100), big.NewInt(200)}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMatrixFull(tt.args.col, tt.args.row, tt.args.arr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatrixFull() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatrixFull() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		a Matrix
		b Matrix
	}
	tests := []struct {
		name    string
		args    args
		want    Matrix
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				a: [][]*big.Int{{big.NewInt(1), big.NewInt(2)},
					{big.NewInt(3), big.NewInt(4)}},
				b: [][]*big.Int{{big.NewInt(10), big.NewInt(20)},
					{big.NewInt(30), big.NewInt(40)}},
			},
			want: [][]*big.Int{{big.NewInt(11), big.NewInt(22)},
				{big.NewInt(33), big.NewInt(44)}},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				a: [][]*big.Int{{big.NewInt(1), big.NewInt(2)},
					{big.NewInt(3), big.NewInt(4)}},
				b: [][]*big.Int{{big.NewInt(-10), big.NewInt(-20)},
					{big.NewInt(-30), big.NewInt(-40)}},
			},
			want: [][]*big.Int{{big.NewInt(-9), big.NewInt(-18)},
				{big.NewInt(-27), big.NewInt(-36)}},
			wantErr: false,
		},
		{
			name: "err 1: a.RowLen() != b.RowLen()",
			args: args{
				a: [][]*big.Int{{big.NewInt(1), big.NewInt(2)},
					{big.NewInt(3), big.NewInt(4)}},
				b: [][]*big.Int{{big.NewInt(10)},
					{big.NewInt(30)}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMul(t *testing.T) {
	type args struct {
		a Matrix
		b Matrix
	}
	tests := []struct {
		name    string
		args    args
		want    Matrix
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				a: [][]*big.Int{{big.NewInt(1), big.NewInt(2)},
					{big.NewInt(3), big.NewInt(4)}},
				b: [][]*big.Int{{big.NewInt(10), big.NewInt(20)},
					{big.NewInt(30), big.NewInt(40)}},
			},
			want: [][]*big.Int{{big.NewInt(70), big.NewInt(100)},
				{big.NewInt(150), big.NewInt(220)}},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				a: [][]*big.Int{{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
					{big.NewInt(4), big.NewInt(5), big.NewInt(6)}},
				b: [][]*big.Int{{big.NewInt(7), big.NewInt(8)},
					{big.NewInt(9), big.NewInt(10)},
					{big.NewInt(11), big.NewInt(12)}},
			},
			want: [][]*big.Int{{big.NewInt(58), big.NewInt(64)},
				{big.NewInt(139), big.NewInt(154)}},
			wantErr: false,
		},
		{
			name: "err 1: a.RowLen() != b.ColLen()",
			args: args{
				a: [][]*big.Int{{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
					{big.NewInt(4), big.NewInt(5), big.NewInt(6)}},
				b: [][]*big.Int{{big.NewInt(7), big.NewInt(8)},
					{big.NewInt(9), big.NewInt(10)}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mul() got = %v, want %v", got, tt.want)
			}
		})
	}
}

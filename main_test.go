package main

import (
	"reflect"
	"testing"
)

func Test_encrypt(t *testing.T) {
	type args struct {
		data []byte
		pub  pubKey
	}
	tests := []struct {
		name string
		args args
		want []uint
	}{
		{
			name: "test",
			args: args{
				data: []byte{byte('B'), byte('a'), byte('t')},
				pub:  pubKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: []uint{736, 852, 719},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encrypt(tt.args.data, tt.args.pub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eea(t *testing.T) {
	type args struct {
		a int
		m int
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "test",
			args: args{
				a: 13,
				m: 672,
			},
			want: 517,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eea(tt.args.a, tt.args.m); got != tt.want {
				t.Errorf("eea() = %v, want %v", got, tt.want)
			}
		})
	}
}

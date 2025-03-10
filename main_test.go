package main

import (
	"reflect"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	type args struct {
		data []byte
		pub  PublicKey
	}
	tests := []struct {
		name string
		args args
		want []uint
	}{
		{
			name: "test Encrypt 1",
			args: args{
				data: []byte{byte('B'), byte('a'), byte('t')},
				pub:  PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: []uint{736, 852, 719},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encrypt(tt.args.data, tt.args.pub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EEA(t *testing.T) {
	type args struct {
		a uint
		m uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "test EEA 1",
			args: args{
				a: 13,
				m: 672,
			},
			want: 517,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EEA(tt.args.a, tt.args.m); got != tt.want {
				t.Errorf("EEA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Decrypt(t *testing.T) {
	type args struct {
		data    []uint
		private PrivateKey
		public  PublicKey
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test Decrypt 1",
			args: args{
				data: []uint{736, 852, 719},
				private: PrivateKey{
					U: 672,
					V: 13,
				},
				public: PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: []byte{byte('B'), byte('a'), byte('t')},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decrypt(tt.args.data, tt.args.private, tt.args.public); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	type args struct {
		a uint
		b uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "test GCD",
			args: args{
				a: 13,
				b: 672,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GCD(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("GCD() = %v, want %v", got, tt.want)
			}
		})
	}
}

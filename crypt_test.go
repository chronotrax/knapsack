package main

import (
	"reflect"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	text := []byte("Bat")

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
			name: "valid 1",
			args: args{
				data: text,
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

func Test_Decrypt(t *testing.T) {
	text := []byte("Bat")

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
			name: "valid 1",
			args: args{
				data: []uint{736, 852, 719},
				private: PrivateKey{
					U: 672,
					V: 13,
				},
				public: PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
		},
		{
			name: "valid 2",
			args: args{
				data: []uint{736, 852, 719},
				private: PrivateKey{
					U: 113,
					V: 13,
				},
				public: PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
		},
		{
			name: "valid 3",
			args: args{
				data: []uint{736, 852, 719},
				private: PrivateKey{
					U: 902,
					V: 464,
				},
				public: PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
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

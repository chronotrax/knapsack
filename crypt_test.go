package main

import (
	"reflect"
	"testing"
)

func Test_Cryptosystem(t *testing.T) {
	text := []byte("CRYPTO test!!?")

	type args struct {
		data []byte
		u    uint
		v    uint
		s    S
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "valid 1",
			args: args{
				data: text,
				u:    881,
				v:    588,
				s:    S{2, 7, 11, 21, 42, 89, 180, 354},
			},
			want: text,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv, err := NewPrivateKey(tt.args.u, tt.args.v)
			if err != nil {
				t.Errorf("NewPrivateKey() error = %v", err)
			}

			pub := NewPublicKey(tt.args.s, priv)

			cipher := Encrypt(tt.args.data, pub)
			if got := Decrypt(cipher, priv, pub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Encrypt(t *testing.T) {
	text := []byte("Bat")

	wikiPriv, _ := NewPrivateKey(881, 588)
	wikiPub := NewPublicKey(S{2, 7, 11, 21, 42, 89, 180, 354}, wikiPriv)
	wikiText := []byte{byte(97)}

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
		{
			name: "valid 2", // https://en.wikipedia.org/wiki/Merkle%E2%80%93Hellman_knapsack_cryptosystem#Example
			args: args{
				data: wikiText,
				pub:  wikiPub,
			},
			want: []uint{1129},
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

	wikiPriv, _ := NewPrivateKey(881, 588)
	wikiPub := NewPublicKey(S{2, 7, 11, 21, 42, 89, 180, 354}, wikiPriv)
	wikiText := []byte{byte(97)}

	type args struct {
		cipher  []uint
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
				cipher: []uint{736, 852, 719},
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
				cipher: []uint{736, 852, 719},
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
				cipher: []uint{736, 852, 719},
				private: PrivateKey{
					U: 902,
					V: 464,
				},
				public: PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
		},
		{
			name: "valid 4", // https://en.wikipedia.org/wiki/Merkle%E2%80%93Hellman_knapsack_cryptosystem#Example
			args: args{
				cipher:  []uint{1129},
				private: wikiPriv,
				public:  wikiPub,
			},
			want: wikiText,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decrypt(tt.args.cipher, tt.args.private, tt.args.public); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

package crypt

import (
	"github.com/chronotrax/knapsack/types"
	"reflect"
	"testing"
)

func Test_Cryptosystem(t *testing.T) {
	text := types.Plaintext("CRYPTO test!!?")

	type args struct {
		data types.Plaintext
		u    uint64
		v    uint64
		s    types.S
	}
	tests := []struct {
		name string
		args args
		want types.Plaintext
	}{
		{
			name: "valid 1",
			args: args{
				data: text,
				u:    881,
				v:    588,
				s:    types.S{2, 7, 11, 21, 42, 89, 180, 354},
			},
			want: text,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			private, err := types.NewPrivateKey(tt.args.u, tt.args.v)
			if err != nil {
				t.Errorf("NewPrivateKey() error = %v", err)
			}

			pub := types.NewPublicKey(tt.args.s, private)

			cipher := Encrypt(tt.args.data, pub)
			if got := Decrypt(cipher, private, pub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Encrypt(t *testing.T) {
	text := types.Plaintext("Bat")

	wikiPrivate, _ := types.NewPrivateKey(881, 588)
	wikiPub := types.NewPublicKey(types.S{2, 7, 11, 21, 42, 89, 180, 354}, wikiPrivate)
	wikiText := types.Plaintext{byte(97)}

	type args struct {
		data []byte
		pub  types.PublicKey
	}
	tests := []struct {
		name string
		args args
		want types.Ciphertext
	}{
		{
			name: "valid 1",
			args: args{
				data: text,
				pub:  types.PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: []uint64{736, 852, 719},
		},
		{
			name: "valid 2", // https://en.wikipedia.org/wiki/Merkle%E2%80%93Hellman_knapsack_cryptosystem#Example
			args: args{
				data: wikiText,
				pub:  wikiPub,
			},
			want: []uint64{1129},
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
	text := types.Plaintext("Bat")

	wikiPrivate, _ := types.NewPrivateKey(881, 588)
	wikiPub := types.NewPublicKey(types.S{2, 7, 11, 21, 42, 89, 180, 354}, wikiPrivate)
	wikiText := types.Plaintext{byte(97)}

	type args struct {
		cipher  types.Ciphertext
		private types.PrivateKey
		public  types.PublicKey
	}
	tests := []struct {
		name string
		args args
		want types.Plaintext
	}{
		{
			name: "valid 1",
			args: args{
				cipher: []uint64{736, 852, 719},
				private: types.PrivateKey{
					U: 672,
					V: 13,
				},
				public: types.PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
		},
		{
			name: "valid 2",
			args: args{
				cipher: []uint64{736, 852, 719},
				private: types.PrivateKey{
					U: 113,
					V: 13,
				},
				public: types.PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
		},
		{
			name: "valid 3",
			args: args{
				cipher: []uint64{736, 852, 719},
				private: types.PrivateKey{
					U: 902,
					V: 464,
				},
				public: types.PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			},
			want: text,
		},
		{
			name: "valid 4", // https://en.wikipedia.org/wiki/Merkle%E2%80%93Hellman_knapsack_cryptosystem#Example
			args: args{
				cipher:  []uint64{1129},
				private: wikiPrivate,
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

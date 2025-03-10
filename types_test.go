package main

import (
	"reflect"
	"testing"
)

func TestS_IsSuperincreasing(t *testing.T) {
	tests := []struct {
		name string
		s    S
		want bool
	}{
		{
			name: "valid 1",
			s:    S{5, 6, 12, 24, 48, 96, 192, 384},
			want: true,
		},
		{
			name: "invalid 1",
			s:    S{3, 5, 13, 18, 99, 108, 323, 350},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsSuperincreasing(); got != tt.want {
				t.Errorf("IsSuperincreasing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPublicKey(t *testing.T) {
	type args struct {
		s       S
		private PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    PublicKey
		wantErr bool
	}{
		{
			name: "valid 1",
			args: args{
				s: S{3, 5, 9, 18, 38, 75, 155, 310},
				private: PrivateKey{
					U: 672,
					V: 13,
				},
			},
			want:    PublicKey{39, 65, 117, 234, 494, 303, 671, 670},
			wantErr: false,
		},
		{
			name: "invalid 1",
			args: args{
				s: S{3, 5, 9, 18, 38, 75, 155, 336}, // 336*2 = 672
				private: PrivateKey{
					U: 672,
					V: 13,
				},
			},
			want:    PublicKey{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPublicKey(tt.args.s, tt.args.private)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPublicKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

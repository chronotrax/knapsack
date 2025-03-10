package main

import (
	"testing"
)

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
			name: "valid 1",
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
			name: "valid 1",
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

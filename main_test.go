package main

import (
	"testing"
)

func TestFakeLogin(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Testing", args: args{username: "hanan", password: "password"}, want: true},
		{name: "Testing Wrong Password", args: args{username: "hanan", password: "pass"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FakeLogin(tt.args.username, tt.args.password); got != tt.want {
				t.Errorf("FakeLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

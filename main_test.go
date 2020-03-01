package main

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{"Test1", "TO THE CLOUD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			print(tt.message)
		})
	}
}

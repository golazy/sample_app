package main

import "testing"

func TestListenAddr(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{name: "unset", want: ":8080"},
		{name: "port only", addr: "9191", want: ":9191"},
		{name: "all interfaces", addr: ":9191", want: ":9191"},
		{name: "specific interface", addr: "127.0.0.1:9191", want: "127.0.0.1:9191"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("ADDR", test.addr)

			if got := listenAddr(); got != test.want {
				t.Fatalf("listenAddr() = %q, want %q", got, test.want)
			}
		})
	}
}

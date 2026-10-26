package main

import "testing"

func TestAbsoluteURL(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "empty", raw: "", wantErr: true},
		{name: "relative", raw: "/vpn-panel/", wantErr: true},
		{name: "no host", raw: "https://", wantErr: true},
		{name: "malformed", raw: "://bad", wantErr: true},
		{name: "absolute", raw: "https://panel.example.com/vpn-panel/"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := absoluteURL(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatalf("absoluteURL(%q) = nil error, want error", tc.raw)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("absoluteURL(%q) returned error: %v", tc.raw, err)
			}
		})
	}
}

package tappay

import (
	"os"
	"testing"
)

// Creates a client using a combination of custom server and TAPPAY_SERVER environment variable
// and verify that the base path (client.url) are set correctly
func TestWithServer(t *testing.T) {
	originalTappayServer := os.Getenv("TAPPAY_SERVER")
	for _, tc := range []struct {
		name         string
		customServer string
		tappayServer string
		WantServer   string
		WantError    bool
	}{
		{
			name:       "Given unset server option and env variable, returns default server url",
			WantServer: APIURL,
		},
		{
			name:         "Given valid server option, returns server option",
			customServer: "http://valid.tappay.com",
			WantServer:   "http://valid.tappay.com",
		},
		{
			name:         "Given valid url in env variable, returns env variable",
			tappayServer: "http://valid.tappay.com",
			WantServer:   "http://valid.tappay.com",
		},
		{
			name:         "Given both(server option and env variable), returns server option which takes precedence",
			customServer: "http://valid1.tappay.com",
			tappayServer: "http://valid2.tappay.com",
			WantServer:   "http://valid1.tappay.com",
		},
		{
			name:         "Given invalid server option, wants error",
			customServer: "BadURL",
			WantError:    true,
		},
		{
			name:         "Given invalid env variable, wants error",
			tappayServer: "BadURL",
			WantError:    true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("TAPPAY_SERVER", tc.tappayServer)

			var cli *client
			var err error
			if tc.customServer != "" {
				cli, err = NewClient("tappay_key", WithServer(tc.customServer))
			} else {

				cli, err = NewClient("tappay_key")
			}
			if tc.WantError && err == nil {
				t.Errorf("expected an error, but the creation succeeded")
			}
			if !tc.WantError && err != nil {
				t.Errorf("expected the creation succeeded, but got error : %v", err)
			}
			if cli != nil && cli.url != tc.WantServer {
				t.Errorf("expected url :%s, but got :%s", tc.WantServer, cli.url)
			}
		})
	}
	os.Setenv("TAPPAY_SERVER", originalTappayServer)
}

package tappay

import (
	"context"
	"testing"
)

func TestRecord(t *testing.T) {
	for _, tc := range []struct {
		name       string
		setupQuery func(*client) RecordParams
		wantError  bool
	}{
		{
			name: "Given empty filter returns no error",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cli, _ := NewClient("partner_6ID1DoDlaPrfHw6HBZsULfTYtDmWs0q0ZZGKMBpp4YICWBxgK97eK3RM", WithServer(SandboxAPIURL))
			var query RecordParams
			if tc.setupQuery != nil {
				query = tc.setupQuery(cli)
			}
			records, err := cli.Records(context.Background(), query)
			if err != nil {
				t.Errorf("unexpected get record error: %v", err)
			}
			if tc.wantError && records.Status != 0 && records.Status != 2 {
				t.Errorf("unexpected status: %d", records.Status)
			}
		})
	}
}

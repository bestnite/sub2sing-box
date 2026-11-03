package model

import (
	"strings"
	"testing"
)

func TestParseConvertRequestReadsPluralKeys(t *testing.T) {
	data := []byte(`{
		"subscriptions": ["https://example.com/sub"],
		"proxies": ["vless://00000000-0000-0000-0000-000000000000@example.com:443"],
		"template": "sing-box-template.json",
		"group-type": "urltest",
		"sort": "tag",
		"sort-type": "asc",
		"output": "config.json",
		"user-agent": "v2rayN/7.25.3"
	}`)

	cfg, err := ParseConvertRequest(data)
	if err != nil {
		t.Fatalf("ParseConvertRequest() error = %v", err)
	}
	if len(cfg.Subscriptions) != 1 {
		t.Fatalf("Subscriptions = %v, want one entry", cfg.Subscriptions)
	}
	if len(cfg.Proxies) != 1 {
		t.Fatalf("Proxies = %v, want one entry", cfg.Proxies)
	}
	if cfg.UserAgent != "v2rayN/7.25.3" {
		t.Fatalf("UserAgent = %q, want %q", cfg.UserAgent, "v2rayN/7.25.3")
	}
}

func TestParseConvertRequestRejectsUnknownKey(t *testing.T) {
	cases := []string{
		`{"subscription": ["https://example.com/sub"]}`,
		`{"proxy": ["vless://00000000-0000-0000-0000-000000000000@example.com:443"]}`,
		`{"output": "config.json", "whatever": true}`,
	}

	for _, c := range cases {
		_, err := ParseConvertRequest([]byte(c))
		if err == nil {
			t.Fatalf("ParseConvertRequest(%s) = nil error, want unknown field error", c)
		}
		if !strings.Contains(err.Error(), "unknown field") {
			t.Fatalf("ParseConvertRequest(%s) error = %v, want unknown field error", c, err)
		}
	}
}

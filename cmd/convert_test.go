package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigMergesPluralKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub2sing-box.json")
	content := `{
		"subscriptions": ["https://example.com/sub"],
		"proxies": ["vless://00000000-0000-0000-0000-000000000000@example.com:443"],
		"output": "config.json"
	}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	origConfig, origSubscriptions, origProxies, origOutput := config, subscriptions, proxies, output
	defer func() {
		config, subscriptions, proxies, output = origConfig, origSubscriptions, origProxies, origOutput
	}()

	config, subscriptions, proxies, output = path, nil, nil, ""
	if err := loadConfig(); err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	if len(subscriptions) != 1 {
		t.Fatalf("subscriptions = %v, want one entry", subscriptions)
	}
	if len(proxies) != 1 {
		t.Fatalf("proxies = %v, want one entry", proxies)
	}
	if output != "config.json" {
		t.Fatalf("output = %q, want %q", output, "config.json")
	}
}

func TestLoadConfigFailsOnUnknownKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub2sing-box.json")
	content := `{"subscription": ["https://example.com/sub"]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	origConfig, origSubscriptions := config, subscriptions
	defer func() {
		config, subscriptions = origConfig, origSubscriptions
	}()

	config, subscriptions = path, nil
	if err := loadConfig(); err == nil {
		t.Fatal("loadConfig() error = nil, want unknown field error")
	}
	if len(subscriptions) != 0 {
		t.Fatalf("subscriptions = %v, want empty on error", subscriptions)
	}
}

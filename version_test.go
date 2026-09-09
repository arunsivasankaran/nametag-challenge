package main

import "testing"

func TestCompareVersions(t *testing.T) {
	if compareVersions("v1.0.0", "v1.0.1") >= 0 {
		t.Fatal("expected v1.0.0 to be older than v1.0.1")
	}

	if compareVersions("v2.3.4", "v2.3.4") != 0 {
		t.Fatal("expected equal versions to compare as equal")
	}

	if compareVersions("v3.0.0", "v2.9.9") <= 0 {
		t.Fatal("expected v3.0.0 to be newer than v2.9.9")
	}
}

func TestManifestValidatesVersion(t *testing.T) {
	m := manifest{Version: "v1.2.3", DownloadURL: "https://example.com/update", SHA256: "abc123"}
	if !m.isValid() {
		t.Fatal("expected manifest to be valid when version and download URL are set")
	}

	m.Version = ""
	if m.isValid() {
		t.Fatal("expected manifest to be invalid when version is missing")
	}
}

func TestCompareVersionsWithGoTagFormat(t *testing.T) {
	if compareVersions("v1.0.0", "go1.23.1") >= 0 {
		t.Fatal("expected go1.23.1 to be newer than v1.0.0")
	}
}

func TestLatestTagFromTagsResponse(t *testing.T) {
	body := []byte(`[
		{"name": "go1.23.1"},
		{"name": "go1.23.0"}
	]`)

	latest, err := latestTagFromTagsResponse(body)
	if err != nil {
		t.Fatalf("expected tag response to parse without error: %v", err)
	}
	if latest != "go1.23.1" {
		t.Fatalf("expected latest tag to be go1.23.1, got %q", latest)
	}
}

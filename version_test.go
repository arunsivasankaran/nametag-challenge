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

func TestCompareVersionsWithGoTagFormat(t *testing.T) {
	if compareVersions("v1.0.0", "go1.23.1") >= 0 {
		t.Fatal("expected go1.23.1 to be newer than v1.0.0")
	}
}

func TestSelectReleaseAssetURL(t *testing.T) {
	assets := []releaseAsset{
		{Name: "nametag-linux-amd64.tar.gz", BrowserDownloadURL: "https://example.com/linux.tar.gz"},
		{Name: "nametag-darwin-arm64.tar.gz", BrowserDownloadURL: "https://example.com/macos.tar.gz"},
	}

	url, err := selectReleaseAssetURL(assets)
	if err != nil {
		t.Fatalf("expected a matching asset to be returned: %v", err)
	}

	if url == "" {
		t.Fatal("expected selected asset URL to be non-empty")
	}
}

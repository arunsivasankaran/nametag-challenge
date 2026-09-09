package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	currentVersion     = "v1.0.1"
	defaultGitHubOwner = "arunsivasankaran"
	defaultGitHubRepo  = "nametag-challenge"
)

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type githubRelease struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

func compareVersions(local, remote string) int {
	local = normalizeVersion(local)
	remote = normalizeVersion(remote)

	localParts := strings.Split(local, ".")
	remoteParts := strings.Split(remote, ".")
	maxLen := len(localParts)
	if len(remoteParts) > maxLen {
		maxLen = len(remoteParts)
	}

	for i := 0; i < maxLen; i++ {
		var lv, rv int
		if i < len(localParts) {
			lv, _ = strconv.Atoi(strings.TrimPrefix(localParts[i], "v"))
		}
		if i < len(remoteParts) {
			rv, _ = strconv.Atoi(strings.TrimPrefix(remoteParts[i], "v"))
		}
		if lv < rv {
			return -1
		}
		if lv > rv {
			return 1
		}
	}
	return 0
}

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	v = strings.TrimPrefix(v, "go")
	v = strings.TrimPrefix(v, "Go")
	v = strings.TrimPrefix(v, "release-")
	return v
}

func githubReleaseURL() string {
	owner := strings.TrimSpace(os.Getenv("GITHUB_OWNER"))
	if owner == "" {
		owner = defaultGitHubOwner
	}
	repo := strings.TrimSpace(os.Getenv("GITHUB_REPO"))
	if repo == "" {
		repo = defaultGitHubRepo
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
}

func latestReleaseFromResponse(body []byte) (githubRelease, error) {
	var rel githubRelease
	if err := json.Unmarshal(body, &rel); err != nil {
		return githubRelease{}, err
	}
	if strings.TrimSpace(rel.TagName) == "" {
		return githubRelease{}, errors.New("latest release tag not found")
	}
	return rel, nil
}

func getLatestGitHubRelease() (githubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, githubReleaseURL(), nil)
	if err != nil {
		return githubRelease{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "nametag-cli")

	resp, err := client.Do(req)
	if err != nil {
		return githubRelease{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return githubRelease{}, fmt.Errorf("github releases API returned status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return githubRelease{}, err
	}

	return latestReleaseFromResponse(body)
}

func assetNameMatchesPlatform(assetName string) bool {
	lowerName := strings.ToLower(assetName)
	lowerOS := strings.ToLower(runtime.GOOS)
	lowerArch := strings.ToLower(runtime.GOARCH)

	if strings.Contains(lowerName, lowerOS) && strings.Contains(lowerName, lowerArch) {
		return true
	}
	if strings.Contains(lowerName, ".") && strings.Contains(lowerName, lowerOS) {
		return true
	}
	return false
}

func selectReleaseAssetURL(assets []releaseAsset) (string, error) {
	for _, asset := range assets {
		if assetNameMatchesPlatform(asset.Name) {
			if strings.TrimSpace(asset.BrowserDownloadURL) != "" {
				return asset.BrowserDownloadURL, nil
			}
		}
	}
	for _, asset := range assets {
		if strings.TrimSpace(asset.BrowserDownloadURL) != "" {
			return asset.BrowserDownloadURL, nil
		}
	}
	return "", errors.New("no release asset found for this platform")
}

func main() {
	fmt.Printf("Nametag CLI version %s\n", currentVersion)

	updated, err := checkForUpdate(currentVersion)
	if err != nil {
		fmt.Printf("update check failed: %v\n", err)
		return
	}

	if updated {
		fmt.Println("A newer version is available. Updating...")
		if err := runUpdate(); err != nil {
			fmt.Printf("update failed: %v\n", err)
			return
		}
		fmt.Println("Update installed successfully.")
	} else {
		fmt.Println("You are already on the latest version.")
	}
}

func checkForUpdate(current string) (bool, error) {
	release, err := getLatestGitHubRelease()
	if err != nil {
		return false, err
	}
	return compareVersions(current, release.TagName) < 0, nil
}

func runUpdate() error {
	release, err := getLatestGitHubRelease()
	if err != nil {
		return err
	}
	if compareVersions(currentVersion, release.TagName) >= 0 {
		return nil
	}

	assetURL, err := selectReleaseAssetURL(release.Assets)
	if err != nil {
		return err
	}

	downloadPath, err := downloadBinary(assetURL)
	if err != nil {
		return err
	}
	defer os.Remove(downloadPath)

	return replaceExecutable(downloadPath)
}

func downloadBinary(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download status: %s", resp.Status)
	}

	tempDir, err := os.MkdirTemp("", "nametag-update-")
	if err != nil {
		return "", err
	}

	targetPath := filepath.Join(tempDir, executableName())
	out, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	if err := os.Chmod(targetPath, 0o755); err != nil {
		return "", err
	}

	return targetPath, nil
}

func executableName() string {
	name := "nametag"
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func verifySHA256(path, expected string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(content)
	actual := hex.EncodeToString(hash[:])
	if !strings.EqualFold(actual, strings.TrimSpace(expected)) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", strings.TrimSpace(expected), actual)
	}

	return nil
}

func replaceExecutable(newBinary string) error {
	currentExec, err := os.Executable()
	if err != nil {
		return err
	}

	backupPath := currentExec + ".bak"
	if err := os.Rename(currentExec, backupPath); err != nil {
		return err
	}

	if err := os.Rename(newBinary, currentExec); err != nil {
		if restoreFromBackup(currentExec, backupPath) != nil {
			return fmt.Errorf("install failed and rollback also failed: %w", err)
		}
		return err
	}

	if err := os.Remove(backupPath); err != nil {
		return err
	}
	return nil
}

func restoreFromBackup(currentExec, backupPath string) error {
	if _, err := os.Stat(backupPath); err != nil {
		return err
	}
	_ = os.Remove(currentExec)
	return os.Rename(backupPath, currentExec)
}

func init() {
	_ = time.Second
}

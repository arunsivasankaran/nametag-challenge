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
	currentVersion = "v1.0.0"
	updateURL      = "http://127.0.0.1:8080/manifest.json"
)

type manifest struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
	SHA256      string `json:"sha256"`
	Description string `json:"description,omitempty"`
}

func (m manifest) isValid() bool {
	return strings.TrimSpace(m.Version) != "" && strings.TrimSpace(m.DownloadURL) != "" && strings.TrimSpace(m.SHA256) != ""
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
	return v
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
	resp, err := http.Get(updateURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var m manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return false, err
	}
	if !m.isValid() {
		return false, errors.New("remote manifest is invalid")
	}

	return compareVersions(current, m.Version) < 0, nil
}

func runUpdate() error {
	resp, err := http.Get(updateURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var m manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return err
	}
	if !m.isValid() {
		return errors.New("manifest missing required fields")
	}

	if compareVersions(currentVersion, m.Version) >= 0 {
		return nil
	}

	downloadPath, err := downloadBinary(m.DownloadURL)
	if err != nil {
		return err
	}
	defer os.Remove(downloadPath)

	if err := verifySHA256(downloadPath, m.SHA256); err != nil {
		return err
	}

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

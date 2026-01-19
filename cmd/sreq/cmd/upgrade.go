package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	githubRepo   = "Priyans-hu/sreq"
	githubAPIURL = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var (
	upgradeForce bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade sreq to the latest version",
	Long: `Check for and install the latest version of sreq.

This command will:
1. Check for the latest release on GitHub
2. Download the appropriate binary for your OS/architecture
3. Replace the current binary with the new version

Examples:
  sreq upgrade              # Upgrade to latest version
  sreq upgrade --force      # Force upgrade even if already on latest`,
	RunE: runUpgrade,
}

func init() {
	upgradeCmd.Flags().BoolVarP(&upgradeForce, "force", "f", false, "Force upgrade even if already on latest version")
	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	fmt.Printf("Current version: %s\n", Version)
	fmt.Println("Checking for updates...")

	// Get latest release info
	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(Version, "v")

	fmt.Printf("Latest version: %s\n", release.TagName)

	// Check if upgrade is needed
	if currentVersion == latestVersion && !upgradeForce {
		fmt.Println("You are already on the latest version!")
		return nil
	}

	if currentVersion == "dev" {
		fmt.Println("Warning: Running development version, upgrading to release version")
	}

	// Find the right asset for this OS/arch
	assetName := getAssetName(release.TagName)
	var assetURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			assetURL = asset.BrowserDownloadURL
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("no release found for %s/%s (looking for %s)", runtime.GOOS, runtime.GOARCH, assetName)
	}

	fmt.Printf("Downloading %s...\n", assetName)

	// Download to temp file
	tmpDir, err := os.MkdirTemp("", "sreq-upgrade")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	archivePath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(archivePath, assetURL); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	fmt.Println("Extracting...")

	// Extract binary
	binaryPath, err := extractBinary(archivePath, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	fmt.Printf("Installing to %s...\n", execPath)

	// Replace current binary
	if err := replaceBinary(binaryPath, execPath); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	fmt.Printf("Successfully upgraded to %s!\n", release.TagName)
	return nil
}

func getLatestRelease() (*githubRelease, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func getAssetName(version string) string {
	goos := runtime.GOOS
	arch := runtime.GOARCH

	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}

	// Asset name format: sreq_0.1.0_darwin_arm64.tar.gz
	ver := strings.TrimPrefix(version, "v")
	return fmt.Sprintf("sreq_%s_%s_%s%s", ver, goos, arch, ext)
}

func downloadFile(destPath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, destDir string) (string, error) {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	return extractTarGz(archivePath, destDir)
}

func extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)

	var binaryPath string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Look for the sreq binary
		if header.Typeflag == tar.TypeReg && (header.Name == "sreq" || strings.HasSuffix(header.Name, "/sreq")) {
			binaryPath = filepath.Join(destDir, "sreq")
			outFile, err := os.OpenFile(binaryPath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				return "", err
			}
			_, copyErr := io.Copy(outFile, tr)
			closeErr := outFile.Close()
			if copyErr != nil {
				return "", copyErr
			}
			if closeErr != nil {
				return "", closeErr
			}
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("sreq binary not found in archive")
	}

	return binaryPath, nil
}

func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = r.Close() }()

	var binaryPath string
	for _, f := range r.File {
		// Look for the sreq binary
		if f.Name == "sreq.exe" || strings.HasSuffix(f.Name, "/sreq.exe") {
			binaryPath = filepath.Join(destDir, "sreq.exe")

			rc, err := f.Open()
			if err != nil {
				return "", err
			}

			outFile, err := os.OpenFile(binaryPath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				_ = rc.Close()
				return "", err
			}

			_, copyErr := io.Copy(outFile, rc)
			closeOutErr := outFile.Close()
			closeRcErr := rc.Close()

			if copyErr != nil {
				return "", copyErr
			}
			if closeOutErr != nil {
				return "", closeOutErr
			}
			if closeRcErr != nil {
				return "", closeRcErr
			}
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("sreq.exe not found in archive")
	}

	return binaryPath, nil
}

func replaceBinary(newBinary, targetPath string) error {
	// On Windows, we can't replace a running executable directly
	// We need to rename it first
	if runtime.GOOS == "windows" {
		oldPath := targetPath + ".old"
		_ = os.Remove(oldPath) // Remove any existing .old file
		if err := os.Rename(targetPath, oldPath); err != nil {
			return fmt.Errorf("failed to rename old binary: %w", err)
		}
	}

	// Copy new binary to target location
	src, err := os.Open(newBinary)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

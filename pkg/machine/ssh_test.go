package machine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLocalhostSSHArgsKnownHostsIsRealFile guards against regressing to a device name
// for UserKnownHostsFile. os.DevNull is "NUL" on Windows, and because the ssh binary is
// taken from PATH, an MSYS/Cygwin ssh can win the race and create a literal file named
// NUL in the caller's working directory (issue #30026).
func TestLocalhostSSHArgsKnownHostsIsRealFile(t *testing.T) {
	var knownHosts string
	for _, arg := range LocalhostSSHArgs() {
		if strings.HasPrefix(arg, "UserKnownHostsFile=") {
			knownHosts = strings.TrimPrefix(arg, "UserKnownHostsFile=")
		}
	}
	if knownHosts == "" {
		t.Fatal("LocalhostSSHArgs() does not set UserKnownHostsFile")
	}
	if knownHosts == os.DevNull {
		t.Fatalf("UserKnownHostsFile must not be a device name (%q): a non-native ssh creates a file with that name", knownHosts)
	}
	// a relative path would land in the caller's working directory
	if !filepath.IsAbs(filepath.FromSlash(knownHosts)) {
		t.Fatalf("UserKnownHostsFile %q is not an absolute path", knownHosts)
	}
	// Cygwin/MSYS builds do not accept backslash-separated Windows paths
	if strings.Contains(knownHosts, "\\") {
		t.Fatalf("UserKnownHostsFile %q contains backslashes and is not portable to MSYS/Cygwin ssh", knownHosts)
	}
	if fi, err := os.Stat(filepath.FromSlash(knownHosts)); err != nil {
		t.Fatalf("UserKnownHostsFile %q does not exist: %v", knownHosts, err)
	} else if fi.IsDir() {
		t.Fatalf("UserKnownHostsFile %q is a directory", knownHosts)
	}
}

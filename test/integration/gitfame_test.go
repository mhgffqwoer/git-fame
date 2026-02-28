package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const importPath = "github.com/mhgffqwoer/git-fame/cmd/git-fame"

var binaryPath string

func getRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "gitfame-bin-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir) //nolint

	binaryPath = filepath.Join(tmpDir, "gitfame")

	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, importPath) //nolint
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build gitfame: %v\n%s", err, output)
		os.Exit(1) //nolint
	}

	os.Exit(m.Run())
}

func TestGitFame(t *testing.T) {
	require.NotEmpty(t, binaryPath, "binary path should not be empty")

	basePath := getRootPath()
	bundlesDir := filepath.Join(basePath, "testdata", "bundles")
	testsDir := filepath.Join(basePath, "testdata", "tests")

	testDirs := ListTestDirs(t, testsDir)

	for _, dir := range testDirs {
		tc := ReadTestCase(t, filepath.Join(testsDir, dir))

		t.Run(dir+"/"+tc.Name, func(t *testing.T) {
			tempRepoDir, err := os.MkdirTemp("", "gitfame-repo-")
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(tempRepoDir) }()

			args := []string{"--repository", tempRepoDir} //nolint
			args = append(args, tc.Args...)

			Unbundle(t, filepath.Join(bundlesDir, tc.Bundle), tempRepoDir)
			headRef := GetHEADRef(t, tempRepoDir)

			cmd := exec.Command(binaryPath, args...) //nolint
			cmd.Stderr = os.Stderr

			output, err := cmd.Output()
			if !tc.Error {
				require.NoError(t, err)
				CompareResults(t, tc.Expected, output, tc.Format)
			} else {
				require.Error(t, err)
				_, ok := err.(*exec.ExitError)
				require.True(t, ok)
			}

			newHEADRef := GetHEADRef(t, tempRepoDir)
			require.Equal(t, headRef, newHEADRef)
		})
	}
}

func ListTestDirs(t *testing.T, testPath string) []string {
	t.Helper()

	files, err := os.ReadDir(testPath)
	require.NoError(t, err)

	var names []string //nolint
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		names = append(names, f.Name())
	}

	toInt := func(name string) int {
		i, err := strconv.Atoi(name)
		require.NoError(t, err)
		return i
	}

	sort.Slice(names, func(i, j int) bool {
		return toInt(names[i]) < toInt(names[j])
	})

	return names
}

type TestCase struct {
	*TestDescription
	Expected []byte
}

func ReadTestCase(t *testing.T, testPath string) *TestCase {
	t.Helper()

	desc := ReadTestDescription(t, testPath)

	expected, err := os.ReadFile(filepath.Join(testPath, "expected.out")) //nolint
	require.NoError(t, err)

	return &TestCase{TestDescription: desc, Expected: expected}
}

type TestDescription struct {
	Name   string   `yaml:"name"`
	Args   []string `yaml:"args"`
	Bundle string   `yaml:"bundle"`
	Error  bool     `yaml:"error"`
	Format string   `yaml:"format,omitempty"`
}

func ReadTestDescription(t *testing.T, testPath string) *TestDescription {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(testPath, "description.yaml")) //nolint
	require.NoError(t, err)

	var desc TestDescription
	require.NoError(t, yaml.Unmarshal(data, &desc))

	return &desc
}

func Unbundle(t *testing.T, src, dst string) {
	t.Helper()

	cmd := exec.Command("git", "clone", src, dst) //nolint
	require.NoError(t, cmd.Run())
}

func CompareResults(t *testing.T, expected, actual []byte, format string) {
	t.Helper()

	switch format {
	case "json":
		if len(expected) == 0 {
			require.Empty(t, string(actual))
		} else {
			require.JSONEq(t, string(expected), string(actual))
		}
	case "json-lines":
		if len(expected) == 0 {
			require.Empty(t, string(actual))
		} else {
			CompareJSONLines(t, expected, actual)
		}
	default:
		require.Equal(t, string(expected), string(actual))
	}
}

func CompareJSONLines(t *testing.T, expected, actual []byte) {
	t.Helper()

	expectedLines := ParseJSONLines(expected)
	actualLines := ParseJSONLines(actual)
	require.Equal(t, len(expectedLines), len(actualLines))

	for i, l := range expectedLines {
		require.JSONEq(t, string(l), string(actualLines[i]))
	}
}

func ParseJSONLines(data []byte) [][]byte {
	return bytes.Split(bytes.TrimSpace(data), []byte("\n"))
}

func GetHEADRef(t *testing.T, testPath string) string {
	t.Helper()

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = testPath

	out, err := cmd.Output()
	require.NoError(t, err)

	return string(out)
}

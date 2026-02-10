package template

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
)

func getTestCases(t *testing.T) []string {
	entries, err := os.ReadDir("./")
	require.NoError(t, err, "failed to read test cases directory")

	var testCases []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "common" {
			testCases = append(testCases, entry.Name())
		}
	}
	return testCases
}

func TestTemplateRendering(t *testing.T) {
	testCases := getTestCases(t)
	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			chartPath, err := filepath.Abs(relativeChartPath)
			require.NoError(t, err)

			releaseName := testCase
			valuesFile := fmt.Sprintf("%s/values.yaml", testCase)
			compareToGoldenFile(t, renderHelmTemplate(t, releaseName, chartPath, valuesFile, clusterTemplateName), fmt.Sprintf("%s/cluster.yaml", testCase))
			compareToGoldenFile(t, renderHelmTemplate(t, releaseName, chartPath, valuesFile, "secret"), "common/secret.yaml")
			for _, templateName := range disabledByDefaultTemplates {
				targetFilePath := fmt.Sprintf("%s/%s.yaml", testCase, templateName)
				if fileExists(targetFilePath) {
					rendered := renderHelmTemplate(t, releaseName, chartPath, valuesFile, templateName)
					compareToGoldenFile(t, rendered, targetFilePath)
				} else {
					verifyTemplateIsNotRendered(t, releaseName, chartPath, valuesFile, templateName)
				}
			}
		})
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func compareToGoldenFile(t *testing.T, renderedYAML []byte, goldenFilePath string) {
	goldenData, err := os.ReadFile(goldenFilePath)
	require.NoError(t, err, "failed to read golden file: %s", goldenFilePath)

	renderedStr := string(filterYamlComments(unifyLineEndings(renderedYAML)))
	goldenStr := string(unifyLineEndings(goldenData))

	// Compare the content
	if renderedStr == goldenStr {
		return // Files match
	}

	// Generate a unified diff for better error reporting
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(goldenStr),
		B:        difflib.SplitLines(renderedStr),
		FromFile: "Golden",
		ToFile:   "Rendered",
	}

	diffText, err := difflib.GetUnifiedDiffString(diff)
	require.NoError(t, err, "failed to generate diff")

	t.Errorf("Rendered YAML does not match golden file %s:\n%s", goldenFilePath, diffText)
}

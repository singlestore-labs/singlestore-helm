package template

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

var helmPath *string

func getHelmPath(t *testing.T) string {
	if helmPath == nil {
		path, err := exec.LookPath("helm")
		require.NoError(t, err)
		helmPath = &path
	}
	return *helmPath
}

func getRenderTemplateCommand(t *testing.T, releaseName string, chartPath string, valuesFile string, templateName string) *exec.Cmd {
	return exec.Command(getHelmPath(t), "template", releaseName, chartPath, "--values", valuesFile, "--show-only", getTemplatePath(templateName))
}

func renderHelmTemplate(t *testing.T, releaseName string, chartPath string, valuesFile string, templateName string) []byte {
	var output, errorOutput bytes.Buffer
	cmd := getRenderTemplateCommand(t, releaseName, chartPath, valuesFile, templateName)
	cmd.Stdout = &output
	cmd.Stderr = &errorOutput
	err := cmd.Run()
	if err != nil {
		fmt.Print(errorOutput.String())
	}
	require.NoError(t, err)
	return output.Bytes()
}

func verifyTemplateIsNotRendered(t *testing.T, releaseName string, chartPath string, valuesFile string, templateName string) {
	var errorOutput bytes.Buffer
	cmd := getRenderTemplateCommand(t, releaseName, chartPath, valuesFile, templateName)
	cmd.Stderr = &errorOutput
	err := cmd.Run()
	require.Error(t, err)
	require.Contains(t, errorOutput.String(), fmt.Sprintf("Error: could not find template templates/%s.yaml in chart", templateName))
}

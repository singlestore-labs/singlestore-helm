package template

import "bytes"

func filterYamlComments(yamlContent []byte) []byte {
	isSeparatorLine := func(line []byte) bool {
		for _, b := range line {
			if b != '-' {
				return false
			}
		}
		return true
	}
	lines := bytes.Split(yamlContent, []byte("\n"))

	var filteredLines [][]byte
	removedTopSeparator := false
	for _, line := range lines {
		trimmedLine := bytes.TrimSpace(line)
		if len(trimmedLine) == 0 || bytes.HasPrefix(trimmedLine, []byte("# Source")) {
			continue // Skip empty lines and source comments
		}
		if !removedTopSeparator && isSeparatorLine(trimmedLine) {
			removedTopSeparator = true
			continue
		}
		filteredLines = append(filteredLines, line)
	}

	return bytes.Join(filteredLines, []byte("\n"))
}

func unifyLineEndings(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
}

package filelib

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Wrapper to ReadKVFile().
// Looks in specific locations for the designated config file.
func ReadConfig(filename string) (map[string]string, error) {
	var err error = nil
	data := make(map[string]string)
	paths := []string{"/etc", "etc", "."}

	for _, path := range paths {
		data, err = ReadKVFile(filepath.Join(path, filename), ":")
		if err == nil {
			break
		}
	}
	return data, err
}

// Read a text file that contains simple key/value pairs (one per line)
// separated by the string specified in the sep param.
// Returns a map of type map[string]string.
func ReadKVFile(filepath string, sep string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { // Skip empty lines and comments
			continue
		}
		elements := strings.SplitN(line, sep, 2)
		if len(elements) == 2 {
			key := strings.TrimSpace(elements[0])
			value := strings.TrimSpace(elements[1])
			result[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

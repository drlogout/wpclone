package util

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gosimple/slug"
)

func NoTrailingSlash(path string) string {
	if path[len(path)-1:] == "/" {
		return path[:len(path)-1]
	}

	return path
}

func RemoveLineInFile(path string, regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("error compiling regex: %s", err)
	}

	inputFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer inputFile.Close()

	// Create a temporary file to write the filtered content
	tempFilepath := path + ".tmp"
	tempFile, err := os.Create(tempFilepath)
	if err != nil {
		return fmt.Errorf("error creating temp file: %s", err)
	}
	defer tempFile.Close()

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(tempFile)

	// Read the input file line by line
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return fmt.Errorf("error reading file: %s", err)
			}
			break
		}

		// Write the line to the temp file if it doesn't match the line to remove
		if !r.MatchString(strings.TrimSpace(line)) {
			_, err := writer.WriteString(line)
			if err != nil {
				return fmt.Errorf("error writing to temp file: %s", err)
			}
		}
	}

	// Flush the writer to ensure all content is written to the temp file
	writer.Flush()

	// Replace the original file with the temp file
	err = os.Rename(tempFilepath, path)
	if err != nil {
		return fmt.Errorf("error replacing original file: %s", err)
	}

	return nil
}

func Slug(s string) string {
	sl := slug.Make(s)
	sl = strings.ReplaceAll(sl, "-", "_")

	// split the string if longer than 32 characters
	if len(sl) > 32 {
		sl = sl[:32]
	}

	return sl
}

package input

import (
	"bufio"
	"log"
	"os"
)

// LoadTextLines reads a text file and returns its lines as members of a string array
func LoadTextLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 { // doesn't add blank lines.. maybe add bool arg to change that?
			data = append(data, line)
		}
	}

	return data
}

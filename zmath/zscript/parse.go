package zscript

import "os"

// ParseFile will parse a .zs file and execute it
func ParseFile(f os.File) {

}

// GetLines reads an os.File and returns its data as a bunch of strings broken up by lines
func GetLines(f os.File) []string {
	lines := make([]string, 0, 10)

	for {
		var (
			line     string = ""
			fileDone bool   = false
		)
		for {
			b := make([]byte, 1, 1)
			n, _ := f.Read(b)

			// if the end of the file
			if n == 0 {
				lines = append(lines, line)
				fileDone = true
				break
			}

			// if the end of a line
			if b[0] == '\n' {
				lines = append(lines, line)
				break
			}

			// if safe, append just-read byte to the end of the current line
			line += string(b)
		}

		if fileDone {
			break
		}
	}

	return lines
}

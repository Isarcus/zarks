package input

import (
	"fmt"
	"os"
)

// FileExists tells you whether or not a file exists at the specified path!
// Also, this function is partially copied off Stack Exchange... hope it works!
func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true

	} else if os.IsNotExist(err) {
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		fmt.Println("\n[ERROR] File detection got wonky:")
		fmt.Println(err)
	}
	return false
}

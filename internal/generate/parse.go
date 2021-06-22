package generate

import (
	"bufio"
	"os"
)

// getStructString gets the portion of the file containing the struct. It is brittle.
// TODO: do this better
func getStructString(source string) string {

	var structString string

	file, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	first := true
	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		structString += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return structString
}

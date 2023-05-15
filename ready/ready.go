package ready

import (
	"bufio"
	"fmt"
	"os"
)

func PromptAndRead(message string) (string, error) {
	fmt.Fprint(os.Stdout, message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

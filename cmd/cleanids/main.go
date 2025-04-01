package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	ids := make(map[string]bool)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ids[scanner.Text()] = true
	}
	for id := range ids {
		fmt.Fprintln(os.Stdout, id)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to shitdb.")
	for {
		fmt.Print("# ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		fmt.Println("ok")
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/DylanGuedes/weakdb/core"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to weakdb.")
	for {
		fmt.Print("# ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		_, err := core.Parse(text)

		if err != nil {
			panic(err)
		}

		fmt.Println("ok")
	}
}

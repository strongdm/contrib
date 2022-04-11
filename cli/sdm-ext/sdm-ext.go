package main

import (
	"ext/cli"
	"fmt"
	"os"
)

func main() {
	err := cli.NewApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"ext/cli"
	"ext/util"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(fmt.Sprintf("%s/.env", util.GetBasePath()))
	cli.Main()
}

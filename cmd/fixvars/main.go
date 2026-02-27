package main

import (
	fixvars "gofix-vars"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(fixvars.Analyzer)
}

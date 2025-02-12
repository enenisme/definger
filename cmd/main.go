package main

import (
	"os"

	"github.com/enenisme/definger/cli"
	"github.com/enenisme/definger/flag"
)

func main() {
	app := flag.NewFlag()
	app.Action = cli.Run
	app.Run(os.Args)
}

package summarizer

import (
	"flag"
)

func CmdFlags(args []string) *string {
	fs := flag.NewFlagSet("summarizer", flag.ContinueOnError)

	cpath := new(string)
	fs.StringVar(cpath, "c", "", "Path to config file")

	if err := fs.Parse(args[1:]); err != nil {
		return nil
	}

	return cpath
}

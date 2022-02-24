package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/robbyriverside/brief"
)

var SemVer = "unknown"

type options struct {
	Args struct {
		File string `positional-arg-name:"file" description:"brief file"`
	} `positional-args:"true" required:"true"`
	Verbose bool `short:"v" long:"verbose" description:"verbose output"`
	Version bool `long:"version" description:"describe version"`
}

func main() {
	opt := &options{}
	parser := flags.NewParser(opt, flags.Default)
	parser.Name = "brief"

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return
		}
	}
	if opt.Version {
		fmt.Println("brief", SemVer)
		return
	}
	dec, err := brief.NewFileDecoder(opt.Args.File)
	if err != nil {
		log.Fatal(err)
	}
	dec.Debug = opt.Verbose
	nodes, err := dec.Decode()
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range nodes {
		if dec.Debug {
			fmt.Println(node)
			continue
		}
		out := node.Encode()
		fmt.Println(string(out))
	}
}

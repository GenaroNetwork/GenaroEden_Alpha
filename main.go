package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

const version = "v0.1.0"

var app = &cli.Command{
	Name: os.Args[0],
	Desc: "Genaro Client",
	Text: `lastest version of the Genaro Client, Genaro Network is the decentralized network`,
	Argv: func() interface{} { return new(gogoT) },
	Fn:   gogo,
}

type gogoT struct {
	cli.Helper
	Version bool `cli:"v,version" usage:"display version"`
	List    bool `cli:"l,list" usage:"list all sub commands or not" dft:"false"`
}

func gogo(ctx *cli.Context) error {
	argv := ctx.Argv().(*gogoT)
	if argv.Version {
		ctx.String(version + "\n")
		return nil
	}

	if argv.List {
		ctx.String(ctx.Command().ChildrenDescriptions(" ", "  =>  "))
		return nil
	}

	ctx.String("try `%s --help for more information'\n", ctx.Path())
	return nil
}

func jsonIndent(i interface{}) string {
	data, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		return ""
	}
	return string(data) + "\n"
}

func main() {
	cli.SetUsageStyle(cli.ManualStyle)
	//NOTE: You can set any writer implements io.Writer
	// default writer is os.Stdout
	if err := app.RunWith(os.Args[1:], os.Stderr, nil); err != nil {
		fmt.Printf("%v\n", err)
	}
}

var ( bucketCommands = &cli.Command{Name: "bucket"}
	  fileCommands = &cli.Command{Name: "file"})


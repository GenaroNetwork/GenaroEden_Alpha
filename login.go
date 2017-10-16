package main

import (
	"github.com/mkideal/cli"
)

//need to use flag to input username and password

var _ = app.Register(&cli.Command{
	Name: "login",
	Desc: "Genaro User login",
	Argv: func() interface{} { return new(loginT) },
	OnBefore: func(ctx *cli.Context) error {
		ctx.String("user name and password is loaded\n")
		return nil
	},		
	Fn:   login,
})

type loginT struct {
	cli.Helper
	Username    string `cli:"u,username" usage:"login username"`
	Password    string `cli:"p,password" usage:"login password"`
}

func login(ctx *cli.Context) error {
	argv := ctx.Argv().(*loginT)
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
package main

/*
#include <microhttpd.h>
#include <assert.h>
#include <stdio.h>
#include <uv.h>

#cgo CFLAGS: -I.
#cgo LDFLAGS: -L ./ -lstorj -luv

#include "storj.h"
#include "callbacks.h"
*/
import "C"
import "errors"
import "unsafe"
import "fmt"
import (
	"github.com/mkideal/cli"
)

var accountCommand = &cli.Command{
	Name: "account",
	Desc: "Genaro account Menu",
	Argv: func() interface{} { return new(accountT) },
	Fn: func(ctx *cli.Context) error {
		return nil
	},
}

var _ = app.Register(accountCommand)

type accountT struct{
	cli.Helper
	Version string `cli:"v,version" usage:"show version" dft:"v0.1.0"`	
}

var _ =  accountCommand.Register(loginCommand)
var _ = accountCommand.Register(exportCommand)


var loginCommand = &cli.Command{
	Name: "login",
	Desc: "Genaro account login",
	Argv: func() interface{} { return new(loginT) },	
	Fn:   login,
}

type loginT struct {
	cli.Helper
}

func login(ctx *cli.Context) error {

	var env *C.storj_env_t
	
	if err := init_env_n(&env); err != nil{
		fmt.Printf("that is a bad thing ")
	}

	user_opts := user_options{
		nil, 
		nil, 
		// C.CString("localhost"),
		C.CString("101.132.159.197"), 
		nil, 
		nil}

	Sts := make(chan C.int)	

	go func(){

    	status := C.import_keys((*C.user_options_t)(unsafe.Pointer(&user_opts)))	


		Sts<-status

	}()

	<-Sts

    if ( user_opts.user == C.CString("") || user_opts.pass == C.CString("")) {
		return errors.New("env is not set properly")	
    }	

	return nil
}




var exportCommand = &cli.Command{
	Name: "export",
	Desc: "Genaro account export",
	Argv: func() interface{} { return new(exportT) },	
	Fn:   export,
}

type exportT struct {
	cli.Helper
}

func export(ctx *cli.Context) error {

	var env *C.storj_env_t
	
	if err := init_env_n(&env); err != nil{
		fmt.Printf("that is a bad thing ")
	}

	Sts := make(chan C.int)	
	host := C.CString("101.132.159.197")

	go func(){

    	status := C.export_keys(host)


		Sts<-status

	}()

	<-Sts

	return nil
}

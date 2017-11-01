package main

import (
	"github.com/mkideal/cli"
)


/*
#include <microhttpd.h>
#include <assert.h>

#cgo CFLAGS: -I.
#cgo LDFLAGS: -L ./ -lstorj

#include "storj.h"

*/
import "C"
import "unsafe"

var _ = app.Register(&cli.Command{
	Name: "keygen",
	Desc: "Genaro Key Genaration",
	Argv: func() interface{} { return new(keygenT) },
	OnBefore: func(ctx *cli.Context) error {
		ctx.String("Please store your key in private safe place\n")
		return nil
	},	
	Fn:   keygen,
})

type keygenT struct {
	cli.Helper
	// key string `cli:"key,k" usage:"genearte genaro private key"`
}

func keygen(ctx *cli.Context) error {

	buf := make([]string, 16)
	// argv := ctx.Argv().(*keygenT)


	//generate 16 size buffer for generate key
	arg := make([](*_Ctype_char),0);

	for i,_ := range buf{
		char := C.CString(buf[i])
		defer C.free(unsafe.Pointer(char))
		strptr := (*_Ctype_char)(unsafe.Pointer(char))
		arg = append(arg,strptr)
	}

	C.storj_mnemonic_generate(128, (**_Ctype_char)(unsafe.Pointer(&arg[0])));

	//return back to go string


	for i,_ := range arg{
		buf[i]=C.GoString(arg[i])
	}

	ctx.String("%s: %s \n", ctx.Path(), buf)

	// if env !=nil {
	// 	C.storj_destroy_env(env);
	// }

	return nil
}
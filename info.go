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
import "unsafe"
import "errors"
import (
	"github.com/mkideal/cli"
)


var _ = app.Register(&cli.Command{
	Name: "info",
	Desc: "Genaro info ",
	Argv: func() interface{} { return new(infoT) },	
	Fn:   info,
})

type infoT struct {
	cli.Helper
}

func info(ctx *cli.Context) error {

	var env *C.storj_env_t
	
	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not init properly")
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}


	// User := os.Getenv("STORJ_BRIDGE_USER")
	// Pass := os.Getenv("STORJ_BRIDGE_PASS")	
	// Key := os.Getenv("STORJ_ENCRYPTION_KEY")
	// Home := os.Getenv("STORJ_KEYPASS")

	// fmt.Printf("%v\n",Home)
	// fmt.Printf("%v\n",User)
	// fmt.Printf("not here ? \n")
	// fmt.Printf("%v\n",Pass)
	// fmt.Printf("%v\n",Key)

	Sts := make(chan C.int)	

	go func(){

    	status := C.storj_bridge_get_info(env, nil,(C.uv_after_work_cb)(unsafe.Pointer(C.get_info_callback)))	 		

		C.uv_run((env).loop, C.UV_RUN_DEFAULT)

		Sts<-status

	}()

	<-Sts
	return nil
}
package main

/*
#include <microhttpd.h>
#include <assert.h>
#include <stdio.h>
#include <uv.h>

#cgo CFLAGS: -I.
#cgo LDFLAGS: -L ./ -lstorj -luv

#include "../libstorj/src/storj.h"
#include "../libstorj/src/uploader.h"
#include "callbacks.h"

*/
import "C"
import "unsafe"
import "fmt"
import "errors"
import (
	"github.com/mkideal/cli"
	// "github.com/stretchr/testify/assert"
)

type storj_bridge_options struct{
	proto *C.char;
	host  *C.char;
	port  C.int;
	user  *C.char;
	pass  *C.char;
}

type storj_encrypt_options struct{
	mnemonic *C.char;
}

type storj_http_options struct{
	user_agent *C.char;
	low_speed_limit C.uint64_t;
	low_speed_time	C.uint64_t;
	timeout C.uint64_t;
}

type storj_log_options struct{
	level C.int;
}

type storj_upload_opts struct{
	index string;
	bucket_id string;
	file_name string;
	fd *C.FILE;
	rs bool;
}

var fileCommand = &cli.Command{
	Name: "file",
	Desc: "Genaro File Menu",
	Argv: func() interface{} { return new(fileT) },
	Fn: func(ctx *cli.Context) error {
		// assert.Equal(t, ctx.Path(), "fileCommand")
		// argv := ctx.Argv().(*argT)
		// assert.Equal(t, argv.Version, "v0.1.0")
		// assert.Equal(t, ctx.Command().Name, "fileCommand")
		return nil
	},
}

var _ = app.Register(fileCommand)

type fileT struct{
	cli.Helper
	Version string `cli:"v,version" usage:"show version" dft:"v0.1.0"`	
}

var _ =  fileCommand.Register(uploadCommand)
var _ = fileCommand.Register(downloadCommand)
var _ = fileCommand.Register(rmfileCommand)

func init_env(env_ptr **C.storj_env_t) error{

	var env *C.storj_env_t

	Env := make(chan *C.storj_env_t) 

	// proto := C.CString("http")
	// host := C.CString("172.19.0.7")
	// user := C.CString("test@storj.io")
	// pass := C.CString("password")
	// mnemonic := C.CString("resemble scorpion weasel gift retreat pigeon piece liar shuffle mind best arctic slender quiz strong jeans misery wide tobacco pact firm wet success again")
	// user_agent := C.CString("libstorj-1.1.0-beta")


	proto := C.CString("http")
	host := C.CString("localhost")
	user := C.CString("begoingto1@163.com")
	pass := C.CString("ashes8871")
	mnemonic := C.CString("negative brown polar admit eagle return valley host weather simple oak assume")
	user_agent := C.CString("libstorj-1.1.0-beta")

	options := storj_bridge_options{
		proto: proto, 
		host: host,
		port: C.int(8080),
		user: user,
		pass: pass}

	encrypt_options := storj_encrypt_options{
		mnemonic : mnemonic}

	http_options := storj_http_options{
	    user_agent : user_agent,
	    low_speed_limit : C.uint64_t(0),
	    low_speed_time : C.uint64_t(0),
	    timeout : C.uint64_t(0)}

	log_options := storj_log_options{
		level : C.int(0)}

	go func(){
		
		env = C.storj_init_env((*C.storj_bridge_options_t)(unsafe.Pointer(&options)),(*C.storj_encrypt_options_t)(unsafe.Pointer(&encrypt_options)),(*C.storj_http_options_t)(unsafe.Pointer(&http_options)),(*C.storj_log_options_t)(unsafe.Pointer(&log_options)))
	
		C.uv_run(env.loop, C.UV_RUN_DEFAULT)
		
		Env <- env
	}()

	*env_ptr = <- Env	

	return nil
}


var uploadCommand = &cli.Command{
	Name: "set",
	Desc: "Genaro File upload",
	Argv: func() interface{} { return new(uploadT) },	
	Fn: genaro_upload,
}


type uploadT struct {
	cli.Helper
	Bucket string `cli:"id,i" usage:"bucket id to be loaded on Genaro network"`
	Path string `cli:"path,p" usage:"upload file path on Genaro network"`	
}


func genaro_upload(ctx *cli.Context) error {
	argv := ctx.Argv().(*uploadT)

	var env *C.storj_env_t
		
	Sts := make(chan C.int)		

	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly");
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}	

	fmt.Printf("%v\n",C.GoString(env.bridge_options.host))   

	bucket_id := C.CString(argv.Bucket)
	file_path := C.CString(argv.Path)

	// file_path := C.CString("/Users/weilongwu/GenaroCore_Storj/genaro_cli/test_upload.data")

	go func(){


		status := C.upload_file(env, bucket_id, file_path);	

	
		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status	

	}()

	<-Sts

	return nil;

}

//===========================================================================================

var downloadCommand =&cli.Command{
	Name: "get",
	Desc: "Genaro File download",
	Argv: func() interface{} { return new(downloadFileT) },	
	Fn: genaro_download,
}

type downloadFileT struct {
	cli.Helper
	BucketId string `cli:"b" usage:"BucketId from your download file on Genaro network"`
	FileId string `cli:"f" usage:"download FileId on Genaro network"`
	Path string `cli:"p" usage:"download path to your computer"`
}

func genaro_download(ctx *cli.Context) error {

	argv := ctx.Argv().(*downloadFileT)

	var env *C.storj_env_t
		
	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly");
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}		

	Sts := make(chan C.int)
	

	bucket_id := C.CString(argv.BucketId)
	file_id := C.CString(argv.FileId)
	path := C.CString(argv.Path)

	go func(){

		status := C.download_file(env, bucket_id, file_id, path)
	
		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status	

		fmt.Printf("%v\n",status)	

	}()	

	<-Sts

	return nil;

}

//========================================================================================

var rmfileCommand = &cli.Command{
	Name: "rm",
	Desc: "Genaro File remove",
	Argv: func() interface{} { return new(genaro_removeT) },	
	Fn: genaro_removefile,
}

type genaro_removeT struct {
	cli.Helper
	FileId string `cli:"f" usage:"remove file from Genaro network"`
	BucketId string `cli:"b" usage:"remove file from Genaro network"`	
}

func genaro_removefile(ctx *cli.Context) error {

	argv := ctx.Argv().(*genaro_removeT)

	Sts := make(chan C.int)
	
	var env *C.storj_env_t
		
	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly");
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}	

	bucket_id := C.CString(argv.BucketId)
	file_id := C.CString(argv.FileId)


	go func(){

		status := C.storj_bridge_delete_file(env, bucket_id, file_id, nil,(C.uv_after_work_cb)(unsafe.Pointer(C.delete_file_callback)))
		// defer C.fclose(file_go)
	
		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status	

	}()

	<-Sts


	return nil;

}

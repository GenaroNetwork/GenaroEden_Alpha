package main

/*
#include <microhttpd.h>
#include <assert.h>
#include <stdio.h>

#include "callbacks.h"

#cgo CFLAGS: -I.
#cgo LDFLAGS: -L ./ -lstorj -luv -ljson-c
#include "storj.h"

extern void create_bucket(uv_work_t *work_req, int status);
extern void delete_bucket(uv_work_t *work_req, int status); 
*/
import "C"
import "unsafe"
import "fmt"
import "errors"
import (
	"github.com/mkideal/cli"
	// "github.com/stretchr/testify/assert"
)

type callbackFunc *func(total_buckets C.uint32_t, created *C.char, name *C.char,id *C.char)

var bucketCommand = &cli.Command{
	Name: "bucket",
	Desc: "Genaro Bucket Menu",
	Argv: func() interface{} { return new(bucketT) },
	Fn: func(ctx *cli.Context) error {
		// assert.Equal(t, ctx.Path(), "bucketCommand")
		// argv := ctx.Argv().(*argT)
		// assert.Equal(t, argv.Version, "v0.1.0")
		// assert.Equal(t, ctx.Command().Name, "bucketCommand")
		return nil
	},
}

type bucketT struct{
	cli.Helper
	Version string `cli:"v,version" usage:"show version" dft:"v0.1.0"`	
}


	var bkcmd = app.Register(bucketCommand)

	var _ = bkcmd.Register(addCommand)
	var _ = bkcmd.Register(removeCommand)
	var _ = bkcmd.Register(listbucketCommand)
	var _ =bkcmd.Register(listfileCommand)



var addCommand = &cli.Command{
	Name: "addbucket",
	Desc: "Genaro Bucket Adding",
	Argv: func() interface{} { return new(addbucketT) },
	Fn: addbucket,
}

type addbucketT struct {
	cli.Helper
	Bucket string `cli:"name,n" usage:"new bucket name on Genaro network"`
}

type initT struct {
	cli.Helper
	Host string `cli:"host" usage:"set host on Genaro network"`
	User string `cli:"user,u" usage:"username on Genaro network"`
	Password string `cli:"password,p" usage:"user password on Genaro network"`
}

func addbucket(ctx *cli.Context) error {
	argv := ctx.Argv().(*addbucketT)
	bucket_name := C.CString(argv.Bucket)	
	Sts := make(chan C.int)

	var env *C.storj_env_t
		
	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly")
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}	

	if argv.Bucket != "" {
		fmt.Printf("%v bucket is created \n", argv.Bucket)	
	}else{
		return errors.New("bucket name is needed for create bucket");
	}

	go func(){


    	status := C.storj_bridge_create_bucket(env, bucket_name, nil, (C.uv_after_work_cb)(unsafe.Pointer(C.create_bucket)))			
		

		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status

	}()
	
	<-Sts

	return nil;		
}

var removeCommand = &cli.Command{
	Name: "removebucket",
	Desc: "Genaro Bucket remove",
	Argv: func() interface{} { return new(removebucketT) },	
	Fn: removebucket,
}

type removebucketT struct {
	cli.Helper
	Bucket string `cli:"id,i" usage:"remove bucket(id) on Genaro network"`
}

func removebucket(ctx *cli.Context) error {
	argv := ctx.Argv().(*removebucketT)
	bucket_id := C.CString(argv.Bucket) 
	Sts := make(chan C.int)	

	var env *C.storj_env_t
		
	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly")
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}	

	go func(){

		status := C.storj_bridge_delete_bucket(env, bucket_id, nil, (C.uv_after_work_cb)(unsafe.Pointer(C.delete_bucket)))			
		
		fmt.Printf("start loading \n")

		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status

	}()
	
	<-Sts

	// C.storj_destroy_env(env);

	return nil;	

}


//=================================================================================================================

var listbucketCommand =&cli.Command{
	Name: "listbuckets",
	Desc: "Genaro Buskets Listing",
	Argv: func() interface{} { return new(listbucketsT) },	
	Fn: listbuckets,
}

type listbucketsT struct {
	cli.Helper
	List string `cli:"list-buskets" usage:"user Genaro buskets list"`
}

func listbuckets(ctx *cli.Context) error{

	var env *C.storj_env_t
	
	Sts := make(chan C.int)		

	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly");
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}	

	go func(){
		status := C.storj_bridge_get_buckets(env, nil, (C.uv_after_work_cb)(unsafe.Pointer(C.get_buckets_callback)))		

		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status

	}()
	
	<-Sts

	// fmt.Println("%s", env.loop)

	// C.storj_destroy_env(env);

	return nil;	

}

var listfileCommand = &cli.Command{
	Name: "listfiles",
	Desc: "Genaro Busket File Listing",
	Argv: func() interface{} { return new(listfilesT) },	
	Fn: listfiles,
}

type listfilesT struct {
	cli.Helper
	Bucket string `cli:"id,i" usage:"check files in certain busket id in Genaro"`
}

func listfiles(ctx *cli.Context) error {

	argv := ctx.Argv().(*listfilesT)
	var env *C.storj_env_t
	
	Sts := make(chan C.int)

	// if err := init_env(&env); err != nil{
	// 	return errors.New("env is not set properly");
	// }

	if err := set_env(&env); err != nil{
		return errors.New("Unlock passphrase is not correct")
	}	

	bucket_id := C.CString(argv.Bucket) 


	go func(){

		defer C.free(unsafe.Pointer(bucket_id))

		status := C.storj_bridge_list_files(env, bucket_id, nil, (C.uv_after_work_cb)(unsafe.Pointer(C.list_files_callback)))			


		C.uv_run(env.loop, C.UV_RUN_DEFAULT)

		Sts<-status


	}()
	
	<-Sts

	return nil;	
}		
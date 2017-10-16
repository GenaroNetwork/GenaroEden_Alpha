package main

/*
#include <microhttpd.h>
#include <assert.h>
#include <stdio.h>
#include <uv.h>

#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. /usr/local/lib/libstorj.dylib /usr/local/lib/libuv.dylib

#include "../libstorj/src/storj.h"

void check_store_file_progress(double progress,
                               uint64_t uploaded_bytes,
                               uint64_t total_bytes,
                               void *handle)
{
    assert(handle == NULL);
    if (progress == (double)1) {
        printf("success");
    }
}

void check_store_file(int error_code, char *file_id, void *handle)
{
    assert(handle == NULL);
    if (error_code == 0) {
        if (strcmp(file_id, "85fb0ed00de1196dc22e0f6d") == 0 ) {
			printf("success");
        } 
    } else {
        printf("\t\tERROR:   %s\n", storj_strerror(error_code));
    }

    free(file_id);
}


void check_resolve_file_progress(double progress,
                                 uint64_t downloaded_bytes,
                                 uint64_t total_bytes,
                                 void *handle)
{
    assert(handle == NULL);
    if (progress == (double)1) {
        printf("success");
    }

    // TODO check error case
}

void check_resolve_file(int status, FILE *fd, void *handle)
{
    fclose(fd);
    assert(handle == NULL);
    if (status) {
        printf("Download failed: %s\n", storj_strerror(status));
    } 
}


void check_delete_file(uv_work_t *work_req, int status)
{
    assert(status == 0);
    json_request_t *req = work_req->data;
    assert(req->handle == NULL);
    assert(req->response == NULL);
    assert(req->status_code == 200);

    free(req->path);
    free(req);
    free(work_req);
}
*/
import "C"
import "unsafe"
import "fmt"
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

var uploadCommand = &cli.Command{
	Name: "set",
	Desc: "Genaro File upload",
	Argv: func() interface{} { return new(uploadT) },	
	Fn: genaro_upload,
}

// var _ = app.Register(uploadCommand)

// var _ = app.Register(&cli.Command{
// 	Name: "set",
// 	Desc: "Genaro File upload",
// 	Argv: func() interface{} { return new(uploadT) },	
// 	Fn: genaro_upload,
// })

type uploadT struct {
	cli.Helper
	Host string `cli:"host" usage:"set host on Genaro network"`
	User string `cli:"user,u" usage:"username on Genaro network"`
	Password string `cli:"password,p" usage:"user password on Genaro network"`
	File string `cli:"file" usage:"upload file"`	
}


func string_to_char(dst *C.char, src string) error {
	src_b := []byte(src)
	dst = (*C.char)(unsafe.Pointer(&src_b[0]))

	return nil
}

//export init_env_go
func init_env_go(ctx *cli.Context, env *C.storj_env_t) error {
	// argv := ctx.Argv().(*uploadT)

	proto := "http"
	host := "172.19.0.6"
	user := "test@storj.io"
	pass := "password"
	mnemonic := "resemble scorpion weasel gift retreat pigeon piece liar shuffle mind best arctic slender quiz strong jeans misery wide tobacco pact firm wet success again"
	user_agent := "libstorj-1.1.0-beta"

	var proto_c *C.char
	var host_c *C.char
	var user_c *C.char
	var pass_c *C.char

	var mnemonic_c *C.char

	var user_agent_c *C.char

	Env := make(chan *C.storj_env_t) 

	string_to_char(proto_c,proto)
	string_to_char(host_c,host)
	string_to_char(user_c,user)
	string_to_char(pass_c,pass)

	string_to_char(mnemonic_c,mnemonic)

	string_to_char(user_agent_c,user_agent)

	options := storj_bridge_options{
		proto: proto_c, 
		host: host_c,
		port: C.int(8080),
		user: user_c,
		pass: pass_c}

	encrypt_options := storj_encrypt_options{
		mnemonic : mnemonic_c}

	http_options := storj_http_options{
	    user_agent : user_agent_c,
	    low_speed_limit : C.uint64_t(0),
	    low_speed_time : C.uint64_t(0),
	    timeout : C.uint64_t(0)}

	log_options := storj_log_options{
		level : C.int(0)}


	go func(){
		
		env = C.storj_init_env((*C.storj_bridge_options_t)(unsafe.Pointer(&options)),(*C.storj_encrypt_options_t)(unsafe.Pointer(&encrypt_options)),(*C.storj_http_options_t)(unsafe.Pointer(&http_options)),(*C.storj_log_options_t)(unsafe.Pointer(&log_options)))
		
		
		defer C.free((unsafe.Pointer)(&env))

		C.uv_run(env.loop, C.UV_RUN_DEFAULT)
		
		Env <- env

		}()
	// env = C.storj_init_env((*C.storj_bridge_options_t)(unsafe.Pointer(&options)),(*C.storj_encrypt_options_t)(unsafe.Pointer(&encrypt_options)),(*C.storj_http_options_t)(unsafe.Pointer(&http_options)),(*C.storj_log_options_t)(unsafe.Pointer(&log_options)))

	fmt.Printf("%v\n", <-Env)

	return nil
}

func genaro_upload(ctx *cli.Context) error {
	// env := make([](*C.storj_env_t),1)
	var env *C.storj_env_t
	init_env(ctx,env)

	file := "test_upload"
	file_name := "test_upload.data"
	file_opt := "r"
	var file_c *C.char
	var file_opt_c *C.char

	string_to_char(file_c,file_name)
	string_to_char(file_opt_c,file_opt)
	var file_go *C.FILE
	file_go = C.fopen(file_c,file_opt_c)
	C.fclose(file_go);

	upload_opts := storj_upload_opts{
		index : "d2891da46d9c3bf42ad619ceddc1b6621f83e6cb74e6b6b6bc96bdbfaefb8692",
		bucket_id : "368be0816766b28fd5f43af5",
		file_name : file,
		fd : file_go,
		rs : true}

	var state C.storj_upload_state_t

	go func(){

		status := C.storj_bridge_store_file(env,(*C.storj_upload_state_t)(unsafe.Pointer(&state)),(*C.storj_upload_opts_t)(unsafe.Pointer(&upload_opts)),nil,(C.storj_progress_cb)(unsafe.Pointer(C.check_store_file)),(C.storj_finished_upload_cb)(unsafe.Pointer(C.check_store_file_progress)))

		panic(status)
	}()
	// assert.Equal(status,0)
	// assert(status == 0);


	if(C.int(C.uv_run(env.loop,C.UV_RUN_DEFAULT)) != 0){
	return nil;
	}

	C.storj_destroy_env(env);

	return nil;

}

//===========================================================================================

var downloadCommand =&cli.Command{
	Name: "get",
	Desc: "Genaro File download",
	Argv: func() interface{} { return new(dwnloadT) },	
	Fn: genaro_download,
}

// var _ = app.Register(downloadCommand)

// var _ = app.Register(&cli.Command{
// 	Name: "get",
// 	Desc: "Genaro File download",
// 	Argv: func() interface{} { return new(dwnloadT) },	
// 	Fn: genaro_download,
// })

type dwnloadT struct {
	cli.Helper
	Host string `cli:"host" usage:"set host on Genaro network"`
	User string `cli:"user,u" usage:"username on Genaro network"`
	Password string `cli:"password,p" usage:"user password on Genaro network"`
	File string `cli:"file" usage:"download file"`	
}

func genaro_download(ctx *cli.Context) error {
	// env := make([](*C.storj_env_t),1)
	var env *C.storj_env_t
	init_env(ctx,env)

	file_name := "test_download.data"
	file_opt := "w+"
	bucket_id := "368be0816766b28fd5f43af5"
	file_id	:= "998960317b6725a3f8080c2b"

	var file_c *C.char
	var file_opt_c *C.char
	var bucket_id_c *C.char
	var file_id_c *C.char

	string_to_char(file_c, file_name)
	string_to_char(file_opt_c, file_opt)
	string_to_char(bucket_id_c, bucket_id)
	string_to_char(file_id_c, file_id)

	var download_file *C.FILE
	download_file = C.fopen(file_c,file_opt_c)
	C.fclose(download_file)
	var handle int

	var state C.storj_download_state_t
	// problem here the download file should be *FILE but in here it needs void* 
	status := C.storj_bridge_resolve_file(env,(*C.storj_download_state_t)(unsafe.Pointer(&state)),bucket_id_c,file_id_c,download_file,unsafe.Pointer(&handle),(C.storj_progress_cb)(unsafe.Pointer(C.check_resolve_file)),(C.storj_finished_download_cb)(unsafe.Pointer(C.check_resolve_file_progress)))

	panic(status)
	// assert(status == 0);

	if(C.int(C.uv_run(env.loop,C.UV_RUN_DEFAULT)) != 0){
		return nil;
	}

	C.storj_destroy_env(env);

	return nil;

}

//========================================================================================

var rmfileCommand = &cli.Command{
	Name: "rm",
	Desc: "Genaro File remove",
	Argv: func() interface{} { return new(genaro_removeT) },	
	Fn: genaro_removefile,
}


// var _ = app.Register(rmfileCommand)

// var _ = app.Register(&cli.Command{
// 	Name: "rm",
// 	Desc: "Genaro File remove",
// 	Argv: func() interface{} { return new(genaro_removeT) },	
// 	Fn: genaro_remove,
// })

type genaro_removeT struct {
	cli.Helper
	File string `cli:"file" usage:"remove file from Genaro network"`
}

func genaro_removefile(ctx *cli.Context) error {
	var env *C.storj_env_t
	init_env(ctx,env)

	bucket_id := "368be0816766b28fd5f43af5"
	file_id	:= "998960317b6725a3f8080c2b"

	var bucket_id_c *C.char
	var file_id_c *C.char
	string_to_char(bucket_id_c, bucket_id)
	string_to_char(file_id_c, file_id)

	var handle int

	status := C.storj_bridge_delete_file(env, bucket_id_c, file_id_c, unsafe.Pointer(&handle),(C.uv_after_work_cb)(unsafe.Pointer(C.check_delete_file)))
	
	panic(status)
	// assert(status == 0);

	if(C.int(C.uv_run(env.loop,C.UV_RUN_DEFAULT)) != 0){
		return nil;
	}

	C.storj_destroy_env(env);

	return nil;

}

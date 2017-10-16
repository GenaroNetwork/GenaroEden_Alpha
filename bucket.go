package main

/*
#include <microhttpd.h>
#include <assert.h>
#include <stdio.h>

#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. /usr/local/lib/libstorj.dylib /usr/local/lib/libuv.dylib /usr/local/lib/libjson-c.dylib

#include "../libstorj/src/storj.h"

void check_get_buckets(uv_work_t *work_req, int status)
{
    assert(status == 0);
    get_buckets_request_t *req = work_req->data;
    assert(req->handle == NULL);
    assert(json_object_is_type(req->response, json_type_array) == 1);

    struct json_object *bucket = json_object_array_get_idx(req->response, 0);
    struct json_object* value;
    int success = json_object_object_get_ex(bucket, "id", &value);

    storj_free_get_buckets_request(req);
    free(work_req);
}

void check_list_files(uv_work_t *work_req, int status)
{
    assert(status == 0);
    list_files_request_t *req = work_req->data;
    assert(req->handle == NULL);
    assert(req->response != NULL);

    struct json_object *file = json_object_array_get_idx(req->response, 0);
    struct json_object *value;
    int success = json_object_object_get_ex(file, "id", &value);
    assert(success == 1);
    assert(json_object_is_type(value, json_type_string) == 1);

    const char* id = json_object_get_string(value);
    assert(strcmp(id, "f18b5ca437b1ca3daa14969f") == 0);

    storj_free_list_files_request(req);
    free(work_req);
}

*/
import "C"
import "unsafe"

import (
	"github.com/mkideal/cli"
	// "github.com/stretchr/testify/assert"
)

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

// var _ = app.Register(&cli.Command{
// 	Name: "addbucket",
// 	Desc: "Genaro Bucket Adding",
// 	Argv: func() interface{} { return new(addbucketT) },
// 	Fn: addbucket,
// })

type addbucketT struct {
	cli.Helper
	Bucket string `cli:"addbkt,adb" usage:"add bucket on Genaro network"`
}

type initT struct {
	cli.Helper
	Host string `cli:"host" usage:"set host on Genaro network"`
	User string `cli:"user,u" usage:"username on Genaro network"`
	Password string `cli:"password,p" usage:"user password on Genaro network"`
}

func addbucket(ctx *cli.Context) error {
	argv := ctx.Argv().(*addbucketT)
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}

var removeCommand = &cli.Command{
	Name: "removebucket",
	Desc: "Genaro Bucket remove",
	Argv: func() interface{} { return new(removebucketT) },	
	Fn: removebucket,
}

// var _ = app.Register(removeCommand)

// var _ = app.Register(&cli.Command{
// 	Name: "removebucket",
// 	Desc: "Genaro Bucket remove",
// 	Argv: func() interface{} { return new(removebucketT) },	
// 	Fn: removebucket,
// })

type removebucketT struct {
	cli.Helper
	Bucket string `cli:"removebkt,rmb" usage:"remove bucket on Genaro network"`
}

func removebucket(ctx *cli.Context) error {
	argv := ctx.Argv().(*removebucketT)
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}


//=================================================================================================================

var listbucketCommand =&cli.Command{
	Name: "listbuckets",
	Desc: "Genaro Buskets Listing",
	Argv: func() interface{} { return new(listbucketsT) },	
	Fn: listbuckets,
}

// var _ = app.Register(&cli.Command{
// 	Name: "listbuckets",
// 	Desc: "Genaro Buskets Listing",
// 	Argv: func() interface{} { return new(listbucketsT) },	
// 	Fn: listbuckets,
// })

type listbucketsT struct {
	cli.Helper
	List string `cli:"list-buskets" usage:"user Genaro buskets list"`
}

func listbuckets(ctx *cli.Context) error{
	var env *C.storj_env_t
	init_env(ctx,env)

	var handle int

	status := C.storj_bridge_get_buckets(env, unsafe.Pointer(&handle), (C.uv_after_work_cb)(unsafe.Pointer(C.check_get_buckets)))

	panic(status)

	// fmt.Println("%s", env.loop)

	if(C.int(C.uv_run(env.loop,C.UV_RUN_DEFAULT)) != 0){
		return nil;
	}

	C.storj_destroy_env(env);

	return nil;	

}

//=================================================================================================================


var listfileCommand = &cli.Command{
	Name: "listfiles",
	Desc: "Genaro Busket File Listing",
	Argv: func() interface{} { return new(listfilesT) },	
	Fn: listfiles,
}

type listfilesT struct {
	cli.Helper
	List string `cli:"list-buskets" usage:"user Genaro files list"`
}

func listfiles(ctx *cli.Context) error {
	var env *C.storj_env_t
	init_env(ctx,env)

	bucket_id := "368be0816766b28fd5f43af5"

	var bucket_id_c *C.char

	string_to_char(bucket_id_c, bucket_id)

	status := C.storj_bridge_list_files(env, bucket_id_c, nil, (C.uv_after_work_cb)(unsafe.Pointer(C.check_list_files)))

	panic(status)

	if(C.int(C.uv_run(env.loop,C.UV_RUN_DEFAULT)) != 0){
		return nil;
	}

	C.storj_destroy_env(env);

	return nil;	
}
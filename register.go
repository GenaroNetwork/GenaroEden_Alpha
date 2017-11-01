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
import "fmt"
import "os"
import "strconv"
import "bytes"
import (
	"github.com/mkideal/cli"
	"github.com/howeyc/gopass"
)

type user_options struct{
	user *C.char;
	pass  *C.char;
	host  *C.char;
	mnemonic  *C.char;
	key  *C.char;
}

var _ = app.Register(&cli.Command{
	Name: "register",
	Desc: "Genaro account register ",
	Argv: func() interface{} { return new(registerT) },	
	Fn:   register,
})

type registerT struct {
	cli.Helper
}

func ScanLine() string {
	var c byte
	var err error
	var b []byte
	for ; err == nil; {
		_, err = fmt.Scanf("%c", &c)
		
		if c != '\n' {
			b = append(b, c)
		} else {
			break;
		}
	}
	return string(b)
}

func register(ctx *cli.Context) error {

	var env *C.storj_env_t
	
	if err := init_env_n(&env); err != nil{
		fmt.Printf("that is a bad thing ")
	}

	fmt.Printf("please input username(email): ")
	user := ScanLine()

	fmt.Printf("%v\n",user)

	fmt.Printf("Password: ")
	pass, _ := gopass.GetPasswd()

	fmt.Printf("%v\n",string(pass))

	user_opts := user_options{
		C.CString(user), 
		C.CString(string(pass)), 
		// C.CString("localhost"),
		C.CString("101.132.159.197"), 
		nil, 
		nil}

    if ( user_opts.user == C.CString("") || user_opts.pass == C.CString("")) {
		return errors.New("env is not set properly")	
    }	

	Sts := make(chan C.int)	

	go func(){

    	status := C.storj_bridge_register(env, C.CString(user), C.CString(string(pass)), (unsafe.Pointer(&user_opts)), (C.uv_after_work_cb)(unsafe.Pointer(C.register_callback)))	

		C.uv_run((env).loop, C.UV_RUN_DEFAULT)

		Sts<-status

	}()

	<-Sts

	return nil
}

func init_env_n(env_ptr **C.storj_env_t) error{

	var env *C.storj_env_t

	Env := make(chan *C.storj_env_t) 
	
	proto := C.CString("http")
	// host := C.CString("localhost")
	host := C.CString("101.132.159.197")
	user_agent := C.CString("libstorj-1.1.0-beta")

	options := storj_bridge_options{
		proto: proto, 
		host: host,
		port: C.int(8080),
		user: nil,
		pass: nil}

	http_options := storj_http_options{
	    user_agent : user_agent,
	    low_speed_limit : C.uint64_t(0),
	    low_speed_time : C.uint64_t(0),
	    timeout : C.uint64_t(0)}

	log_options := storj_log_options{
		level : C.int(0)}

	go func(){
		
		env = C.storj_init_env((*C.storj_bridge_options_t)(unsafe.Pointer(&options)),nil,(*C.storj_http_options_t)(unsafe.Pointer(&http_options)),(*C.storj_log_options_t)(unsafe.Pointer(&log_options)))
	
		C.uv_run(env.loop, C.UV_RUN_DEFAULT)
		
		Env <- env
	}()

	*env_ptr = <- Env	

	return nil
}

func set_env(env_ptr **C.storj_env_t) error {
	
	// var env *C.storj_env_t

// follow the storj library steps to set environment

	var env *C.storj_env_t

	// host := C.CString("localhost")
	host := C.CString("101.132.159.197")
    var user *C.char
    var	pass *C.char
    var mnemonic *C.char

	var char C.char

	var user_file *C.char
	var root_dir *C.char


    defer C.free(unsafe.Pointer(root_dir))

    if C.get_user_auth_location(host, &root_dir, &user_file) != C.int(0) {
		return errors.New("Unable to determine user auth filepath.\n")
    }


    if os.Getenv("STORJ_BRIDGE_USER") != ""{
    	user = C.strdup(C.CString(os.Getenv("STORJ_BRIDGE_USER")))
    }else{
    	user = C.CString("")   	
    }


    if os.Getenv("STORJ_BRIDGE_PASS") != "" {
    	pass = C.strdup(C.CString(os.Getenv("STORJ_BRIDGE_PASS")))   	
    }else{
    	pass = C.CString("")    	
    }


    if os.Getenv("STORJ_ENCRYPTION_KEY") != "" {
    	mnemonic = C.strdup(C.CString(os.Getenv("STORJ_ENCRYPTION_KEY"))) 	
    }else{
    	mnemonic = C.CString("")     	
    }

    keypass := C.CString(os.Getenv("STORJ_KEYPASS"))


    // Second, try to get from encrypted user file
    if   _, err := os.Stat(C.GoString(user_file)) ; err == nil &&
    	(len(C.GoString(user)) == 0 || len(C.GoString(pass)) == 0 || len(C.GoString(mnemonic)) == 0){
    		var key *C.char
    		if len(C.GoString(keypass)) != 0 {

 				key = (*C.char)(C.calloc(C.strlen(keypass)+1, C.size_t(char)))    			
    			if key == nil{
    				return errors.New("Unable to generate keys.\n") 
    			}
    			C.strcpy(key,keypass)
    		}else{  			
    			key = (*C.char)(C.calloc(C.BUFSIZ, C.size_t(char)))
    			if key == nil{
    				return errors.New("Unable to generate keys.\n") 
    			}
    			fmt.Printf("Unlock passphrase: ")
    			mask , _:= strconv.Atoi("*")
    			C.get_password(key,C.int(mask))
    			fmt.Printf("\n")        			
    		}
  		
    		var file_user *C.char
    		var file_pass *C.char
    		var file_mnemonic *C.char
					

    		if res := C.storj_decrypt_read_auth(user_file, key, &file_user,
                                    &file_pass, &file_mnemonic); res != C.int(0){  
    			defer C.free(unsafe.Pointer(key))
    			defer C.free(unsafe.Pointer(user_file))
    			defer C.free(unsafe.Pointer(file_user))
    			defer C.free(unsafe.Pointer(file_pass))
    			defer C.free(unsafe.Pointer(file_mnemonic))

    			return errors.New("Unable to read user file. Invalid keypass or path.\n")
    		}

    		var userbuffer bytes.Buffer
    		var passbuffer bytes.Buffer
    		var mnemonicbuffer bytes.Buffer

    		userbuffer.WriteString(C.GoString(file_user))
    		passbuffer.WriteString(C.GoString(file_pass))
    		mnemonicbuffer.WriteString(C.GoString(file_mnemonic))

			user = C.CString(userbuffer.String())
			pass = C.CString(passbuffer.String())
			mnemonic = C.CString(mnemonicbuffer.String())			

            defer C.free(unsafe.Pointer(file_user))
			defer C.free(unsafe.Pointer(file_pass))
            defer C.free(unsafe.Pointer(file_mnemonic))								
					// fmt.Printf("working_3_1\n")
	    //         if len(C.GoString(pass)) == 0 && len(C.GoString(file_pass)) != 0{
	    //             // user = file_user   
	    //             Pass := C.GoString(pass)
	    //             File_pass := C.GoString(file_pass)
	    //             copy(Pass,File_pass)  
					// fmt.Printf("%v\n ", Pass)     	            				
	    //         }  


	     //        if len(C.GoString(pass)) == 0 && len(C.GoString(file_pass)) != 0{
	     //            // pass = file_pass     
	     //            C.strcpy(pass,file_pass)
	     //            fmt.Printf("%v\n", C.GoString(pass))
      //       		fmt.Printf("working_2_1\n")	    
			 		// defer C.free(unsafe.Pointer(file_pass))             		            
	     //        }

					// fmt.Printf("working_3_2\n")

	    //         if len(C.GoString(mnemonic)) == 0 && len(C.GoString(file_mnemonic)) != 0 {
	    //             // mnemonic = file_mnemonic               
	    //             C.strcpy(mnemonic,file_mnemonic) 
					// defer C.free(unsafe.Pointer(file_mnemonic))
	    //         }

        // Third, ak for authentication
        if len(C.GoString(user)) == 0{
        	
        	var user_input *C.char

        	user_input = (*C.char)(C.malloc(C.BUFSIZ))

            if (user_input == nil) {
                return errors.New("user_input not set\n")
            }
            fmt.Printf("Bridge username (email): ")
            C.get_input(user_input)

            num_chars := C.strlen(user_input)
            user = (*C.char)(C.calloc(num_chars + 1, C.size_t(char)))
            
            if len(C.GoString(user)) == 0 {
				return errors.New("user_input not set\n")
            }
            
            C.memcpy(unsafe.Pointer(user), unsafe.Pointer(user_input), num_chars)
            C.free(unsafe.Pointer(user_input))
        }

        if len(C.GoString(pass)) == 0 {
            fmt.Printf("Bridge password: ")
            pass = (*C.char)(C.calloc(C.BUFSIZ, C.size_t(unsafe.Sizeof(char))))
            if len(C.GoString(pass)) == 0 {
                return errors.New("user password not set\n")
            }
            C.get_password(pass, '*');
            fmt.Printf("\n");
        }

        if len(C.GoString(mnemonic)) == 0 {
            fmt.Printf("Encryption key: ")
            var mnemonic_input *C.char

            mnemonic_input = (*C.char)(C.malloc(C.BUFSIZ))
            if len(C.GoString(mnemonic_input)) == 0 {
				return errors.New("user mnemonic not set\n")	
            }

            C.get_input(mnemonic_input)


            num_chars := C.strlen(mnemonic_input)

            mnemonic = (*C.char)(C.calloc(num_chars + 1, C.size_t(unsafe.Sizeof(char))))
            

            C.memcpy(unsafe.Pointer(mnemonic), unsafe.Pointer(mnemonic_input), num_chars)
            defer C.free(unsafe.Pointer(mnemonic_input))

            fmt.Printf("\n")
        }

    }

	Env := make(chan *C.storj_env_t) 
	
	proto := C.CString("http")
	// host = C.CString("101.132.159.197")
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

	defer C.free(unsafe.Pointer(user))    
	defer C.free(unsafe.Pointer(pass))       
	defer C.free(unsafe.Pointer(mnemonic))


    return nil
}
	

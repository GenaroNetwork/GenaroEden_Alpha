#include <microhttpd.h>
#include <assert.h>
#include <stdio.h>
// #include <curl/curl.h>
// #include <json-c/json.h>

#include "_cgo_export.h"
#include "../libstorj/src/storj.h"
#include "callbacks.h"


static inline void noop() {};

static char *get_home_dir()
{
#ifdef _WIN32
    return getenv("USERPROFILE");
#else
    return getenv("HOME");
#endif
}

void get_buckets(uv_work_t *work_req, int status)
{

    assert(status == 0);
    get_buckets_request_t *req = work_req->data;
    assert(req->handle == NULL);

   	uint32_t total_buckets = req->total_buckets;
    char *created;
    char *name;
    char *id;

    //show for display, should have one verision for json files

   	for (int i=0; i<total_buckets; i++){
   		created = calloc(strlen(req->buckets[i].created)+1,sizeof( char));
   		strcpy(created, req->buckets[i].created);
   		printf("%s\n", created);

   		name = calloc(strlen(req->buckets[i].name)+1,sizeof( char));
   		strcpy(name, req->buckets[i].name);
   		printf("%s\n", name);

   		id = calloc(strlen(req->buckets[i].id)+1,sizeof( char));
   		strcpy(id, req->buckets[i].id);
   		printf("%s\n", id);   		
   	}

	// bucketcallBack(total_buckets,created,name,id);

    storj_free_get_buckets_request(req);
    free(work_req);
}

void list_files_callback(uv_work_t *work_req, int status)
{
    int ret_status = 0;
    assert(status == 0);
    list_files_request_t *req = work_req->data;

    if (req->status_code == 404) {
        printf("Bucket id [%s] does not exist\n", req->bucket_id);
        goto cleanup;
    } else if (req->status_code == 400) {
        printf("Bucket id [%s] is invalid\n", req->bucket_id);
        goto cleanup;
    } else if (req->status_code == 401) {
        printf("Invalid user credentials.\n");
        goto cleanup;
    } else if (req->status_code != 200) {
        printf("Request failed with status code: %i\n", req->status_code);
    }

    if (req->total_files == 0) {
        printf("No files for bucket.\n");
    }

    for (int i = 0; i < req->total_files; i++) {

        storj_file_meta_t *file = &req->files[i];

        printf("ID: %s \tSize: %" PRIu64 " bytes \tDecrypted: %s \tType: %s \tCreated: %s \tName: %s\n",
               file->id,
               file->size,
               file->decrypted ? "true" : "false",
               file->mimetype,
               file->created,
               file->filename);
    }

cleanup:
    json_object_put(req->response);
    storj_free_list_files_request(req);
    free(work_req);
    exit(ret_status);
}

void get_buckets_callback(uv_work_t *work_req, int status)
{
    assert(status == 0);
    get_buckets_request_t *req = work_req->data;

    if (req->status_code == 401) {
       printf("Invalid user credentials.\n");
    } else if (req->status_code != 200 && req->status_code != 304) {
        printf("Request failed with status code: %i\n", req->status_code);
    } else if (req->total_buckets == 0) {
        printf("No buckets.\n");
    }

    for (int i = 0; i < req->total_buckets; i++) {
        storj_bucket_meta_t *bucket = &req->buckets[i];
        printf("ID: %s \tDecrypted: %s \tCreated: %s \tName: %s\n",
               bucket->id, bucket->decrypted ? "true" : "false",
               bucket->created, bucket->name);
    }

    json_object_put(req->response);
    storj_free_get_buckets_request(req);
    free(work_req);
}

void create_bucket(uv_work_t *work_req, int status)
{

    assert(status == 0);
    create_bucket_request_t *req = work_req->data;
    assert(req->handle == NULL);


    struct json_object* value;
    int success = json_object_object_get_ex(req->response, "name", &value);
    assert(success == 1);
    assert(json_object_is_type(value, json_type_string) == 1);

    const char* name = json_object_get_string(value);

    json_object_put(req->response);
    free((char *)req->encrypted_bucket_name);
    free(req->bucket);
    free(req);
    free(work_req);

}

void delete_bucket(uv_work_t *work_req, int status)
{
    assert(status == 0);
    json_request_t *req = work_req->data;
    assert(req->handle == NULL);
    assert(req->response == NULL);
    assert(req->status_code == 204);

    json_object_put(req->response);
    free(req->path);
    free(req);
    free(work_req);
}

//=============================================================================
//	get info  
//=============================================================================
void get_info_callback(uv_work_t *work_req, int status)
{
    assert(status == 0);
    json_request_t *req = work_req->data;

    if (req->error_code || req->response == NULL) {
        free(req);
        free(work_req);
        if (req->error_code) {
            printf("Request failed, reason: %i\n",
                   (req->error_code));
        } else {
            printf("Failed to get info.\n");
        }
        exit(1);
    }


    struct json_object *info;
    json_object_object_get_ex(req->response, "info", &info);

    struct json_object *title;
    json_object_object_get_ex(info, "title", &title);
    struct json_object *description;
    json_object_object_get_ex(info, "description", &description);
    struct json_object *version;
    json_object_object_get_ex(info, "version", &version);
    struct json_object *host;
    json_object_object_get_ex(req->response, "host", &host);


    printf("Title:  %s\n", json_object_get_string(title));
    printf("Description: %s\n", json_object_get_string(description));
    printf("Version:     %s\n", json_object_get_string(version));
    printf("Host:        %s\n", json_object_get_string(host));

    json_object_put(req->response);
    free(req);
    free(work_req);
}


//=============================================================================
//	upload file  
//=============================================================================
void check_store_file_progress(double progress,
                               uint64_t uploaded_bytes,
                               uint64_t total_bytes,
                               void *handle)
{
    assert(handle == NULL);
    if (progress == (double)1) {
    	printf("progress finished\n");
    }
}

void check_store_file(int error_code, char *file_id, void *handle)
{
    assert(handle == NULL);
    if (error_code == 0) {
   		printf("%s\n", file_id);
    } else {
        printf("\t\tERROR:   %s\n", storj_strerror(error_code));
    }

    free(file_id);
}

int upload_file(storj_env_t *env, char *bucket_id, const char *file_path)
{
    FILE *fd = fopen(file_path, "r");

    if (!fd) {
        printf("Invalid file path: %s\n", file_path);
    }

    const char *file_name = get_filename_separator(file_path);

    printf("file name : %s\n", file_name);

    if (!file_name) {
        file_name = file_path;
    }

    // Upload opts env variables:
    char *prepare_frame_limit = getenv("STORJ_PREPARE_FRAME_LIMIT");
    char *push_frame_limit = getenv("STORJ_PUSH_FRAME_LIMIT");
    char *push_shard_limit = getenv("STORJ_PUSH_SHARD_LIMIT");
    char *rs = getenv("STORJ_REED_SOLOMON");

    storj_upload_opts_t upload_opts = {
        .prepare_frame_limit = (prepare_frame_limit) ? atoi(prepare_frame_limit) : 1,
        .push_frame_limit = (push_frame_limit) ? atoi(push_frame_limit) : 64,
        .push_shard_limit = (push_shard_limit) ? atoi(push_shard_limit) : 64,
        .rs = (!rs) ? true : (strcmp(rs, "false") == 0) ? false : true,
        .bucket_id = bucket_id,
        .file_name = file_name,
        .fd = fd
    };

    uv_signal_t *sig = malloc(sizeof(uv_signal_t));
    if (!sig) {
        return 1;
    }
    uv_signal_init(env->loop, sig);
    uv_signal_start(sig, upload_signal_handler, SIGINT);

    storj_upload_state_t *state = malloc(sizeof(storj_upload_state_t));
    if (!state) {
        return 1;
    }

    sig->data = state;

    storj_progress_cb progress_cb = (storj_progress_cb)noop;
    if (env->log_options->level == 0) {
        progress_cb = file_progress;
    }

    int status = storj_bridge_store_file(env,
                                         state,
                                         &upload_opts,
                                         NULL,
                                         progress_cb,
                                         upload_file_complete);

    return status;
}

static void file_progress(double progress,
                          uint64_t downloaded_bytes,
                          uint64_t total_bytes,
                          void *handle)
{
    int bar_width = 70;

    if (progress == 0 && downloaded_bytes == 0) {
        printf("Preparing File...");
        fflush(stdout);
        return;
    }

    printf("\r[");
    int pos = bar_width * progress;
    for (int i = 0; i < bar_width; ++i) {
        if (i < pos) {
            printf("=");
        } else if (i == pos) {
            printf(">");
        } else {
            printf(" ");
        }
    }
    printf("] %.*f%%", 2, progress * 100);

    fflush(stdout);
}

static const char *get_filename_separator(const char *file_path)
{
    const char *file_name = NULL;
#ifdef _WIN32
    file_name = strrchr(file_path, '\\');
    if (!file_name) {
        file_name = strrchr(file_path, '/');
    }
    if (!file_name && file_path) {
        file_name = file_path;
    }
    if (!file_name) {
        return NULL;
    }
    if (file_name[0] == '\\' || file_name[0] == '/') {
        file_name++;
    }
#else
    file_name = strrchr(file_path, '/');
    if (!file_name && file_path) {
        file_name = file_path;
    }
    if (!file_name) {
        return NULL;
    }
    if (file_name[0] == '/') {
        file_name++;
    }
#endif
    return file_name;
}

static void upload_file_complete(int status, char *file_id, void *handle)
{
    printf("\n");
    if (status != 0) {
        printf("Upload failure: %s\n", storj_strerror(status));
        exit(status);
    }

    printf("Upload Success! File ID: %s\n", file_id);

    free(file_id);

    exit(0);
}

void upload_signal_handler(uv_signal_t *req, int signum)
{
    storj_upload_state_t *state = req->data;
    storj_bridge_store_file_cancel(state);
    if (uv_signal_stop(req)) {
        printf("Unable to stop signal\n");
    }
    uv_close((uv_handle_t *)req, close_signal);
}

void close_signal(uv_handle_t *handle)
{
    ((void)0);
}


//=============================================================================
//	download file 
//=============================================================================

int download_file(storj_env_t *env, char *bucket_id,
                         char *file_id, char *path)
{
    FILE *fd = NULL;

    if (path) {
        char user_input[BUFSIZ];
        memset(user_input, '\0', BUFSIZ);

        if(access(path, F_OK) != -1 ) {
            printf("Warning: File already exists at path [%s].\n", path);
            while (strcmp(user_input, "y") != 0 && strcmp(user_input, "n") != 0)
            {
                memset(user_input, '\0', BUFSIZ);
                printf("Would you like to overwrite [%s]: [y/n] ", path);
                get_input(user_input);
            }

            if (strcmp(user_input, "n") == 0) {
                printf("\nCanceled overwriting of [%s].\n", path);
                return 1;
            }

            unlink(path);
        }

        fd = fopen(path, "w+");
    } else {
        fd = stdout;
    }

    if (fd == NULL) {
        // TODO send to stderr
        printf("Unable to open %s: %s\n", path, strerror(errno));
        return 1;
    }

    uv_signal_t *sig = malloc(sizeof(uv_signal_t));
    uv_signal_init(env->loop, sig);
    uv_signal_start(sig, download_signal_handler, SIGINT);

    storj_download_state_t *state = malloc(sizeof(storj_download_state_t));
    if (!state) {
        return 1;
    }

    sig->data = state;

    storj_progress_cb progress_cb = (storj_progress_cb)noop;
    if (path && env->log_options->level == 0) {
        progress_cb = file_progress;
    }

    int status = storj_bridge_resolve_file(env, state, bucket_id,
                                           file_id, fd, NULL,
                                           progress_cb,
                                           download_file_complete);

    // printf("status is %d \n", status);

    return status;
}

static void download_file_complete(int status, FILE *fd, void *handle)
{
    printf("\n");
    fclose(fd);
    if (status) {
        // TODO send to stderr
        switch(status) {
            case STORJ_FILE_DECRYPTION_ERROR:
                printf("Unable to properly decrypt file, please check " \
                       "that the correct encryption key was " \
                       "imported correctly.\n\n");
                break;
            default:
                printf("Download failure: %s\n", storj_strerror(status));
        }

        exit(status);
    }
    printf("Download Success!\n");
    exit(0);
}

void download_signal_handler(uv_signal_t *req, int signum)
{
    storj_download_state_t *state = req->data;
    storj_bridge_resolve_file_cancel(state);
    if (uv_signal_stop(req)) {
        printf("Unable to stop signal\n");
    }
    uv_close((uv_handle_t *)req, close_signal);
}

void get_input(char *line)
{
    if (fgets(line, BUFSIZ, stdin) == NULL) {
        line[0] = '\0';
    } else {
        int len = strlen(line);
        if (len > 0) {
            char *last = strrchr(line, '\n');
            if (last) {
                last[0] = '\0';
            }
            last = strrchr(line, '\r');
            if (last) {
                last[0] = '\0';
            }
        }
    }
}
//=============================================================================
//	delete file 
//=============================================================================
void delete_file_callback(uv_work_t *work_req, int status)
{
    assert(status == 0);
    json_request_t *req = work_req->data;

    if (req->status_code == 200 || req->status_code == 204) {
        printf("File was successfully removed from bucket.\n");
    } else if (req->status_code == 401) {
        printf("Invalid user credentials.\n");
    } else {
        printf("Failed to remove file from bucket. (%i)\n", req->status_code);
    }

    json_object_put(req->response);
    free(req->path);
    free(req);
    free(work_req);
}

//================================================================================
// register 
//================================================================================

int import_keys(user_options_t *options)
{
    int status = 0;
    // char *host = options->host ? strdup(options->host) : NULL;
    char *host = malloc(64 * sizeof(char));
    strcpy(host, "genaro_eden_alpha");
    // char *host = "genaro_eden_alpha";
    char *user = options->user ? strdup(options->user) : NULL;
    char *pass = options->pass ? strdup(options->pass) : NULL;
    char *key = options->key ? strdup(options->key) : NULL;
    char *mnemonic = options->mnemonic ? strdup(options->mnemonic): NULL;
    char *mnemonic_input = NULL;
    char *user_file = NULL;
    char *root_dir = NULL;
    int num_chars;

    char *user_input = calloc(BUFSIZ, sizeof(char));
    if (user_input == NULL) {
        printf("Unable to allocate buffer\n");
        status = 1;
        goto clear_variables;
    }

    if (get_user_auth_location(host, &root_dir, &user_file)) {
        printf("Unable to determine user auth filepath.\n");
        status = 1;
        goto clear_variables;
    }

    struct stat st;
    if (stat(user_file, &st) == 0) {
        printf("Would you like to overwrite the current settings?: [y/n] ");
        get_input(user_input);
        while (strcmp(user_input, "y") != 0 && strcmp(user_input, "n") != 0)
        {
            printf("Would you like to overwrite the current settings?: [y/n] ");
            get_input(user_input);
        }

        if (strcmp(user_input, "n") == 0) {
            printf("\nCanceled overwriting of stored credentials.\n");
            status = 1;
            goto clear_variables;
        }
    }

    if (!user) {
        printf("Bridge username (email): ");
        get_input(user_input);
        num_chars = strlen(user_input);
        user = calloc(num_chars + 1, sizeof(char));
        if (!user) {
            status = 1;
            goto clear_variables;
        }
        memcpy(user, user_input, num_chars * sizeof(char));
    }

    if (!pass) {
        printf("Bridge password: ");
        pass = calloc(BUFSIZ, sizeof(char));
        if (!pass) {
            status = 1;
            goto clear_variables;
        }
        get_password(pass, '*');
        printf("\n");
    }

    if (!mnemonic) {
        mnemonic_input = calloc(BUFSIZ, sizeof(char));
        if (!mnemonic_input) {
            status = 1;
            goto clear_variables;
        }

        printf("\nIf you've previously uploaded files, please enter your" \
               " existing encryption key (12 to 24 words). \nOtherwise leave" \
               " the field blank to generate a new key.\n\n");

        printf("Encryption key: ");
        get_input(mnemonic_input);
        num_chars = strlen(mnemonic_input);

        if (num_chars == 0) {
            printf("\n");
            generate_mnemonic(&mnemonic);
            printf("\n");

            printf("Encryption key: %s\n", mnemonic);
            printf("\n");
            printf("Please make sure to backup this key in a safe location. " \
                   "If the key is lost, the data uploaded will also be lost.\n\n");
        } else {
            mnemonic = calloc(num_chars + 1, sizeof(char));
            if (!mnemonic) {
                status = 1;
                goto clear_variables;
            }
            memcpy(mnemonic, mnemonic_input, num_chars * sizeof(char));
        }

        if (!storj_mnemonic_check(mnemonic)) {
            printf("Encryption key integrity check failed.\n");
            status = 1;
            goto clear_variables;
        }
    }

    if (!key) {
        key = calloc(BUFSIZ, sizeof(char));
        printf("We now need to save these settings. Please enter a passphrase" \
               " to lock your settings.\n\n");
        if (get_password_verify("Unlock passphrase: ", key, 0)) {
            printf("Unable to store encrypted authentication.\n");
            status = 1;
            goto clear_variables;
        }
        printf("\n");
    }

    if (make_user_directory(root_dir)) {
        status = 1;
        goto clear_variables;
    }

    if (storj_encrypt_write_auth(user_file, key, user, pass, mnemonic)) {
        status = 1;
        printf("Failed to write to disk\n");
        goto clear_variables;
    }

    printf("Successfully stored bridge username, password, and encryption "\
           "key to %s\n\n",
           user_file);

clear_variables:
    if (user) {
        free(user);
    }
    if (user_input) {
        free(user_input);
    }
    if (pass) {
        free(pass);
    }
    if (mnemonic) {
        free(mnemonic);
    }
    if (mnemonic_input) {
        free(mnemonic_input);
    }
    if (key) {
        free(key);
    }
    if (root_dir) {
        free(root_dir);
    }
    if (user_file) {
        free(user_file);
    }
    if (host) {
        free(host);
    }

    return status;
}

void register_callback(uv_work_t *work_req, int status)
{
    assert(status == 0);
    json_request_t *req = work_req->data;

    if (req->status_code != 201) {
        printf("Request failed with status code: %i\n",
               req->status_code);
        struct json_object *error;
        json_object_object_get_ex(req->response, "error", &error);
        printf("Error: %s\n", json_object_get_string(error));

        user_options_t *handle = (user_options_t *) req->handle;
        free(handle->user);
        free(handle->host);
        free(handle->pass);
    } else {
        struct json_object *email;
        json_object_object_get_ex(req->response, "email", &email);
        printf("\n");
        printf("Successfully registered %s, please check your email "\
               "to confirm.\n", json_object_get_string(email));

        // save credentials
        char *mnemonic = NULL;
        printf("\n");
        generate_mnemonic(&mnemonic);
        printf("\n");

        printf("Encryption key: %s\n", mnemonic);
        printf("\n");
        printf("Please make sure to backup this key in a safe location. " \
               "If the key is lost, the data uploaded will also be lost.\n\n");

        user_options_t *user_opts = req->handle;

        user_opts->mnemonic = mnemonic;
        import_keys(user_opts);

        if (mnemonic) {
            free(mnemonic);
        }
        if (user_opts->pass) {
            free(user_opts->pass);
        }
        if (user_opts->user) {
            free(user_opts->user);
        }
        if (user_opts->host) {
            free(user_opts->host);
        }
    }

    json_object_put(req->response);
    json_object_put(req->body);
    free(req);
    free(work_req);
}

static int generate_mnemonic(char **mnemonic)
{
    char *strength_str = NULL;
    int strength = 0;
    int status = 0;

    printf("We now need to create an secret key used for encrypting " \
           "files.\nPlease choose strength from: 128, 160, 192, 224, 256\n\n");

    while (strength % 32 || strength < 128 || strength > 256) {
        strength_str = calloc(BUFSIZ, sizeof(char));
        printf("Strength: ");
        get_input(strength_str);

        if (strength_str != NULL) {
            strength = atoi(strength_str);
        }

        free(strength_str);
    }

    if (*mnemonic) {
        free(*mnemonic);
    }

    *mnemonic = NULL;

    int generate_code = storj_mnemonic_generate(strength, mnemonic);
    if (*mnemonic == NULL || generate_code == 0) {
        printf("Failed to generate encryption key.\n");
        status = 1;
        status = generate_mnemonic(mnemonic);
    }

    return status;
}

int get_user_auth_location(char *host, char **root_dir, char **user_file)
{
    char *home_dir = get_home_dir();
    if (home_dir == NULL) {
        return 1;
    }

    // int len = strlen(home_dir) + strlen("/.genaro/");
    int len = strlen(home_dir) + strlen("/.genaro/");    
    *root_dir = calloc(len + 1, sizeof(char));
    if (!*root_dir) {
        return 1;
    }

    strcpy(*root_dir, home_dir);

    strcat(*root_dir, "/.genaro/");
    // strcat(*root_dir, "/.genaro/");    

    len = strlen(*root_dir) + strlen(host) + strlen(".json");
    *user_file = calloc(len + 1, sizeof(char));
    if (!*user_file) {
        return 1;
    }

    strcpy(*user_file, *root_dir);
    strcat(*user_file, host);
    strcat(*user_file, ".json");

    return 0;
}

static int make_user_directory(char *path)
{
    struct stat st = {0};
    if (stat(path, &st) == -1) {
#if _WIN32
        int mkdir_status = _mkdir(path);
        if (mkdir_status) {
            printf("Unable to create directory %s: code: %i.\n",
                   path,
                   mkdir_status);
            return 1;
        }
#else
        if (mkdir(path, 0700)) {
            printf("Unable to create directory %s: reason: %s\n",
                   path,
                   strerror(errno));
            return 1;
        }
#endif
    }
    return 0;
}
// export _keys 
//
//
int export_keys(char *host)
{
    int status = 0;
    char *user_file = NULL;
    char *root_dir = NULL;
    char *user = NULL;
    char *pass = NULL;
    char *mnemonic = NULL;
    char *key = NULL;

    if (get_user_auth_location(host, &root_dir, &user_file)) {
        printf("Unable to determine user auth filepath.\n");
        status = 1;
        goto clear_variables;
    }

    if (access(user_file, F_OK) != -1) {
        key = calloc(BUFSIZ, sizeof(char));
        printf("Unlock passphrase: ");
        get_password(key, '*');
        printf("\n\n");

        if (storj_decrypt_read_auth(user_file, key, &user, &pass, &mnemonic)) {
            printf("Unable to read user file.\n");
            status = 1;
            goto clear_variables;
        }

        printf("Username:\t%s\nPassword:\t%s\nEncryption key:\t%s\n",
               user, pass, mnemonic);
    }

clear_variables:
    if (user) {
        free(user);
    }
    if (pass) {
        free(pass);
    }
    if (mnemonic) {
        free(mnemonic);
    }
    if (root_dir) {
        free(root_dir);
    }
    if (user_file) {
        free(user_file);
    }
    if (key) {
        free(key);
    }
    return status;
}

//
//
//
int get_password(char *password, int mask)
{
    int max_pass_len = 512;

#ifdef _WIN32
    HANDLE hstdin = GetStdHandle(STD_INPUT_HANDLE);
    DWORD mode = 0;
    DWORD prev_mode = 0;
    GetConsoleMode(hstdin, &mode);
    GetConsoleMode(hstdin, &prev_mode);
    SetConsoleMode(hstdin, mode & ~(ENABLE_LINE_INPUT | ENABLE_ECHO_INPUT));
#else
    static struct termios prev_terminal;
    static struct termios terminal;

    tcgetattr(STDIN_FILENO, &prev_terminal);

    memcpy (&terminal, &prev_terminal, sizeof(struct termios));
    terminal.c_lflag &= ~(ICANON | ECHO);
    terminal.c_cc[VTIME] = 0;
    terminal.c_cc[VMIN] = 1;
    tcsetattr(STDIN_FILENO, TCSANOW, &terminal);
#endif

    size_t idx = 0;         /* index, number of chars in read   */
    int c = 0;

    const char BACKSPACE = 8;
    const char RETURN = 13;

    /* read chars from fp, mask if valid char specified */
#ifdef _WIN32
    long unsigned int char_read = 0;
    while ((ReadConsole(hstdin, &c, 1, &char_read, NULL) && c != '\n' && c != RETURN && c != EOF && idx < max_pass_len - 1) ||
            (idx == max_pass_len - 1 && c == BACKSPACE))
#else
    while (((c = fgetc(stdin)) != '\n' && c != EOF && idx < max_pass_len - 1) ||
            (idx == max_pass_len - 1 && c == 127))
#endif
    {
        if (c != 127 && c != BACKSPACE) {
            if (31 < mask && mask < 127)    /* valid ascii char */
                fputc(mask, stdout);
            password[idx++] = c;
        } else if (idx > 0) {         /* handle backspace (del)   */
            if (31 < mask && mask < 127) {
                fputc(0x8, stdout);
                fputc(' ', stdout);
                fputc(0x8, stdout);
            }
            password[--idx] = 0;
        }
    }
    password[idx] = 0; /* null-terminate   */

    // go back to the previous settings
#ifdef _WIN32
    SetConsoleMode(hstdin, prev_mode);
#else
    tcsetattr(STDIN_FILENO, TCSANOW, &prev_terminal);
#endif

    return idx; /* number of chars in passwd    */
}

int get_password_verify(char *prompt, char *password, int count)
{
    printf("%s", prompt);
    char first_password[BUFSIZ];
    get_password(first_password, '*');

    printf("\nAgain to verify: ");
    char second_password[BUFSIZ];
    get_password(second_password, '*');

    int match = strcmp(first_password, second_password);
    strncpy(password, first_password, BUFSIZ);

    if (match == 0) {
        return 0;
    } else {
        printf("\nPassphrases did not match. ");
        count++;
        if (count > 3) {
            printf("\n");
            return 1;
        }
        printf("Try again...\n");
        return get_password_verify(prompt, password, count);
    }
}


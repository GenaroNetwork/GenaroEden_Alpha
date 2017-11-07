#include <uv.h>

#include "storj.h"


typedef struct {
    char *user;
    char *pass;
    char *host;
    char *mnemonic;
    char *key;
} user_options_t;

int upload_file(storj_env_t *env, char *bucket_id, const char *file_path);

static void file_progress(double progress,uint64_t downloaded_bytes,uint64_t total_bytes,void *handle);

static void upload_file_complete(int status, char *file_id, void *handle);

void upload_signal_handler(uv_signal_t *req, int signum);

void close_signal(uv_handle_t *handle);

static const char *get_filename_separator(const char *file_path);

int download_file(storj_env_t *env, char *bucket_id,
                         char *file_id, char *path);

static void download_file_complete(int status, FILE *fd, void *handle);

void download_signal_handler(uv_signal_t *req, int signum);

void get_input(char *line);

void get_info_callback(uv_work_t *work_req, int status);

void delete_file_callback(uv_work_t *work_req, int status);

//============================================================
void list_files_callback(uv_work_t *work_req, int status);
void get_buckets_callback(uv_work_t *work_req, int status);

int import_keys(user_options_t *options);
int export_keys(char *host);
void register_callback(uv_work_t *work_req, int status);
static int generate_mnemonic(char **mnemonic);
int get_user_auth_location(char *host, char **root_dir, char **user_file);
static int make_user_directory(char *path);

int get_password(char *password, int mask);
int get_password_verify(char *prompt, char *password, int count);


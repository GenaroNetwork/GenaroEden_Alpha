# GenaroEden Alpha
Genaro Eden Project, unix Alpha version 

> Genaro Network based on golang and C, accept commands in unix/macos system.

## Setup

``` bash
# install to local

go build -o genaro

# start genaro 

./genaro

# check all commands

./genaro -h

# register 

./genaro register 

# list all buckets 

./genaro bucket listbuckets

# list files in bucket with certain id

./genaro bucket listfiles -i (your bucket id)

# add bucket

./genaro bucket addbucket -n (your bucket name)

# upload file to a given bucket (Genaro set command)

./genaro file set -i (your bucket id) -p (your file path)

# download file from a given bucket (Genaro get command)

./genaro file get -b (your bucket id) -f (your file id in bucket) -p (your download file path)

```

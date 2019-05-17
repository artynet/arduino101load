#!/bin/bash -x

checkgopath () {

    GOPATH=$(printenv GOPATH)

    if [ -z $GOPATH ]
    then
        mkdir -p ${HOME}/go
        export GOPATH=${HOME}/go
    fi

}

# check if GOPATH variable is blank or not
checkgopath

# cleaning all go packages
rm -rf $GOPATH/{pkg,src}/*

# downloading dependencies
go get github.com/mattn/go-shellwords

# cd $GOPATH/src/github.com/arduino/arduino-builder
# git checkout ${VERSION}

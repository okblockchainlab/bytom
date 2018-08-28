#!/bin/sh

if [ -z "$JAVA_HOME" ]; then
  printf "No JAVA_HOME detected! "
  printf "Setup JAVA_HOME before build: export JAVA_HOME=/path/to/java\\n"
  exit 1
fi

LDFLAGS="-X github.com/bytom/version.GitCommit=`git rev-parse HEAD`"

EXT=so
NM_FLAGS=
TARGET_OS=`uname -s`
case "$TARGET_OS" in
  Darwin)
    EXT=dylib
    ;;
  Linux)
    EXT=so
    NM_FLAGS=-D
    ;;
  *)
  echo "Unknown platform!" >&2
  exit 1
esac


go build -o libbytom.${EXT} -buildmode=c-shared -ldflags="${ldflags}" ./okwallet/libbytom
[ $? -ne 0 ] && exit 1
nm ${NM_FLAGS} libbytom.${EXT} |grep "[ _]Java_com_okcoin"

#!/bin/bash

REL_BASE=build
APP=box2perm

for D in linux-arm7 linux-ia32 linux-x86_64 freebsd win32 darwin
do
  [ -d $REL_BASE/$D ] || mkdir -p $REL_BASE/$D/box2perm
  pushd $REL_BASE/$D/box2perm/ >/dev/null
  ln -sf ../../../README.md ./
  popd >/dev/null
done

GOOS=linux GOARCH=arm GOARM=7 go build -o $REL_BASE/linux-arm7/box2perm/$APP
GOOS=linux GOARCH=386 go build -o $REL_BASE/linux-ia32/box2perm/$APP
go build && cp $APP $REL_BASE/linux-x86_64/box2perm
GOOS=freebsd GOARCH=amd64 go build -o $REL_BASE/freebsd/box2perm/$APP
GOOS=windows GOARCH=386 go build -o $REL_BASE/win32/box2perm/$APP.exe
GOOS=darwin GOARCH=amd64 go build -o $REL_BASE/darwin/box2perm/$APP

mkdir -p /tmp/box2perm-rel
for D in linux-arm7 linux-ia32 linux-x86_64 freebsd win32 darwin
do
  pushd $REL_BASE/$D >/dev/null
  zip -9r /tmp/box2perm-rel/box2perm-$D.zip box2perm
  popd>/dev/null
done

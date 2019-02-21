#!/bin/bash
export GOPATH="$PWD/"
go build
if [ $? -eq 0 ]; then
  echo "Done building."
else
  return 3;
fi

#!/bin/bash
nc -v -z -w 1 -n 127.0.0.1 $1 >null 2>&1
if [ $? -ne 0 ]; then
    exit 1
else
#    nc -v -z -w 1 -n 127.0.0.1 $2 >null 2>&1
#    if [ $? -ne 0 ]; then
#      exit 1
#    else
      exit 0
#    fi
fi

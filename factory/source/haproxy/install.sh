#!/bin/bash

source config.cfg

#lua
cd source/haproxy && \
tar xvzf lua-${lua_version}.tar.gz && \
cd lua-${lua_version} && \
make all test && \
mkdir -p /usr/local/src/lua-${lua_version}/src && \
\cp -f -fa src/* /usr/local/src/lua-${lua_version}/src && \
cd .. && \
rm -rf lua-${lua_version}

#haproxy
tar xvzf haproxy-${haproxy_version}.tar.gz && \
cd haproxy-${haproxy_version} && \
make -j 2 ARCH=x86_64 TARGET=linux-glibc USE_PCRE=1 USE_OPENSSL=1 USE_ZLIB=1 USE_SYSTEMD=1 USE_LUA=1 LUA_INC=/usr/local/src/lua-${lua_version}/src LUA_LIB=/usr/local/src/lua-${lua_version}/src && \
make install PREFIX=/usr/local/haproxy && \
\cp -f  /usr/local/haproxy/sbin/haproxy /${bin_dir} && \
cd .. && \
rm -rf haproxy-${haproxy_version} && cd .. && cd ..
echo "install.sh(haproxy) ok"

#!/bin/bash

source config.cfg

#openssl
# k8s 1.24版本不允许证书采用 rsa-sha1，必须采用 rsa-sha512. 低版本低openssl不支持rsa-sha512.
# k8s版本 >1.23.x,openssl版本必须大于等于1.1.0
# k8s版本 <=1.23.x,openssl版本默认不变（一般系统自带为1.0.x）
#if [[ $k8s_common_ver > "v1.23.x" ]];then
#  openssl version |grep 1.1.0
# if [[ $? -eq 0 ]];then
#     exit 0
#  fi

#centos8 版本高于1.1.0，退出
if [[ $target_system == "centos8" ]]; then
   exit 0
fi
#其他版本linux 检测版本情况
curr_ver=`openssl version | awk '{print $2}'`
curr_ver=${curr_ver:0:${#curr_ver}-1}

ver=`echo "$openssl_version" |  sed 's/_/./g'`
ver=${ver:0:${#ver}-1}

echo "curr_ver=$curr_ver"
echo "ver=$ver"

seelp 3s
#当前版本高于或等于定义版本，退出
if [[ $curr_ver > $ver || $curr_ver == $ver ]];then
  exit 0
fi
echo "begin update core ... "
cd source/openssl && \
tar xvzf OpenSSL_${openssl_version}.tar.gz && \
cd openssl-OpenSSL_${openssl_version} && \
./config --prefix=/usr/local/openssl && \
make clean && make && make install  && \
rm -rf /usr/bin/openssl  && \
rm -rf /usr/include/openssl  && \
rm -rf /usr/lib64/libssl.so  && \
rm -rf /usr/lib64/libcrypto.so  && \
rm -rf /usr/lib/libssl.so  && \
rm -rf /usr/lib/libcrypto.so  && \
ln -s /usr/local/openssl/bin/openssl /usr/bin/openssl && \
ln -s /usr/local/openssl/include/openssl /usr/include/openssl && \
ln -s /usr/local/openssl/lib/libssl.so /usr/lib64/libssl.so && \
ln -s /usr/local/openssl/lib/libcrypto.so /usr/lib64/libcrypto.so && \
ln -s /usr/local/openssl/lib/libssl.so /usr/lib/libssl.so && \
ln -s /usr/local/openssl/lib/libcrypto.so /usr/lib/libcrypto.so
cat > /etc/ld.so.conf << EOF
include ld.so.conf.d/*.conf
/usr/local/openssl/lib
EOF
ldconfig -v && \
openssl version && \
cd .. && \
rm -rf openssl-OpenSSL_${openssl_version} && \
cd .. &&  cd ..
echo "install.sh(openssl) ok"

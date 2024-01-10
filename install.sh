#!/usr/bin/env bash
#
# for installing le

set -o errexit
set -o nounset

install_path="/usr/bin"

if [[ ":${PATH}:" == *":${install_path}:"* ]]; then
    cp ./bin/le ${install_path}/ 
    chown $(id -u):$(id -g) ${install_path}/le
else
    echo "error: ${install_path} not in ${PATH}"
    exit 1
fi

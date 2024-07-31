#!/bin/bash

Version=$1
#Version=${Version:-0.2.0}
PRODUCT_NAME="AI-Demo"
TARGET_ROOT=`cd "$(dirname "$0")"; pwd`
RELEASE_NAME=release_$Version
RELEASE_ROOT=${TARGET_ROOT}/${RELEASE_NAME}

echo RELEASE_NAME=$RELEASE_NAME

cd ${TARGET_ROOT}/

echo "build x86_64 ${PRODUCT_NAME}... "
make $PRODUCT_NAME
echo "build arm64 ${PRODUCT_NAME}... "
make $PRODUCT_NAME-ARM64
echo "build arm ${PRODUCT_NAME}... "
make $PRODUCT_NAME-ARM
echo "build windows exe ${PRODUCT_NAME}..."
make $PRODUCT_NAME.exe

if [ ! -d ${RELEASE_ROOT} ]; then
	rm -rf ${RELEASE_ROOT}		
fi

[ ! -e ${RELEASE_ROOT} ] && mkdir -p ${RELEASE_ROOT}

cp -a frontend ${RELEASE_ROOT}/
cp -a conf ${RELEASE_ROOT}/
cp -a ffmpeg ${RELEASE_ROOT}/
cp -a ffmpeg-ARM64 ${RELEASE_ROOT}/
cp  -a $PRODUCT_NAME ${PRODUCT_NAME}-ARM ${PRODUCT_NAME}-ARM64 ${PRODUCT_NAME}.exe ${RELEASE_ROOT}/
[ -a "startup.sh" ] && cp -a startup.sh ${RELEASE_ROOT}/
[ -a "AI-Demo.service" ] && cp -a ${PRODUCT_NAME}.service  ${RELEASE_ROOT}/AI-Demo.service

echo "Syncing...."
sync;sync;
sync
sync

#tar cJvf ${RELEASE_NAME}.tar.xz ${RELEASE_NAME}
tar zcvf ${RELEASE_NAME}.tar.gz ${RELEASE_NAME}

rm -rf ${RELEASE_ROOT}
echo "[Done]"


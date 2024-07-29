#!/bin/bash

if [ $# -ne 1 ];
then
    echo "Please select os type"
	echo "./release.sh windows or ./release.sh linux"
    exit
fi

ClientVersion=$(cat common/constants/default.go  | grep ClientVersion | sed 's/\"//g')
ClientVersion=${ClientVersion#*=}

V=$(echo $ClientVersion |awk '{print $1}' |sed 's/ //g' | sed 's/\t//g')
V=${V:-Beta}
Version=$(echo $ClientVersion |awk '{print $2}' |sed 's/ //g' | sed 's/\t//g')
#Version=${Version:-0.2.0}

OS=$1
echo "$OS"

TARGET_ROOT=`cd "$(dirname "$0")"; pwd`
#RELEASE_NAME=release_$V-$Version
RELEASE_NAME=release_${OS}_$V
RELEASE_ROOT=${TARGET_ROOT}/${RELEASE_NAME}

echo "$RELEASE_NAME"

if [ ! -d ${RELEASE_ROOT} ]; then
	rm -rf ${RELEASE_ROOT}		
fi

[ ! -e ${RELEASE_ROOT} ] && mkdir -p ${RELEASE_ROOT}

cd ${TARGET_ROOT}/

if [ $OS == 'windows' ]
	then
		echo "build AppHub-Agent.exe... "
		make AppHub-Agent.exe

		cd ${TARGET_ROOT}/

		# cp -a ffmpeg ${RELEASE_ROOT}/
		cp -a conf ${RELEASE_ROOT}/
		cp -a frontend ${RELEASE_ROOT}/
		cp  -a AppHub-Agent.exe  ${RELEASE_ROOT}/
		cp -a winpackage ${RELEASE_ROOT}/
elif [ $OS == 'linux' ]
	then
		echo "build x86_64 AppHub-Agent... "
		make AppHub-Agent

		if [ ! -d ${RELEASE_ROOT} ]; then
			rm -rf ${RELEASE_ROOT}		
		fi

		[ ! -e ${RELEASE_ROOT} ] && mkdir -p ${RELEASE_ROOT}

		# cp -a ffmpeg ${RELEASE_ROOT}/
		cp -a conf ${RELEASE_ROOT}/
		cp -a frontend ${RELEASE_ROOT}/
		cp  -a AppHub-Agent ${RELEASE_ROOT}/
		cp -a linuxpackage ${RELEASE_ROOT}/
		[ -a "startup.sh" ] && cp -a startup.sh ${RELEASE_ROOT}/

fi

echo "Syncing...."
sync;sync;
sync
sync

zip -r ${RELEASE_NAME}.zip ${RELEASE_NAME}

rm -rf ${RELEASE_ROOT}
echo "[Done]"

#!/bin/bash -x

# Unable to proceed if zfs.ko is absent.
if ! modprobe zfs; then
    echo "Unable to load zfs.ko on the host" >&2
    exit 1
fi

# Setup a loopback pool if pool is non-existent
if ! zpool list -H -o name ${NEMO_POOLNAME}; then
    if ! [[ -d ${NEMO_POOLDIR}/${NEMO_POOLNAME} ]]; then
        mkdir -p ${NEMO_POOLDIR}/${NEMO_POOLNAME}
    fi
    NEMO_ZPOOL=${NEMO_POOLDIR}/${NEMO_POOLNAME}
    if ! [[ -e ${NEMO_ZPOOL}/zpool ]]; then
        dd if=/dev/zero of=${NEMO_ZPOOL}/zpool \
           count=1 seek=$(expr ${NEMO_POOLSIZE} \* 2000000) bs=512
        LOOP_DEV=$(losetup -f)
        losetup ${LOOP_DEV} ${NEMO_ZPOOL}/zpool
        zpool create ${NEMO_POOLNAME} ${LOOP_DEV}
        zpool set cachefile=${NEMO_ZPOOL}/zpool.cache ${NEMO_POOLNAME}
    else
        LOOP_DEV=$(losetup -a | grep 'zpool (deleted)' | awk -F: '{ print $0 }')
        if [[ -z ${LOOP_DEV} ]]; then
            LOOP_DEV=$(losetup -f)
        else
            losetup -d ${LOOP_DEV}
            LOOP_DEV=$(losetup -f)
        fi
        losetup ${LOOP_DEV} ${NEMO_ZPOOL}/zpool
        zpool import -c ${NEMO_ZPOOL}/zpool.cache ${NEMO_POOLNAME}
    fi
fi

# Setup dory and k8s
DRIVER_NAME=${NEMO_DRIVER:-nemo}
DRIVER_PREFIX=${DORY_PREFIX:-dev.hpe.com}
DRIVER_DIR="/usr/libexec/kubernetes/kubelet-plugins/volume/exec/${DRIVER_PREFIX}~${DRIVER_NAME}"
DRIVER_PATH=${DRIVER_DIR}/${DRIVER_NAME}
DRIVER_SOCKET=/run/docker/plugins/${DRIVER_NAME}.sock
KUBELET_MOUNT=/var/lib/kubelet/plugins/${DRIVER_PREFIX}/mounts

if ! [[ -e ${KUBELET_MOUNT} ]]; then
	mkdir -p ${KUBELET_MOUNT}
fi

rm -rf ${DRIVER_DIR}
mkdir -p ${DRIVER_DIR}
cp /dory ${DRIVER_PATH}
sed -e 's#%DRIVER_SOCKET%#'${DRIVER_SOCKET}'#g' dory.json > ${DRIVER_PATH}.json

nemo -n ${DRIVER_NAME} -m ${KUBELET_MOUNT} -b ${NEMO_POOLNAME}/k8s

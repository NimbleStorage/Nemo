#!/bin/bash -x
if ! modprobe zfs; then
    echo "Unable to load zfs.ko on the host" >&2
    exit 1
fi
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

nemo -b ${NEMO_POOLNAME}/v2

BASEDS=tank/default
MOUNT=/opt/nemo
DRIVER=nemo
VOL:=$(shell date | md5sum | cut -f1 -d' ')
all:
	go get github.com/drajen/go-libzfs \
		github.com/docker/go-plugins-helpers/volume \
		github.com/Sirupsen/logrus \
		github.com/urfave/cli 
	go build -o nemo

test:
	echo $(shell date)
	docker volume create -d ${DRIVER} ${VOL} -o destroyOnRm=true
	docker run --rm -it -v ${VOL}:/data busybox touch /data/myfile.txt
	docker run --rm -it -v ${VOL}:/data busybox ls /data/myfile.txt
	docker volume create -d ${DRIVER} cloneof-${VOL} -o cloneOf=${VOL} -o description="a clone of ${VOL}" -o destroyOnRm=true
	docker volume inspect cloneof-${VOL} | grep '"Description": "a clone of ${VOL}'
	docker run --rm -it -v cloneof-${VOL}:/data busybox ls /data/myfile.txt
	docker volume rm cloneof-${VOL}
	docker volume rm ${VOL}
	docker volume create -d ${DRIVER} ${VOL} -o destroyOnRm=false
	docker volume rm ${VOL}
	docker volume create -d ${DRIVER} imported-${VOL} -o importVol=$$(zfs list -o name | grep $(VOL) | sed -e 's#$(BASEDS)/##'g) -o destroyOnRm=true
	docker volume rm imported-${VOL}
	zfs create ${BASEDS}/dataset1
	touch ${MOUNT}/dataset1/myfile.txt
	zfs create ${BASEDS}/dataset2
	zfs snapshot ${BASEDS}/dataset1@now
	zfs snapshot ${BASEDS}/dataset2@later
	touch ${MOUNT}/dataset2/myfile.txt
	docker volume create -d ${DRIVER} imported-${VOL} -o importVolAsClone=dataset1 -o destroyOnRm=true 
	docker run --rm -it -v imported-${VOL}:/data busybox ls /data/myfile.txt
	docker volume rm imported-${VOL}
	docker volume create -d ${DRIVER} imported-${VOL} -o importVolAsClone=dataset1 -o destroyOnRm=true -o snapshot=now
	docker run --rm -it -v imported-${VOL}:/data busybox ls /data/myfile.txt
	docker volume rm imported-${VOL}
	zfs unmount ${BASEDS}/dataset1
	zfs unmount ${BASEDS}/dataset2
	docker volume create -d ${DRIVER} imported-${VOL} -o importVol=dataset1 -o destroyOnRm=true -o restore=true
	docker run --rm -it -v imported-${VOL}:/data busybox ls /data/myfile.txt
	docker volume rm imported-${VOL}
	docker volume create -d ${DRIVER} imported-${VOL} -o importVol=dataset2 -o destroyOnRm=true -o restore=true -o snapshot=later
	if docker run --rm -it -v imported-${VOL}:/data busybox ls /data/myfile.txt; then false; fi
	docker volume rm imported-${VOL}
	if zfs list ${BASEDS} | grep ${BASEDS}/; then false; fi

ifeq ($(VERSION),)
  $(error VERSION is not set)
endif
all:
	docker build -t nemo-v2 .
	docker run --name nemo-v2 nemo-v2 true
	mkdir rootfs
	docker export nemo-v2 | tar x -C rootfs
	docker plugin create nimblestorage/nemo:$(VERSION) .
	docker plugin create nimblestorage/nemo:latest .
clean:
	docker rm nemo-v2 || true
	docker rmi nemo-v2 || true
	docker plugin rm nimblestorage/nemo:$(VERSION) || true
	docker plugin rm nimblestorage/nemo:latest || true
	rm -rf rootfs || true
push:
	docker plugin push nimblestorage/nemo:$(VERSION)
	docker plugin push nimblestorage/nemo:latest

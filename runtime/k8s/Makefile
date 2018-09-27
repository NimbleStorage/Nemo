ifeq ($(VERSION),)
  $(error VERSION is not set)
endif
all:
	docker build -t nemod:latest .
	docker tag nemod nimblestorage/nemod:latest
	docker tag nemod nimblestorage/nemod:$(VERSION) 
clean:
	docker rmi nemod:latest || true
	docker rmi nemod:$(VERSION) || true
push:
	docker push nimblestorage/nemod:$(VERSION)
	docker push nimblestorage/nemod:latest

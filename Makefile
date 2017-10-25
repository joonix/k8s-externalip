APP?=externalip_init
VERSION?=v0.1.0
REGISTRY?=eu.gcr.io/foobar
GOLANG?=1.9.1

docker:
	docker run --rm -v $(CURDIR):/go/src/joonix/k8s-externalip-init \
		-w /go/src/joonix/k8s-externalip-init \
		-e CGO_ENABLED=0 golang:$(GOLANG) go build .
	docker build -t $(REGISTRY)/$(APP):$(VERSION) .
.PHONY: docker

release: docker
	gcloud docker -- push $(REGISTRY)/$(APP):$(VERSION)
.PHONY: release

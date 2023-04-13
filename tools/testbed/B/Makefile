IMG ?= pluviophile225/test_b_param
docker-build:
	docker buildx build -t ${IMG} . --load

docker-push:
	docker push ${IMG} 
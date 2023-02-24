services := auth_service products_service

build_fn = docker build -t ghcr.io/manoamaro/$(1):latest -f $(1)/Dockerfile .
publish_fn = docker push ghcr.io/manoamaro/$(1):latest

build:
	$(foreach service,$(services), $(call build_fn,$(service));)

publish: build
	echo $(DOCKER_PASSWORD) | docker login ghcr.io -u $(DOCKER_USERNAME) --password-stdin
	$(foreach service,$(services),$(call publish_fn,$(service));)

monolith:
	go run store-monolith/cmd/main.go
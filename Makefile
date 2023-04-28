VERSION ?= 0.1.0
REPO = github.com/manoamaro/microservice-store
BRANCH ?= develop
MAIN_BRANCH ?= main

build-%:
	cd $(@:build-%=%) && go build -o ../dist/$(notdir $(@:/=)) ./cmd/main.go

run-%:
	cd $(@:run-%=%) && go run ./cmd/main.go

test-%:
	cd $(@:test-%=%) && go test ./...

release-start:
	git checkout $(BRANCH)
	git pull origin $(BRANCH)
	git checkout -b release-$(VERSION)

release-finish:
	git checkout $(MAIN_BRANCH)
	git pull origin $(MAIN_BRANCH)
	git merge release-$(VERSION) --no-ff -m "Merge release-$(VERSION) into $(MAIN_BRANCH)"
	git tag -a v$(VERSION) -m "Version $(VERSION)"
	git branch -d release-$(VERSION)
	git push origin $(MAIN_BRANCH) v$(VERSION)
	git checkout $(BRANCH)
	git merge $(MAIN_BRANCH) --no-ff -m "Merge $(MAIN_BRANCH) into $(BRANCH)"
	git push origin $(BRANCH)

update-version:
	go mod edit -version $(REPO)/v$(VERSION)

changelog:
	@echo "Updating CHANGELOG.md..."
	@echo "" > CHANGELOG.tmp.md
	@echo "## [$(VERSION)] - $$(date +'%Y-%m-%d')" >> CHANGELOG.tmp.md
	@echo "" >> CHANGELOG.tmp.md
	@echo "### Added" >> CHANGELOG.tmp.md
	@git log --pretty=format:'- %s' --grep='^Added' $(PREV_VERSION)..HEAD >> CHANGELOG.tmp.md || true
	@echo "" >> CHANGELOG.tmp.md
	@echo "### Removed" >> CHANGELOG.tmp.md
	@git log --pretty=format:'- %s' --grep='^Removed' $(PREV_VERSION)..HEAD >> CHANGELOG.tmp.md || true
	@echo "" >> CHANGELOG.tmp.md
	@echo "### Changed" >> CHANGELOG.tmp.md
	@git log --pretty=format:'- %s' --grep='^Changed' $(PREV_VERSION)..HEAD >> CHANGELOG.tmp.md || true
	@echo "" >> CHANGELOG.tmp.md
	@cat CHANGELOG.md >> CHANGELOG.tmp.md
	@mv CHANGELOG.tmp.md CHANGELOG.md
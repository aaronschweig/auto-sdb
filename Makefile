VERSION=dev

.PHONY: release
release:
	docker build -t aaronschweig/auto-sdb:$(VERSION) .
	docker push aaronschweig/auto-sdb:$(VERSION)
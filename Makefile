VERSION=dev

build-frontend:
	cd frontend && pnpm run build

build: build-frontend
	go mod vendor
	docker build -t aaronschweig/auto-sdb:$(VERSION) . -f Containerfile
	rm -rf vendor

release: build
	docker push aaronschweig/auto-sdb:$(VERSION)
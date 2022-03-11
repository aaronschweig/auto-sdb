VERSION=dev

build-frontend:
	cd frontend && pnpm run build

build: build-frontend
	docker build -t aaronschweig/auto-sdb:$(VERSION) . -f Containerfile

release: build
	docker push aaronschweig/auto-sdb:$(VERSION)
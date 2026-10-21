BIN := bk
BIN_PATH := /usr/local/bin/bk
ARCHES := amd64 arm64 386 s390x
ARMS := 5 6 7

.PHONY: all prepare build check release-linux release version install sbom
all: build
prepare:
	go run ./tools/preparecore
build:
	go run ./tools/build -o $(BIN)
version:
	@grep -q '"distribution_version": "'"$$(cat VERSION | tr -d '\r\n')"'"' upstream.lock.json
check: prepare
	go test ./... -race
	go vet ./...
	python3 -m unittest discover -s scripts -p 'test_*.py'
	cd addons && go vet ./... && go test ./... -race -timeout 20m
	cd .build/core && go vet ./... && go test ./... -race -timeout 20m
	@if command -v node >/dev/null; then node --test .build/core/internal/webui/paneltest/*.test.js; else echo 'Node tests require Node.js; run them before publishing.'; fi
release-linux:
	mkdir -p dist
	@for arch in $(ARCHES); do go run ./tools/build -goos linux -goarch $$arch -o dist/bk-linux-$$arch || exit 1; done
	@for arm in $(ARMS); do go run ./tools/build -goos linux -goarch arm -goarm $$arm -o dist/bk-linux-armv$$arm || exit 1; done
release: version release-linux
	mkdir -p release
	@for arch in $(ARCHES) $(addprefix armv,$(ARMS)); do cp dist/bk-linux-$$arch dist/bk && tar -czf release/bk_linux_$$arch.tar.gz -C dist bk && rm dist/bk || exit 1; done
	cd release && sha256sum bk_linux_*.tar.gz > SHA256SUMS
	go -C .build/core run ./tools/signsums "$(CURDIR)/release/SHA256SUMS" "$${RELEASE_TAG:-$$(cat VERSION | tr -d '\r\n')}"
	$(MAKE) sbom
sbom:
	mkdir -p release
	@cat .build/sbom/*.txt > release/SBOM.txt
install: build
	install -m 0755 $(BIN) $(BIN_PATH)

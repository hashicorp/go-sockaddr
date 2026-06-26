TOOLS= golang.org/x/tools/cover
GOCOVER_TMPFILE?=	$(GOCOVER_FILE).tmp
GOCOVER_FILE?=	.cover.out
GOCOVERHTML?=	coverage.html
FIND=`/usr/bin/which 2> /dev/null gfind find | /usr/bin/grep -v ^no | /usr/bin/head -n 1`
XARGS=`/usr/bin/which 2> /dev/null gxargs xargs | /usr/bin/grep -v ^no | /usr/bin/head -n 1`

test:: $(GOCOVER_FILE)
	@$(MAKE) -C cmd/sockaddr test

cover:: coverage_report

$(GOCOVER_FILE)::
	@${FIND} . -type d ! -path '*cmd*' ! -path '*.git*' -print0 | ${XARGS} -0 -I % sh -ec "cd % && rm -f $(GOCOVER_TMPFILE) && go test -coverprofile=$(GOCOVER_TMPFILE)"

	@echo 'mode: set' > $(GOCOVER_FILE)
	@${FIND} . -type f ! -path '*cmd*' ! -path '*.git*' -name "$(GOCOVER_TMPFILE)" -print0 | ${XARGS} -0 -n1 cat $(GOCOVER_TMPFILE) | grep -v '^mode: ' >> ${PWD}/$(GOCOVER_FILE)

$(GOCOVERHTML): $(GOCOVER_FILE)
	go tool cover -html=$(GOCOVER_FILE) -o $(GOCOVERHTML)

coverage_report:: $(GOCOVER_FILE)
	go tool cover -html=$(GOCOVER_FILE)

deps::
	@go install github.com/hashicorp/copywrite@b3e6599f43beff698f471c6f46888045453fa030 # v0.25.3
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@c0d3ddc9cf3faa61a4e378e879ece580256d76e5 # v2.12.2

.PHONY: copywriteheaders
copywriteheaders:
	@echo "==> Running copywrite headers plan..."
	@copywrite headers --plan
	@echo "==> Done"

lint::
	@echo "==> Running linters..."
	@golangci-lint run ./...
	@echo "==> Done"

clean::
	rm -f $(GOCOVER_FILE) $(GOCOVERHTML)

dev::
	@go build
	@$(MAKE) -B -C cmd/sockaddr sockaddr

install::
	@go install
	@$(MAKE) -C cmd/sockaddr install

doc::
	@echo Visit: http://127.0.0.1:6161/pkg/github.com/hashicorp/go-sockaddr/
	godoc -http=:6161 -goroot $GOROOT

world::
	@set -e; \
	for os in solaris darwin freebsd linux windows android; do \
		for arch in amd64; do \
			printf "Building on %s-%s\n" "$${os}" "$${arch}" ; \
			env GOOS="$${os}" GOARCH="$${arch}" go build -o /dev/null; \
		done; \
	done

	$(MAKE) -C cmd/sockaddr world

default: build

build b:
	go build

fmt f:
	golangci-lint fmt ./... 
 
lint l:
	golangci-lint run ./...

test t:
	go test ./...

cover test_cover:
	go test ./... --cover

semgrep:
	docker run --rm -v ${PWD}:/app:ro -w /app semgrep/semgrep semgrep scan --error --metrics=off \
		--config=p/command-injection \
		--config=p/comment \
		--config=p/cwe-top-25 \
		--config=p/default \
		--config=p/gitleaks \
		--config=p/golang \
		--config=p/gosec \
		--config=p/insecure-transport \
		--config=p/owasp-top-ten \
		--config=p/r2c-best-practices \
		--config=p/r2c-bug-scan \
		--config=p/r2c-security-audit \
		--config=p/secrets \
		--config=p/security-audit \
		--config=p/sql-injection \
		--config=p/xss \
		.

.PHONY: lint cover build test fmt semgrep
# See PLAN.md for the roadmap and PORTING.md for the translation rules.
.PHONY: all build vet test oracle-tools

all: build vet test

build:
	go build ./...

# Vet everything under internal/ (the whole port).
vet:
	go vet ./internal/...

test:
	go test ./...

# Pinned Closure Compiler jar + sources, and the Java oracle tools (needs a JDK).
oracle-tools:
	bash tools/oracle/build.sh

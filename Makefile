# See PLAN.md for the roadmap and PORTING.md for the translation rules.
.PHONY: all build vet test oracle-tools

all: build vet test

build:
	go build ./...

# Vet everything under internal/ (the whole port), check the licence
# notices (tools/licenses; see PORTING.md, top) and that every Go map
# iteration says why its order cannot matter (tools/maprangecheck).
vet:
	go vet ./internal/...
	go run ./tools/licenses -check
	go run ./tools/maprangecheck

test:
	go test ./...

# Pinned Closure Compiler jar + sources, and the Java oracle tools (needs a JDK).
oracle-tools:
	bash tools/oracle/build.sh

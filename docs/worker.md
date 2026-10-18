# Persistent command-line worker

`jsgo` keeps its ordinary one-shot command-line behavior. `jsgo-worker` is an
optional persistent process for build tools and watch loops. Each request creates
a fresh compiler and rereads/recompiles the user inputs. It adds a bounded cache
of built-in extern parse trees; templates are cloned before use. It does not
cache compiled output or user ASTs.

Build from the repository root:

```sh
go build -o bin/jsgo ./cmd/jsgo
go build -pgo=cmd/jsgo/default.pgo -o bin/jsgo-worker ./cmd/jsgo-worker
```

The explicit PGO argument gives the worker the same profile as the CLI. Go 1.24
or newer is supported. Java character classification/case mappings stay pinned
to Unicode 15, matching the reference JDK 21, when using a newer Go compiler.

Send JSON objects on stdin, conventionally one per line. The worker emits one
JSON response per request on stdout. Requests run sequentially in the worker's
working directory. Arguments are the normal Closure/jsgo CLI arguments; output
files are written normally.

```json
{"id":1,"args":["--compilation_level=WHITESPACE_ONLY","--js","input.js"]}
```

The response contains `id`, the shell `exit` code, `stdout_base64`, and
`stderr_base64`. The base64 fields contain the exact output bytes, including
non-UTF-8 encodings. Supply optional `stdin_base64` for `--js -` or JSON-stream
input. An empty base64 string means empty output. A compilation error produces
a response and the worker continues accepting requests. Malformed protocol JSON
terminates the worker with status 2 and a message on its stderr.

Keep the process open for subsequent requests. Closing stdin cleanly shuts it
down. The normal CLI remains useful for one-off builds; the worker requires a
client that sends requests and decodes responses.

See the [performance report](performance.md) for the workload-specific comparison. Worker comparisons must give both runtimes equivalent extern reuse. The earlier
small WHITESPACE comparison against stock Java used unequal cache policies;
that result does not establish a win against a comparable Java worker. The
benchmark-only Java worker and matched comparison are documented in
[tools/bench/java/README.md](../tools/bench/java/README.md).

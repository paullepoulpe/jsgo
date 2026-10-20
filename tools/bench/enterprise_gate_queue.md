# Sequential enterprise correctness gates

`enterprise_gate_queue.py --plan <frozen-plan.json> --out <new-directory>`
executes one stock Java compilation, then the matching Go compilation, then
independent Node checks on both outputs when a generated-runtime oracle exists.
It stops at the first failure and keeps all subprocess commands, streams, exits
and partial outputs. A retry needs a new output directory.

The plan pins the Go executable, unchanged stock jar, driver classes, Node,
helper source and exact task specs. Existing capped helpers enforce an actual
4 GiB process cgroup, zero swap, two CPU quota and affinity 0,1. Java uses the
same explicit 3 GiB maximum heap as the paired benchmark. There is no worker,
extern AST cache, compiler overlay, runtime GC tuning or performance claim.

Every compiler invocation uses identical complete arguments and frozen source
and default-extern bytes. All stdout, stderr, exit, output JavaScript, source
maps, manifests and expected absences must agree with stock Java. Chunk tasks
also require every chunk artifact, the synthetic weak manifest and dependency
JSON, with the declared input count covering every original source once.

J2CL and Soy require the independent generated-runtime result on both compiled
outputs. Browser application/chunk and mega tasks currently have no independent
application runtime oracle. Their report states this limitation explicitly;
matching bytes do not invent that coverage. A failed Java fit under the cap
stops the queue and is a failure, not an infinite Go speed improvement.

Only run the queue in a scheduled correctness slot, outside timing windows.
Raw subprocess clocks are diagnostic data. Successful correctness gates must
precede a separate warm-Java/fresh-Go paired benchmark.

Pure model tests are in `enterprise_gate_queue_test.py`; they do not start
compilers, runtimes or containers.

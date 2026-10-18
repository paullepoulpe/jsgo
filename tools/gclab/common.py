"""Shared helpers: the benchmark's tasks, workspaces and run_proc (tools/bench/bench.py)."""
import json, os, sys

S = "/tmp/claude-0/-home-user-jsgo/6e91e3d3-452b-57d0-b724-37154263b2f7/scratchpad"
WT = os.path.join(S, "wt-gclab")
LAB = os.path.join(S, "gclab")
REF = os.path.join(S, "bench2", "results")
sys.path.insert(0, os.path.join(WT, "tools", "bench"))
sys.dont_write_bytecode = True
import bench  # noqa: E402

CACHE = bench.default_cache()
TASKS = {t["name"]: t for t in bench.load_tasks()}
SINGLE_STEP = 10  # the fixed 100-file subset: every 10th file of tasks/single-files.txt


def go_env(extra=None):
    env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP and not k.startswith("JSGO_GCLAB")}
    env["LC_ALL"] = "C.UTF-8"
    env.update(extra or {})
    return env


def jobs_of(task, ws):
    """[(index in the stored results' job list, argv)]"""
    if task == "single-100":
        jobs = bench.task_jobs(TASKS["single-files"], ws)
        return [(i, jobs[i]) for i in range(0, len(jobs), SINGLE_STEP)]
    return [(0, bench.task_jobs(TASKS[task], ws)[0])]


def ref_task(task):
    return "single-files" if task == "single-100" else task


def ref_digests(task, who):
    d = json.load(open(os.path.join(REF, ref_task(task), who + ".json")))["runs"][0]["digests"]
    return d


def compare_names(task):
    return TASKS[ref_task(task)].get("compare", [])

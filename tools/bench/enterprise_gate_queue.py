#!/usr/bin/env python3
"""Sequential capped compiler/runtime correctness gates, never a benchmark.

Preparation writes a plan with exact executable/source hashes. Execution stops
at the first failure, preserving every prior artifact and failed subprocess log.
Run only in an explicitly scheduled correctness slot, never during timed trials.
"""
import argparse
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

sys.dont_write_bytecode = True
import paired_cli as control
HERE = Path(__file__).resolve().parent
PINNED_JAR = 'd6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6'
SCOPE = 'Sequential diagnostic correctness/setup only; clocks are not performance evidence'


def sha(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2)+'\n')


def check_pins(pins):
    for path, digest in pins.items():
        if sha(path) != digest:
            raise ValueError('frozen dependency changed: '+path)


def validate_spec(spec, semantics):
    if spec.get('schema') != 'jsgo-enterprise-cli-v1':
        raise ValueError('unsupported enterprise spec')
    argv = spec['argv']
    if any(not isinstance(arg, str) or any(c in arg for c in '\r\n\x00') for arg in argv):
        raise ValueError('invalid compiler argv')
    if not {'out.js', 'out.js.map', 'inputs.manifest'} <= set(spec['artifacts']) and not spec['same_work'].get('chunks'):
        raise ValueError('need source/map/manifest artifact gates')
    if semantics and not all(spec.get('semantic_check', {}).get(key) for key in ('script', 'entry', 'expected')):
        raise ValueError('generated fixture needs its independent runtime oracle')
    if semantics and spec['semantic_check'].get('native_crosscheck_required'):
        raise ValueError('this API fixture also requires its separate native protobuf crosscheck; use the complete staged workflow')
    inputs = [item for item in spec['inputs'] if item['flag'] == '--js']
    paths = [item['link']+'/'+item['relative'] for item in inputs]
    actual = [argv[i+1] for i, arg in enumerate(argv) if arg == '--js']
    if actual != paths or len(paths) != len(set(paths)):
        raise ValueError('compiler must receive every frozen source exactly once in order')
    chunks = spec['same_work'].get('chunks', [])
    if chunks:
        flags = [argv[i+1] for i, arg in enumerate(argv) if arg == '--chunk']
        expected = [c['name']+':'+str(c['input_count'])+(':'+','.join(c['dependencies']) if c['dependencies'] else '') for c in chunks]
        if flags != expected or sum(c['input_count'] for c in chunks) != len(paths):
            raise ValueError('chunk definitions do not cover the exact frozen source list')
        required = {f'chunks/{c["name"]}.js{suffix}' for c in chunks for suffix in ('', '.map', '.inputs')}
        required |= {'chunks/$weak$.js.inputs', 'chunk-dependencies.json'}
        if set(spec['artifacts']) != required or set(spec.get('absent_artifacts', [])) != {'chunks/$weak$.js', 'chunks/$weak$.js.map'}:
            raise ValueError('incomplete chunk/maps/manifests/dependency/absence gate')


def run_step(command, root, label, pins):
    check_pins(pins)
    before = control.outer_state()
    save(root/(label+'.outer-before.json'), before)
    save(root/(label+'.command.json'), {'argv': command, 'scope': SCOPE})
    with (root/(label+'.stdout')).open('wb') as stdout, (root/(label+'.stderr')).open('wb') as stderr:
        completed = subprocess.run(command, stdout=stdout, stderr=stderr)
    save(root/(label+'.exit.json'), {'exit': completed.returncode, 'scope': SCOPE})
    after = control.outer_state()
    save(root/(label+'.outer-after.json'), after)
    assert_no_outer_oom(before, after)
    check_pins(pins)
    if completed.returncode:
        raise RuntimeError(label+' failed; original artifacts/logs retained')


def assert_no_outer_oom(before, after):
    for event in ('oom', 'oom_kill'):
        if control.event_count(before.get('memory.events', ''), event) != control.event_count(after.get('memory.events', ''), event):
            raise RuntimeError('outer workspace cgroup '+event+' changed; correctness gate invalidated')


def kernel(memory, swap, cpu, cpuset):
    quota, period = cpu.split()
    if str(memory) != '4294967296' or str(swap) != '0' or quota == 'max' or int(quota) != 2*int(period) or cpuset not in ('0-1', '0,1'):
        raise ValueError('actual kernel4GiB/swap0/CPU2/affinity0,1 required')


def verify_task(root, spec, semantic, node_dirs):
    java = json.loads((root/'java-oracle/result.json').read_text())
    go = json.loads((root/'go-bytegate/result.json').read_text())
    jm = json.loads((root/'java-oracle/metadata.json').read_text())
    gm = json.loads((root/'go-bytegate/metadata.json').read_text())
    row = java['row']
    if row['exit'] or go['exit'] or go.get('oracle_gate') != 'MATCH' or row['digests'] != go['digests']:
        raise ValueError('every exit/stream/artifact/absence must match successful stock Java')
    if jm['spec'] != spec or gm['spec'] != spec or jm['input_manifest'] != gm['input_manifest'] or jm['jar_sha256'] != PINNED_JAR:
        raise ValueError('stock Java and Go must use exactly the same frozen spec and inputs')
    fields = json.loads((root/'java-oracle/java/ready.json').read_text())['fields']
    # validate_ready() already enforces these in the stock oracle subprocess;
    # preserve that actual kernel report rather than infer from a heap setting.
    if not fields:
        raise ValueError('missing actual Java kernel READY report')
    cg = go['cgroup']; kernel(cg['memory.max'], cg['memory.swap.max'], cg['cpu.max'], cg['cpuset.cpus.effective'])
    proof = {'status': 'PASS', 'scope': SCOPE, 'all_artifact_digests': row['digests'],
             'input_identity_equal': True, 'argv_equal': True, 'java_kernel_ready': fields,
             'go_kernel': cg, 'semantic_coverage': 'not an independent app runtime oracle; exact compiler/artifact checks only'}
    if spec['same_work'].get('chunks'):
        for work in (root/'java-oracle/java', root/'go-bytegate/run'):
            dependencies = json.loads((work/'chunk-dependencies.json').read_text())
            if not dependencies:
                raise ValueError('empty chunk dependency artifact')
        proof['chunk_count'] = len(spec['same_work']['chunks'])
        proof['source_count'] = sum(item['flag'] == '--js' for item in spec['inputs'])
    if semantic:
        proof['runtime_semantics'] = {}
        for name, directory in node_dirs.items():
            metrics = json.loads((directory/'result.json').read_text())
            result = json.loads((directory/'run/stdout').read_text())
            cg = metrics['cgroup']; kernel(cg['memory.max'], cg['memory.swap.max'], cg['cpu.max'], cg['cpuset.cpus.effective'])
            if metrics['exit'] or result.get('status') != 'PASS' or result.get('result') != spec['semantic_check']['expected']:
                raise ValueError('independent '+name+' compiled runtime semantics failed')
            proof['runtime_semantics'][name] = result
        proof['semantic_coverage'] = 'both independently executed compiled outputs match the frozen generated-runtime oracle'
    return proof


def execute(plan, out):
    if plan.get('schema') != 'jsgo-enterprise-gate-queue-v1' or out.exists():
        raise ValueError('need supported plan and a new evidence output directory')
    names = [task['name'] for task in plan['tasks']]
    if not names or len(set(names)) != len(names) or any(not re.fullmatch('[a-z0-9][a-z0-9-]*', name) for name in names):
        raise ValueError('need distinct safe task names')
    if any(type(task['semantics']) is not bool for task in plan['tasks']):
        raise ValueError('runtime semantic coverage must be explicit')
    if not re.fullmatch(r'.+@sha256:[0-9a-f]{64}', plan['image']):
        raise ValueError('pinned image digest required')
    required_pins = {str(Path(__file__).resolve()), plan['binary'], plan['jar'], plan['node'],
                     str(Path(plan['java_home'])/'bin/java'), str(Path(plan['driver'])/'CliPaired.class'),
                     plan['semantic_runner'], plan['semantic_setup_runner']}
    required_pins |= {str(HERE/name) for name in ('capped_java_oracle.py', 'capped_cli.py', 'paired_cli.py',
                                               'paired_factors.py', 'limited_cli.py', 'cgroup_exec.py', 'go_factor_screen.py')}
    required_pins |= {task['spec'] for task in plan['tasks']}
    required_pins |= {str(path) for path in Path(plan['driver']).rglob('*.class')}
    for task in plan['tasks']:
        if task['semantics']:
            required_pins.add(json.loads(Path(task['spec']).read_text())['semantic_check']['script'])
    if Path(plan['semantic_setup_runner']) != Path(plan['semantic_runner']).with_name('run_capped_setup.py'):
        raise ValueError('the pinned semantic setup helper must match the actual runner import')
    if not required_pins <= set(plan['pins']):
        raise ValueError('plan must pin every executable, helper and task spec')
    check_pins(plan['pins'])
    if sha(plan['jar']) != PINNED_JAR:
        raise ValueError('unchanged repository-pinned stock Java jar required')
    out.mkdir(parents=True)
    save(out/'plan.json', plan)
    before = control.outer_state()
    save(out/'outer-before.json', before)
    results = []
    try:
        for task in plan['tasks']:
            root = out/task['name']; root.mkdir()
            spec = json.loads(Path(task['spec']).read_text())
            validate_spec(spec, task['semantics'])
            save(root/'spec.json', spec)
            common = ['--spec', str(root/'spec.json'), '--binary', plan['binary'], '--image', plan['image']]
            run_step([sys.executable, str(HERE/'capped_java_oracle.py')]+common+
                     ['--jar', plan['jar'], '--driver', plan['driver'], '--java-home', plan['java_home'], '--out', str(root/'java-oracle')], root, '01-java', plan['pins'])
            run_step([sys.executable, str(HERE/'capped_cli.py')]+common+
                     ['--oracle', str(root/'java-oracle/stock-java-oracle-digests.json'), '--out', str(root/'go-bytegate')], root, '02-go', plan['pins'])
            node_dirs = {}
            if task['semantics']:
                for name, compiled in [('java', root/'java-oracle/java/out.js'), ('go', root/'go-bytegate/run/out.js')]:
                    directory = root/('node-'+name); node_dirs[name] = directory
                    run_step([sys.executable, plan['semantic_runner'], '--spec', str(root/'spec.json'), '--compiled-js', str(compiled), '--node', plan['node'], '--out', str(directory)], root, '03-node-'+name, plan['pins'])
            proof = verify_task(root, spec, task['semantics'], node_dirs)
            save(root/'validation.json', proof)
            results.append({'name': task['name'], 'status': 'PASS', 'validation': str(root/'validation.json')})
            save(out/'progress.json', {'scope': SCOPE, 'completed': results})
        after = control.outer_state()
        save(out/'outer-after.json', after)
        assert_no_outer_oom(before, after)
    except BaseException as error:
        save(out/'failure.json', {'type': type(error).__name__, 'message': str(error), 'completed': results,
                                 'outer_state': control.outer_state(), 'scope': SCOPE})
        raise
    save(out/'result.json', {'status': 'PASS', 'scope': SCOPE, 'tasks': results})


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--plan', required=True); p.add_argument('--out', required=True)
    a = p.parse_args()
    execute(json.loads(Path(a.plan).read_text()), Path(a.out).resolve())


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
"""Independent, read-only receipt checks; never launches Java or reads archives."""
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import re
from urllib.parse import unquote, urlparse

HERE = Path(__file__).resolve().parent
BASE = HERE.parent.parent
PINNED = Path('/workspace/jsgo/.cache/closure-src')
JAR = '/workspace/jsgo/.cache/closure-compiler-v20261015-maven.jar'
FLAGS = ['-source', '21', '-target', '21', '-g', '-proc:none', '-implicit:none', '-encoding', 'UTF-8']
RELOCATION = r'(?<![\w.])com\.google\.(common|auto|errorprone|gson|protobuf|j2objc|thirdparty)(?=\.)'
INPUTS = {}


def require(condition, detail):
    if not condition:
        raise ValueError(detail)


def data(path):
    path = Path(path)
    require(path.suffix not in ('.jar', '.jfr', '.zip'), str(path))
    value = path.read_bytes()
    digest = hashlib.sha256(value).hexdigest()
    previous = INPUTS.setdefault(str(path), digest)
    require(previous == digest, 'Changed input: ' + str(path))
    return value


def sha(path):
    data(path)
    return INPUTS[str(Path(path))]


def obj(path):
    return json.loads(data(path))


def relocate(value):
    return re.sub(RELOCATION, r'com.google.javascript.jscomp.jarjar.com.google.\1',
                  value.decode('utf-8')).encode('utf-8')


def class_files(directory):
    directory = Path(directory)
    return {str(p.relative_to(directory)): sha(p) for p in sorted(directory.rglob('*.class'))}


def origins(log, classname):
    pattern = re.compile(r'\[class,load\]\s+' + re.escape(classname) + r'\s+source:\s+(\S+)')
    found = set(pattern.findall(log))
    require(bool(found), 'Missing class load: ' + classname)
    result = set()
    for value in found:
        parsed = urlparse(value)
        require(parsed.scheme == 'file', 'Unexpected class source: ' + value)
        result.add(str(Path(unquote(parsed.path)).resolve()))
    return result


def build_review(name, relative, control, candidate):
    manifest_path = BASE/'build'/name/'manifest.json'
    manifest = obj(manifest_path)
    require(manifest['status'] == 'PASS_BUILD_VALIDATION', name)
    require(manifest['javac_flags'] == FLAGS, name)
    require(manifest['compiler_source_families'] == [relative[:-5]], name)
    require(set(manifest['variants']) == {'STOCK', 'CONTROL', name}, name)
    require(data(control) == data(PINNED/relative), 'Control is not pinned original: ' + name)
    entries = {}
    for variant, source in [('CONTROL', control), (name, candidate)]:
        entry = manifest['variants'][variant]
        receipt_path = Path(entry['build_proof'])
        require(sha(receipt_path) == entry['build_proof_sha256'], str(receipt_path))
        require(sha(receipt_path) == manifest['compiler_build_receipt_sha256'][variant], name)
        receipt = obj(receipt_path)
        require(receipt['status'] == 'PASS' and receipt['exit'] == 0, str(receipt_path))
        require(receipt['variant'] == variant and receipt['javac_flags'] == FLAGS, variant)
        require(receipt['original_sources'] == {relative: sha(control)}, variant)
        require(receipt['effective_sources'] == {relative: sha(source)}, variant)
        relocated = receipt_path.parent/'src'/relative
        require(data(relocated) == relocate(data(source)), 'Build shading: ' + variant)
        require(receipt['relocated_sources'] == {relative: sha(relocated)}, variant)
        require(receipt['jar_sha256'] == manifest['jar_sha256'], variant)
        require(receipt['source_jar_sha256'] == manifest['source_jar_sha256'], variant)
        require(entry['classes'] == receipt['classes'] == class_files(entry['overlay']), variant)
        require(all(p[:-6] == relative[:-5] or p[:-6].startswith(relative[:-5]+'$')
                    for p in entry['classes']), 'Unexpected production class')
        expected = [manifest['java_home']+'/bin/java', '-Xmx1g', '-m',
                    'jdk.compiler/com.sun.tools.javac.Main', *FLAGS, '-cp', JAR,
                    '-sourcepath', str(receipt_path.parent/'src'), '-d', entry['overlay'],
                    str(relocated)]
        require(receipt['argv'] == expected, 'Build argv: '+variant)
        entries[variant] = {'build_proof': str(receipt_path),
                            'build_proof_sha256': sha(receipt_path),
                            'effective_source_sha256': sha(source),
                            'classes_verified': len(entry['classes'])}
    return manifest, {'manifest': str(manifest_path), 'sha256': sha(manifest_path),
                      'status': 'PASS', 'variants': entries}


def probe_review(directory_name, candidate, classname, source, expected_args, build):
    directory = BASE/'results'/directory_name
    receipt = obj(directory/'manifest.json')
    require(receipt['status'] == 'PASS' and receipt['execution_complete'] is True, directory_name)
    require(receipt['candidate'] == candidate and receipt['expected_runs'] == 3, directory_name)
    require(receipt['preflight_runtimes'] == receipt['postflight_runtimes'] == [], directory_name)
    require(receipt['build_manifest_unchanged'] is True, directory_name)
    require(all(receipt['source_integrity'].values()) and all(receipt['overlay_integrity'].values()), directory_name)
    require(data(directory/'build-manifest.json') == data(receipt['build_manifest']), directory_name)
    require(sha(receipt['build_manifest']) == receipt['build_manifest_sha256'], directory_name)
    require(sha(directory/'runner.py') == receipt['runner_sha256'], directory_name)
    require(receipt['jar'] == JAR and receipt['jar_sha256'] == build['jar_sha256'], directory_name)
    require(receipt['javac_flags'] == FLAGS and receipt['java_version_from_build_manifest'] == build['java_version'], directory_name)
    require(receipt['java'] == build['java_home']+'/bin/java', directory_name)
    relative = classname.replace('.', '/')+'.java'
    require(set(receipt['sources']) == {relative}, directory_name)
    source_receipt = receipt['sources'][relative]
    require(source_receipt['original_path'] == str(source), directory_name)
    require(sha(source) == source_receipt['original_sha256'], directory_name)
    require(data(directory/'original-sources'/relative) == data(source), directory_name)
    require(data(directory/'relocated-sources'/relative) == relocate(data(source)), directory_name)
    require(sha(directory/'relocated-sources'/relative) == source_receipt['relocated_sha256'], directory_name)
    require(class_files(directory/'classes') == receipt['probe_classes'], directory_name)
    stem = classname.replace('.', '/')
    require(all(p[:-6] == stem or p[:-6].startswith(stem+'$') for p in receipt['probe_classes']), directory_name)
    require(list((directory/'empty-sourcepath').iterdir()) == [], directory_name)
    compile_result = receipt['compile']
    expected_compile = [receipt['java'], '-Xmx1g', '-m', 'jdk.compiler/com.sun.tools.javac.Main',
                        *FLAGS, '-cp', JAR, '-sourcepath', str(directory/'empty-sourcepath'),
                        '-d', str(directory/'classes'), str(directory/'relocated-sources'/relative)]
    require(compile_result['argv'] == expected_compile, directory_name)
    require(compile_result['status'] == 'PASS' and compile_result['exit'] == 0, directory_name)
    require(compile_result['started_unix'] >= receipt['started_unix'], directory_name)
    for stream in ('stdout', 'stderr'):
        require(sha(directory/'compile'/stream) == compile_result[stream+'_sha256'], directory_name)
        require(data(directory/'compile'/stream) == b'', directory_name)
    compile_command = obj(directory/'compile/command.json')
    require(all(compile_result[k] == v for k,v in compile_command.items()), directory_name)
    require([r['variant'] for r in receipt['runs']] == ['STOCK', 'CONTROL', candidate], directory_name)
    expected_origins = [family.replace('/', '.') for family in build['compiler_source_families']]
    require(receipt['origin_classes'] == expected_origins, directory_name)
    previous_end = compile_result['ended_unix']
    outputs = []
    reviews = []
    for run in receipt['runs']:
        variant = run['variant']
        run_directory = Path(run['cwd'])
        require(run_directory.parent == directory, str(run_directory))
        require(run == obj(run_directory/'result.json'), str(run_directory))
        command = obj(run_directory/'command.json')
        require(all(run[k] == v for k,v in command.items()), str(run_directory))
        require(run['status'] == 'PASS' and run['exit'] == 0 and run['probe'] == classname, variant)
        require(previous_end <= run['started_unix'] <= run['ended_unix'], variant)
        previous_end = run['ended_unix']
        require(run['environment_variables_removed'] == ['JAVA_TOOL_OPTIONS', '_JAVA_OPTIONS', 'JDK_JAVA_OPTIONS'], variant)
        overlay = build['variants'][variant]['overlay']
        classpath = ':'.join([str(directory/'classes')] + ([overlay] if overlay else []) + [JAR])
        args = expected_args if variant == candidate else []
        expected = [receipt['java'], '-Xmx1g', '-Xlog:class+load=info:file='+str(run_directory/'class-load.log'),
                    '-cp', classpath, classname, *args]
        require(run['argv'] == expected, 'Probe argv: '+str(run_directory))
        for stream in ('stdout', 'stderr'):
            require(sha(run_directory/stream) == run[stream+'_sha256'], str(run_directory))
        require(data(run_directory/'stderr') == b'', str(run_directory))
        output = data(run_directory/'stdout')
        require(output.rstrip().endswith((classname.rsplit('.',1)[1]+' OK').encode()), str(run_directory))
        outputs.append(output)
        require(sha(run_directory/'class-load.log') == run['class_load_log_sha256'], str(run_directory))
        log = data(run_directory/'class-load.log').decode()
        expected_origin = str(Path(overlay or JAR).resolve())
        for cls in expected_origins:
            require(origins(log, cls) == {expected_origin}, 'Production origin '+cls)
            origin_receipt = run['class_origins'][cls]
            require(origin_receipt['pass'] is True and origin_receipt['normalized_origins'] == [expected_origin], cls)
        require(origins(log, classname) == {str((directory/'classes').resolve())}, classname)
        if variant != 'STOCK':
            require(run['stdout_equals_stock'] is True, variant)
        reviews.append({'variant':variant, 'exit_status':run['exit'], 'argv':run['argv'],
                        'stdout_sha256':run['stdout_sha256'], 'stdout_bytes':len(output),
                        'stderr_empty':True, 'production_origin':expected_origin,
                        'class_origin_log_independently_checked':True})
    require(previous_end <= receipt['ended_unix'], directory_name)
    require(outputs[0] == outputs[1] == outputs[2], 'Output mismatch: '+directory_name)
    return {'results':str(directory), 'receipt_sha256':sha(directory/'manifest.json'),
            'status':'PASS_INDEPENDENT_RECEIPT_REVIEW', 'source_sha256':sha(source),
            'probe':classname, 'candidate_args':expected_args,
            'stdout_byte_identical':True, 'runs':reviews}


def main():
    require(not (HERE/'review.json').exists(), 'Keep previous reviews; choose a fresh directory')
    fn_rel = 'com/google/javascript/jscomp/FunctionArgumentInjector.java'
    scope_rel = 'com/google/javascript/jscomp/AbstractScope.java'
    empty_root = BASE/'candidates/scope-maps'
    scope_root = BASE/'candidates/scope-declare-once'
    empty_source = empty_root/'empty-parameters-only/src'/fn_rel
    scope_source = scope_root/'src'/scope_rel
    frozen = obj(scope_root/'frozen-source-manifest.json')['files_sha256']
    for relative in ['src/'+scope_rel, 'control-src/'+scope_rel,
                     'probe/com/google/javascript/jscomp/ScopeDeclareOnceProbe.java']:
        require(sha(scope_root/relative) == frozen[relative], relative)
    audit = obj(BASE/'audits/empty-parameters-callsite-v1/source-review.json')
    for entry in audit['files']:
        require(sha(entry['path']) == entry['sha256'], entry['path'])
    empty_build, empty_build_review = build_review('empty-parameters', fn_rel,
        empty_root/'control-src'/fn_rel, empty_source)
    scope_build, scope_build_review = build_review('scope-declare-once', scope_rel,
        scope_root/'control-src'/scope_rel, scope_source)
    specs = [
        ('empty-parameters-modified-probes-attempt-1', 'empty-parameters',
         'com.google.javascript.jscomp.ModifiedParametersProbe',
         empty_root/'probe/com/google/javascript/jscomp/ModifiedParametersProbe.java', [], empty_build),
        ('empty-parameters-boundary-probes-attempt-1', 'empty-parameters',
         'com.google.javascript.jscomp.EmptyParametersProbe',
         empty_root/'empty-parameters-only/probe/com/google/javascript/jscomp/EmptyParametersProbe.java',
         ['--expect-shortcut'], empty_build),
        ('scope-declare-once-probes-attempt-1', 'scope-declare-once',
         'com.google.javascript.jscomp.ScopeDeclareOnceProbe',
         scope_root/'probe/com/google/javascript/jscomp/ScopeDeclareOnceProbe.java', [], scope_build),
    ]
    reviews = [probe_review(*spec) for spec in specs]
    for path, digest in list(INPUTS.items()):
        require(sha(path) == digest, 'Input changed: '+path)
    result = {'status':'PASS_ALL_THREE_INDEPENDENT_RECEIPT_REVIEWS',
              'reviewed_at_utc':datetime.now(timezone.utc).isoformat(),
              'scope':'Source/build/probe receipt integrity, exact argv, recorded success, byte-identical stdout and independently parsed class origins; no new runtime or performance review.',
              'jar_hash_scope':'Recorded hashes cross-checked only; no archive opened or rehashed.',
              'performance_status':'UNMEASURED_AT_TASK_ASSIGNMENT_NO_ACTIVE_SCREEN_READ',
              'builds':[empty_build_review, scope_build_review], 'probe_reviews':reviews,
              'input_sha256':INPUTS, 'review_script_sha256':hashlib.sha256(Path(__file__).read_bytes()).hexdigest()}
    (HERE/'review.json').write_text(json.dumps(result,indent=2)+'\n')
    for review in reviews:
        print(review['probe']+': PASS STOCK/CONTROL/candidate; exact stdout, sources/builds and class origins verified')


if __name__ == '__main__':
    main()

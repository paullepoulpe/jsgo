`JfrPublicationAudit.java` is closed source for root to compile/run after performance tests. It has not been compiled or executed by this agent. It uses only public JDK 21 JFR consumer APIs, and does not change recordings or export observed values/value hashes.

Exact compilation and one-recording audit commands (root only, when quiet):

```sh
mkdir -p /tmp/jfr-publication-audit-classes
/usr/lib/jvm/java-21-openjdk-amd64/bin/java -m jdk.compiler/com.sun.tools.javac.Main --add-modules jdk.jfr -source 21 -target 21 -d /tmp/jfr-publication-audit-classes /workspace/jsgo-java-experiments/profile-v2/publication-preparation/JfrPublicationAudit.java
/usr/lib/jvm/java-21-openjdk-amd64/bin/java -Xmx256m --add-modules jdk.jfr -cp /tmp/jfr-publication-audit-classes JfrPublicationAudit /workspace/jsgo-java-experiments/profile-v2/results/profile-stock-bigquery-steady-v2/work/profile.jfr
```

Root may redirect stdout to a new audit receipt. Repeat independently for `results/profile-stock-bigquery-v1/work/profile.jfr`; do not infer its contents from the corrected capture. Exit 0 means the standard consumer string audit completed with no unresolved fields; exit 2 means unresolved or incomplete; exit 3 means the receipt could not be written. Exit 2 is a publication disposition, not a benchmark failure. Preserve the receipt and raw recording locally.

The optional second argument is a **private exact-value allowlist outside the experiment/publication trees**. Keep it private (for example, a mode-0600 file under `/tmp/jfr-publication-private/`); do not publish it or its contents. Each noncomment line contains a selector, one literal tab, and an exact allowed value. Repeated selectors allow several explicitly reviewed values. Supported escapes in values are `\\`, `\n`, `\r`, `\t`; other escapes fail closed. Example benign entries:

```text
PROP:os.name<TAB>Linux
PROP:os.arch<TAB>amd64
PROP:java.home<TAB>/usr/lib/jvm/java-21-openjdk-amd64
```

Replace `<TAB>` with an actual tab. `ENV:key` checks a particular environment value, `PROP:key` a particular system-property value, and `FIELD:ownerType.fieldName` a particular other recorded string. Only `ENV:LC_ALL=C.UTF-8` is built in from the explicit container construction. Empty strings are safe. Every nonempty environment/property value, command line, path, thread name, thread dump or other unrecognized string needs an exact allowlist match. Generate expected entries from already-reviewed source/command configuration, not by automatically admitting recording values. The utility never prints unknown values to help construct a permissive list.

The receipt inventories event counts and every string/nested field schema exposed by `readEventTypes`, including nested class/method/stack metadata; fields not observed have zero occurrence counts. It emits only fixed benign environment/property key names, never unknown keys, values, hashes, exception messages or input paths. It groups repeated strings by schema field and reports exact-match, known diagnostic metadata and unresolved counts. Standard nested class/method/package/module names and selected JDK enum descriptions are classified as diagnostic metadata under the pinned-JDK/compiler provenance assumption; arbitrary classloader names are not. JDK schema labels/annotations are treated as producer metadata, not user environment payloads. This is a standard consumer audit, not a raw-binary forensic claim.

Bounds: recording <=512 MiB, policy <=16 MiB/100,000 entries, 5 million events, 100 million visited value nodes, depth 64, 20,000 schema fields, 1 MiB per matched string and a 120-second scan budget. An unsupported representation, exceeded bound or read error produces INCOMPLETE without its exception message. Known diagnostic strings must also be printable ASCII and <=4 KiB; larger/unusual strings remain unresolved. The external `-Xmx256m` bounds heap use.

Publication decision: try this bounded audit once when quiet. If unresolved/incomplete fields remain, root's stated fallback is to publish all source/results/aggregate summaries and explicitly exclude raw `.jfr` files, retaining originals locally. Do not prolong measurements or introduce repeated audit setup. Record the exclusion reason before any broad archive inventory; do not claim the public archive contains every raw recording.

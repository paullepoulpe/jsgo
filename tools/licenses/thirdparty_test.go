package main

import (
	"strings"
	"testing"
)

const guavaHeader = `/*
 * Copyright (C) 2010 The Guava Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 */
`

const mitHeader = `/*
 * Copyright (c) 2013 Kohsuke Kawaguchi and other contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software ...
 */
`

func TestCheckNotices(t *testing.T) {
	for _, tc := range []struct {
		name, src string
		problem   string // substring of the only problem; "" for none
	}{
		{"closure port", "// Port of com/google/javascript/jscomp/Foo.java @ closure-compiler v20261015\n\npackage x\n", ""},
		{"closure test port", "// Port of test/com/google/debugging/sourcemap/FooTest.java @ closure-compiler v20261015\n\npackage x\n", ""},
		{"no marker", "package x\n\n// Go-only: language-adapter — Guava's Ascii.\n", ""},
		{"guava port", "// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)\n\n" + guavaHeader + "\npackage x\n", ""},
		{"guava port, no notice", "// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)\n\npackage x\n", "lacks Guava's copyright line and Guava's licence"},
		{"guava notice after package", "// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)\n\npackage x\n\n" + guavaHeader, "lacks Guava's"},
		{"marker after package", "package x\n\n// Port of com/google/gson/stream/JsonWriter.java @ gson 2.9.1\n", "translates Gson code"},
		{"derived, short block", "// Derived from com/google/common/cache/LocalCache.java (Guava 33.4.6-jre):\n// a part.\n//\n// Copyright (C) 2009 The Guava Authors\n// Licensed under the Apache License, Version 2.0; see LICENSES/Apache-2.0.txt.\n\npackage x\n", ""},
		{"gson notice for a guava port", "// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)\n\n/*\n * Copyright (C) 2010 Google Inc.\n *\n * Licensed under the Apache License, Version 2.0 (the \"License\");\n */\n\npackage x\n", "lacks Guava's copyright line"},
		{"protobuf with pointer", "// Derived from com/google/protobuf/TextFormatEscaper.java (protobuf-java 4.30.2)\n//\n// Copyright 2008 Google Inc.  All rights reserved.\n//\n// Use of this source code is governed by a BSD-style\n// license that can be found in the LICENSE file.\n// Here that LICENSE file is LICENSES/BSD-3-Clause-protobuf.txt.\n\npackage x\n", ""},
		{"protobuf without pointer", "// Port of com/google/protobuf/TextFormat.java @ protobuf-java 4.30.2\n\n// Copyright 2008 Google Inc.  All rights reserved.\n//\n// Use of this source code is governed by a BSD-style\n// license that can be found in the LICENSE file.\n\npackage x\n", "a pointer to LICENSES/BSD-3-Clause-protobuf.txt"},
		{"args4j port", "// Port of org/kohsuke/args4j/CmdLineParser.java @ args4j 2.33\n\n" + mitHeader + "\npackage x\n", ""},
		{"args4j port, no notice", "// Port of org/kohsuke/args4j/CmdLineParser.java @ args4j 2.33\n\npackage x\n", "lacks args4j's"},
		{"unknown upstream", "// Port of java/util/HashMap.java\n\npackage x\n", "not a Closure source nor one of a known third-party project"},
		{"free-form marker in a body", "// Port of com/google/javascript/jscomp/Foo.java @ closure-compiler v20261015\n\npackage x\n\n// Port of TypeInference.java#traverseAdd.\nfunc f() {}\n", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			problems, _ := checkNotices("x.go", []byte(tc.src))
			switch {
			case tc.problem == "" && len(problems) > 0:
				t.Errorf("unexpected problems: %q", problems)
			case tc.problem != "" && (len(problems) != 1 || !strings.Contains(problems[0], tc.problem)):
				t.Errorf("problems = %q, want one containing %q", problems, tc.problem)
			}
		})
	}
}

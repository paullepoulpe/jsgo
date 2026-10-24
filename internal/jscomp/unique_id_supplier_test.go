// Port of test/com/google/javascript/jscomp/UniqueIdSupplierTest.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func uniqueIdSupplierTest_hashString(name string) string {
	fileHashCode := jutil.StringHashCode(name)
	if fileHashCode < 0 {
		return "m" + strconv.Itoa(int(-fileHashCode))
	}
	return strconv.Itoa(int(fileHashCode))
}

func TestUniqueIdSupplier_SingleCompilerInputGeneratesUniqueIds(t *testing.T) {
	uniqueIdSupplier := NewUniqueIdSupplier()
	input := NewCompilerInput(SourceFile_FromCode("tmp", "function foo() {}"))
	inputHashString := uniqueIdSupplierTest_hashString(input.GetSourceFile().GetName())

	uniqueId1 := uniqueIdSupplier.GetUniqueId(input)
	uniqueId2 := uniqueIdSupplier.GetUniqueId(input)
	if uniqueId1 == "" || uniqueId2 == "" || uniqueId1 == uniqueId2 {
		t.Errorf("ids %q %q", uniqueId1, uniqueId2)
	}
	if !strings.Contains(uniqueId1, inputHashString+"$0") || !strings.Contains(uniqueId2, inputHashString+"$1") {
		t.Errorf("ids %q %q, hash %s", uniqueId1, uniqueId2, inputHashString)
	}
}

func TestUniqueIdSupplier_FilePathsHavingSameHashCode_generatesUniqueIds(t *testing.T) {
	uniqueIdSupplier := NewUniqueIdSupplier()
	// Strings "FB" and "Ea" generate the same hashcode.
	input1 := NewCompilerInput(SourceFile_FromCode("FB", "function foo() {}"))
	input2 := NewCompilerInput(SourceFile_FromCode("Ea", "function foo() {}"))

	uniqueId1 := uniqueIdSupplier.GetUniqueId(input1)
	uniqueId2 := uniqueIdSupplier.GetUniqueId(input2)
	if uniqueId1 == "" || uniqueId2 == "" || uniqueId1 == uniqueId2 {
		t.Errorf("ids %q %q", uniqueId1, uniqueId2)
	}

	inputHashString1 := uniqueIdSupplierTest_hashString(input1.GetSourceFile().GetName())
	inputHashString2 := uniqueIdSupplierTest_hashString(input2.GetSourceFile().GetName())
	// The hash strings for the two files are the same, still the supplier generates unique IDs for
	// them
	assertEqual(t, "hash strings", inputHashString1, inputHashString2)
	if !strings.Contains(uniqueId1, inputHashString1+"$0") || !strings.Contains(uniqueId2, inputHashString2+"$1") {
		t.Errorf("ids %q %q", uniqueId1, uniqueId2)
	}
}

func TestUniqueIdSupplier_MultipleCompilerInputsGenerateUniqueIds(t *testing.T) {
	uniqueIdSupplier := NewUniqueIdSupplier()
	input1 := NewCompilerInput(SourceFile_FromCode("tmp1", "function foo() {}"))
	inputHashString1 := uniqueIdSupplierTest_hashString(input1.GetSourceFile().GetName())

	input2 := NewCompilerInput(SourceFile_FromCode("tmp2", "function foo() {}"))
	inputHashString2 := uniqueIdSupplierTest_hashString(input2.GetSourceFile().GetName())

	uniqueId1 := uniqueIdSupplier.GetUniqueId(input1)
	uniqueId2 := uniqueIdSupplier.GetUniqueId(input2)
	if uniqueId1 == "" || uniqueId2 == "" || uniqueId1 == uniqueId2 {
		t.Errorf("ids %q %q", uniqueId1, uniqueId2)
	}
	if !strings.Contains(uniqueId1, inputHashString1+"$0") || !strings.Contains(uniqueId2, inputHashString2+"$0") {
		t.Errorf("ids %q %q", uniqueId1, uniqueId2)
	}
}

// Ensures that unique Ids generated are the same across compiler runs
func TestUniqueIdSupplier_GeneratedIdsAreDeterministicAcrossRuns(t *testing.T) {
	compiler1 := NewCompiler()
	compiler2 := NewCompiler()

	input1 := NewCompilerInput(SourceFile_FromCode("tmp", "function foo() {}"))
	input2 := NewCompilerInput(SourceFile_FromCode("tmp", "function foo() {}"))

	uniqueId1 := compiler1.GetUniqueIdSupplier().GetUniqueId(input1)
	differentUniqueId1 := compiler1.GetUniqueIdSupplier().GetUniqueId(input1)

	uniqueId2 := compiler2.GetUniqueIdSupplier().GetUniqueId(input2)

	if uniqueId1 == "" || uniqueId2 == "" {
		t.Errorf("empty ids")
	}
	assertEqual(t, "same input, different compilers", uniqueId1, uniqueId2)
	if uniqueId1 == differentUniqueId1 {
		t.Errorf("same input, same compiler: %q twice", uniqueId1)
	}
}

// Java: UniqueIdSupplierTest.testProtoRoundTripPreservesInsertionOrder()
func TestUniqueIdSupplier_ProtoRoundTripPreservesInsertionOrder(t *testing.T) {
	uniqueIdSupplier := NewUniqueIdSupplier()
	input1 := NewCompilerInput(SourceFile_FromCode("fileA.js", ""))
	input2 := NewCompilerInput(SourceFile_FromCode("fileB.js", ""))
	input3 := NewCompilerInput(SourceFile_FromCode("fileC.js", ""))

	_ = uniqueIdSupplier.GetUniqueId(input1)
	_ = uniqueIdSupplier.GetUniqueId(input2)
	_ = uniqueIdSupplier.GetUniqueId(input3)
	_ = uniqueIdSupplier.GetUniqueId(input1) // increment counter for fileA

	protos := uniqueIdSupplier.ToProto()

	assertEqual(t, "protos size", len(protos), 3)
	assertEqual(t, "protos[0] hash", protos[0].GetHash(), jutil.StringHashCode("fileA.js"))
	assertEqual(t, "protos[0] counter", protos[0].GetCounter(), int32(2))
	assertEqual(t, "protos[1] hash", protos[1].GetHash(), jutil.StringHashCode("fileB.js"))
	assertEqual(t, "protos[1] counter", protos[1].GetCounter(), int32(1))
	assertEqual(t, "protos[2] hash", protos[2].GetHash(), jutil.StringHashCode("fileC.js"))
	assertEqual(t, "protos[2] counter", protos[2].GetCounter(), int32(1))

	restored := UniqueIdSupplier_FromProto(protos)
	restoredProtos := restored.ToProto()
	if len(restoredProtos) != len(protos) {
		t.Fatalf("restored.toProto(): got %d protos, want %d", len(restoredProtos), len(protos))
	}
	for i := range protos {
		if !restoredProtos[i].Equals(protos[i]) {
			t.Errorf("restored.toProto()[%d]: got %v, want %v", i, restoredProtos[i], protos[i])
		}
	}
}

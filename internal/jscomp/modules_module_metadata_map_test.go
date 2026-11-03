/*
 * Copyright 2018 The Closure Compiler Authors.
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

import "testing"

func TestModuleMetadataMap_NilReceiverSafe(t *testing.T) {
	var nilMap *ModuleMetadataMap
	if got := nilMap.ModuleByPath("foo.js"); got != nil {
		t.Fatalf("nilMap.ModuleByPath = %v, want nil", got)
	}
	if got := nilMap.ModuleByGoogNamespace("foo.bar"); got != nil {
		t.Fatalf("nilMap.ModuleByGoogNamespace = %v, want nil", got)
	}
	if got := nilMap.AllMetadataSlice(); got != nil {
		t.Fatalf("nilMap.AllMetadataSlice = %v, want nil", got)
	}
	if got := nilMap.GetModulesByPath(); got != nil {
		t.Fatalf("nilMap.GetModulesByPath = %v, want nil", got)
	}
	if got := nilMap.GetModulesByGoogNamespace(); got != nil {
		t.Fatalf("nilMap.GetModulesByGoogNamespace = %v, want nil", got)
	}
	if got := nilMap.GetAllModuleMetadata(); got.Size() != 0 {
		t.Fatalf("nilMap.GetAllModuleMetadata.Size() = %d, want 0", got.Size())
	}
}

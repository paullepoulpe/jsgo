/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

func nodeStringID(n rhino.Node) strpool.ID {
	return strpool.ID(n.GetStringID())
}

// Pre-interned well-known JavaScript and Closure identifier IDs shared across
// compiler passes.
var (
	idPrototype        = strpool.Intern("prototype")
	idArguments        = strpool.Intern("arguments")
	idThis             = strpool.Intern("this")
	idSuper            = strpool.Intern("super")
	idCall             = strpool.Intern("call")
	idApply            = strpool.Intern("apply")
	idEval             = strpool.Intern("eval")
	idConstructor      = strpool.Intern("constructor")
	idLength           = strpool.Intern("length")
	idToString         = strpool.Intern("toString")
	idValueOf          = strpool.Intern("valueOf")
	idDefineProperties = strpool.Intern("defineProperties")

	idUndefined = strpool.Intern("undefined")
	idNaN       = strpool.Intern("NaN")
	idInfinity  = strpool.Intern("Infinity")
	idObject    = strpool.Intern("Object")
	idArray     = strpool.Intern("Array")
	idString    = strpool.Intern("String")
	idNumber    = strpool.Intern("Number")
	idBoolean   = strpool.Intern("Boolean")
	idMath      = strpool.Intern("Math")
	idError     = strpool.Intern("Error")
	idRegExp    = strpool.Intern("RegExp")
	idSymbol    = strpool.Intern("Symbol")
	idWindow    = strpool.Intern("window")

	idGoog              = strpool.Intern("goog")
	idWeakUsage         = strpool.Intern("weakUsage")
	idGoogWeakUsage     = strpool.Intern("goog$weakUsage")
	idJscPropertyNameFn = strpool.Intern(NodeUtil_JSC_PROPERTY_NAME_FN)
)

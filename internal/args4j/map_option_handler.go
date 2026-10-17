// Port of org/kohsuke/args4j/spi/MapOptionHandler.java @ args4j 2.33

/*
 * Copyright (c) 2013 Kohsuke Kawaguchi and other contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */

package args4j

import "github.com/paullepoulpe/jsgo/internal/jutil"

// MapOptionHandler parses options into a Map.
//
//	class Foo {
//	 @Option(name="-P",handler=MapOptionHandler.class)
//	  Map<String,String> args;
//	}
//
// With this, -P x=1 -P y=2 parses to map of size 2. This option handler can
// be subtyped if you want to convert values to different types or to handle
// key=value in other formats, like key:=value.
//
// Java: class spi.MapOptionHandler extends OptionHandler<Map<?,?>>
type MapOptionHandler struct {
	OptionHandlerBase
}

// NewMapOptionHandler returns the handler as an OptionHandler, so that it is
// an OptionHandlerClass (Go-only: language-adapter).
//
// Java: MapOptionHandler(CmdLineParser, OptionDef, Setter)
func NewMapOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &MapOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	if setter.AsFieldSetter() == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("MapOptionHandler can only work with fields"))
	}
	return h
}

// Java: MapOptionHandler.getDefaultMetaVariable()
func (h *MapOptionHandler) GetDefaultMetaVariable() *string {
	return nil
}

// Java: MapOptionHandler.parseArguments(Parameters)
func (h *MapOptionHandler) ParseArguments(params Parameters) int {
	fs := h.Setter.AsFieldSetter()
	var v Map
	if o := fs.GetValue(); o != nil {
		v = o.(Map)
	}
	if v == nil {
		v = h.CreateNewCollection(fs.GetType())
		fs.AddValue(v)
	}

	h.AddToMap(params.GetParameter(0), v)

	return 1
}

// CreateNewCollection creates a new instance of the collection.
//
// Java: MapOptionHandler.createNewCollection(Class)
func (h *MapOptionHandler) CreateNewCollection(type_ *Class) Map {
	return jutil.NewJavaHashMap[string, *string]()
}

// AddToMap encapsulates how a single string argument gets converted into key
// and value.
//
// Java: MapOptionHandler.addToMap(String, Map)
func (h *MapOptionHandler) AddToMap(argument string, m Map) {
	if jutil.StringIndexOf(argument, "=") == -1 {
		panic(NewCmdLineExceptionWithLocalizable(h.Owner, SpiMessages_FORMAT_ERROR_FOR_MAP))
	}

	var mapKey string
	var mapValue *string

	//Splitting off the key from the value
	idx := jutil.StringIndexOf(argument, "=")
	if idx >= 0 {
		mapKey = jutil.StringSubstring(argument, 0, idx)
		value := jutil.StringSubstringFrom(argument, idx+1)
		mapValue = &value
		if jutil.StringLength(value) == 0 {
			// Kohsuke: I think this is a bad choice, but this is needed to
			// remain backward compatible
			mapValue = nil
		}
	} else {
		mapKey = argument
		mapValue = nil
	}

	if jutil.StringLength(mapKey) == 0 {
		panic(NewCmdLineExceptionWithLocalizable(h.Owner, SpiMessages_MAP_HAS_NO_KEY))
	}

	h.AddToMapWithKeyAndValue(m, mapKey, mapValue)
}

// AddToMapWithKeyAndValue is the opportunity to convert values to some typed
// objects.
//
// Java: MapOptionHandler.addToMap(Map, String, String)
func (h *MapOptionHandler) AddToMapWithKeyAndValue(m Map, key string, value *string) {
	m.Put(key, value)
}

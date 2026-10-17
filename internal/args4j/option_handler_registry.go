// Port of org/kohsuke/args4j/OptionHandlerRegistry.java @ args4j 2.33

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

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// OptionHandlerRegistry manages the registration of option handlers. This is
// good for registering custom handlers for specific parameter classes not
// yet implemented. The registry is a singleton that can be retrieved with
// the getRegistry() call.
//
// Java: class OptionHandlerRegistry
type OptionHandlerRegistry struct {
	// All OptionHandlers known to the CmdLineParser.
	//
	// Constructors of OptionHandler-derived class keyed by their supported
	// types.
	//
	// Go-only: Collections.synchronizedMap(new HashMap<Class, ...>()); the
	// map is only looked up, and its keys are identity-hashed classes.
	// handlersMonitor is the wrapper's monitor, which each of its methods
	// holds.
	handlers        *jutil.LinkedHashMap[*Class, OptionHandlerRegistry_OptionHandlerFactory]
	handlersMonitor jutil.Monitor
}

// The shared reference.
//
// See getRegistry().
var optionHandlerRegistry_instance *OptionHandlerRegistry

// Go-only: language-adapter — the monitor of OptionHandlerRegistry.class,
// which its static synchronized methods hold.
var optionHandlerRegistry_classMonitor jutil.Monitor

// OptionHandlerRegistry_GetRegistry gets the option handler registry
// singleton instance.
//
// Java: OptionHandlerRegistry.getRegistry()
func OptionHandlerRegistry_GetRegistry() *OptionHandlerRegistry {
	optionHandlerRegistry_classMonitor.Lock()
	defer optionHandlerRegistry_classMonitor.Unlock()
	if optionHandlerRegistry_instance == nil {
		optionHandlerRegistry_instance = newOptionHandlerRegistry()
	}
	return optionHandlerRegistry_instance
}

// newOptionHandlerRegistry constructs an option handler manager with the
// default handlers initialized.
//
// Java: OptionHandlerRegistry()
func newOptionHandlerRegistry() *OptionHandlerRegistry {
	r := &OptionHandlerRegistry{handlers: jutil.NewLinkedHashMap[*Class, OptionHandlerRegistry_OptionHandlerFactory]()}
	r.initHandlers()
	return r
}

// initHandlers registers the default handlers.
//
// Java: OptionHandlerRegistry.initHandlers()
func (r *OptionHandlerRegistry) initHandlers() {
	r.RegisterHandler(Class_Boolean, NewBooleanOptionHandler)
	r.RegisterHandler(Class_boolean, NewBooleanOptionHandler)
	// Go-only: the handlers that are not ported (see the package
	// documentation) are not registered: File, URL and URI here, ...
	r.RegisterHandler(Class_Integer, NewIntOptionHandler)
	r.RegisterHandler(Class_int, NewIntOptionHandler)
	// ... Double and double here, ...
	r.RegisterHandler(Class_String, NewStringOptionHandler)
	// ... Byte, byte, Character, char, Float, float, Long, long, Short,
	// short, InetAddress and Pattern here, ...
	// enum is a special case
	r.RegisterHandler(Class_Map, NewMapOptionHandler)

	// ... and java.nio.file.Path here.
}

// optionHandlerRegistry_getConstructor finds the constructor for an option
// handler.
//
// Go-only: language-adapter — reflection: an OptionHandlerClass is the
// constructor itself, so Java's NoSuchMethodException (reported as
// IllegalArgumentException(NO_CONSTRUCTOR_ON_HANDLER)) cannot happen.
//
// Java: OptionHandlerRegistry.getConstructor(Class)
func optionHandlerRegistry_getConstructor(handlerClass OptionHandlerClass) OptionHandlerClass {
	return handlerClass
}

// RegisterHandler registers a user-defined OptionHandler class with args4j.
//
// This method allows users to extend the behavior of args4j by writing their
// own OptionHandler implementation.
//
// valueType: the specified handler is used when the field/method annotated by
// Option is of this type. handlerClass: this class must have the constructor
// that has the same signature as OptionHandler(CmdLineParser, OptionDef,
// Setter). It throws NullPointerException if valueType or handlerClass is
// nil.
//
// Java: OptionHandlerRegistry.registerHandler(Class, Class)
func (r *OptionHandlerRegistry) RegisterHandler(valueType *Class, handlerClass OptionHandlerClass) {
	Utilities_CheckNonNull(valueType, "valueType")
	Utilities_CheckNonNull(handlerClass, "handlerClass")

	// Go-only: Java first checks that handlerClass is a subtype of
	// OptionHandler (IllegalArgumentException(NO_OPTIONHANDLER)); an
	// OptionHandlerClass always constructs one.

	r.putHandler(valueType, newOptionHandlerRegistry_DefaultConstructorHandlerFactory(handlerClass))
}

// RegisterHandlerWithFactory registers a user-defined OptionHandler class
// with args4j.
//
// This method allows users to extend the behavior of args4j by writing their
// own OptionHandler implementation.
//
// valueType: the specified handler is used when the field/method annotated by
// Option is of this type. factory: factory to instantiate handler upon
// request. It throws NullPointerException if valueType or factory is nil.
//
// Java: OptionHandlerRegistry.registerHandler(Class, OptionHandlerFactory)
func (r *OptionHandlerRegistry) RegisterHandlerWithFactory(valueType *Class, factory OptionHandlerRegistry_OptionHandlerFactory) {
	Utilities_CheckNonNull(valueType, "valueType")
	Utilities_CheckNonNull(factory, "factory")

	r.putHandler(valueType, factory)
}

func (r *OptionHandlerRegistry) putHandler(valueType *Class, factory OptionHandlerRegistry_OptionHandlerFactory) {
	r.handlersMonitor.Lock()
	defer r.handlersMonitor.Unlock()
	r.handlers.Put(valueType, factory)
}

// CreateOptionHandler creates an OptionHandler that handles the given Option
// annotation and the Setter instance.
//
// Java: OptionHandlerRegistry.createOptionHandler(CmdLineParser, OptionDef, Setter)
func (r *OptionHandlerRegistry) CreateOptionHandler(parser *CmdLineParser, o OptionDef, setter Setter) OptionHandler {
	Utilities_CheckNonNull(o, "CmdLineParser")
	Utilities_CheckNonNull(o, "OptionDef")
	Utilities_CheckNonNull(setter, "Setter")

	h := o.Handler()
	if h == nil {
		// infer the type
		t := setter.GetType()

		// enum is the special case
		if Class_Enum.IsAssignableFrom(t) {
			return NewEnumOptionHandler(parser, o, setter, t)
		}

		r.handlersMonitor.Lock()
		factory := r.handlers.Get(t)
		r.handlersMonitor.Unlock()
		if factory == nil {
			panic(NewIllegalAnnotationError(Messages_UNKNOWN_HANDLER.Format(t)))
		}

		return factory.GetHandler(parser, o, setter)
	} else {
		// explicit handler specified
		return newOptionHandlerRegistry_DefaultConstructorHandlerFactory(h).GetHandler(parser, o, setter)
	}
}

// OptionHandlerRegistry_OptionHandlerFactory provides custom logic for
// creating OptionHandler implementation.
//
// Java: interface OptionHandlerRegistry.OptionHandlerFactory
type OptionHandlerRegistry_OptionHandlerFactory interface {
	// GetHandler provides a handler instance to use.
	//
	// Java: OptionHandlerFactory.getHandler(CmdLineParser, OptionDef, Setter)
	GetHandler(parser *CmdLineParser, o OptionDef, setter Setter) OptionHandler
}

// Java: class OptionHandlerRegistry.DefaultConstructorHandlerFactory
type OptionHandlerRegistry_DefaultConstructorHandlerFactory struct {
	cons OptionHandlerClass
}

// Java: DefaultConstructorHandlerFactory(Class)
func newOptionHandlerRegistry_DefaultConstructorHandlerFactory(type_ OptionHandlerClass) *OptionHandlerRegistry_DefaultConstructorHandlerFactory {
	f := &OptionHandlerRegistry_DefaultConstructorHandlerFactory{}
	f.cons = optionHandlerRegistry_getConstructor(type_)
	return f
}

// Java: DefaultConstructorHandlerFactory.getHandler(CmdLineParser, OptionDef, Setter)
func (f *OptionHandlerRegistry_DefaultConstructorHandlerFactory) GetHandler(parser *CmdLineParser, o OptionDef, setter Setter) OptionHandler {
	// Go-only: language-adapter — reflection: calling the constructor cannot
	// fail with InstantiationException or IllegalAccessException. Java wraps
	// an exception thrown by the constructor in an InvocationTargetException
	// and that in an IllegalAnnotationError; here it propagates as thrown.
	return f.cons(parser, o, setter)
}

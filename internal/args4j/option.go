// Port of org/kohsuke/args4j/Option.java @ args4j 2.33

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

// Option marks a field/setter that receives a command line switch value.
//
// This annotation can be placed on a field of type T or the method of the
// form void methodName(T value). Its access modified can be anything, but if
// it's not public, your application needs to run in a security context that
// allows args4j to access the field/method (see
// AccessibleObject.setAccessible(boolean)).
//
// The behavior of the annotation differs depending on T --- the type of the
// field or the parameter of the method.
//
// Boolean Switch: When T is boolean, it represents a boolean option that
// takes the form of -OPT. When this option is set, the property will be set
// to true.
//
// String Switch: When T is String, it represents an option that takes one
// operand. The value of the operand is set to the property.
//
// Enum Switch: When T is derived from Enum, it represents an option that
// takes an operand, which must be one of the enum constant. The comparion
// between the operand and the enum constant name is done in a case
// insensitive fashion.
//
// For example, the following definition will represent command line options
// like -coin penny or -coin DIME, but things like -coin or -coin abc are
// errors.
//
//	enum Coin { PENNY,NICKEL,DIME,QUARTER }
//
//	class Option {
//	  @Option(name="-coin")
//	  public Coin coin;
//	}
//
// File Switch: When T is a File, it represents an option that takes a
// file/directory name as an operand.
//
// Go-only: language-adapter — annotations: the annotation is a struct whose
// fields are its elements; a field left at its zero value has the
// element's default (see the package documentation).
//
// Java: @interface Option
type Option struct {
	// Name of the option, such as -foo or -bar.
	Name string

	// Aliases for the options, such as --long-option-name.
	Aliases []string

	// Help string used to display the usage screen.
	//
	// This parameter works in two ways. For a simple use, you can just encode
	// the human-readable help string directly, and that will be used as the
	// message. This is easier, but it doesn't support localization.
	//
	// For more advanced use, this property is set to a key of a
	// ResourceBundle. The actual message is obtained by querying a
	// ResourceBundle instance supplied to CmdLineParser by this key. This
	// allows the usage screen to be properly localized.
	//
	// If this value is empty, the option will not be displayed in the usage
	// screen.
	Usage string

	// When the option takes an operand, the usage screen will show something
	// like this
	//
	//	-x FOO  : blah blah blah
	//
	// You can replace the FOO token by using this parameter.
	//
	// If left unspecified, this value is infered from the type of the option.
	//
	// Just like Usage, normally, this value is printed as is. But if a
	// ResourceBundle is given to the CmdLineParser, it will be used to obtain
	// the locale-specific value.
	MetaVar string

	// Specify that the option is mandatory.
	//
	// At the end of CmdLineParser.parseArgument(String...), a
	// CmdLineException will be thrown if a required option is not present.
	//
	// Note that in most of the command line interface design principles,
	// options should be really optional. So use caution when using this flag.
	Required bool

	// Specify that the option is a help option.
	//
	// When flagging an option being the help option, required arguments or
	// options that are missing in an actual command line don't cause an
	// exception to be thrown.
	Help bool

	// Specify that the option is hidden from the usage, by default.
	//
	// You can still have CmdLineParser show hidden options by using
	// OptionHandlerFilter.ALL, which allows you to create an option that shows
	// hidden options.
	//
	// If you need more complicated filtering logic, define your own
	// annotations and check them in Setter.asAnnotatedElement().
	Hidden bool

	// Specify the OptionHandler that processes the command line arguments.
	//
	// The default value OptionHandler (nil) indicates that the OptionHandler
	// will be infered from the type of the field/method where a Option
	// annotation is placed.
	//
	// If this annotation element is used, it overrides the inference and
	// determines the handler to be used. This is convenient for defining a
	// non-standard option parsing semantics.
	//
	//	// this is a normal "-r" option
	//	@Option(name="-r")
	//	boolean value;
	//
	//	// this causes arg4j to use MyHandler, not the default
	//	// handler provided for boolean
	//	@Option(name="-b",handler=MyHandler.class)
	//	boolean value;
	Handler OptionHandlerClass

	// List of other options that this option depends on.
	//
	//	@Option(name="-a")
	//	int a;
	//	//-b is not required but if it's provided, then a becomes required
	//	@Option(name="-b", depends={"-a"})
	//	int b;
	//
	// At the end of CmdLineParser.parseArgument(String...), a
	// CmdLineException will be thrown if options required by another one are
	// not present.
	Depends []string

	// List of other options that this option is incompatible with..
	//
	//	@Option(name="-a")
	//	int a;
	//	// -h and -a cannot be specified together
	//	@Option(name="-h", forbids={"-a"})
	//	boolean h;
	//
	// At the end of CmdLineParser.parseArgument(String...), a
	// CmdLineException will be thrown if forbidden option combinations are
	// present.
	Forbids []string
}

// OptionHandlerClass is the Go form of Class<? extends OptionHandler>, as an
// annotation's handler element and OptionHandlerRegistry use it: the
// handler's constructor with the signature of
// OptionHandler(CmdLineParser, OptionDef, Setter). nil is
// OptionHandler.class, the default.
//
// Go-only: language-adapter — reflection (args4j instantiates the class
// through that constructor).
type OptionHandlerClass func(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler

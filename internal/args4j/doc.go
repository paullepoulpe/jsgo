// Package args4j is a port of args4j 2.33, the command line parser of
// Closure Compiler's command-line tools (CommandLineRunner, LinterMain,
// MapDiagnosticToSuppression, ProductionInstrumentationReporter). The pinned
// closure-compiler.jar shades it as com.google.javascript.jscomp.jarjar.
// org.kohsuke.args4j. The Java packages org.kohsuke.args4j and
// org.kohsuke.args4j.spi import each other, so they form this one Go
// package; org.kohsuke.args4j.spi.Messages, whose simple name collides with
// org.kohsuke.args4j.Messages, is SpiMessages here.
//
// Only what Closure uses is ported: CmdLineParser (parsing, @-files,
// required/help/depends/forbids checks, printUsage, printExample,
// printSingleLineUsage), CmdLineException, ParserProperties, the Option and
// Argument metadata (OptionDef, NamedOptionDef, ClassParser,
// OptionHandlerFilter, ExampleMode), OptionHandlerRegistry, the message
// bundles, and from the spi package OptionHandler, Parameters, Setter,
// Getter, Setters, FieldSetter, MultiValueFieldSetter and the handlers
// OneArgumentOptionHandler, BooleanOptionHandler, StringOptionHandler,
// IntOptionHandler, EnumOptionHandler, MapOptionHandler and
// StringArrayOptionHandler. Not ported, because Closure uses none of them:
// annotated methods (MethodSetter), array fields (ArrayFieldSetter), the XML
// configuration (Config, ConfigElement, XmlParser, AnnotationImpl,
// OptionImpl, ArgumentImpl), Starter, FieldParser, sub-commands, and the
// handlers for File, URL, URI, Double, Byte, Character, Float, Long, Short,
// InetAddress, Pattern, Path, UUID, MAC addresses, explicit booleans, the
// rest of the arguments, StopOptionHandler and the delimited/multi-file
// handlers. OptionHandlerRegistry registers the ported handlers only.
//
// # Go-only: language-adapter — annotations
//
// args4j discovers options by reflection: ClassParser walks the fields of
// the bean's class and its superclasses and reads their @Option and
// @Argument annotations, and Setters creates a setter over each such
// java.lang.reflect.Field. Go has neither annotations nor field reflection of
// that kind, so the metadata is registered explicitly:
//
//   - Option and Argument are structs whose fields are the annotation
//     elements, with the same names. A field left at its zero value takes the
//     element's default (Handler nil is the default OptionHandler.class:
//     infer the handler from the setter's type). An OptionHandlerClass is a
//     handler's (CmdLineParser, OptionDef, Setter) constructor, which is all
//     args4j uses of a handler class.
//   - An option bean implements Bean. Its Args4jFields method lists the
//     annotated fields in the order ClassParser.parse visits them: the
//     bean's own class first, in field declaration order, then each
//     superclass. That order is the order in which options are registered;
//     options are then sorted for display (ParserProperties' default
//     comparator), arguments keep their index.
//   - Field stands for java.lang.reflect.Field: the field's declared type
//     (and the type argument of a List), and accessors bound to the bean.
//     Class stands for java.lang.Class, as far as args4j inspects it.
//     Constructors such as BooleanField, StringField, EnumField, ListField
//     and StringMapField build a Field over a Go variable, boxing its values
//     as args4j sees them: a Java boolean/Boolean is a bool, an int/Integer
//     an int32, a String a string, null nil.
//
// For example, Java's
//
//	@Option(name = "--js", handler = JsOptionHandler.class, usage = "...")
//	private List<String> js = new ArrayList<>();
//	@Option(name = "--debug", usage = "...")
//	private boolean debug = false;
//	@Argument private List<String> arguments = new ArrayList<>();
//
// becomes
//
//	func (f *Flags) Args4jFields() []args4j.AnnotatedField {
//		const c = "com.google.javascript.jscomp.CommandLineRunner$Flags"
//		return []args4j.AnnotatedField{
//			{Option: &args4j.Option{Name: "--js", Handler: NewJsOptionHandler, Usage: "..."},
//				Field: args4j.ListField(c, "js", args4j.Class_String, &f.js)},
//			{Option: &args4j.Option{Name: "--debug", Usage: "..."},
//				Field: args4j.BooleanField(c, "debug", &f.debug)},
//			{Argument: &args4j.Argument{},
//				Field: args4j.ListField(c, "arguments", args4j.Class_String, &f.arguments)},
//		}
//	}
//
// and NewCmdLineParser(f) then registers and sorts them as args4j does.
//
// # Go-only: language-adapter — generics, locales, resource bundles
//
// Setter<T>, Getter<T> and OptionHandler<T> are not generic: values are
// passed as any, boxed as above. The message bundles are the jar's
// Messages.properties files, resolved as for an English default locale (the
// _en bundles are empty and fall back to them; the _de and _ru bundles are
// not ported), so Messages.format and Messages.formatWithLocale(ENGLISH)
// agree. A ResourceBundle given to printUsage or printExample is any value
// with GetString.
package args4j

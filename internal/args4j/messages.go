// Port of org/kohsuke/args4j/Messages.java @ args4j 2.33

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
	_ "embed"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Messages are the messages of org.kohsuke.args4j.
//
// Java: enum Messages implements Localizable
type Messages int32

const (
	Messages_MISSING_OPERAND Messages = iota
	Messages_UNDEFINED_OPTION
	Messages_NO_ARGUMENT_ALLOWED
	Messages_REQUIRED_OPTION_MISSING
	Messages_TOO_MANY_ARGUMENTS
	Messages_REQUIRED_ARGUMENT_MISSING
	Messages_METADATA_ERROR
	Messages_MULTIPLE_USE_OF_ARGUMENT
	Messages_MULTIPLE_USE_OF_OPTION
	Messages_UNKNOWN_HANDLER
	Messages_NO_OPTIONHANDLER
	Messages_NO_CONSTRUCTOR_ON_HANDLER
	Messages_REQUIRES_OPTION_MISSING
	Messages_FORBIDDEN_OPTION_PRESENT
	Messages_NO_SUCH_FILE
	Messages_DEFAULT_VALUE
)

var messages_names = [...]string{
	"MISSING_OPERAND",
	"UNDEFINED_OPTION",
	"NO_ARGUMENT_ALLOWED",
	"REQUIRED_OPTION_MISSING",
	"TOO_MANY_ARGUMENTS",
	"REQUIRED_ARGUMENT_MISSING",
	"METADATA_ERROR",
	"MULTIPLE_USE_OF_ARGUMENT",
	"MULTIPLE_USE_OF_OPTION",
	"UNKNOWN_HANDLER",
	"NO_OPTIONHANDLER",
	"NO_CONSTRUCTOR_ON_HANDLER",
	"REQUIRES_OPTION_MISSING",
	"FORBIDDEN_OPTION_PRESENT",
	"NO_SUCH_FILE",
	"DEFAULT_VALUE",
}

// Java: Messages.values()
func Messages_values() []Messages {
	v := make([]Messages, len(messages_names))
	for i := range v {
		v[i] = Messages(i)
	}
	return v
}

// String is Messages.name().
func (m Messages) String() string { return messages_names[m] }

// Go-only: language-adapter — classpath resource via go:embed: the bundle
// org/kohsuke/args4j/Messages.properties (see locale.go).
//
//go:embed messages.properties
var messages_properties string

var messages_bundle = &propertyBundle{text: messages_properties}

// Java: Messages.formatWithLocale(Locale, Object...)
func (m Messages) FormatWithLocale(locale *Locale, args ...any) string {
	localized := messages_bundle
	return jutil.MessageFormat_Format(localized.GetString(m.String()), args)
}

// Java: Messages.format(Object...)
func (m Messages) Format(args ...any) string {
	return m.FormatWithLocale(Locale_GetDefault(), args...)
}

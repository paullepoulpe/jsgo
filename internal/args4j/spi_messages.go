// Port of org/kohsuke/args4j/spi/Messages.java @ args4j 2.33

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

// SpiMessages are the messages of org.kohsuke.args4j.spi.
//
// Go-only: naming — org.kohsuke.args4j.spi.Messages is SpiMessages, since
// org.kohsuke.args4j.Messages is in the same Go package.
//
// Java: enum spi.Messages implements Localizable
type SpiMessages int32

const (
	SpiMessages_ILLEGAL_OPERAND SpiMessages = iota
	SpiMessages_ILLEGAL_CHAR
	SpiMessages_ILLEGAL_BOOLEAN
	SpiMessages_ILLEGAL_METHOD_SIGNATURE
	SpiMessages_ILLEGAL_FIELD_SIGNATURE
	SpiMessages_ILLEGAL_LIST
	SpiMessages_FORMAT_ERROR_FOR_MAP
	SpiMessages_MAP_HAS_NO_KEY
	SpiMessages_ILLEGAL_IP_ADDRESS
	SpiMessages_ILLEGAL_PATTERN
	SpiMessages_ILLEGAL_MAC_ADDRESS
	SpiMessages_ILLEGAL_UUID
	SpiMessages_ILLEGAL_PATH
	SpiMessages_DEFAULT_META_EXPLICIT_BOOLEAN_OPTION_HANDLER
	SpiMessages_DEFAULT_META_FILE_OPTION_HANDLER
	SpiMessages_DEFAULT_META_INET_ADDRESS_OPTION_HANDLER
	SpiMessages_DEFAULT_META_MAC_ADDRESS_OPTION_HANDLER
	SpiMessages_DEFAULT_META_PATH_OPTION_HANDLER
	SpiMessages_DEFAULT_META_PATTERN_OPTION_HANDLER
	SpiMessages_DEFAULT_META_REST_OF_ARGUMENTS_HANDLER
	SpiMessages_DEFAULT_META_STRING_ARRAY_OPTION_HANDLER
	SpiMessages_DEFAULT_META_STRING_OPTION_HANDLER
	SpiMessages_DEFAULT_META_SUB_COMMAND_HANDLER
	SpiMessages_DEFAULT_META_URI_OPTION_HANDLER
	SpiMessages_DEFAULT_META_URL_OPTION_HANDLER
	SpiMessages_DEFAULT_META_UUID_OPTION_HANDLER
)

var spiMessages_names = [...]string{
	"ILLEGAL_OPERAND",
	"ILLEGAL_CHAR",
	"ILLEGAL_BOOLEAN",
	"ILLEGAL_METHOD_SIGNATURE",
	"ILLEGAL_FIELD_SIGNATURE",
	"ILLEGAL_LIST",
	"FORMAT_ERROR_FOR_MAP",
	"MAP_HAS_NO_KEY",
	"ILLEGAL_IP_ADDRESS",
	"ILLEGAL_PATTERN",
	"ILLEGAL_MAC_ADDRESS",
	"ILLEGAL_UUID",
	"ILLEGAL_PATH",
	"DEFAULT_META_EXPLICIT_BOOLEAN_OPTION_HANDLER",
	"DEFAULT_META_FILE_OPTION_HANDLER",
	"DEFAULT_META_INET_ADDRESS_OPTION_HANDLER",
	"DEFAULT_META_MAC_ADDRESS_OPTION_HANDLER",
	"DEFAULT_META_PATH_OPTION_HANDLER",
	"DEFAULT_META_PATTERN_OPTION_HANDLER",
	"DEFAULT_META_REST_OF_ARGUMENTS_HANDLER",
	"DEFAULT_META_STRING_ARRAY_OPTION_HANDLER",
	"DEFAULT_META_STRING_OPTION_HANDLER",
	"DEFAULT_META_SUB_COMMAND_HANDLER",
	"DEFAULT_META_URI_OPTION_HANDLER",
	"DEFAULT_META_URL_OPTION_HANDLER",
	"DEFAULT_META_UUID_OPTION_HANDLER",
}

// Java: spi.Messages.values()
func SpiMessages_values() []SpiMessages {
	v := make([]SpiMessages, len(spiMessages_names))
	for i := range v {
		v[i] = SpiMessages(i)
	}
	return v
}

// String is spi.Messages.name().
func (m SpiMessages) String() string { return spiMessages_names[m] }

// Go-only: language-adapter — classpath resource via go:embed: the bundle
// org/kohsuke/args4j/spi/Messages.properties (see locale.go).
//
//go:embed spi_messages.properties
var spiMessages_properties string

var spiMessages_bundle = &propertyBundle{text: spiMessages_properties}

// Java: spi.Messages.formatWithLocale(Locale, Object...)
func (m SpiMessages) FormatWithLocale(locale *Locale, args ...any) string {
	localized := spiMessages_bundle
	return jutil.MessageFormat_Format(localized.GetString(m.String()), args)
}

// Java: spi.Messages.format(Object...)
func (m SpiMessages) Format(args ...any) string {
	return m.FormatWithLocale(Locale_GetDefault(), args...)
}

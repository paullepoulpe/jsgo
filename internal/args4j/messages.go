// Port of org/kohsuke/args4j/Localizable.java @ args4j 2.33
// Port of org/kohsuke/args4j/Messages.java @ args4j 2.33
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ResourceBundle localizes option usage strings and meta variables when passed
// to CmdLineParser usage printers.
type ResourceBundle interface {
	GetString(key string) string
}

// messageTemplate is an English args4j MessageFormat pattern from
// org/kohsuke/args4j/Messages.properties or spi/Messages.properties.
type messageTemplate string

func (m messageTemplate) format(args ...any) string {
	s := string(m)
	var sb strings.Builder
	sb.Grow(len(s) + 16)
	quoted := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			if i+1 < len(s) && s[i+1] == '\'' {
				sb.WriteByte('\'')
				i++
			} else {
				quoted = !quoted
			}
			continue
		}
		if c == '{' && !quoted {
			if end := strings.IndexByte(s[i+1:], '}'); end >= 0 {
				if idx, err := strconv.Atoi(s[i+1 : i+1+end]); err == nil && idx >= 0 {
					if idx < len(args) {
						sb.WriteString(jutil.ToJavaString(args[idx]))
					} else {
						sb.WriteByte('{')
						sb.WriteString(strconv.Itoa(idx))
						sb.WriteByte('}')
					}
					i += 1 + end
					continue
				}
			}
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

const (
	msgMissingOperand          messageTemplate = `Option "{0}" takes an operand`
	msgUndefinedOption         messageTemplate = `"{0}" is not a valid option`
	msgNoArgumentAllowed       messageTemplate = `No argument is allowed: {0}`
	msgRequiredOptionMissing   messageTemplate = `Option "{0}" is required`
	msgRequiredArgumentMissing messageTemplate = `Argument "{0}" is required`
	msgTooManyArguments        messageTemplate = `Too many arguments: {0}`
	msgMultipleUseOfArgument   messageTemplate = `Argument with index {0} is used more than once`
	msgMultipleUseOfOption     messageTemplate = `Option name {0} is used more than once`
	msgUnknownHandler          messageTemplate = `No OptionHandler is registered to handle {0}`
	msgRequiresOptionMissing   messageTemplate = `option "{0}" requires the option(s) {1}`
	msgForbiddenOptionPresent  messageTemplate = `option "{0}" cannot be used with the option(s) {1}`
	msgNoSuchFile              messageTemplate = `No such file: {0}`
	msgDefaultValue            messageTemplate = `(default: {0})`

	spiMsgIllegalOperand        messageTemplate = `"{1}" is not a valid value for "{0}"`
	spiMsgIllegalBoolean        messageTemplate = `"{0}" is not a legal boolean value`
	spiMsgIllegalFieldSignature messageTemplate = `Field of type {0} isn't supported by args4j`
	spiMsgIllegalList           messageTemplate = `type of {0} is not a List`
	spiMsgFormatErrorForMap     messageTemplate = `An argument for setting a Map must contain a "="`
	spiMsgMapHasNoKey           messageTemplate = `A key must be set.`

	defaultMetaStringArray = "STRING[]"
	defaultMetaString      = "VAL"
)

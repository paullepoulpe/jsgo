// Port of org/kohsuke/args4j/NamedOptionDef.java @ args4j 2.33

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

// NamedOptionDef is the immutable run-time copy of the Option annotation.
//
// Java: final class NamedOptionDef extends OptionDef
type NamedOptionDef struct {
	OptionDefBase
	name    string
	aliases []string
	depends []string
	forbids []string
}

// Deprecated: multi-valuedness as option definition does not make sense.
// It's driven by the setter.
//
// Java: NamedOptionDef(Option, boolean)
func NewNamedOptionDefWithForceMultiValued(o *Option, forceMultiValued bool) *NamedOptionDef {
	return NewNamedOptionDef(o)
}

// Java: NamedOptionDef(Option)
func NewNamedOptionDef(o *Option) *NamedOptionDef {
	d := &NamedOptionDef{}
	d.InitOptionDefBase(d, o.Usage, o.MetaVar, o.Required, o.Help, o.Hidden, o.Handler, false)

	d.name = o.Name
	d.aliases = namedOptionDef_createZeroSizedArrayIfNull(o.Aliases)
	d.depends = namedOptionDef_createZeroSizedArrayIfNull(o.Depends)
	d.forbids = namedOptionDef_createZeroSizedArrayIfNull(o.Forbids)
	return d
}

// Java: NamedOptionDef.createZeroSizedArrayIfNull(String[])
func namedOptionDef_createZeroSizedArrayIfNull(in []string) []string {
	if in == nil {
		return []string{}
	} else {
		return in
	}
}

// Java: NamedOptionDef.name()
func (d *NamedOptionDef) Name() string {
	return d.name
}

// Java: NamedOptionDef.aliases()
func (d *NamedOptionDef) Aliases() []string {
	return append([]string{}, d.aliases...)
}

// Java: NamedOptionDef.depends()
func (d *NamedOptionDef) Depends() []string {
	return append([]string{}, d.depends...)
}

// Java: NamedOptionDef.forbids()
func (d *NamedOptionDef) Forbids() []string {
	return append([]string{}, d.forbids...)
}

// Java: NamedOptionDef.toString()
func (d *NamedOptionDef) String() string {
	if len(d.aliases) > 0 {
		str := ""
		for _, alias := range d.aliases {
			if len(str) > 0 {
				str += ", "
			}
			str += alias
		}
		return d.Name() + " (" + str + ")"
	}
	return d.Name()
}

// Java: NamedOptionDef.isArgument()
func (d *NamedOptionDef) IsArgument() bool {
	return false
}

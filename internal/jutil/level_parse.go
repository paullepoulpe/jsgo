package jutil

// levelKnown are the levels the Level class defines, which Level.parse
// knows by name and by value.
var levelKnown = []*Level{Level_OFF, Level_SEVERE, Level_WARNING, Level_INFO, Level_CONFIG, Level_FINE, Level_FINER, Level_FINEST, Level_ALL}

// Level_Parse is Level.parse(String): a known level name (case-sensitive),
// or an integer value, which yields the known level of that value or a new
// Level named by the string. Anything else throws
// IllegalArgumentException("Bad level \"<name>\"").
func Level_Parse(name string) *Level {
	for _, l := range levelKnown {
		if l.name == name {
			return l
		}
	}
	var value int32
	parsed := false
	Try(func() {
		value = IntegerParseInt(name, 10)
		parsed = true
	}, func(e INumberFormatException) {})
	if parsed {
		for _, l := range levelKnown {
			if l.value == value {
				return l
			}
		}
		return &Level{name: name, value: value}
	}
	panic(NewIllegalArgumentExceptionMsg("Bad level \"" + name + "\""))
}

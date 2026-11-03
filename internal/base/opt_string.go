package base

// OptString is an unboxed optional string that distinguishes an absent value
// (Java null) from a present empty string (""). The zero value is absent.
type OptString struct {
	val string
	ok  bool
}

// SomeString returns a present OptString holding s.
func SomeString(s string) OptString {
	return OptString{val: s, ok: true}
}

// NoString returns an absent OptString.
func NoString() OptString {
	return OptString{}
}

// OptStringFromPtr converts a nil-able *string into an OptString.
func OptStringFromPtr(p *string) OptString {
	if p == nil {
		return OptString{}
	}
	return OptString{val: *p, ok: true}
}

// Get returns the string value and whether it is present.
func (o OptString) Get() (string, bool) {
	return o.val, o.ok
}

// ValueOrEmpty returns the string value if present, or "" if absent.
func (o OptString) ValueOrEmpty() string {
	if !o.ok {
		return ""
	}
	return o.val
}

// ValueOr returns the string value if present, or fallback if absent.
func (o OptString) ValueOr(fallback string) string {
	if !o.ok {
		return fallback
	}
	return o.val
}

// IsPresent reports whether a string value is present.
func (o OptString) IsPresent() bool {
	return o.ok
}

// Ptr returns a *string when present, or nil when absent.
func (o OptString) Ptr() *string {
	if !o.ok {
		return nil
	}
	s := o.val
	return &s
}

// String returns the string value if present, or "null" if absent (matching
// Java String.valueOf).
func (o OptString) String() string {
	if !o.ok {
		return "null"
	}
	return o.val
}

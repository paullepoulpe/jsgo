package ast

import "github.com/paullepoulpe/jsgo/internal/rhino"

// TokenSet is a set of rhino tokens.
type TokenSet [4]uint64

// Add adds t to s.
func (s *TokenSet) Add(t rhino.Token) { s[t>>6] |= 1 << (t & 63) }

// Has reports whether t is in s.
func (s *TokenSet) Has(t rhino.Token) bool { return s[t>>6]&(1<<(t&63)) != 0 }

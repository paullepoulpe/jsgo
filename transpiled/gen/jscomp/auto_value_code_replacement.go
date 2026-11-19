// Port of com/google/javascript/refactoring/AutoValue_CodeReplacement.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.refactoring.AutoValue_CodeReplacement
type AutoValue_CodeReplacement struct {
	CodeReplacementBase
	newContent    string `gson:"newContent"`
	sortKey       string `gson:"sortKey"`
	startPosition int32  `gson:"startPosition"`
	length        int32  `gson:"length"`
}

func (this_ *AutoValue_CodeReplacement) BindSelf_(s any) {
	_ = s
	this_.CodeReplacementBase.BindSelf_(s)
}
func (this_ *AutoValue_CodeReplacement) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_CodeReplacement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.AutoValue_CodeReplacement", "AutoValue_CodeReplacement"
}
func (this_ *AutoValue_CodeReplacement) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodeReplacementBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.AutoValue_CodeReplacement", "com.google.javascript.refactoring.AutoValue_CodeReplacement", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_CodeReplacement](v_)
	})
}
func (*AutoValue_CodeReplacement) Implements_jscomp_CodeReplacement() {}

func (this_ *AutoValue_CodeReplacement) initFields_AutoValue_CodeReplacement() {
	_ = this_
	this_.newContent = j2gort.NullString
	this_.sortKey = j2gort.NullString
}

// Java constructor: AutoValue_CodeReplacement(int startPosition, int length, java.lang.String newContent, java.lang.String sortKey)
func NewAutoValue_CodeReplacement(startPosition int32, length int32, newContent string, sortKey string) *AutoValue_CodeReplacement {
	this_ := &AutoValue_CodeReplacement{}
	this_.BindSelf_(this_)
	InitAutoValue_CodeReplacement(this_, startPosition, length, newContent, sortKey)
	return this_
}

func InitAutoValue_CodeReplacement(this_ *AutoValue_CodeReplacement, startPosition int32, length int32, newContent string, sortKey string) {
	_ = startPosition
	_ = length
	_ = newContent
	_ = sortKey
	InitCodeReplacement(&this_.CodeReplacementBase)
	this_.initFields_AutoValue_CodeReplacement()
	this_.startPosition = startPosition
	this_.length = length
	if newContent == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null newContent")))
	}
	this_.newContent = newContent
	if sortKey == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null sortKey")))
	}
	this_.sortKey = sortKey
}

// Java method: AutoValue_CodeReplacement.getStartPosition()
func (this_ *AutoValue_CodeReplacement) GetStartPosition() int32 {
	return this_.startPosition
}

// Java method: AutoValue_CodeReplacement.getLength()
func (this_ *AutoValue_CodeReplacement) GetLength() int32 {
	return this_.length
}

// Java method: AutoValue_CodeReplacement.getNewContent()
func (this_ *AutoValue_CodeReplacement) GetNewContent() string {
	return this_.newContent
}

// Java method: AutoValue_CodeReplacement.getSortKey()
func (this_ *AutoValue_CodeReplacement) GetSortKey() string {
	return this_.sortKey
}

// Java method: AutoValue_CodeReplacement.toString()
func (this_ *AutoValue_CodeReplacement) String() string {
	return ("CodeReplacement{startPosition=" + j2gort.StrValInt32(this_.startPosition) + ", length=" + j2gort.StrValInt32(this_.length) + ", newContent=" + j2gort.StrValString(this_.newContent) + ", sortKey=" + j2gort.StrValString(this_.sortKey) + "}")
}

// Java method: AutoValue_CodeReplacement.equals(java.lang.Object o)
func (this_ *AutoValue_CodeReplacement) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[CodeReplacement](o) {
		var that CodeReplacement = j2gort.AsIface[CodeReplacement](o)
		_ = that
		return ((((this_.startPosition == that.GetStartPosition()) && (this_.length == that.GetLength())) && func() bool {
			sr_2 := this_.newContent
			se_3 := that.GetNewContent()
			return j2gort.CheckNotNullString(sr_2) == se_3
		}()) && func() bool {
			sr_4 := this_.sortKey
			se_5 := that.GetSortKey()
			return j2gort.CheckNotNullString(sr_4) == se_5
		}())
	}
	return false
}

// Java method: AutoValue_CodeReplacement.hashCode()
func (this_ *AutoValue_CodeReplacement) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.startPosition)
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.length)
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.newContent)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.sortKey)))
	return h_d
}

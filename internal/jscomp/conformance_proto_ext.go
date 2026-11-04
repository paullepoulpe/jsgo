package jscomp

import (
	"fmt"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ParseConformanceConfigBytes parses binary-encoded ConformanceConfig protobuf
// bytes, returning an error instead of panicking on malformed wire input.
func ParseConformanceConfigBytes(data []byte) (cfg *ConformanceConfig, err error) {
	defer func() {
		if r := recover(); r != nil {
			cfg = nil
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("invalid ConformanceConfig protobuf: %v", r)
			}
		}
	}()
	return ConformanceConfig_ParseFrom(data), nil
}

// ConformanceConfig_ParseFrom parses binary-encoded ConformanceConfig bytes.
//
// Java: ConformanceConfig.parseFrom(byte[])
func ConformanceConfig_ParseFrom(data []byte) *ConformanceConfig {
	m := &ConformanceConfig{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

// Requirement_ParseFrom parses binary-encoded Requirement bytes.
//
// Java: Requirement.parseFrom(byte[])
func Requirement_ParseFrom(data []byte) *Requirement {
	m := &Requirement{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

// Clone returns a deep copy of m and its requirements.
func (m *ConformanceConfig) Clone() *ConformanceConfig {
	if m == nil {
		return nil
	}
	c := *m
	c.cachedSize = 0
	if len(m.requirement_) > 0 {
		c.requirement_ = make([]*Requirement, len(m.requirement_))
		for i, req := range m.requirement_ {
			if req != nil {
				c.requirement_[i] = req.Clone()
			}
		}
	}
	c.unknownFields = slices.Clone(m.unknownFields)
	return &c
}

// Java: ConformanceConfig.getRequirementCount()
func (m *ConformanceConfig) GetRequirementCount() int { return len(m.requirement_) }

// Java: ConformanceConfig.getRequirement(int)
func (m *ConformanceConfig) GetRequirement(index int) *Requirement { return m.requirement_[index] }

// Java: ConformanceConfig.Builder.clearRequirement()
func (m *ConformanceConfig) ClearRequirement() *ConformanceConfig {
	m.requirement_ = nil
	m.cachedSize = 0
	return m
}

// Java: Requirement.Builder.clearWhitelistEntry()
func (m *Requirement) ClearWhitelistEntry() *Requirement {
	m.whitelistEntry_ = nil
	m.cachedSize = 0
	return m
}

// Java: Requirement.Builder.clearAllowlistEntry()
func (m *Requirement) ClearAllowlistEntry() *Requirement {
	m.allowlistEntry_ = nil
	m.cachedSize = 0
	return m
}

// Clone returns a deep copy of m.
func (m *RequirementScopeEntry) Clone() *RequirementScopeEntry {
	if m == nil {
		return nil
	}
	c := *m
	c.cachedSize = 0
	c.prefix_ = slices.Clone(m.prefix_)
	c.regexp_ = slices.Clone(m.regexp_)
	c.comment_ = slices.Clone(m.comment_)
	c.unknownFields = slices.Clone(m.unknownFields)
	return &c
}

// Java: RequirementScopeEntry.getReason()
func (m *RequirementScopeEntry) GetReason() RequirementScopeEntry_Reason {
	if m.HasReason() {
		if r, ok := RequirementScopeEntry_Reason_ForNumber(m.reason_); ok {
			return r
		}
	}
	return RequirementScopeEntry_Reason_UNSPECIFIED
}

// Java: RequirementScopeEntry.getExplanation()
func (m *RequirementScopeEntry) GetExplanation() string { return m.explanation_ }

// Java: RequirementScopeEntry.getCommentList()
func (m *RequirementScopeEntry) GetCommentList() []string { return m.comment_ }

// Java: RequirementScopeEntry.getAutomaticallyPrune()
func (m *RequirementScopeEntry) GetAutomaticallyPrune() bool { return m.automaticallyPrune_ }

// Java: RequirementScopeEntry.getDoNotRecordViolationsInSummaryBecauseThisIsGencodeAndAlwaysSafe()
func (m *RequirementScopeEntry) GetDoNotRecordViolationsInSummaryBecauseThisIsGencodeAndAlwaysSafe() bool {
	return m.doNotRecordViolationsInSummaryBecauseThisIsGencodeAndAlwaysSafe_
}

// Java: RequirementScopeEntry.toByteArray()
func (m *RequirementScopeEntry) ToByteArray() []byte {
	return m.appendTo(make([]byte, 0, m.serializedSize()))
}

// Java: ConformanceViolationsProto.newBuilder()
func ConformanceViolationsProto_NewBuilder() *ConformanceViolationsProto {
	return &ConformanceViolationsProto{}
}

// ConformanceViolationsProto_ParseFrom parses binary-encoded ConformanceViolationsProto bytes.
func ConformanceViolationsProto_ParseFrom(data []byte) *ConformanceViolationsProto {
	m := &ConformanceViolationsProto{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

// ParseConformanceViolationsBytes parses binary-encoded ConformanceViolationsProto bytes safely.
func ParseConformanceViolationsBytes(data []byte) (out *ConformanceViolationsProto, err error) {
	defer func() {
		if r := recover(); r != nil {
			out = nil
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("invalid ConformanceViolationsProto protobuf: %v", r)
			}
		}
	}()
	return ConformanceViolationsProto_ParseFrom(data), nil
}

// Clone returns a deep copy of m and its ViolationProto entries.
func (m *ConformanceViolationsProto) Clone() *ConformanceViolationsProto {
	if m == nil {
		return nil
	}
	c := *m
	c.cachedSize = 0
	if len(m.violations_) > 0 {
		c.violations_ = make([]*ViolationProto, len(m.violations_))
		for i, v := range m.violations_ {
			if v != nil {
				c.violations_[i] = v.Clone()
			}
		}
	}
	return &c
}

// Java: ConformanceViolationsProto.getInputName()
func (m *ConformanceViolationsProto) GetInputName() string { return m.inputName_ }

// Java: ConformanceViolationsProto.Builder.setInputName(String)
func (m *ConformanceViolationsProto) SetInputName(v string) *ConformanceViolationsProto {
	m.inputName_ = v
	m.cachedSize = 0
	return m
}

// Java: ConformanceViolationsProto.getViolationsList()
func (m *ConformanceViolationsProto) GetViolationsList() []*ViolationProto { return m.violations_ }

// Java: ConformanceViolationsProto.getViolationsCount()
func (m *ConformanceViolationsProto) GetViolationsCount() int { return len(m.violations_) }

// Java: ConformanceViolationsProto.getViolations(int)
func (m *ConformanceViolationsProto) GetViolations(index int) *ViolationProto {
	return m.violations_[index]
}

// Java: ConformanceViolationsProto.Builder.addViolations(ViolationProto)
func (m *ConformanceViolationsProto) AddViolations(v *ViolationProto) *ConformanceViolationsProto {
	if v == nil {
		panic(jutil.NewNullPointerException())
	}
	m.violations_ = append(m.violations_, v)
	m.cachedSize = 0
	return m
}

// Java: ConformanceViolationsProto.Builder.clearViolations()
func (m *ConformanceViolationsProto) ClearViolations() *ConformanceViolationsProto {
	m.violations_ = nil
	m.cachedSize = 0
	return m
}

// Java: ConformanceViolationsProto.Builder.build()
func (m *ConformanceViolationsProto) Build() *ConformanceViolationsProto {
	m.cachedSize = 0
	return m
}

// Java: ConformanceViolationsProto.toByteArray()
func (m *ConformanceViolationsProto) ToByteArray() []byte {
	return m.appendTo(make([]byte, 0, m.serializedSize()))
}

// Java: ViolationProto.newBuilder()
func ViolationProto_NewBuilder() *ViolationProto { return &ViolationProto{} }

// ViolationProto_ParseFrom parses binary-encoded ViolationProto bytes.
func ViolationProto_ParseFrom(data []byte) *ViolationProto {
	m := &ViolationProto{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

// ParseViolationBytes parses binary-encoded ViolationProto bytes safely.
func ParseViolationBytes(data []byte) (out *ViolationProto, err error) {
	defer func() {
		if r := recover(); r != nil {
			out = nil
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("invalid ViolationProto protobuf: %v", r)
			}
		}
	}()
	return ViolationProto_ParseFrom(data), nil
}

// Clone returns a deep copy of m.
func (m *ViolationProto) Clone() *ViolationProto {
	if m == nil {
		return nil
	}
	c := *m
	c.cachedSize = 0
	if m.requirement_ != nil {
		c.requirement_ = m.requirement_.Clone()
	}
	if m.allowlistEntry_ != nil {
		c.allowlistEntry_ = m.allowlistEntry_.Clone()
	}
	return &c
}

// Java: ViolationProto.hasRequirement()
func (m *ViolationProto) HasRequirement() bool { return m.requirement_ != nil }

// Java: ViolationProto.getRequirement()
func (m *ViolationProto) GetRequirement() *Requirement { return m.requirement_ }

// Java: ViolationProto.Builder.setRequirement(Requirement)
func (m *ViolationProto) SetRequirement(v *Requirement) *ViolationProto {
	if v == nil {
		panic(jutil.NewNullPointerException())
	}
	m.requirement_ = v
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.hasAllowlistEntry()
func (m *ViolationProto) HasAllowlistEntry() bool { return m.allowlistEntry_ != nil }

// Java: ViolationProto.getAllowlistEntry()
func (m *ViolationProto) GetAllowlistEntry() *RequirementScopeEntry { return m.allowlistEntry_ }

// Java: ViolationProto.Builder.setAllowlistEntry(RequirementScopeEntry)
func (m *ViolationProto) SetAllowlistEntry(v *RequirementScopeEntry) *ViolationProto {
	if v == nil {
		panic(jutil.NewNullPointerException())
	}
	m.allowlistEntry_ = v
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.hasSourceName()
func (m *ViolationProto) HasSourceName() bool { return m.bitField0_&(1<<0) != 0 }

// Java: ViolationProto.getSourceName()
func (m *ViolationProto) GetSourceName() string { return m.sourceName_ }

// Java: ViolationProto.Builder.setSourceName(String)
func (m *ViolationProto) SetSourceName(v string) *ViolationProto {
	m.sourceName_ = v
	m.bitField0_ |= 1 << 0
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.getLineNumber()
func (m *ViolationProto) GetLineNumber() int32 { return m.lineNumber_ }

// Java: ViolationProto.Builder.setLineNumber(int)
func (m *ViolationProto) SetLineNumber(v int32) *ViolationProto {
	m.lineNumber_ = v
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.getCharNumber()
func (m *ViolationProto) GetCharNumber() int32 { return m.charNumber_ }

// Java: ViolationProto.Builder.setCharNumber(int)
func (m *ViolationProto) SetCharNumber(v int32) *ViolationProto {
	m.charNumber_ = v
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.getShouldReportDuringLibraryDepsConformanceChecker()
func (m *ViolationProto) GetShouldReportDuringLibraryDepsConformanceChecker() bool {
	return m.shouldReportDuringLibraryDepsConformanceChecker_
}

// Java: ViolationProto.Builder.setShouldReportDuringLibraryDepsConformanceChecker(boolean)
func (m *ViolationProto) SetShouldReportDuringLibraryDepsConformanceChecker(v bool) *ViolationProto {
	m.shouldReportDuringLibraryDepsConformanceChecker_ = v
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.Builder.build()
func (m *ViolationProto) Build() *ViolationProto {
	m.cachedSize = 0
	return m
}

// Java: ViolationProto.toByteArray()
func (m *ViolationProto) ToByteArray() []byte {
	return m.appendTo(make([]byte, 0, m.serializedSize()))
}

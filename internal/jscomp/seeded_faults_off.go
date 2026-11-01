//go:build !jsgo_faults

package jscomp

// Go-only: the seeded faults of seeded_faults.go are not built in (the
// default): every fault check is false at compile time.

const seededFault = faultNone

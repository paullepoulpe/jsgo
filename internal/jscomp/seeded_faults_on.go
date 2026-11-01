//go:build jsgo_faults

package jscomp

// Go-only: the seeded faults of seeded_faults.go are built in. A test-only
// build (`go build -tags jsgo_faults`): JSGO_FAST_FAULT names the fault to turn
// on in fast mode, and an unknown name stops the program at start-up.

import "os"

var seededFault = parseSeededFault(os.Getenv("JSGO_FAST_FAULT"))

package jscomp

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// sourceFileCacheLock uses the ordinary Go mutex only when the loader's concrete
// implementation is known not to call back into SourceFile. An embedding loader
// has a different concrete type and keeps a Java reentrant monitor, even if it
// embeds one of the native implementations.
//
// The native call graph while locked is:
//   - getCodeStateLocked -> native loader -> setCodeAndDoBookkeepingLocked;
//   - getLineOffsets/getCodeAndLineOffsets -> getCodeStateLocked and/or
//     findLineOffsetsLocked;
//   - GetCodeReader -> native OpenUncachedReader (GetCode is called after unlock);
//   - SetCodeAndDoBookkeeping/restoreCachedCounts -> immutable-state publication.
//
// Preloaded returns a string; Stub only throws; OnDisk invokes os file I/O and
// charset decoding; AtZip invokes the private ZIP cache and reader path. None
// invokes a SourceFile method. Selection is fixed before first publication.
type sourceFileCacheLock struct {
	mu        sync.Mutex
	reentrant *jutil.Monitor
}

func (l *sourceFileCacheLock) init(loader SourceFile_CodeLoader) {
	switch loader.(type) {
	case *SourceFile_CodeLoader_Preloaded, *SourceFile_CodeLoader_StubSourceFileCodeLoader,
		*SourceFile_CodeLoader_OnDisk, *SourceFile_CodeLoader_AtZip:
		// Exact native types only. Unknown/user loaders may reenter.
	default:
		l.reentrant = &jutil.Monitor{}
	}
}

func (l *sourceFileCacheLock) Lock() {
	if l.reentrant != nil {
		l.reentrant.Lock()
		return
	}
	l.mu.Lock()
}

func (l *sourceFileCacheLock) Unlock() {
	if l.reentrant != nil {
		l.reentrant.Unlock()
		return
	}
	l.mu.Unlock()
}

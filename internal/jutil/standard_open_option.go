package jutil

// Go-only: language-adapter — java.nio.file.StandardOpenOption.

// StandardOpenOption is java.nio.file.StandardOpenOption.
type StandardOpenOption int32

const (
	StandardOpenOption_READ StandardOpenOption = iota
	StandardOpenOption_WRITE
	StandardOpenOption_APPEND
	StandardOpenOption_TRUNCATE_EXISTING
	StandardOpenOption_CREATE
	StandardOpenOption_CREATE_NEW
	StandardOpenOption_DELETE_ON_CLOSE
	StandardOpenOption_SPARSE
	StandardOpenOption_SYNC
	StandardOpenOption_DSYNC
)

var standardOpenOption_names = [...]string{
	"READ", "WRITE", "APPEND", "TRUNCATE_EXISTING", "CREATE", "CREATE_NEW",
	"DELETE_ON_CLOSE", "SPARSE", "SYNC", "DSYNC",
}

// String is StandardOpenOption.name().
func (o StandardOpenOption) String() string { return standardOpenOption_names[o] }

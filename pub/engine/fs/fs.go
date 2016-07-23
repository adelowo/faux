package fs

import (
	"io"

	"github.com/influx6/faux/pub"
)

// FileSystem provides a core interface for FileSystem operations, where each
// instance of this type is series of operation for a single task or objective.
// This allows you to construct a nice chain of events, that can be recalled
// over and over again to do the same thing purely and nothing else.
type FileSystem interface {
	pub.Node

	ReadFile(string) FileSystem
	ReadReader(io.Reader) FileSystem

	WriteBytes([]byte) FileSystem
	WriteWritter(io.Writer)

	CreateFile(string) FileSystem
	WriteFile(string) FileSystem

	ReadDir(string) FileSystem
	MkDir(string) FileSystem

	DeleteFile(string) FileSystem
	DeleteDir(string) FileSystem
}

// New returns a new FileSystem
func New() FileSystem {
	var fu fs
	fu.Node = pub.InverseMagic(pub.IdentityHandler())
	return &fu
}

// DeferFileSystem provides a core interface for FileSystem operations, where each
// instance of this type is series of operation for a single task or objective.
// It defers in the sense that it takes in functions as argument. This may seem
// somewhat an inversion or an over-engineered but it becomes useful when you
// realise that FileSystem itself once created can not change, nor does its
// parameters. Hence DeferFileSystem allows you to pass functions in its critical
// areas that require inputs. This allows you to create an interesting permanent
// series of operation but with changing values.

//==============================================================================

type fs struct {
	pub.Node
}

// newFS returns a new instance pointer of a fs struct, which attaches the
// inverse of the provided Node to itself. The inverse ensures this is a
// serially chained link.
func newFS(node pub.Node) *fs {
	var fu fs
	fu.Node = node.Inverse()
	return &fu
}

// ReadFile adds a readFile operation whoes contents get passed to the next
// event/Node/Task in the link.
func (f *fs) ReadFile(path string) FileSystem {

	return f
}

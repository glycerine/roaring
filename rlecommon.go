package roaring

import (
	"fmt"
)

// common to rle32.go and rle16.go

// rleVerbose controls whether p() prints show up.
// The testing package sets this based on
// testing.Verbose().
var rleVerbose bool

// p is a shorthand for fmt.Printf with beginning and
// trailing newlines. p() makes it easy
// to add diagnostic print statements.
func p(format string, args ...interface{}) {
	if rleVerbose {
		fmt.Printf("\n"+format+"\n", args...)
	}
}

// MaxUint32 is only used internally for the endx
// value when UpperLimit32 is stored; users should
// only ever store up to UpperLimit32.
const MaxUint32 = 4294967295

// UpperLimit32 is the largest
// integer we can store in an RunContainer32. As
// we need to reserve one value for the open
// interval endpoint endx, this is MaxUint32 - 1.
const UpperLimit32 = MaxUint32 - 1

// MaxUint16 is only used internally for the endx
// value when UpperLimit16 is stored; users should
// only ever store up to UpperLimit16.
const MaxUint16 = 65535

// UpperLimit16 is the largest
// integer we can store in an RunContainer16. As
// we need to reserve one value for the open
// interval endpoint endx, this is MaxUint16 - 1.
const UpperLimit16 = MaxUint16 - 1

// SearchOptions allows us to accelerate RunContainer16.Search with
// prior knowledge of (mostly lower) bounds. This is used by Union
// and Intersect.
type SearchOptions struct {

	// start here instead of at 0
	StartIndex int

	// upper bound instead of len(rc.iv);
	// EndxIndex == 0 means ignore the bound and use
	// EndxIndex == n ==len(rc.iv) which is also
	// naturally the default for Search()
	// when opt = nil.
	EndxIndex int
}

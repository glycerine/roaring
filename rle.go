package roaring

//
// Copyright (c) 2016 by the roaring authors.
// Licensed under the Apache License, Version 2.0.
//
// We derive a few lines of code from the sort.Search
// funcion in the golang standard library. That function
// is Copyright 2009 The Go Authors, and licensed
// under the following BSD-style license.
/*
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"fmt"
	"sort"
	"unsafe"
)

// MaxUint32 is only used internally for the endx
// value when UpperLimit32 is stored; users should
// only ever store up to UpperLimit32.
const MaxUint32 = 4294967295

// UpperLimit32 is the largest
// integer we can store in an RunContainer32. As
// we need to reserve one value for the open
// interval endpoint endx, this is MaxUint32 - 1.
const UpperLimit32 = MaxUint32 - 1

// RunContainer32 does run-length encoding of sets of
// uint32 integers.
type RunContainer32 struct {
	iv   []interval32
	card int

	// avoid allocation during Search
	myOpts SearchOptions
}

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

// interval32 is the internal to RunContainer32
// structure that maintains the individual [start, endx)
// half-open intervals.
type interval32 struct {
	start uint32
	endx  uint32
}

// runlen returns iv.endx - iv.start, i.e. a
// count of integers in the [start, endx) half-open interval.
func (iv interval32) runlen() uint32 {
	return iv.endx - iv.start
}

// String produces a human viewable string of the contents.
func (iv interval32) String() string {
	return fmt.Sprintf("[%d, %d)", iv.start, iv.endx)
}

// String produces a human viewable string of the contents.
func (rc *RunContainer32) String() string {
	if len(rc.iv) == 0 {
		return "RunContainer32{}"
	}
	var s string
	for j, p := range rc.iv {
		s += fmt.Sprintf("%v:[%d, %d), ", j, p.start, p.endx)
	}
	return `RunContainer32{` + s + `}`
}

// And finds the intersection of rc and b.
func (rc *RunContainer32) And(b *Bitmap) *Bitmap {
	out := NewBitmap()
	for _, p := range rc.iv {
		for i := p.start; i < p.endx; i++ {
			if b.Contains(i) {
				out.Add(i)
			}
		}
	}
	return out
}

// Xor returns the exclusive-or of rc and b.
func (rc *RunContainer32) Xor(b *Bitmap) *Bitmap {
	out := b.Clone()
	for _, p := range rc.iv {
		for v := p.start; v < p.endx; v++ {
			if out.Contains(v) {
				out.RemoveRange(uint64(v), uint64(v+1))
			} else {
				out.Add(v)
			}
		}
	}
	return out
}

// Or returns the union of rc and b.
func (rc *RunContainer32) Or(b *Bitmap) *Bitmap {
	out := b.Clone()
	for _, p := range rc.iv {
		for v := p.start; v < p.endx; v++ {
			out.Add(v)
		}
	}
	return out
}

// uint32Slice is a sort.Sort convenience method
type uint32Slice []uint32

// Len returns the length of p.
func (p uint32Slice) Len() int { return len(p) }

// Less returns p[i] < p[j]
func (p uint32Slice) Less(i, j int) bool { return p[i] < p[j] }

// Swap swaps elements i and j.
func (p uint32Slice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// NewRunContainer32FromVals makes a new container from vals.
//
// For efficiency, vals should be sorted in ascending order.
// Ideally vals should not contain duplicates, but we detect and
// ignore them. If vals is already sorted in ascending order, then
// pass alreadySorted = true. Otherwise, for !alreadySorted,
// we will sort vals before creating a RunContainer32 of them.
// We sort the original vals, so this will change what the
// caller sees in vals as a side effect.
func NewRunContainer32FromVals(alreadySorted bool, vals ...uint32) *RunContainer32 {
	// keep this in sync with NewRunContainer32FromArray below

	if !alreadySorted {
		sort.Sort(uint32Slice(vals))
	}
	m := make([]interval32, 0)
	n := len(vals)
	actuallyAdded := 0
	var cur, prev uint32
	switch {
	case n == 0:
		// nothing more
	case n == 1:
		m = append(m, interval32{start: vals[0], endx: vals[0] + 1})
		actuallyAdded++
	default:
		runstart := vals[0]
		actuallyAdded++
		var runlen uint32
		for i := 1; i < n; i++ {
			prev = vals[i-1]
			cur = vals[i]
			if cur == prev+1 {
				runlen++
				actuallyAdded++
			} else {
				if cur < prev {
					panic(fmt.Sprintf("NewRunContainer32FromVals sees "+
						"unsorted vals; vals[%v]=cur=%v < prev=%v. Sort your vals"+
						" before calling us with alreadySorted == true.", i, cur, prev))
				}
				if cur == prev {
					// ignore duplicates
				} else {
					actuallyAdded++
					m = append(m, interval32{start: runstart, endx: runstart + 1 + runlen})
					runstart = cur
					runlen = 0
				}
			}
		}
		m = append(m, interval32{start: runstart, endx: runstart + 1 + runlen})
	}
	return &RunContainer32{iv: m, card: actuallyAdded}
}

//
// NewRunContainer32FromArray populates a new
// RunContainer32 from the contents of arr.
//
func NewRunContainer32FromArray(arr *arrayContainer) *RunContainer32 {
	// keep this in sync with NewRunContainer32FromVals above

	n := arr.getCardinality()
	actuallyAdded := 0
	m := make([]interval32, 0)
	var cur, prev uint32
	switch {
	case n == 0:
		// nothing more
	case n == 1:
		m = append(m, interval32{start: uint32(arr.content[0]), endx: uint32(arr.content[0]) + 1})
		actuallyAdded++
	default:
		runstart := uint32(arr.content[0])
		actuallyAdded++
		var runlen uint32
		for i := 1; i < n; i++ {
			prev = uint32(arr.content[i-1])
			cur = uint32(arr.content[i])
			if cur == prev+1 {
				runlen++
				actuallyAdded++
			} else {
				if cur < prev {
					panic(fmt.Sprintf("NewRunContainer32FromVals sees "+
						"unsorted vals; vals[%v]=cur=%v < prev=%v. Sort your vals"+
						" before calling us with alreadySorted == true.", i, cur, prev))
				}
				if cur == prev {
					//p("ignore duplicates")
				} else {
					actuallyAdded++
					m = append(m, interval32{start: runstart, endx: runstart + 1 + runlen})
					runstart = cur
					runlen = 0
				}
			}
		}
		m = append(m, interval32{start: runstart, endx: runstart + 1 + runlen})
	}
	return &RunContainer32{iv: m, card: actuallyAdded}
}

// Set adds the integers in vals to the set. Vals
// must be sorted in increasing order; if not, you should set
// alreadySorted to false, and we will sort them in place for you.
// (Be aware of this side effect -- it will affect the callers
// view of vals).
//
// If you have a small number of additions to an already
// big RunContainer32, calling Add() may be faster.
func (rc *RunContainer32) Set(alreadySorted bool, vals ...uint32) {

	rc2 := NewRunContainer32FromVals(alreadySorted, vals...)
	//p("Set: rc2 is %s", rc2)
	un := rc.Union(rc2)
	rc.iv = un.iv
	rc.card = 0
}

// canMerge returns true iff the intervals
// a and b either overlap or they are
// contiguous and so can be merged into
// a single interval.
func canMerge(a, b interval32) bool {
	if a.endx < b.start {
		return false
	}
	return b.endx >= a.start
}

// haveOverlap differs from canMerge in that
// it tells you if the intersection of a
// and b would contain an element (otherwise
// it would be the empty set, and we return
// false).
func haveOverlap(a, b interval32) bool {
	if a.endx <= b.start {
		return false
	}
	return b.endx > a.start
}

// mergeInterval32s joins a and b into a
// new interval, and panics if it cannot.
func mergeInterval32s(a, b interval32) (res interval32) {
	if !canMerge(a, b) {
		panic(fmt.Sprintf("cannot merge %#v and %#v", a, b))
	}
	if b.start < a.start {
		res.start = b.start
	} else {
		res.start = a.start
	}
	if b.endx > a.endx {
		res.endx = b.endx
	} else {
		res.endx = a.endx
	}
	return
}

// intersectInterval32s returns the intersection
// of a and b. The isEmpty flag will be true if
// a and b were disjoint.
func intersectInterval32s(a, b interval32) (res interval32, isEmpty bool) {
	if !haveOverlap(a, b) {
		isEmpty = true
		return
	}
	if b.start > a.start {
		res.start = b.start
	} else {
		res.start = a.start
	}
	if b.endx < a.endx {
		res.endx = b.endx
	} else {
		res.endx = a.endx
	}
	return
}

// Union merges two RunContainer32s, producing
// a new RunContainer32 with the union of rc and b.
func (rc *RunContainer32) Union(b *RunContainer32) *RunContainer32 {

	// rc is also known as 'a' here, but golint insisted we
	// call it rc for consistency with the rest of the methods.

	m := make([]interval32, 0)

	alim := len(rc.iv)
	blim := len(b.iv)

	na := 0 // next from a
	nb := 0 // next from b

	// merged holds the current merge output, which might
	// get additional merges before being appended to m.
	var merged interval32
	var mergedUsed bool // is merged being used at the moment?

	var cura interval32 // currently considering this interval32 from a
	var curb interval32 // currently considering this interval32 from b

	pass := 0
	for na < alim && nb < blim {
		pass++
		cura = rc.iv[na]
		curb = b.iv[nb]

		//p("pass=%v, cura=%v, curb=%v, merged=%v, mergedUsed=%v m=%v", pass, cura, curb, merged, mergedUsed, m)

		if mergedUsed {
			//p("mergedUsed is true")
			mergedUpdated := false
			if canMerge(cura, merged) {
				//p("canMerge(cura=%s, merged=%s) is true", cura, merged)
				merged = mergeInterval32s(cura, merged)
				na = rc.indexOfIntervalAtOrAfter(merged.endx, na+1)
				mergedUpdated = true
			}
			if canMerge(curb, merged) {
				//p("canMerge(curb=%s, merged=%s) is true", curb, merged)
				merged = mergeInterval32s(curb, merged)
				nb = b.indexOfIntervalAtOrAfter(merged.endx, nb+1)
				mergedUpdated = true
			}
			if !mergedUpdated {
				//p("!mergedUpdated")
				// we know that merged is disjoint from cura and curb
				m = append(m, merged)
				mergedUsed = false
			}
			continue

		} else {
			//p("!mergedUsed")
			// !mergedUsed
			if !canMerge(cura, curb) {
				if cura.start < curb.start {
					//p("cura is before curb")
					m = append(m, cura)
					na++
				} else {
					//p("curb is before cura")
					m = append(m, curb)
					nb++
				}
			} else {
				//p("intervals are not disjoint, we can merge them. cura=%s, curb=%s", cura, curb)
				merged = mergeInterval32s(cura, curb)
				mergedUsed = true
				na = rc.indexOfIntervalAtOrAfter(merged.endx, na+1)
				nb = b.indexOfIntervalAtOrAfter(merged.endx, nb+1)
			}
		}
	}
	var aDone, bDone bool
	if na >= alim {
		aDone = true
		//p("na(%v) >= alim=%v, the 'a' sequence is done, finish up on 'merged' and 'b'.", na, alim)
	}
	if nb >= blim {
		bDone = true
		//p("nb(%v) >= blim=%v, the 'b' sequence is done, finish up on 'merged' and 'a'.", nb, blim)
	}
	// finish by merging anything remaining into merged we can:
	if mergedUsed {
		if !aDone {
		aAdds:
			for na < alim {
				cura = rc.iv[na]
				if canMerge(cura, merged) {
					//p("canMerge(cura=%s, merged=%s) is true. na=%v", cura, merged, na)
					merged = mergeInterval32s(cura, merged)
					na = rc.indexOfIntervalAtOrAfter(merged.endx, na+1)
				} else {
					break aAdds
				}
			}

		}

		if !bDone {
		bAdds:
			for nb < blim {
				curb = b.iv[nb]
				if canMerge(curb, merged) {
					//p("canMerge(curb=%s, merged=%s) is true. nb=%v", curb, merged, nb)
					merged = mergeInterval32s(curb, merged)
					nb = b.indexOfIntervalAtOrAfter(merged.endx, nb+1)
				} else {
					break bAdds
				}
			}

		}

		//p("mergedUsed==true, before adding merged=%s, m=%v", merged, sliceToString(m))
		m = append(m, merged)
		//p("added mergedUsed, m=%v", sliceToString(m))
	}
	if na < alim {
		//p("adding the rest of a.vi[na:] = %v", sliceToString(rc.iv[na:]))
		m = append(m, rc.iv[na:]...)
		//p("after the rest of a.vi[na:] to m, now m = %v", sliceToString(m))
	}
	if nb < blim {
		//p("adding the rest of b.vi[nb:] = %v", sliceToString(b.iv[nb:]))
		m = append(m, b.iv[nb:]...)
		//p("after the rest of a.vi[nb:] to m, now m = %v", sliceToString(m))
	}

	//p("making res out of m = %v", sliceToString(m))
	res := &RunContainer32{iv: m}
	//p("Union returning %s", res)
	return res
}

// indexOfIntervalAtOrAfter is a helper for Union. We check
// for already and panic, as this is dedicated to use
// by Union() and that should always be the case when
// Union calls.
func (rc *RunContainer32) indexOfIntervalAtOrAfter(key uint32, startIndex int) int {
	rc.myOpts.StartIndex = startIndex
	rc.myOpts.EndxIndex = 0

	w, already, _ := rc.Search(key, &rc.myOpts)
	if already {
		return w
	}
	return w + 1
}

// Intersect returns a new RunContainer32 holding the
// intersection of rc (also known as 'a')  and b.
func (rc *RunContainer32) Intersect(b *RunContainer32) *RunContainer32 {

	a := rc
	numa := len(a.iv)
	numb := len(b.iv)
	res := &RunContainer32{}
	if numa == 0 || numb == 0 {
		//p("intersection is empty, returning early")
		return res
	}

	if numa == 1 && numb == 1 {
		if !haveOverlap(a.iv[0], b.iv[0]) {
			//p("intersection is empty, returning early")
			return res
		}
	}

	output := make([]interval32, 0)

	acuri := 0
	bcuri := 0

	astart := a.iv[acuri].start
	bstart := b.iv[bcuri].start

	var intersection interval32
	var leftoverStart uint32
	var isOverlap, isLeftoverA, isLeftoverB bool
	var done bool
	pass := 0
toploop:
	for acuri < numa && bcuri < numb {
		//p("============     top of loop, pass = %v", pass)
		pass++

		isOverlap, isLeftoverA, isLeftoverB, leftoverStart, intersection = intersectWithLeftover(astart, a.iv[acuri].endx, bstart, b.iv[bcuri].endx)

		//p("acuri=%v, astart=%v, a.iv[acuri].endx=%v,   bcuri=%v, bstart=%v, b.iv[bcuri].endx=%v, isOverlap=%v, isLeftoverA=%v, isLeftoverB=%v, leftoverStart=%v, intersection = %#v", acuri, astart, a.iv[acuri].endx, bcuri, bstart, b.iv[bcuri].endx, isOverlap, isLeftoverA, isLeftoverB, leftoverStart, intersection)

		if !isOverlap {
			switch {
			case astart < bstart:
				//p("no overlap, astart < bstart ... acuri = %v, key=bstart= %v", acuri, bstart)
				acuri, done = a.findNextIntervalThatIntersectsStartingFrom(acuri+1, bstart)
				//p("b.findNextIntervalThatIntersectsStartingFrom(startIndex=%v, key=%v) returned: acuri = %v, done=%v", acuri+1, bstart, acuri, done)
				if done {
					break toploop
				}
				astart = a.iv[acuri].start

			case astart > bstart:
				//p("no overlap, astart > bstart ... bcuri = %v, key=astart= %v", bcuri, astart)
				bcuri, done = b.findNextIntervalThatIntersectsStartingFrom(bcuri+1, astart)
				//p("b.findNextIntervalThatIntersectsStartingFrom(startIndex=%v, key=%v) returned: bcuri = %v, done=%v", bcuri+1, astart, bcuri, done)
				if done {
					break toploop
				}
				bstart = b.iv[bcuri].start

				//default:
				//	panic("impossible that astart == bstart, since !isOverlap")
			}

		} else {
			// isOverlap
			//p("isOverlap == true, intersection = %#v", intersection)
			output = append(output, intersection)
			switch {
			case isLeftoverA:
				//p("isLeftoverA true... new astart = leftoverStart = %v", leftoverStart)
				// note that we change astart without advancing acuri,
				// since we need to capture any 2ndary intersections with a.iv[acuri]
				astart = leftoverStart
				bcuri++
				if bcuri >= numb {
					break toploop
				}
				bstart = b.iv[bcuri].start
				//p("new bstart is %v", bstart)
			case isLeftoverB:
				//p("isLeftoverB true... new bstart = leftoverStart = %v", leftoverStart)
				// note that we change bstart without advancing bcuri,
				// since we need to capture any 2ndary intersections with b.iv[bcuri]
				bstart = leftoverStart
				acuri++
				if acuri >= numa {
					break toploop
				}
				astart = a.iv[acuri].start
				//p(" ... and new astart is %v", astart)
			default:
				//p("no leftovers after intersection")
				// neither had leftover, both completely consumed
				// optionally, assert for sanity:
				//if a.iv[acuri].endx != b.iv[bcuri].endx {
				//	panic("huh? should only be possible that endx agree now!")
				//}

				// advance to next a interval
				acuri++
				if acuri >= numa {
					//p("out of 'a' elements, breaking out of loop")
					break toploop
				}
				astart = a.iv[acuri].start

				// advance to next b interval
				bcuri++
				if bcuri >= numb {
					//p("out of 'b' elements, breaking out of loop")
					break toploop
				}
				bstart = b.iv[bcuri].start
				//p("no leftovers after intersection, new acuri=%v, astart=%v, bcuri=%v, bstart=%v", acuri, astart, bcuri, bstart)
			}
		}
	} // end for toploop

	if len(output) == 0 {
		return res
	}

	res.iv = output
	//p("Intersect returning %#v", res)
	return res
}

// Get returns true iff key is in the container.
func (rc *RunContainer32) Get(key uint32) bool {
	_, in, _ := rc.Search(key, nil)
	return in
}

// NumIntervals returns the count of intervals in the container.
func (rc *RunContainer32) NumIntervals() int {
	return len(rc.iv)
}

// SearchOptions allows us to accelerate RunContainer32.Search with
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

// Search returns alreadyPresent to indicate if the
// key is already in one of our interval32s.
//
// If key is alreadyPresent, then whichInterval32 tells
// you where.
//
// If key is not already present, then whichInterval32 is
// set as follows:
//
//  a) whichInterval32 == len(rc.iv)-1 if key is beyond our
//     last interval32 in rc.iv;
//
//  b) whichInterval32 == -1 if key is before our first
//     interval32 in rc.iv;
//
//  c) whichInterval32 is set to the minimum index of rc.iv
//     which comes strictly before the key;
//     so  key >= rc.iv[whichInterval32].endx,
//     and  if whichInterval32+1 exists, then key < rc.iv[whichInterval32+1].startx
//     (Note that whichInterval32+1 won't exist when
//     whichInterval32 is the last interval.)
//
// RunContainer32.Search always returns whichInterval32 < len(rc.iv).
//
// If not nil, opts can be used to further restrict
// the search space.
//
func (rc *RunContainer32) Search(key uint32, opts *SearchOptions) (whichInterval32 int, alreadyPresent bool, numCompares int) {

	n := len(rc.iv)
	if n == 0 {
		return -1, false, 0
	}

	startIndex := 0
	endxIndex := n
	if opts != nil {
		startIndex = opts.StartIndex

		// let EndxIndex == 0 mean no effect
		if opts.EndxIndex > 0 {
			endxIndex = opts.EndxIndex
		}
	}

	// sort.Search returns the smallest index i
	// in [0, n) at which f(i) is true, assuming that on the range [0, n),
	// f(i) == true implies f(i+1) == true.
	// If there is no such index, Search returns n.

	// For correctness, this began as verbatim snippet from
	// sort.Search in the Go standard lib.
	// We inline our comparison function for speed, and
	// annotate with numCompares
	// to observe and test that extra bounds are utilized.
	i, j := startIndex, endxIndex
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h as the bisector
		// i <= h < j
		numCompares++
		if !(key < rc.iv[h].start) {
			i = h + 1
		} else {
			j = h
		}
	}
	below := i
	// end std lib snippet.

	// The above is a simple in-lining and annotation of:
	/*	below := sort.Search(n,
		func(i int) bool {
			return key < rc.iv[i].start
		})
	*/
	whichInterval32 = below - 1

	if below == n {
		// all falses => key is >= start of all interval32s
		// ... so does it belong to the last interval32?
		if key < rc.iv[n-1].endx {
			// yes, it belongs to the last interval32
			alreadyPresent = true
			return
		}
		// no, it is beyond the last interval32.
		// leave alreadyPreset = false
		return
	}

	// INVAR: key is below rc.iv[below]
	if below == 0 {
		// key is before the first first interval32.
		// leave alreadyPresent = false
		return
	}

	// INVAR: key is >= rc.iv[below-1].start and
	//        key is <  rc.iv[below].start

	// is key in below-1 interval32?
	if key >= rc.iv[below-1].start && key < rc.iv[below-1].endx {
		// yes, it is. key is in below-1 interval32.
		alreadyPresent = true
		return
	}

	// INVAR: key >= rc.iv[below-1].endx && key < rc.iv[below].start
	//p("Search, INVAR: key >= rc.iv[below-1].endx && key < rc.iv[below].start, where key=%v, below=%v, below-1=%v, rc.iv[below-1]=%v, rc.iv[below]=%v", key, below, below-1, rc.iv[below-1], rc.iv[below])
	// leave alreadyPresent = false
	return
}

// Cardinality returns the count of the integers stored in the
// RunContainer32.
func (rc *RunContainer32) Cardinality() int {
	if len(rc.iv) == 0 {
		rc.card = 0
		return 0
	}
	if rc.card > 0 {
		return rc.card // already cached
	}
	// have to compute it
	n := 0
	for _, p := range rc.iv {
		n += int(p.runlen())
	}
	rc.card = n // cache it
	return n
}

// AsSlice decompresses the contents into a []uint32 slice.
func (rc *RunContainer32) AsSlice() []uint32 {
	s := make([]uint32, rc.Cardinality())
	j := 0
	for _, p := range rc.iv {
		for i := p.start; i < p.endx; i++ {
			s[j] = uint32(i)
			j++
		}
	}
	return s
}

// NewRunContainer32 creates an empty run container.
func NewRunContainer32() *RunContainer32 {
	return &RunContainer32{}
}

// NewRunContainer32CopyIv creates a run container, initializing
// with a copy of the supplied iv slice.
//
func NewRunContainer32CopyIv(iv []interval32) *RunContainer32 {
	rc := &RunContainer32{
		iv: make([]interval32, len(iv)),
	}
	copy(rc.iv, iv)
	return rc
}

// NewRunContainer32TakeOwnership returns a new RunContainer32
// backed by the provided iv slice, which we will
// assume exclusive control over from now on.
//
func NewRunContainer32TakeOwnership(iv []interval32) *RunContainer32 {
	rc := &RunContainer32{
		iv: iv,
	}
	return rc
}

const baseRc32Size = int(unsafe.Sizeof(RunContainer32{}))
const perIntervalRc32Size = int(unsafe.Sizeof(interval32{}))

// RunContainer32SerializedSizeInBytes returns the number of bytes of memory
// required to hold numRuns in a RunContainer32.
func RunContainer32SerializedSizeInBytes(numRuns int) int {
	return baseRc32Size + perIntervalRc32Size*numRuns
}

// SerializedSizeInBytes returns the number of bytes of memory
// required by this RunContainer32.
func (rc *RunContainer32) SerializedSizeInBytes() int {
	return baseRc32Size + perIntervalRc32Size*rc.Cardinality()
}

// Add adds a single value k to the set.
func (rc *RunContainer32) Add(k uint32) {
	// TODO comment from RunContainer32.java:
	// it might be better and simpler to do return
	// toBitmapOrArrayContainer(getCardinality()).add(k)
	// but note that some unit tests use this method to build up test
	// runcontainers without calling runOptimize
	index, present, _ := rc.Search(k, nil)
	if present {
		return // already there
	}
	// increment card if it is cached already
	if rc.card > 0 {
		rc.card++
	}
	n := len(rc.iv)
	if index == -1 {
		// we may need to extend the first run
		if n > 0 {
			if rc.iv[0].start == k+1 {
				rc.iv[0].start = k
				return
			}
		}
		// nope, k stands alone, starting the new first interval32.
		rc.iv = append([]interval32{interval32{start: k, endx: k + 1}}, rc.iv...)
		return
	}

	// are we off the end? handle both index == n and index == n-1:
	if index >= n-1 {
		if rc.iv[n-1].endx == k {
			rc.iv[n-1].endx++
			return
		}
		rc.iv = append(rc.iv, interval32{start: k, endx: k + 1})
		return
	}

	// INVAR: index and index+1 both exist, and k goes between them.
	//
	// Now: add k into the middle,
	// possibly fusing with index or index+1 interval32
	// and possibly resulting in fusing of two interval32s
	// that had a one integer gap.

	left := index
	right := index + 1

	// are we fusing left and right by adding k?
	if rc.iv[left].endx == k && rc.iv[right].start == k+1 {
		// fuse into left
		rc.iv[left].endx = rc.iv[right].endx
		// remove redundant right
		rc.iv = append(rc.iv[:left+1], rc.iv[right+1:]...)
		return
	}

	// are we an addition to left?
	if rc.iv[left].endx == k {
		// yes
		rc.iv[left].endx++
		return
	}

	// are we an addition to right?
	if rc.iv[right].start == k+1 {
		// yes
		rc.iv[right].start = k
		return
	}

	// k makes a standalone new interval32, inserted in the middle
	tail := append([]interval32{interval32{start: k, endx: k + 1}}, rc.iv[right:]...)
	rc.iv = append(rc.iv[:left+1], tail...)
}

// RunIterator32 advice: you must call Next() at least once
// before calling Cur(); and you should call HasNext()
// before calling Next() to insure there are contents.
type RunIterator32 struct {
	rc            *RunContainer32
	curIndex      int
	curPosInIndex uint32
	curSeq        int
}

// NewRunIterator32 returns a new empty run container.
func (rc *RunContainer32) NewRunIterator32() *RunIterator32 {
	return &RunIterator32{rc: rc, curIndex: -1}
}

// HasNext returns false if calling Next will panic. It
// returns true when there is at least one more value
// available in the iteration sequence.
func (ri *RunIterator32) HasNext() bool {
	if len(ri.rc.iv) == 0 {
		return false
	}
	if ri.curIndex == -1 {
		return true
	}
	return ri.curSeq+1 < ri.rc.Cardinality()
}

// Cur returns the current value pointed to by the iterator.
func (ri *RunIterator32) Cur() uint32 {
	//p("in Cur, curIndex=%v, curPosInIndex=%v", ri.curIndex, ri.curPosInIndex)
	return ri.rc.iv[ri.curIndex].start + ri.curPosInIndex
}

// Next returns the next value in the iteration sequence.
func (ri *RunIterator32) Next() uint32 {
	if !ri.HasNext() {
		panic("no Next available")
	}
	if ri.curIndex >= len(ri.rc.iv) {
		panic("RunIterator.Next() going beyond what is available")
	}
	if ri.curIndex == -1 {
		// first time is special
		ri.curIndex = 0
	} else {
		ri.curPosInIndex++
		if ri.rc.iv[ri.curIndex].start+ri.curPosInIndex == ri.rc.iv[ri.curIndex].endx {
			//p("rolling from ri.curIndex==%v to ri.curIndex=%v", ri.curIndex, ri.curIndex+1)
			ri.curPosInIndex = 0
			ri.curIndex++
		} else {
			//p("no roll ... ri.curPosInIndex is now %v, ri.rc.iv[ri.curIndex].endx=%v", ri.curPosInIndex, ri.rc.iv[ri.curIndex].endx)
		}
		ri.curSeq++
	}
	return ri.Cur()
}

// Remove removes the element that the iterator
// is on from the run container. You can use
// Cur if you want to double check what is about
// to be deleted.
func (ri *RunIterator32) Remove() uint32 {
	n := ri.rc.Cardinality()
	if n == 0 {
		panic("RunIterator.Remove called on empty RunContainer32")
	}
	cur := ri.Cur()

	ri.rc.deleteAt(&ri.curIndex, &ri.curPosInIndex, &ri.curSeq)
	return cur
}

// Remove removes key from the container.
func (rc *RunContainer32) Remove(key uint32) (wasPresent bool) {
	var index, curSeq int
	index, wasPresent, _ = rc.Search(key, nil)
	if !wasPresent {
		return // already removed, nothing to do.
	}
	pos := key - rc.iv[index].start
	rc.deleteAt(&index, &pos, &curSeq)
	return
}

// internal helper functions

func (rc *RunContainer32) deleteAt(curIndex *int, curPosInIndex *uint32, curSeq *int) {
	rc.card--
	(*curSeq)--
	ci := *curIndex
	pos := *curPosInIndex

	// are we first, last, or in the middle of our interval32?
	switch {
	case pos == 0:
		// first
		rc.iv[ci].start++
		// does our interval32 disappear?
		if rc.iv[ci].start == rc.iv[ci].endx {
			// yes, delete it
			rc.iv = append(rc.iv[:ci], rc.iv[ci+1:]...)
			// curIndex stays the same, since the delete did
			// the advance for us.
			*curPosInIndex = 0
		}
	case pos == rc.iv[ci].runlen()-1:
		// last
		rc.iv[ci].endx--
		// our interval32 cannot disappear, else we would have been pos == 0, case first above.
		//p("deleteAt: pos is last case, curIndex=%v, curPosInIndex=%v", *curIndex, *curPosInIndex)
		(*curPosInIndex)--
		// if we leave *curIndex alone, then Next() will work properly even after the delete.
		//p("deleteAt: pos is last case, after update: curIndex=%v, curPosInIndex=%v", *curIndex, *curPosInIndex)
	default:
		//p("middle...split")
		//middle
		// split into two, adding an interval32
		new0 := interval32{
			start: rc.iv[ci].start,
			endx:  rc.iv[ci].start + *curPosInIndex}

		new1 := interval32{
			start: rc.iv[ci].start + *curPosInIndex + 1,
			endx:  rc.iv[ci].endx}

		//p("new0 = %#v", new0)
		//p("new1 = %#v", new1)

		tail := append([]interval32{new0, new1}, rc.iv[ci+1:]...)
		rc.iv = append(rc.iv[:ci], tail...)
		// update curIndex and curPosInIndex
		(*curIndex)++
		*curPosInIndex = 0
	}

}

func haveOverlap4(astart, aendx, bstart, bendx uint32) bool {
	if aendx <= bstart {
		return false
	}
	return bendx > astart
}

func intersectWithLeftover(astart, aendx, bstart, bendx uint32) (isOverlap, isLeftoverA, isLeftoverB bool, leftoverStart uint32, intersection interval32) {
	if !haveOverlap4(astart, aendx, bstart, bendx) {
		return
	}
	isOverlap = true

	// do the intersection:
	if bstart > astart {
		intersection.start = bstart
	} else {
		intersection.start = astart
	}
	switch {
	case bendx < aendx:
		isLeftoverA = true
		leftoverStart = bendx
		intersection.endx = bendx
	case aendx < bendx:
		isLeftoverB = true
		leftoverStart = aendx
		intersection.endx = aendx
	default:
		// aendx == bendx
		intersection.endx = aendx
	}

	return
}

func (rc *RunContainer32) findNextIntervalThatIntersectsStartingFrom(startIndex int, key uint32) (index int, done bool) {

	rc.myOpts.StartIndex = startIndex
	rc.myOpts.EndxIndex = 0

	w, _, _ := rc.Search(key, &rc.myOpts)
	// rc.Search always returns w < len(rc.iv)
	if w < startIndex {
		// not found and comes before lower bound startIndex,
		// so just use the lower bound.
		if startIndex == len(rc.iv) {
			// also this bump up means that we are done
			return startIndex, true
		}
		return startIndex, false
	}

	return w, false
}

func sliceToString(m []interval32) string {
	s := ""
	for i := range m {
		s += fmt.Sprintf("%v: %s, ", i, m[i])
	}
	return s
}

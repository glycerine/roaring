package roaring

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	rleVerbose = testing.Verbose()
}

func TestRleInterval32s(t *testing.T) {

	Convey("canMerge, and mergeInterval32s should do what they say", t, func() {
		a := interval32{start: 0, endx: 10}
		msg := a.String()
		p("a is %v", msg)
		b := interval32{start: 0, endx: 2}
		report := sliceToString32([]interval32{a, b})
		_ = report
		p("a and b together are: %s", report)
		c := interval32{start: 2, endx: 5}
		d := interval32{start: 2, endx: 6}
		e := interval32{start: 0, endx: 5}
		f := interval32{start: 9, endx: 10}
		g := interval32{start: 8, endx: 10}
		h := interval32{start: 5, endx: 7}
		i := interval32{start: 6, endx: 7}

		aIb, empty := intersectInterval32s(a, b)
		So(empty, ShouldBeFalse)
		So(aIb, ShouldResemble, b)

		So(canMerge32(b, c), ShouldBeTrue)
		So(canMerge32(c, b), ShouldBeTrue)
		So(canMerge32(a, h), ShouldBeTrue)

		So(canMerge32(d, e), ShouldBeTrue)
		So(canMerge32(f, g), ShouldBeTrue)
		So(canMerge32(c, h), ShouldBeTrue)

		So(canMerge32(b, h), ShouldBeFalse)
		So(canMerge32(h, b), ShouldBeFalse)
		So(canMerge32(c, i), ShouldBeFalse)

		So(mergeInterval32s(b, c), ShouldResemble, e)
		So(mergeInterval32s(c, b), ShouldResemble, e)

		So(mergeInterval32s(h, i), ShouldResemble, h)
		So(mergeInterval32s(i, h), ShouldResemble, h)

		So(mergeInterval32s(interval32{start: 0, endx: 1}, interval32{start: 1, endx: 2}), ShouldResemble, interval32{start: 0, endx: 2})
		So(mergeInterval32s(interval32{start: 1, endx: 2}, interval32{start: 0, endx: 1}), ShouldResemble, interval32{start: 0, endx: 2})
		So(mergeInterval32s(interval32{start: 0, endx: 4}, interval32{start: 3, endx: 5}), ShouldResemble, interval32{start: 0, endx: 5})
		So(mergeInterval32s(interval32{start: 0, endx: 4}, interval32{start: 3, endx: 4}), ShouldResemble, interval32{start: 0, endx: 4})

		So(mergeInterval32s(interval32{start: 0, endx: 8}, interval32{start: 1, endx: 7}), ShouldResemble, interval32{start: 0, endx: 8})
		So(mergeInterval32s(interval32{start: 1, endx: 7}, interval32{start: 0, endx: 8}), ShouldResemble, interval32{start: 0, endx: 8})

		So(func() { _ = mergeInterval32s(interval32{start: 0, endx: 1}, interval32{start: 2, endx: 3}) }, ShouldPanic)

	})
}

func TestRleRunIterator(t *testing.T) {

	Convey("RunIterator unit tests for Cur, Next, HasNext, and Remove should pass", t, func() {
		{
			rc := newRunContainer32()
			msg := rc.String()
			_ = msg
			p("an empty container: '%s'\n", msg)
			So(rc.Cardinality(), ShouldEqual, 0)
			it := rc.NewRunIterator32()
			So(it.HasNext(), ShouldBeFalse)
		}
		{
			rc := newRunContainer32TakeOwnership([]interval32{{start: 4, endx: 5}})
			So(rc.Cardinality(), ShouldEqual, 1)
			it := rc.NewRunIterator32()
			So(it.HasNext(), ShouldBeTrue)
			So(it.Next(), ShouldResemble, uint32(4))
			So(it.Cur(), ShouldResemble, uint32(4))
		}
		{
			rc := newRunContainer32CopyIv([]interval32{{start: 4, endx: 10}})
			So(rc.Cardinality(), ShouldEqual, 6)
			it := rc.NewRunIterator32()
			So(it.HasNext(), ShouldBeTrue)
			for i := 4; i < 10; i++ {
				So(it.Next(), ShouldEqual, uint32(i))
			}
			So(it.HasNext(), ShouldBeFalse)
		}

		{
			rc := newRunContainer32TakeOwnership([]interval32{{start: 4, endx: 10}})
			So(rc.Cardinality(), ShouldEqual, 6)
			So(rc.serializedSizeInBytes(), ShouldEqual, 96)

			it := rc.NewRunIterator32()
			So(it.HasNext(), ShouldBeTrue)
			for i := 4; i < 6; i++ {
				So(it.Next(), ShouldEqual, uint32(i))
			}
			So(it.Cur(), ShouldEqual, uint32(5))

			p("before Remove of 5, rc = '%s'", rc)

			So(it.Remove(), ShouldEqual, uint32(5))

			p("after Remove of 5, rc = '%s'", rc)
			So(rc.Cardinality(), ShouldEqual, 5)

			it2 := rc.NewRunIterator32()
			So(rc.Cardinality(), ShouldEqual, 5)
			So(it2.Next(), ShouldEqual, uint32(4))
			for i := 6; i < 10; i++ {
				So(it2.Next(), ShouldEqual, uint32(i))
			}
		}
		{
			rc := newRunContainer32TakeOwnership([]interval32{
				{start: 0, endx: 1},
				{start: 2, endx: 3},
				{start: 4, endx: 5},
			})
			rc1 := newRunContainer32TakeOwnership([]interval32{
				{start: 6, endx: 8},
				{start: 10, endx: 12},
				{start: UpperLimit32 - 1, endx: UpperLimit32},
			})

			rc = rc.Union(rc1)

			So(rc.Cardinality(), ShouldEqual, 8)
			it := rc.NewRunIterator32()
			So(it.Next(), ShouldEqual, uint32(0))
			So(it.Next(), ShouldEqual, uint32(2))
			So(it.Next(), ShouldEqual, uint32(4))
			So(it.Next(), ShouldEqual, uint32(6))
			So(it.Next(), ShouldEqual, uint32(7))
			So(it.Next(), ShouldEqual, uint32(10))
			So(it.Next(), ShouldEqual, uint32(11))
			So(it.Next(), ShouldEqual, uint32(UpperLimit32-1))
			So(it.HasNext(), ShouldEqual, false)

			rc2 := newRunContainer32TakeOwnership([]interval32{
				{start: 0, endx: UpperLimit32},
			})

			p("Union with a full [0,2^32-1) container should yield that same single interval run container")
			rc2 = rc2.Union(rc)
			So(rc2.NumIntervals(), ShouldEqual, 1)
		}
	})
}

func TestRleRunSearch(t *testing.T) {

	Convey("RunContainer.Search should respect the prior bounds we provide for efficiency of searching through a subset of the intervals", t, func() {
		{
			vals := []uint32{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, UpperLimit32 - 3, UpperLimit32 - 1}
			expWhere := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
			absent := []uint32{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, UpperLimit32 - 2}

			rc := newRunContainer32FromVals(true, vals...)

			So(rc.Cardinality(), ShouldEqual, 12)

			var where int
			var present bool

			for i, v := range vals {
				where, present, _ = rc.Search(v, nil)
				So(present, ShouldBeTrue)
				So(where, ShouldEqual, expWhere[i])
			}

			for i, v := range absent {
				p("absent check on i=%v, v=%v in rc=%v", i, v, rc)
				where, present, _ = rc.Search(v, nil)
				So(present, ShouldBeFalse)
				So(where, ShouldEqual, i)
			}

			// delete the UpperLimit32 -1 so we can test
			// the behavior when searching near upper limit.

			p("before removing UpperLimit32-1: %v", rc)

			So(rc.Cardinality(), ShouldEqual, 12)
			So(rc.NumIntervals(), ShouldEqual, 12)

			rc.Remove(UpperLimit32 - 1)
			p("after removing UpperLimit32-1: %v", rc)
			So(rc.Cardinality(), ShouldEqual, 11)
			So(rc.NumIntervals(), ShouldEqual, 11)

			p("Search for absent UpperLimit32-1 should return the interval before our key")
			where, present, _ = rc.Search(UpperLimit32-1, nil)
			So(present, ShouldBeFalse)
			So(where, ShouldEqual, 10)

			var numCompares int
			where, present, numCompares = rc.Search(UpperLimit32, nil)
			So(present, ShouldBeFalse)
			So(where, ShouldEqual, 10)
			p("numCompares = %v", numCompares)
			So(numCompares, ShouldEqual, 3)

			p("confirm that opts SearchOptions to Search reduces search time")
			opts := &SearchOptions{
				StartIndex: 5,
			}
			where, present, numCompares = rc.Search(UpperLimit32, opts)
			So(present, ShouldBeFalse)
			So(where, ShouldEqual, 10)
			p("numCompares = %v", numCompares)
			So(numCompares, ShouldEqual, 2)

			p("confirm that opts SearchOptions to Search is respected")
			where, present, _ = rc.Search(UpperLimit32-3, opts)
			So(present, ShouldBeTrue)
			So(where, ShouldEqual, 10)

			// with the bound in place, UpperLimit32-3 should not be found
			opts.EndxIndex = 10
			where, present, _ = rc.Search(UpperLimit32-3, opts)
			So(present, ShouldBeFalse)
			So(where, ShouldEqual, 9)

		}
	})

}

func TestRleIntersection(t *testing.T) {

	Convey("RunContainer.Intersect of two RunContainers should return their intersection", t, func() {
		{
			vals := []uint32{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, UpperLimit32 - 3, UpperLimit32 - 1}

			a := newRunContainer32FromVals(true, vals[:5]...)
			b := newRunContainer32FromVals(true, vals[2:]...)

			p("a is %v", a)
			p("b is %v", b)

			So(haveOverlap32(interval32{0, 3}, interval32{2, 3}), ShouldBeTrue)
			So(haveOverlap32(interval32{0, 3}, interval32{3, 4}), ShouldBeFalse)

			isect := a.Intersect(b)

			p("isect is %v", isect)

			So(isect.Cardinality(), ShouldEqual, 3)
			So(isect.Get(4), ShouldBeTrue)
			So(isect.Get(6), ShouldBeTrue)
			So(isect.Get(8), ShouldBeTrue)

			d := newRunContainer32TakeOwnership([]interval32{{start: 0, endx: UpperLimit32}})

			isect = isect.Intersect(d)
			p("isect is %v", isect)
			So(isect.Cardinality(), ShouldEqual, 3)
			So(isect.Get(4), ShouldBeTrue)
			So(isect.Get(6), ShouldBeTrue)
			So(isect.Get(8), ShouldBeTrue)

			p("test breaking apart intervals")
			e := newRunContainer32TakeOwnership([]interval32{{2, 5}, {8, 10}, {14, 17}, {20, 23}})
			f := newRunContainer32TakeOwnership([]interval32{{3, 19}, {22, 24}})

			p("e = %v", e)
			p("f = %v", f)

			{
				isect = e.Intersect(f)
				p("isect of e and f is %v", isect)
				So(isect.Cardinality(), ShouldEqual, 8)
				So(isect.Get(3), ShouldBeTrue)
				So(isect.Get(4), ShouldBeTrue)
				So(isect.Get(8), ShouldBeTrue)
				So(isect.Get(9), ShouldBeTrue)
				So(isect.Get(14), ShouldBeTrue)
				So(isect.Get(15), ShouldBeTrue)
				So(isect.Get(16), ShouldBeTrue)
				So(isect.Get(22), ShouldBeTrue)
			}

			{
				// check for symmetry
				isect = f.Intersect(e)
				p("isect of f and e is %v", isect)
				So(isect.Cardinality(), ShouldEqual, 8)
				So(isect.Get(3), ShouldBeTrue)
				So(isect.Get(4), ShouldBeTrue)
				So(isect.Get(8), ShouldBeTrue)
				So(isect.Get(9), ShouldBeTrue)
				So(isect.Get(14), ShouldBeTrue)
				So(isect.Get(15), ShouldBeTrue)
				So(isect.Get(16), ShouldBeTrue)
				So(isect.Get(22), ShouldBeTrue)
			}

		}
	})
}

type trial struct {
	n           int
	percentFill float64
	ntrial      int

	// only in the Union test
	percentDelete float64
}

func TestRleRandomIntersection(t *testing.T) {

	Convey("RunContainer.Intersect of two RunContainers should return their intersection, and this should hold over randomized container content when compared to intersection done with hash maps", t, func() {

		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 100, percentFill: .80, ntrial: 10},
			trial{n: 1000, percentFill: .20, ntrial: 20},
			trial{n: 10000, percentFill: .01, ntrial: 10},
			trial{n: 1000, percentFill: .99, ntrial: 10},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRleRandomIntersection on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint32{}
				b := []uint32{}

				var first, second int

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint32(r0))
					ma[r0] = true
					if i == 0 {
						first = r0
						second = r0 + 1
						p("i is 0, so appending also to a the r0+1 == %v value", second)
						a = append(a, uint32(second))
						ma[second] = true
					}

					r1 := rand.Intn(n)
					b = append(b, uint32(r1))
					mb[r1] = true
				}

				// print a; very likely it has dups
				sort.Sort(uint32Slice(a))
				stringA := ""
				for i := range a {
					stringA += fmt.Sprintf("%v, ", a[i])
				}
				p("a is '%v'", stringA)

				// hash version of intersect:
				hashi := make(map[int]bool)
				for k := range ma {
					if mb[k] {
						hashi[k] = true
					}
				}

				// RunContainer's Intersect
				brle := newRunContainer32FromVals(false, b...)

				//arle := newRunContainer32FromVals(false, a...)
				// instead of the above line, create from array
				// get better test coverage:
				arr := newArrayContainerRange(int(first), int(second))
				arle := newRunContainer32FromArray(arr)
				p("after newRunContainer32FromArray(arr), arle is %v", arle)
				arle.Set(false, a...)
				p("after Set(false, a), arle is %v", arle)

				p("arle is %v", arle)
				p("brle is %v", brle)

				isect := arle.Intersect(brle)

				p("isect is %v", isect)

				//showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in isect", k)
					So(isect.Get(uint32(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: isect is %v, len(hashi) is %v", isect.Cardinality(), len(hashi))
				So(isect.Cardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized Intersect() checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRleRandomUnion(t *testing.T) {

	Convey("RunContainer.Union of two RunContainers should return their union, and this should hold over randomized container content when compared to union done with hash maps", t, func() {

		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 100, percentFill: .80, ntrial: 10},
			trial{n: 1000, percentFill: .20, ntrial: 20},
			trial{n: 10000, percentFill: .01, ntrial: 10},
			trial{n: 1000, percentFill: .99, ntrial: 10, percentDelete: .04},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRleRandomUnion on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint32{}
				b := []uint32{}

				draw := int(float64(n) * tr.percentFill)
				numDel := int(float64(n) * tr.percentDelete)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint32(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint32(r1))
					mb[r1] = true
				}

				// hash version of union:
				hashu := make(map[int]bool)
				for k := range ma {
					hashu[k] = true
				}
				for k := range mb {
					hashu[k] = true
				}

				//showHash("hashu", hashu)

				// RunContainer's Union
				arle := newRunContainer32()
				for i := range a {
					arle.Add(a[i])
				}
				brle := newRunContainer32()
				brle.Set(false, b...)

				p("arle is %v", arle)
				p("brle is %v", brle)

				union := arle.Union(brle)

				p("union is %v", union)

				p("union.Cardinality(): %v, versus len(hashu): %v", union.Cardinality(), len(hashu))

				un := union.AsSlice()
				sort.Sort(uint32Slice(un))

				for kk, v := range un {
					p("kk:%v, RunContainer.Union has %v, checking hashmap: %v", kk, v, hashu[int(v)])
					_ = kk
					So(hashu[int(v)], ShouldBeTrue)
				}

				for k := range hashu {
					p("hashu has %v, checking in union", k)
					So(union.Get(uint32(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement:")
				So(union.Cardinality(), ShouldEqual, len(hashu))

				// do the deletes, exercising the Remove functionality
				for i := 0; i < numDel; i++ {
					r1 := rand.Intn(len(a))
					goner := a[r1]
					union.Remove(goner)
					delete(hashu, int(goner))
				}
				// verify the same as in the hashu
				So(union.Cardinality(), ShouldEqual, len(hashu))
				for k := range hashu {
					p("hashu has %v, checking in union", k)
					So(union.Get(uint32(k)), ShouldBeTrue)
				}

			}
			p("done with randomized Union() checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRleAndOrXor(t *testing.T) {

	Convey("RunContainer And, Or, Xor tests", t, func() {
		{
			rc := newRunContainer32TakeOwnership([]interval32{
				{start: 0, endx: 1},
				{start: 2, endx: 3},
				{start: 4, endx: 5},
			})
			b0 := NewBitmap()
			b0.Add(2)
			b0.Add(6)
			b0.Add(8)

			and := rc.And(b0)
			or := rc.Or(b0)
			xor := rc.Xor(b0)

			So(and.GetCardinality(), ShouldEqual, 1)
			So(or.GetCardinality(), ShouldEqual, 5)
			So(xor.GetCardinality(), ShouldEqual, 4)

			// test creating size 0 and 1 from array
			arr := newArrayContainerCapacity(0)
			empty := newRunContainer32FromArray(arr)
			onceler := newArrayContainerCapacity(1)
			onceler.content = append(onceler.content, uint16(0))
			oneZero := newRunContainer32FromArray(onceler)
			So(empty.Cardinality(), ShouldEqual, 0)
			So(oneZero.Cardinality(), ShouldEqual, 1)
			So(empty.And(b0).GetCardinality(), ShouldEqual, 0)
			So(empty.Or(b0).GetCardinality(), ShouldEqual, 3)

			// exercise newRunContainer32FromVals() with 0 and 1 inputs.
			empty2 := newRunContainer32FromVals(false, []uint32{}...)
			So(empty2.Cardinality(), ShouldEqual, 0)
			one2 := newRunContainer32FromVals(false, []uint32{1}...)
			So(one2.Cardinality(), ShouldEqual, 1)
		}
	})
}

func TestRlePanics(t *testing.T) {

	Convey("Some RunContainer calls/methods should panic if misused", t, func() {

		// newRunContainer32FromVals
		So(func() { newRunContainer32FromVals(true, 1, 0) }, ShouldPanic)

		arr := newArrayContainerRange(1, 3)
		arr.content = []uint16{2, 3, 3, 2, 1}
		So(func() { newRunContainer32FromArray(arr) }, ShouldPanic)
	})
}

func TestRleCoverageOddsAndEnds(t *testing.T) {

	Convey("Some RunContainer code paths that don't otherwise get coverage -- these should be tested to increase percentage of code coverage in testing", t, func() {

		// p() code path
		cur := rleVerbose
		rleVerbose = true
		p("")
		rleVerbose = cur

		// RunContainer.String()
		rc := &runContainer32{}
		So(rc.String(), ShouldEqual, "runContainer32{}")
		rc.iv = make([]interval32, 1)
		rc.iv[0] = interval32{start: 3, endx: 5}
		So(rc.String(), ShouldEqual, "runContainer32{0:[3, 5), }")

		a := interval32{start: 5, endx: 10}
		b := interval32{start: 0, endx: 2}
		c := interval32{start: 1, endx: 3}

		// intersectInterval32s(a, b interval32)
		isect, isEmpty := intersectInterval32s(a, b)
		So(isEmpty, ShouldBeTrue)
		So(isect.runlen(), ShouldEqual, 0)
		isect, isEmpty = intersectInterval32s(b, c)
		So(isEmpty, ShouldBeFalse)
		So(isect.runlen(), ShouldEqual, 1)

		// runContainer32.Union
		{
			ra := newRunContainer32FromVals(false, 4, 5)
			rb := newRunContainer32FromVals(false, 4, 6, 8, 9, 10)
			ra.Union(rb)
			So(rb.indexOfIntervalAtOrAfter(4, 2), ShouldEqual, 2)
			So(rb.indexOfIntervalAtOrAfter(3, 2), ShouldEqual, 2)
		}

		// runContainer.Intersect
		{
			ra := newRunContainer32()
			rb := newRunContainer32()
			So(ra.Intersect(rb).Cardinality(), ShouldEqual, 0)
		}
		{
			ra := newRunContainer32FromVals(false, 1)
			rb := newRunContainer32FromVals(false, 4)
			So(ra.Intersect(rb).Cardinality(), ShouldEqual, 0)
		}

		// runContainer.Add
		{
			ra := newRunContainer32FromVals(false, 1)
			rb := newRunContainer32FromVals(false, 4)
			So(ra.Cardinality(), ShouldEqual, 1)
			So(rb.Cardinality(), ShouldEqual, 1)
			ra.Add(5)
			So(ra.Cardinality(), ShouldEqual, 2)

			// NewRunIterator32()
			empty := newRunContainer32()
			it := empty.NewRunIterator32()
			So(func() { it.Next() }, ShouldPanic)
			it2 := ra.NewRunIterator32()
			it2.curIndex = len(it2.rc.iv)
			So(func() { it2.Next() }, ShouldPanic)

			// RunIterator32.Remove()
			emptyIt := empty.NewRunIterator32()
			So(func() { emptyIt.Remove() }, ShouldPanic)

			// newRunContainer32FromArray
			arr := newArrayContainerRange(1, 6)
			arr.content = []uint16{5, 5, 5, 6, 9}
			rc3 := newRunContainer32FromArray(arr)
			So(rc3.Cardinality(), ShouldEqual, 3)

			// runContainer32SerializedSizeInBytes
			// runContainer32.SerializedSizeInBytes
			_ = runContainer32SerializedSizeInBytes(3)
			_ = rc3.serializedSizeInBytes()

			// findNextIntervalThatIntersectsStartingFrom
			idx, _ := rc3.findNextIntervalThatIntersectsStartingFrom(0, 100)
			So(idx, ShouldEqual, 1)

			// deleteAt / Remove
			rc3.Add(10)
			rc3.Remove(10)
			rc3.Remove(9)
			So(rc3.Cardinality(), ShouldEqual, 2)
			rc3.Add(9)
			rc3.Add(10)
			rc3.Add(12)
			So(rc3.Cardinality(), ShouldEqual, 5)
			it3 := rc3.NewRunIterator32()
			it3.Next()
			it3.Next()
			it3.Next()
			it3.Next()
			rleVerbose = true
			//p("rc3 = %v", rc3) // 5, 6, 9, 10, 12
			So(it3.Cur(), ShouldEqual, uint32(10))
			it3.Remove()
			//p("after Remove of 10, rc3 = %v", rc3) // 5, 6, 9, 12
			So(it3.Next(), ShouldEqual, uint32(12))
		}
	})
}

func showHash(name string, h map[int]bool) {
	hv := []int{}
	for k := range h {
		hv = append(hv, k)
	}
	sort.Sort(sort.IntSlice(hv))
	stringH := ""
	for i := range hv {
		stringH += fmt.Sprintf("%v, ", hv[i])
	}

	p("%s is (len %v): %s", name, len(h), stringH)
}

package roaring

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"sort"
	"testing"
)

func TestRle16RandomIntersectAgainstOtherContainers010(t *testing.T) {

	Convey("runContainer16 `and` operation against other container types should correctly do the intersection", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 100, percentFill: .95, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRleRandomAndAgainstOtherContainers on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of intersect:
				hashi := make(map[int]bool)
				for k := range ma {
					if mb[k] {
						hashi[k] = true
					}
				}

				// RunContainer's Intersect
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)

				// vs bitmapContainer
				bc := newBitmapContainer()
				for _, bv := range b {
					bc.iadd(bv)
				}

				// vs arrayContainer
				ac := newArrayContainer()
				for _, bv := range b {
					ac.iadd(bv)
				}

				// vs runContainer
				rcb := newRunContainer16FromVals(false, b...)

				rc_vs_bc_isect := rc.and(bc)
				rc_vs_ac_isect := rc.and(ac)
				rc_vs_rcb_isect := rc.and(rcb)

				p("rc_vs_bc_isect is %v", rc_vs_bc_isect)
				p("rc_vs_ac_isect is %v", rc_vs_ac_isect)
				p("rc_vs_rcb_isect is %v", rc_vs_rcb_isect)

				showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in rc_vs_bc_isect", k)
					So(rc_vs_bc_isect.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_ac_isect", k)
					So(rc_vs_ac_isect.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_rcb_isect", k)
					So(rc_vs_rcb_isect.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc_vs_bc_isect is %v, len(hashi) is %v", rc_vs_bc_isect.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_ac_isect is %v, len(hashi) is %v", rc_vs_ac_isect.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_rcb_isect is %v, len(hashi) is %v", rc_vs_rcb_isect.getCardinality(), len(hashi))
				So(rc_vs_bc_isect.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_ac_isect.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_rcb_isect.getCardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized and() vs bitmapContainer and arrayContainer checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16RandomUnionAgainstOtherContainers011(t *testing.T) {

	Convey("runContainer16 `or` operation against other container types should correctly do the intersection", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 100, percentFill: .95, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRleRandomAndAgainstOtherContainers on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of union
				hashi := make(map[int]bool)
				for k := range ma {
					hashi[k] = true
				}
				for k := range mb {
					hashi[k] = true
				}

				// RunContainer's 'or'
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)

				// vs bitmapContainer
				bc := newBitmapContainer()
				for _, bv := range b {
					bc.iadd(bv)
				}

				// vs arrayContainer
				ac := newArrayContainer()
				for _, bv := range b {
					ac.iadd(bv)
				}

				// vs runContainer
				rcb := newRunContainer16FromVals(false, b...)

				rc_vs_bc_union := rc.or(bc)
				rc_vs_ac_union := rc.or(ac)
				rc_vs_rcb_union := rc.or(rcb)

				p("rc_vs_bc_union is %v", rc_vs_bc_union)
				p("rc_vs_ac_union is %v", rc_vs_ac_union)
				p("rc_vs_rcb_union is %v", rc_vs_rcb_union)

				showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in rc_vs_bc_union", k)
					So(rc_vs_bc_union.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_ac_union", k)
					So(rc_vs_ac_union.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_rcb_union", k)
					So(rc_vs_rcb_union.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc_vs_bc_union is %v, len(hashi) is %v", rc_vs_bc_union.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_ac_union is %v, len(hashi) is %v", rc_vs_ac_union.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_rcb_union is %v, len(hashi) is %v", rc_vs_rcb_union.getCardinality(), len(hashi))
				So(rc_vs_bc_union.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_ac_union.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_rcb_union.getCardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized or() vs bitmapContainer and arrayContainer checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16RandomInplaceUnionAgainstOtherContainers012(t *testing.T) {

	Convey("runContainer16 `ior` inplace union operation against other container types should correctly do the intersection", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 10, percentFill: .95, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRleRandomInplaceUnionAgainstOtherContainers on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of union
				hashi := make(map[int]bool)
				for k := range ma {
					hashi[k] = true
				}
				for k := range mb {
					hashi[k] = true
				}

				// RunContainer's 'or'
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)

				rc_vs_bc_union := rc.Clone()
				rc_vs_ac_union := rc.Clone()
				rc_vs_rcb_union := rc.Clone()

				// vs bitmapContainer
				bc := newBitmapContainer()
				for _, bv := range b {
					bc.iadd(bv)
				}

				// vs arrayContainer
				ac := newArrayContainer()
				for _, bv := range b {
					ac.iadd(bv)
				}

				// vs runContainer
				rcb := newRunContainer16FromVals(false, b...)

				rc_vs_bc_union.ior(bc)
				rc_vs_ac_union.ior(ac)
				rc_vs_rcb_union.ior(rcb)

				p("rc_vs_bc_union is %v", rc_vs_bc_union)
				p("rc_vs_ac_union is %v", rc_vs_ac_union)
				p("rc_vs_rcb_union is %v", rc_vs_rcb_union)

				showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in rc_vs_bc_union", k)
					So(rc_vs_bc_union.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_ac_union", k)
					So(rc_vs_ac_union.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_rcb_union", k)
					So(rc_vs_rcb_union.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc_vs_bc_union is %v, len(hashi) is %v", rc_vs_bc_union.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_ac_union is %v, len(hashi) is %v", rc_vs_ac_union.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_rcb_union is %v, len(hashi) is %v", rc_vs_rcb_union.getCardinality(), len(hashi))
				So(rc_vs_bc_union.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_ac_union.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_rcb_union.getCardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized or() vs bitmapContainer and arrayContainer checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16RandomInplaceIntersectAgainstOtherContainers014(t *testing.T) {

	Convey("runContainer16 `iand` inplace-and operation against other container types should correctly do the intersection", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 100, percentFill: .95, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRleRandomAndAgainstOtherContainers on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of intersect:
				hashi := make(map[int]bool)
				for k := range ma {
					if mb[k] {
						hashi[k] = true
					}
				}

				// RunContainer's Intersect
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)

				// vs bitmapContainer
				bc := newBitmapContainer()
				for _, bv := range b {
					bc.iadd(bv)
				}

				// vs arrayContainer
				ac := newArrayContainer()
				for _, bv := range b {
					ac.iadd(bv)
				}

				// vs runContainer
				rcb := newRunContainer16FromVals(false, b...)

				rc_vs_bc_isect := rc.Clone()
				rc_vs_ac_isect := rc.Clone()
				rc_vs_rcb_isect := rc.Clone()

				rc_vs_bc_isect.iand(bc)
				rc_vs_ac_isect.iand(ac)
				rc_vs_rcb_isect.iand(rcb)

				p("rc_vs_bc_isect is %v", rc_vs_bc_isect)
				p("rc_vs_ac_isect is %v", rc_vs_ac_isect)
				p("rc_vs_rcb_isect is %v", rc_vs_rcb_isect)

				showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in rc_vs_bc_isect", k)
					So(rc_vs_bc_isect.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_ac_isect", k)
					So(rc_vs_ac_isect.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_rcb_isect", k)
					So(rc_vs_rcb_isect.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc_vs_bc_isect is %v, len(hashi) is %v", rc_vs_bc_isect.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_ac_isect is %v, len(hashi) is %v", rc_vs_ac_isect.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_rcb_isect is %v, len(hashi) is %v", rc_vs_rcb_isect.getCardinality(), len(hashi))
				So(rc_vs_bc_isect.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_ac_isect.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_rcb_isect.getCardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized and() vs bitmapContainer and arrayContainer checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16RemoveApi015(t *testing.T) {

	Convey("runContainer16 `remove` (a minus b) should work", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 100, percentFill: .95, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRle16RemoveApi015 on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of remove:
				hashrm := make(map[int]bool)
				for k := range ma {
					hashrm[k] = true
				}
				for k := range mb {
					delete(hashrm, k)
				}

				// RunContainer's remove
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a, pre-remove, is %v", rc)

				for k := range mb {
					rc.iremove(uint16(k))
				}

				p("rc from a, post-iremove, is %v", rc)

				showHash("correct answer is hashrm", hashrm)

				for k := range hashrm {
					p("hashrm has %v, checking in rc", k)
					So(rc.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc is %v, len(hashrm) is %v", rc.getCardinality(), len(hashrm))
				So(rc.getCardinality(), ShouldEqual, len(hashrm))
			}
			p("done with randomized remove() checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func showArray16(a []uint16, name string) {
	sort.Sort(uint16Slice(a))
	stringA := ""
	for i := range a {
		stringA += fmt.Sprintf("%v, ", a[i])
	}
	p("%s is '%v'", name, stringA)
}

func TestRle16RandomAndNot16(t *testing.T) {

	Convey("runContainer16 `andNot` operation against other container types should correctly do the and-not operation", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 1000, percentFill: .95, ntrial: 2},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRle16RandomAndNot16 on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of and-not
				hashi := make(map[int]bool)
				for k := range ma {
					hashi[k] = true
				}
				for k := range mb {
					delete(hashi, k)
				}

				// RunContainer's and-not
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)

				// vs bitmapContainer
				bc := newBitmapContainer()
				for _, bv := range b {
					bc.iadd(bv)
				}

				// vs arrayContainer
				ac := newArrayContainer()
				for _, bv := range b {
					ac.iadd(bv)
				}

				// vs runContainer
				rcb := newRunContainer16FromVals(false, b...)

				rc_vs_bc_andnot := rc.andNot(bc)
				rc_vs_ac_andnot := rc.andNot(ac)
				rc_vs_rcb_andnot := rc.andNot(rcb)

				p("rc_vs_bc_andnot is %v", rc_vs_bc_andnot)
				p("rc_vs_ac_andnot is %v", rc_vs_ac_andnot)
				p("rc_vs_rcb_andnot is %v", rc_vs_rcb_andnot)

				showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in rc_vs_bc_andnot", k)
					So(rc_vs_bc_andnot.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_ac_andnot", k)
					So(rc_vs_ac_andnot.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_rcb_andnot", k)
					So(rc_vs_rcb_andnot.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc_vs_bc_andnot is %v, len(hashi) is %v", rc_vs_bc_andnot.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_ac_andnot is %v, len(hashi) is %v", rc_vs_ac_andnot.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_rcb_andnot is %v, len(hashi) is %v", rc_vs_rcb_andnot.getCardinality(), len(hashi))
				So(rc_vs_bc_andnot.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_ac_andnot.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_rcb_andnot.getCardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized andNot() vs bitmapContainer and arrayContainer checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16RandomInplaceAndNot017(t *testing.T) {

	Convey("runContainer16 `iandNot` operation against other container types should correctly do the inplace-and-not operation", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 1000, percentFill: .95, ntrial: 2},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRle16RandomAndNot16 on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				showArray16(a, "a")
				showArray16(b, "b")

				// hash version of and-not
				hashi := make(map[int]bool)
				for k := range ma {
					hashi[k] = true
				}
				for k := range mb {
					delete(hashi, k)
				}

				// RunContainer's and-not
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)

				// vs bitmapContainer
				bc := newBitmapContainer()
				for _, bv := range b {
					bc.iadd(bv)
				}

				// vs arrayContainer
				ac := newArrayContainer()
				for _, bv := range b {
					ac.iadd(bv)
				}

				// vs runContainer
				rcb := newRunContainer16FromVals(false, b...)

				rc_vs_bc_iandnot := rc.Clone()
				rc_vs_ac_iandnot := rc.Clone()
				rc_vs_rcb_iandnot := rc.Clone()

				rc_vs_bc_iandnot.iandNot(bc)
				rc_vs_ac_iandnot.iandNot(ac)
				rc_vs_rcb_iandnot.iandNot(rcb)

				p("rc_vs_bc_iandnot is %v", rc_vs_bc_iandnot)
				p("rc_vs_ac_iandnot is %v", rc_vs_ac_iandnot)
				p("rc_vs_rcb_iandnot is %v", rc_vs_rcb_iandnot)

				showHash("hashi", hashi)

				for k := range hashi {
					p("hashi has %v, checking in rc_vs_bc_iandnot", k)
					So(rc_vs_bc_iandnot.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_ac_iandnot", k)
					So(rc_vs_ac_iandnot.contains(uint16(k)), ShouldBeTrue)

					p("hashi has %v, checking in rc_vs_rcb_iandnot", k)
					So(rc_vs_rcb_iandnot.contains(uint16(k)), ShouldBeTrue)
				}

				p("checking for cardinality agreement: rc_vs_bc_iandnot is %v, len(hashi) is %v", rc_vs_bc_iandnot.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_ac_iandnot is %v, len(hashi) is %v", rc_vs_ac_iandnot.getCardinality(), len(hashi))
				p("checking for cardinality agreement: rc_vs_rcb_iandnot is %v, len(hashi) is %v", rc_vs_rcb_iandnot.getCardinality(), len(hashi))
				So(rc_vs_bc_iandnot.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_ac_iandnot.getCardinality(), ShouldEqual, len(hashi))
				So(rc_vs_rcb_iandnot.getCardinality(), ShouldEqual, len(hashi))
			}
			p("done with randomized andNot() vs bitmapContainer and arrayContainer checks for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16InversionOfIntervals018(t *testing.T) {

	Convey("runContainer `invert` operation should do a NOT on the set of intervals, in-place", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 1000, percentFill: .90, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRle16InversinoOfIntervals018 on check# j=%v", j)
				ma := make(map[int]bool)
				hashNotA := make(map[int]bool)

				n := tr.n
				a := []uint16{}

				// hashNotA will be NOT ma
				//for i := 0; i < n; i++ {
				for i := 0; i < MaxUint16+1; i++ {
					hashNotA[i] = true
				}

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true
					delete(hashNotA, r0)
				}

				//showArray16(a, "a")
				// too big to print: showHash("hashNotA is not a:", hashNotA)

				// RunContainer's invert
				rc := newRunContainer16FromVals(false, a...)

				p("rc from a is %v", rc)
				p("rc.cardinality = %v", rc.cardinality())
				inv := rc.invert()

				p("inv of a (card=%v) is %v", inv.cardinality(), inv)

				So(inv.cardinality(), ShouldEqual, 1+MaxUint16-rc.cardinality())

				for k := 0; k < n; k++ {
					if hashNotA[k] {
						//p("hashNotA has %v, checking inv", k)
						So(inv.contains(uint16(k)), ShouldBeTrue)
					}
				}

				// skip for now, too big to do 2^16-1
				p("checking for cardinality agreement: inv is %v, len(hashNotA) is %v", inv.getCardinality(), len(hashNotA))
				So(inv.getCardinality(), ShouldEqual, len(hashNotA))
			}
			p("done with randomized invert() check for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

func TestRle16SubtractionOfIntervals019(t *testing.T) {

	Convey("runContainer `subtract` operation removes an interval in-place", t, func() {
		seed := int64(42)
		p("seed is %v", seed)
		rand.Seed(seed)

		trials := []trial{
			trial{n: 10, percentFill: .90, ntrial: 1},
		}

		tester := func(tr trial) {
			for j := 0; j < tr.ntrial; j++ {
				p("TestRle16SubtractionOfIntervals019 on check# j=%v", j)
				ma := make(map[int]bool)
				mb := make(map[int]bool)

				n := tr.n
				a := []uint16{}
				b := []uint16{}

				// hashAminusB will be  ma - mb
				hashAminusB := make(map[int]bool)

				draw := int(float64(n) * tr.percentFill)
				for i := 0; i < draw; i++ {
					r0 := rand.Intn(n)
					a = append(a, uint16(r0))
					ma[r0] = true
					hashAminusB[r0] = true

					r1 := rand.Intn(n)
					b = append(b, uint16(r1))
					mb[r1] = true
				}

				for k := range mb {
					delete(hashAminusB, k)
				}

				rleVerbose = true
				showHash("hash a is:", ma)
				showHash("hash b is:", mb)
				showHash("hashAminusB is:", hashAminusB)

				// RunContainer's subtract A - B
				rc := newRunContainer16FromVals(false, a...)
				rcb := newRunContainer16FromVals(false, b...)

				p("rc from a is %v", rc)
				p("rc.cardinality = %v", rc.cardinality())
				p("rcb from b is %v", rcb)
				p("rcb.cardinality = %v", rcb.cardinality())
				it := rcb.NewRunIterator16()
				for it.HasNext() {
					nx := it.Next()
					rc.isubtract(interval16{start: nx, last: nx})
				}

				p("rc = a - b; has (card=%v), is %v", rc.cardinality(), rc)

				for k := range hashAminusB {
					p("hashAminusB has element %v, checking rc (which is A - B)", k)
					So(rc.contains(uint16(k)), ShouldBeTrue)
				}
				p("checking for cardinality agreement: sub is %v, len(hashAminusB) is %v", rc.getCardinality(), len(hashAminusB))
				So(rc.getCardinality(), ShouldEqual, len(hashAminusB))
			}
			p("done with randomized subtract() check for trial %#v", tr)
		}

		for i := range trials {
			tester(trials[i])
		}

	})
}

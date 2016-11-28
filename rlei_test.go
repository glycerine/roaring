package roaring

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"sort"
	"testing"
)

func TestRle16RandomIntersectAgainstOtherContainers010(t *testing.T) {

	Convey("runIterator16 `and` operation against other container types should correctly do the intersection", t, func() {
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

	Convey("runIterator16 `or` operation against other container types should correctly do the intersection", t, func() {
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

	Convey("runIterator16 `ior` inplace union operation against other container types should correctly do the intersection", t, func() {
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

	Convey("runIterator16 `iand` inplace-and operation against other container types should correctly do the intersection", t, func() {
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

	Convey("runIterator16 `remove` (a minus b) should work", t, func() {
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

	Convey("runIterator16 `andNot` operation against other container types should correctly do the and-not operation", t, func() {
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

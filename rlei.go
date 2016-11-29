package roaring

///////////////////////////////////////////////////
//
// container interface methods for runContainer16
//
///////////////////////////////////////////////////

// compile time verify we meet interface requirements
var _ container = &runContainer16{}

func (rc *runContainer16) clone() container {
	return newRunContainer16CopyIv(rc.iv)
}

func (rc *runContainer16) and(a container) container {
	switch c := a.(type) {
	case *runContainer16:
		return rc.intersect(c)
	case *arrayContainer:
		return rc.andArray(c)
	case *bitmapContainer:
		return rc.andBitmapContainer(c)
	}
	panic("unsupported container type")
}

// andBitmapContainer finds the intersection of rc and b.
func (rc *runContainer16) andBitmapContainer(bc *bitmapContainer) container {
	out := newRunContainer16()
	for _, p := range rc.iv {
		for i := p.start; i <= p.last; i++ {
			if bc.contains(i) {
				out.Add(i)
			}
		}
	}
	return out
}

func (rc *runContainer16) andArray(ac *arrayContainer) container {
	out := newRunContainer16()
	for _, p := range rc.iv {
		for i := p.start; i <= p.last; i++ {
			if ac.contains(i) {
				out.Add(i)
			}
		}
	}
	return out
}

func (rc *runContainer16) iand(a container) container {
	switch c := a.(type) {
	case *runContainer16:
		return rc.inplaceIntersect(c)
	case *arrayContainer:
		return rc.iandArray(c)
	case *bitmapContainer:
		return rc.iandBitmapContainer(c)
	}
	panic("unsupported container type")
}

func (rc *runContainer16) inplaceIntersect(rc2 *runContainer16) container {
	// TODO: optimize by doing less allocation, possibly?
	sect := rc.intersect(rc2)
	*rc = *sect
	return rc
}

func (rc *runContainer16) iandBitmapContainer(bc *bitmapContainer) container {
	// TODO: optimize by doing less allocation, possibly?
	out := newRunContainer16()
	for _, p := range rc.iv {
		for i := p.start; i <= p.last; i++ {
			if bc.contains(i) {
				out.Add(i)
			}
		}
	}
	*rc = *out
	return rc
}

func (rc *runContainer16) iandArray(ac *arrayContainer) container {
	// TODO: optimize by doing less allocation, possibly?
	out := newRunContainer16()
	for _, p := range rc.iv {
		for i := p.start; i <= p.last; i++ {
			if ac.contains(i) {
				out.Add(i)
			}
		}
	}
	*rc = *out
	return rc
}

func (rc *runContainer16) andNot(a container) container {
	switch c := a.(type) {
	case *arrayContainer:
		return rc.andNotArray(c)
	case *bitmapContainer:
		return rc.andNotBitmap(c)
	case *runContainer16:
		return rc.andNotRunContainer16(c)
	}
	panic("unsupported container type")
}

func (rc *runContainer16) fillLeastSignificant16bits(x []uint32, i int, mask uint32) {
	k := 0
	var val int64
	for _, p := range rc.iv {
		n := p.runlen()
		for j := int64(0); j < n; j++ {
			val = int64(p.start) + j
			x[k+i] = uint32(val) | mask
			k++
		}
	}
}

func (rc *runContainer16) getShortIterator() shortIterable {
	return rc.NewRunIterator16()
}

func (rc *runContainer16) getSizeInBytes() int {
	return baseRc16Size + perIntervalRc16Size*len(rc.iv)
}

// add the values in the range [firstOfRange,lastofRange). lastofRange
// is still abe to express 2^16 because it is an int not an uint16.
func (rc *runContainer16) iaddRange(firstOfRange, lastOfRange int) container {
	addme := newRunContainer16TakeOwnership([]interval16{
		{
			start: uint16(firstOfRange),
			last:  uint16(lastOfRange - 1),
		},
	})
	*rc = *rc.union(addme)
	return rc
}

// remove the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) iremoveRange(firstOfRange, lastOfRange int) container {
	x := interval16{start: uint16(firstOfRange), last: uint16(lastOfRange - 1)}
	rc.isubtract(x)
	return rc
}

// not flip the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) not(firstOfRange, lastOfRange int) container {
	return rc.Not(firstOfRange, lastOfRange)
}

// Not flips the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) Not(firstOfRange, lastOfRange int) *runContainer16 {

	if firstOfRange >= lastOfRange {
		return rc.Clone()
	}
	x := interval16{start: uint16(firstOfRange), last: uint16(lastOfRange - 1)}
	xs := []interval16{x}

	isect := rc.intersect(newRunContainer16TakeOwnership(xs))
	rc2 := rc.Clone()
	rc2.isubtract(x)
	invertedIsect := isect.invert()
	rc2 = rc2.union(invertedIsect)
	return rc2
}

// equals is not logical equals,  apparently from the array implementation,
// equals also requires the same type of container.
func (rc *runContainer16) equals(o interface{}) bool {
	srb, ok := o.(*runContainer16)
	if ok {
		// Check if the containers are the same object.
		if rc == srb {
			return true
		}

		if len(srb.iv) != len(rc.iv) {
			return false
		}

		for i, v := range rc.iv {
			if v != srb.iv[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (rc *runContainer16) iaddReturnMinimized(x uint16) container {
	rc.Add(x)
	return rc
}

func (rc *runContainer16) iadd(x uint16) (wasNew bool) {
	return rc.Add(x)
}

func (rc *runContainer16) iremoveReturnMinimized(x uint16) container {
	rc.removeKey(x)
	return rc
}

func (rc *runContainer16) iremove(x uint16) bool {
	return rc.removeKey(x)
}

func (rc *runContainer16) or(a container) container {
	switch c := a.(type) {
	case *runContainer16:
		return rc.union(c)
	case *arrayContainer:
		return rc.orArray(c)
	case *bitmapContainer:
		return rc.orBitmapContainer(c)
	}
	panic("unsupported container type")
}

// orBitmapContainer finds the union of rc and bc.
func (rc *runContainer16) orBitmapContainer(bc *bitmapContainer) container {
	out := bc.clone()
	for _, p := range rc.iv {
		for i := p.start; i <= p.last; i++ {
			out.iadd(i)
		}
	}
	return out
}

// orArray finds the union of rc and ac.
func (rc *runContainer16) orArray(ac *arrayContainer) container {
	out := ac.clone()
	for _, p := range rc.iv {
		for i := p.start; i <= p.last; i++ {
			out.iadd(i)
		}
	}
	return out
}

func (rc *runContainer16) ior(a container) container {
	switch c := a.(type) {
	case *runContainer16:
		return rc.inplaceUnion(c)
	case *arrayContainer:
		return rc.iorArray(c)
	case *bitmapContainer:
		return rc.iorBitmapContainer(c)
	}
	panic("unsupported container type")
}

func (rc *runContainer16) inplaceUnion(rc2 *runContainer16) container {
	for _, p := range rc2.iv {
		for i := p.start; i <= p.last; i++ {
			rc.Add(i)
		}
	}
	return rc
}

func (rc *runContainer16) iorBitmapContainer(bc *bitmapContainer) container {

	it := bc.getShortIterator()
	for it.hasNext() {
		rc.Add(it.next())
	}
	return rc
}

func (rc *runContainer16) iorArray(ac *arrayContainer) container {
	it := ac.getShortIterator()
	for it.hasNext() {
		rc.Add(it.next())
	}
	return rc
}

func (rc *runContainer16) lazyIOR(a container) container {
	// TODO part of container interface, must implement.
	/*	switch a.(type) {
		case *runContainer16:
			return ac.lazyorArray(a.(*runContainer16))
		case *bitmapContainer:
			return a.lazyOR(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (rc *runContainer16) lazyOR(a container) container {
	// TODO part of container interface, must implement.
	/*	switch a.(type) {
		case *runContainer16:
			return ac.lazyorArray(a.(*runContainer16))
		case *bitmapContainer:
			return a.lazyOR(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (rc *runContainer16) intersects(a container) bool {
	// TODO: optimize by doing inplace/less allocation, possibly?
	isect := rc.and(a)
	return isect.getCardinality() > 0
}

func (rc *runContainer16) xor(a container) container {
	switch c := a.(type) {
	case *arrayContainer:
		return rc.xorArray(c)
	case *bitmapContainer:
		return rc.xorBitmap(c)
	case *runContainer16:
		return rc.xorRunContainer16(c)
	}
	panic("unsupported container type")
}

func (rc *runContainer16) iandNot(a container) container {
	switch c := a.(type) {
	case *arrayContainer:
		return rc.iandNotArray(c)
	case *bitmapContainer:
		return rc.iandNotBitmap(c)
	case *runContainer16:
		return rc.iandNotRunContainer16(c)
	}
	panic("unsupported container type")
}

// flip the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) inot(firstOfRange, lastOfRange int) container {
	// TODO: minimize copies, do it all inplace; not() makes a copy.
	rc = rc.Not(firstOfRange, lastOfRange)
	return rc
}

func (rc *runContainer16) getCardinality() int {
	return int(rc.cardinality())
}

func (rc *runContainer16) rank(x uint16) int {
	n := int64(len(rc.iv))
	xx := int64(x)
	w, already, _ := rc.search(xx, nil)
	if w < 0 {
		return 0
	}
	if w == n-1 {
		return rc.getCardinality()
	}

	var rnk int64
	if !already {
		for i := int64(0); i <= w; i++ {
			rnk += rc.iv[i].runlen()
		}
		return int(rnk)
	}
	for i := int64(0); i < w; i++ {
		rnk += rc.iv[i].runlen()
	}
	rnk += int64(x-rc.iv[w].start) + 1
	return int(rnk)
}

func (rc *runContainer16) selectInt(x uint16) int {
	return rc.selectInt16(x)
}

func (rc *runContainer16) andNotRunContainer16(b *runContainer16) container {
	return rc.AndNotRunContainer16(b)
}

func (rc *runContainer16) andNotArray(ac *arrayContainer) container {
	rcb := rc.toBitmapContainer()
	acb := ac.toBitmapContainer()
	return rcb.andNotBitmap(acb)
}

func (rc *runContainer16) andNotBitmap(bc *bitmapContainer) container {
	rcb := rc.toBitmapContainer()
	return rcb.andNotBitmap(bc)
}

func (rc *runContainer16) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	n := rc.getCardinality()
	bc.cardinality = n
	it := rc.NewRunIterator16()
	for it.HasNext() {
		x := it.Next()
		i := int(x) / 64
		bc.bitmap[i] |= (uint64(1) << uint(x%64))
	}
	return bc
}

func (rc *runContainer16) iandNotRunContainer16(x2 *runContainer16) container {
	rcb := rc.toBitmapContainer()
	x2b := x2.toBitmapContainer()
	rcb.iandNotBitmapSurely(x2b)
	// TODO: check size and optimize the return value
	// TODO: is inplace modification really required? If not, elide the copy.
	rc2 := newRunContainer16FromBitmapContainer(rcb)
	*rc = *rc2
	return rc
}

func (rc *runContainer16) iandNotArray(ac *arrayContainer) container {
	rcb := rc.toBitmapContainer()
	acb := ac.toBitmapContainer()
	rcb.iandNotBitmapSurely(acb)
	// TODO: check size and optimize the return value
	// TODO: is inplace modification really required? If not, elide the copy.
	rc2 := newRunContainer16FromBitmapContainer(rcb)
	*rc = *rc2
	return rc
}

func (rc *runContainer16) iandNotBitmap(bc *bitmapContainer) container {
	rcb := rc.toBitmapContainer()
	rcb.iandNotBitmapSurely(bc)
	// TODO: check size and optimize the return value
	// TODO: is inplace modification really required? If not, elide the copy.
	rc2 := newRunContainer16FromBitmapContainer(rcb)
	*rc = *rc2
	return rc
}

func (rc *runContainer16) xorRunContainer16(x2 *runContainer16) container {
	rcb := rc.toBitmapContainer()
	x2b := x2.toBitmapContainer()
	return rcb.xorBitmap(x2b)
}

func (rc *runContainer16) xorArray(ac *arrayContainer) container {
	rcb := rc.toBitmapContainer()
	acb := ac.toBitmapContainer()
	return rcb.xorBitmap(acb)
}

func (rc *runContainer16) xorBitmap(bc *bitmapContainer) container {
	rcb := rc.toBitmapContainer()
	return rcb.xorBitmap(bc)
}

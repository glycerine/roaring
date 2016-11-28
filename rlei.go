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
		return rc.andNot(c)
	}
	panic("unsupported container type")
}

func (rc *runContainer16) fillLeastSignificant16bits(x []uint32, i int, mask uint32) {
	/*	for k := 0; k < len(ac.content); k++ {
			x[k+i] = uint16(ac.content[k]) | mask
		}
	*/
}

func (rc *runContainer16) getShortIterator() shortIterable {
	// TODO
	return &shortIterator{}
	//	return &shortIterator{ac.content, 0}
}

func (rc *runContainer16) getSizeInBytes() int {
	/*
		// unsafe.Sizeof calculates the memory used by the top level of the slice
		// descriptor - not including the size of the memory referenced by the slice.
		// http://golang.org/pkg/unsafe/#Sizeof
		return ac.getCardinality()*2 + int(unsafe.Sizeof(ac.content))
	*/
	return 0 // TODO
}

// serializedSizeInBytes returns the number of bytes of memory
// required by this runContainer16.
//func (rc *runContainer16) serializedSizeInBytes() int

// add the values in the range [firstOfRange,lastofRange)
func (rc *runContainer16) iaddRange(firstOfRange, lastOfRange int) container {
	/*	if firstOfRange >= lastOfRange {
			return ac
		}
		indexstart := binarySearch(ac.content, uint16(firstOfRange))
		if indexstart < 0 {
			indexstart = -indexstart - 1
		}
		indexend := binarySearch(ac.content, uint16(lastOfRange-1))
		if indexend < 0 {
			indexend = -indexend - 1
		} else {
			indexend++
		}
		rangelength := lastOfRange - firstOfRange
		newcardinality := indexstart + (ac.getCardinality() - indexend) + rangelength
		if newcardinality > arrayDefaultMaxSize {
			a := ac.toBitmapContainer()
			return a.iaddRange(firstOfRange, lastOfRange)
		}
		if cap(ac.content) < newcardinality {
			tmp := make([]uint16, newcardinality, newcardinality)
			copy(tmp[:indexstart], ac.content[:indexstart])
			copy(tmp[indexstart+rangelength:], ac.content[indexend:])

			ac.content = tmp
		} else {
			ac.content = ac.content[:newcardinality]
			copy(ac.content[indexstart+rangelength:], ac.content[indexend:])

		}
		for k := 0; k < rangelength; k++ {
			ac.content[k+indexstart] = uint16(firstOfRange + k)
		}
		return ac
	*/
	return nil
}

// remove the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) iremoveRange(firstOfRange, lastOfRange int) container {
	/*	if firstOfRange >= lastOfRange {
			return ac
		}
		indexstart := binarySearch(ac.content, uint16(firstOfRange))
		if indexstart < 0 {
			indexstart = -indexstart - 1
		}
		indexend := binarySearch(ac.content, uint16(lastOfRange-1))
		if indexend < 0 {
			indexend = -indexend - 1
		} else {
			indexend++
		}
		rangelength := indexend - indexstart
		answer := ac
		copy(answer.content[indexstart:], ac.content[indexstart+rangelength:])
		answer.content = answer.content[:ac.getCardinality()-rangelength]
		return answer
	*/
	return nil

}

// flip the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) not(firstOfRange, lastOfRange int) container {
	/*	if firstOfRange >= lastOfRange {
			return ac.clone()
		}
		return ac.notClose(firstOfRange, lastOfRange-1) // remove everything in [firstOfRange,lastOfRange-1]
	*/
	return nil
}

func (rc *runContainer16) equals(o interface{}) bool {
	/*
		srb, ok := o.(*runContainer16)
		if ok {
			// Check if the containers are the same object.
			if ac == srb {
				return true
			}

			if len(srb.content) != len(ac.content) {
				return false
			}

			for i, v := range ac.content {
				if v != srb.content[i] {
					return false
				}
			}
			return true
		}
		return false
	*/
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
	/*
		switch a.(type) {
		case *runContainer16:
			return ac.intersectsArray(a.(*runContainer16))
		case *bitmapContainer:
			return a.intersects(ac)
		}
		panic("unsupported container type")
	*/
	return false
}

func (rc *runContainer16) xor(a container) container {
	/*	switch a.(type) {
		case *runContainer16:
			return ac.xorArray(a.(*runContainer16))
		case *bitmapContainer:
			return a.xor(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (rc *runContainer16) iandNot(a container) container {
	/*	switch a.(type) {
		case *runContainer16:
			return ac.iandNotArray(a.(*runContainer16))
		case *bitmapContainer:
			return ac.iandNotBitmap(a.(*bitmapContainer))
		}
		panic("unsupported container type")
	*/
	return nil
}

// flip the values in the range [firstOfRange,lastOfRange)
func (rc *runContainer16) inot(firstOfRange, lastOfRange int) container {
	/*	if firstOfRange >= lastOfRange {
			return ac
		}
		return ac.inotClose(firstOfRange, lastOfRange-1) // remove everything in [firstOfRange,lastOfRange-1]
	*/
	return nil
}

func (rc *runContainer16) getCardinality() int {
	return int(rc.cardinality())
}

func (rc *runContainer16) rank(x uint16) int {
	/*
		answer := binarySearch(ac.content, x)
		if answer >= 0 {
			return answer + 1
		}
		return -answer - 1

	*/
	return 0 // TODO
}

func (rc *runContainer16) selectInt(x uint16) int {
	return 0 // TODO
	//	return int(ac.content[x])
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

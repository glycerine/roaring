package roaring

///////////////////////////////////////////////////
//
// container interface methods for runContainer32
//
///////////////////////////////////////////////////

// compile time verify we meet interface requirements
var _ container = &runContainer32{}

func (ac *runContainer32) clone() container {
	/*
		ptr := runContainer32{make([]uint16, len(ac.content))}
		copy(ptr.content, ac.content[:])
		return &ptr
	*/
	return nil
}

func (ac *runContainer32) and(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.andArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.and(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) iand(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.iandArray(a.(*runContainer32))
		case *bitmapContainer:
			return ac.iandBitmap(a.(*bitmapContainer))
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) andNot(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.andNotArray(a.(*runContainer32))
		case *bitmapContainer:
			return ac.andNotBitmap(a.(*bitmapContainer))
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) fillLeastSignificant16bits(x []uint32, i int, mask uint32) {
	/*	for k := 0; k < len(ac.content); k++ {
			x[k+i] = uint32(ac.content[k]) | mask
		}
	*/
}

func (ac *runContainer32) getShortIterator() shortIterable {
	// TODO
	return &shortIterator{}
	//	return &shortIterator{ac.content, 0}
}

func (ac *runContainer32) getSizeInBytes() int {
	/*
		// unsafe.Sizeof calculates the memory used by the top level of the slice
		// descriptor - not including the size of the memory referenced by the slice.
		// http://golang.org/pkg/unsafe/#Sizeof
		return ac.getCardinality()*2 + int(unsafe.Sizeof(ac.content))
	*/
	return 0 // TODO
}

// serializedSizeInBytes returns the number of bytes of memory
// required by this runContainer32.
//func (rc *runContainer32) serializedSizeInBytes() int

// add the values in the range [firstOfRange,lastofRange)
func (ac *runContainer32) iaddRange(firstOfRange, lastOfRange int) container {
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
func (ac *runContainer32) iremoveRange(firstOfRange, lastOfRange int) container {
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
func (ac *runContainer32) not(firstOfRange, lastOfRange int) container {
	/*	if firstOfRange >= lastOfRange {
			return ac.clone()
		}
		return ac.notClose(firstOfRange, lastOfRange-1) // remove everything in [firstOfRange,lastOfRange-1]
	*/
	return nil
}

func (ac *runContainer32) equals(o interface{}) bool {
	/*
		srb, ok := o.(*runContainer32)
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

func (ac *runContainer32) add(x uint16) container {
	/*	// Special case adding to the end of the container.
		l := len(ac.content)
		if l > 0 && l < arrayDefaultMaxSize && ac.content[l-1] < x {
			ac.content = append(ac.content, x)
			return ac
		}

		loc := binarySearch(ac.content, x)

		if loc < 0 {
			if len(ac.content) >= arrayDefaultMaxSize {
				a := ac.toBitmapContainer()
				a.add(x)
				return a
			}
			s := ac.content
			i := -loc - 1
			s = append(s, 0)
			copy(s[i+1:], s[i:])
			s[i] = x
			ac.content = s
		}
		return ac
	*/
	return nil
}

func (ac *runContainer32) remove(x uint16) container {
	/*	loc := binarySearch(ac.content, x)
		if loc >= 0 {
			s := ac.content
			s = append(s[:loc], s[loc+1:]...)
			ac.content = s
		}
		return ac
	*/
	return nil

}

func (ac *runContainer32) or(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.orArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.or(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) ior(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.orArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.ior(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) lazyIOR(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.lazyorArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.lazyOR(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) lazyOR(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.lazyorArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.lazyOR(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) intersects(a container) bool {
	/*
		switch a.(type) {
		case *runContainer32:
			return ac.intersectsArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.intersects(ac)
		}
		panic("unsupported container type")
	*/
	return false
}

func (ac *runContainer32) xor(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.xorArray(a.(*runContainer32))
		case *bitmapContainer:
			return a.xor(ac)
		}
		panic("unsupported container type")
	*/
	return nil
}

func (ac *runContainer32) iandNot(a container) container {
	/*	switch a.(type) {
		case *runContainer32:
			return ac.iandNotArray(a.(*runContainer32))
		case *bitmapContainer:
			return ac.iandNotBitmap(a.(*bitmapContainer))
		}
		panic("unsupported container type")
	*/
	return nil
}

// flip the values in the range [firstOfRange,lastOfRange)
func (ac *runContainer32) inot(firstOfRange, lastOfRange int) container {
	/*	if firstOfRange >= lastOfRange {
			return ac
		}
		return ac.inotClose(firstOfRange, lastOfRange-1) // remove everything in [firstOfRange,lastOfRange-1]
	*/
	return nil
}

func (ac *runContainer32) getCardinality() int {
	return int(ac.cardinality())
}

func (ac *runContainer32) rank(x uint16) int {
	/*
		answer := binarySearch(ac.content, x)
		if answer >= 0 {
			return answer + 1
		}
		return -answer - 1

	*/
	return 0 // TODO
}

func (ac *runContainer32) selectInt(x uint16) int {
	return 0 // TODO
	//	return int(ac.content[x])
}

func (ac *runContainer32) contains(x uint16) bool {
	return false // TODO
	//	return binarySearch(ac.content, x) >= 0
}

package roaring

/*
func (ac *runContainer32) loadData(bitmapContainer *bitmapContainer) {
	ac.content = make([]uint16, bitmapContainer.cardinality, bitmapContainer.cardinality)
	bitmapContainer.fillArray(ac.content)
}
func newArrayContainer() *runContainer32 {
	p := new(runContainer32)
	return p
}

func newArrayContainerCapacity(size int) *runContainer32 {
	p := new(runContainer32)
	p.content = make([]uint16, 0, size)
	return p
}

func newArrayContainerSize(size int) *runContainer32 {
	p := new(runContainer32)
	p.content = make([]uint16, size, size)
	return p
}

func newRunContainer32Range(firstOfRun, lastOfRun int) *runContainer32 {
	valuesInRange := lastOfRun - firstOfRun + 1
	this := newArrayContainerCapacity(valuesInRange)
	for i := 0; i < valuesInRange; i++ {
		this.content = append(this.content, uint16(firstOfRun+i))
	}
	return this
}
*/
/*
// flip the values in the range [firstOfRange,lastOfRange]
func (ac *runContainer32) notClose(firstOfRange, lastOfRange int) container {
	if firstOfRange > lastOfRange { // unlike add and remove, not uses an inclusive range [firstOfRange,lastOfRange]
		return ac.clone()
	}

	// determine the span of array indices to be affected^M
	startIndex := binarySearch(ac.content, uint16(firstOfRange))
	if startIndex < 0 {
		startIndex = -startIndex - 1
	}
	lastIndex := binarySearch(ac.content, uint16(lastOfRange))
	if lastIndex < 0 {
		lastIndex = -lastIndex - 2
	}
	currentValuesInRange := lastIndex - startIndex + 1
	spanToBeFlipped := lastOfRange - firstOfRange + 1
	newValuesInRange := spanToBeFlipped - currentValuesInRange
	cardinalityChange := newValuesInRange - currentValuesInRange
	newCardinality := len(ac.content) + cardinalityChange

	if newCardinality > arrayDefaultMaxSize {
		return ac.toBitmapContainer().not(firstOfRange, lastOfRange+1)
	}
	answer := newArrayContainer()
	answer.content = make([]uint16, newCardinality, newCardinality) //a hack for sure

	copy(answer.content, ac.content[:startIndex])
	outPos := startIndex
	inPos := startIndex
	valInRange := firstOfRange
	for ; valInRange <= lastOfRange && inPos <= lastIndex; valInRange++ {
		if uint16(valInRange) != ac.content[inPos] {
			answer.content[outPos] = uint16(valInRange)
			outPos++
		} else {
			inPos++
		}
	}

	for ; valInRange <= lastOfRange; valInRange++ {
		answer.content[outPos] = uint16(valInRange)
		outPos++
	}

	for i := lastIndex + 1; i < len(ac.content); i++ {
		answer.content[outPos] = ac.content[i]
		outPos++
	}
	answer.content = answer.content[:newCardinality]
	return answer

}
*/
/*
func (ac *runContainer32) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	bc.loadData(ac)
	return bc

}
*/
/*
func (ac *runContainer32) orArray(value2 *runContainer32) container {
	value1 := ac
	maxPossibleCardinality := value1.getCardinality() + value2.getCardinality()
	if maxPossibleCardinality > arrayDefaultMaxSize { // it could be a bitmap!^M
		bc := newBitmapContainer()
		for k := 0; k < len(value2.content); k++ {
			v := value2.content[k]
			i := uint(v) >> 6
			mask := uint64(1) << (v % 64)
			bc.bitmap[i] |= mask
		}
		for k := 0; k < len(ac.content); k++ {
			v := ac.content[k]
			i := uint(v) >> 6
			mask := uint64(1) << (v % 64)
			bc.bitmap[i] |= mask
		}
		bc.cardinality = int(popcntSlice(bc.bitmap))
		if bc.cardinality <= arrayDefaultMaxSize {
			return bc.toArrayContainer()
		}
		return bc
	}
	answer := newArrayContainerCapacity(maxPossibleCardinality)
	nl := union2by2(value1.content, value2.content, answer.content)
	answer.content = answer.content[:nl] // reslice to match actual used capacity
	return answer
}
*/

/*
func (ac *runContainer32) lazyorArray(value2 *runContainer32) container {
	value1 := ac
	maxPossibleCardinality := value1.getCardinality() + value2.getCardinality()
	if maxPossibleCardinality > arrayLazyLowerBound { // it could be a bitmap!^M
		bc := newBitmapContainer()
		for k := 0; k < len(value2.content); k++ {
			v := value2.content[k]
			i := uint(v) >> 6
			mask := uint64(1) << (v % 64)
			bc.bitmap[i] |= mask
		}
		for k := 0; k < len(ac.content); k++ {
			v := ac.content[k]
			i := uint(v) >> 6
			mask := uint64(1) << (v % 64)
			bc.bitmap[i] |= mask
		}
		bc.cardinality = invalidCardinality
		return bc
	}
	answer := newArrayContainerCapacity(maxPossibleCardinality)
	nl := union2by2(value1.content, value2.content, answer.content)
	answer.content = answer.content[:nl] // reslice to match actual used capacity
	return answer
}
*/
/*
func (ac *runContainer32) iandBitmap(bc *bitmapContainer) *runContainer32 {
	pos := 0
	c := ac.getCardinality()
	for k := 0; k < c; k++ {
		if bc.contains(ac.content[k]) {
			ac.content[pos] = ac.content[k]
			pos++
		}
	}
	ac.content = ac.content[:pos]
	return ac

}
*/
/*
func (ac *runContainer32) xorArray(value2 *runContainer32) container {
	value1 := ac
	totalCardinality := value1.getCardinality() + value2.getCardinality()
	if totalCardinality > arrayDefaultMaxSize { // it could be a bitmap!
		bc := newBitmapContainer()
		for k := 0; k < len(value2.content); k++ {
			v := value2.content[k]
			i := uint(v) >> 6
			bc.bitmap[i] ^= (uint64(1) << (v % 64))
		}
		for k := 0; k < len(ac.content); k++ {
			v := ac.content[k]
			i := uint(v) >> 6
			bc.bitmap[i] ^= (uint64(1) << (v % 64))
		}
		bc.computeCardinality()
		if bc.cardinality <= arrayDefaultMaxSize {
			return bc.toArrayContainer()
		}
		return bc
	}
	desiredCapacity := totalCardinality
	answer := newArrayContainerCapacity(desiredCapacity)
	length := exclusiveUnion2by2(value1.content, value2.content, answer.content)
	answer.content = answer.content[:length]
	return answer

}
*/
/*
func (ac *runContainer32) andNotArray(value2 *runContainer32) container {
	value1 := ac
	desiredcapacity := value1.getCardinality()
	answer := newArrayContainerCapacity(desiredcapacity)
	length := difference(value1.content, value2.content, answer.content)
	answer.content = answer.content[:length]
	return answer
}

func (ac *runContainer32) iandNotArray(value2 *runContainer32) container {
	length := difference(ac.content, value2.content, ac.content)
	ac.content = ac.content[:length]
	return ac
}

func (ac *runContainer32) andNotBitmap(value2 *bitmapContainer) container {
	desiredcapacity := ac.getCardinality()
	answer := newArrayContainerCapacity(desiredcapacity)
	answer.content = answer.content[:desiredcapacity]
	pos := 0
	for _, v := range ac.content {
		if !value2.contains(v) {
			answer.content[pos] = v
			pos++
		}
	}
	answer.content = answer.content[:pos]
	return answer
}

func (ac *runContainer32) andBitmap(value2 *bitmapContainer) container {
	desiredcapacity := ac.getCardinality()
	answer := newArrayContainerCapacity(desiredcapacity)
	answer.content = answer.content[:desiredcapacity]
	pos := 0
	for _, v := range ac.content {
		if value2.contains(v) {
			answer.content[pos] = v
			pos++
		}
	}
	answer.content = answer.content[:pos]
	return answer
}

func (ac *runContainer32) iandNotBitmap(value2 *bitmapContainer) container {
	pos := 0
	for _, v := range ac.content {
		if !value2.contains(v) {
			ac.content[pos] = v
			pos++
		}
	}
	ac.content = ac.content[:pos]
	return ac
}

func copyOf(array []uint16, size int) []uint16 {
	result := make([]uint16, size)
	for i, x := range array {
		if i == size {
			break
		}
		result[i] = x
	}
	return result
}
*/
/*
// flip the values in the range [firstOfRange,lastOfRange]
func (ac *runContainer32) inotClose(firstOfRange, lastOfRange int) container {
	if firstOfRange > lastOfRange { // unlike add and remove, not uses an inclusive range [firstOfRange,lastOfRange]
		return ac
	}
	// determine the span of array indices to be affected
	startIndex := binarySearch(ac.content, uint16(firstOfRange))
	if startIndex < 0 {
		startIndex = -startIndex - 1
	}
	lastIndex := binarySearch(ac.content, uint16(lastOfRange))
	if lastIndex < 0 {
		lastIndex = -lastIndex - 1 - 1
	}
	currentValuesInRange := lastIndex - startIndex + 1
	spanToBeFlipped := lastOfRange - firstOfRange + 1

	newValuesInRange := spanToBeFlipped - currentValuesInRange
	buffer := make([]uint16, newValuesInRange)
	cardinalityChange := newValuesInRange - currentValuesInRange
	newCardinality := len(ac.content) + cardinalityChange
	if cardinalityChange > 0 {
		if newCardinality > len(ac.content) {
			if newCardinality > arrayDefaultMaxSize {
				return ac.toBitmapContainer().inot(firstOfRange, lastOfRange+1)
			}
			ac.content = copyOf(ac.content, newCardinality)
		}
		base := lastIndex + 1
		copy(ac.content[lastIndex+1+cardinalityChange:], ac.content[base:base+len(ac.content)-1-lastIndex])
		ac.negateRange(buffer, startIndex, lastIndex, firstOfRange, lastOfRange+1)
	} else { // no expansion needed
		ac.negateRange(buffer, startIndex, lastIndex, firstOfRange, lastOfRange+1)
		if cardinalityChange < 0 {

			for i := startIndex + newValuesInRange; i < newCardinality; i++ {
				ac.content[i] = ac.content[i-cardinalityChange]
			}
		}
	}
	ac.content = ac.content[:newCardinality]
	return ac
}


func (ac *runContainer32) negateRange(buffer []uint16, startIndex, lastIndex, startRange, lastRange int) {
	// compute the negation into buffer
	outPos := 0
	inPos := startIndex // value here always >= valInRange,
	// until it is exhausted
	// n.b., we can start initially exhausted.

	valInRange := startRange
	for ; valInRange < lastRange && inPos <= lastIndex; valInRange++ {
		if uint16(valInRange) != ac.content[inPos] {
			buffer[outPos] = uint16(valInRange)
			outPos++
		} else {
			inPos++
		}
	}

	// if there are extra items (greater than the biggest
	// pre-existing one in range), buffer them
	for ; valInRange < lastRange; valInRange++ {
		buffer[outPos] = uint16(valInRange)
		outPos++
	}

	if outPos != len(buffer) {
		panic("negateRange: internal bug")
	}

	for i, item := range buffer {
		ac.content[i+startIndex] = item
	}
}
*/

/*
func (ac *runContainer32) andArray(value2 *runContainer32) *runContainer32 {
	desiredcapacity := min(ac.getCardinality(), value2.getCardinality())
	answer := newArrayContainerCapacity(desiredcapacity)
	length := intersection2by2(
		ac.content,
		value2.content,
		answer.content)
	answer.content = answer.content[:length]
	return answer
}

func (ac *runContainer32) intersectsArray(value2 *runContainer32) bool {
	return intersects2by2(
		ac.content,
		value2.content)
}

func (ac *runContainer32) iandArray(value2 *runContainer32) *runContainer32 {
	length := intersection2by2(
		ac.content,
		value2.content,
		ac.content)
	ac.content = ac.content[:length]
	return ac
}
*/

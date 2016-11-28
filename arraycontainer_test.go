package roaring

// to run just these tests: go test -run TestArrayContainer*

import (
	"reflect"
	"testing"
)

func TestArrayContainerTransition(t *testing.T) {
	v := container(newArrayContainer())
	arraytype := reflect.TypeOf(v)
	for i := 0; i < arrayDefaultMaxSize; i++ {
		v = v.iaddReturnMinimized(uint16(i))
	}
	if v.getCardinality() != arrayDefaultMaxSize {
		t.Errorf("Bad cardinality.")
	}
	if reflect.TypeOf(v) != arraytype {
		t.Errorf("Should be an array.")
	}
	for i := 0; i < arrayDefaultMaxSize; i++ {
		v = v.iaddReturnMinimized(uint16(i))
	}
	if v.getCardinality() != arrayDefaultMaxSize {
		t.Errorf("Bad cardinality.")
	}
	if reflect.TypeOf(v) != arraytype {
		t.Errorf("Should be an array.")
	}
	v = v.iaddReturnMinimized(uint16(arrayDefaultMaxSize))
	if v.getCardinality() != arrayDefaultMaxSize+1 {
		t.Errorf("Bad cardinality.")
	}
	if reflect.TypeOf(v) == arraytype {
		t.Errorf("Should be a bitmap.")
	}
	v = v.iremoveReturnMinimized(uint16(arrayDefaultMaxSize))
	if v.getCardinality() != arrayDefaultMaxSize {
		t.Errorf("Bad cardinality.")
	}
	if reflect.TypeOf(v) != arraytype {
		t.Errorf("Should be an array.")
	}
}

func TestArrayContainerRank(t *testing.T) {
	v := container(newArrayContainer())
	v = v.iaddReturnMinimized(10)
	v = v.iaddReturnMinimized(100)
	v = v.iaddReturnMinimized(1000)
	if v.getCardinality() != 3 {
		t.Errorf("Bogus cardinality.")
	}
	for i := 0; i <= arrayDefaultMaxSize; i++ {
		thisrank := v.rank(uint16(i))
		if i < 10 {
			if thisrank != 0 {
				t.Errorf("At %d should be zero but is %d ", i, thisrank)
			}
		} else if i < 100 {
			if thisrank != 1 {
				t.Errorf("At %d should be zero but is %d ", i, thisrank)
			}
		} else if i < 1000 {
			if thisrank != 2 {
				t.Errorf("At %d should be zero but is %d ", i, thisrank)
			}
		} else {
			if thisrank != 3 {
				t.Errorf("At %d should be zero but is %d ", i, thisrank)
			}
		}
	}
}

func TestArrayContainerMassiveSetAndGet(t *testing.T) {
	v := container(newArrayContainer())
	for j := 0; j <= arrayDefaultMaxSize; j++ {

		v = v.iaddReturnMinimized(uint16(j))
		if v.getCardinality() != 1+j {
			t.Errorf("Bogus cardinality %d %d. ", v.getCardinality(), j)
		}
		for i := 0; i <= arrayDefaultMaxSize; i++ {
			if i <= j {
				if v.contains(uint16(i)) != true {
					t.Errorf("I added a number in vain.")
				}
			} else {
				if v.contains(uint16(i)) != false {
					t.Errorf("Ghost content")
					break
				}
			}
		}
	}
}

type FakeContainer struct {
	arrayContainer
}

func TestArrayContainerUnsupportedType(t *testing.T) {
	a := container(newArrayContainer())
	testContainerPanics(t, a)
	b := container(newBitmapContainer())
	testContainerPanics(t, b)
}

func testContainerPanics(t *testing.T, c container) {
	f := &FakeContainer{}
	assertPanic(t, func() {
		c.or(f)
	})
	assertPanic(t, func() {
		c.ior(f)
	})
	assertPanic(t, func() {
		c.lazyIOR(f)
	})
	assertPanic(t, func() {
		c.lazyOR(f)
	})
	assertPanic(t, func() {
		c.and(f)
	})
	assertPanic(t, func() {
		c.intersects(f)
	})
	assertPanic(t, func() {
		c.iand(f)
	})
	assertPanic(t, func() {
		c.xor(f)
	})
	assertPanic(t, func() {
		c.andNot(f)
	})
	assertPanic(t, func() {
		c.iandNot(f)
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

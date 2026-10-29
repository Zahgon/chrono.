package js

import (
	"math"
	"slices"
)

// minGallopWins is V8's kMinGallopWins, the initial galloping threshold.
const minGallopWins = 7

type sortRun struct{ base, length int }

// sorter is V8's SortState for one Array.prototype.sort call (third_party/v8/builtins/array-sort.tq).
type sorter[T any] struct {
	a         []T
	cmp       func(a, b T) float64
	minGallop int
	runs      []sortRun
	tmp       []T
}

// SortSlice sorts s in place exactly as V8's Array.prototype.sort(comparefn) does in Node.js 24
// (V8 13.6): the same TimSort with the same sequence of comparator calls, so the final order
// matches V8 even for inconsistent comparators. A NaN comparator result counts as 0.
func SortSlice[T any](s []T, cmp func(a, b T) float64) {
	sortSlice(s, cmp, 0)
}

// v8v14InsertionSortMax is the largest length that V8 14 (Node.js 26) sorts with a plain binary
// insertion sort before falling back to TimSort; V8 13.6 (Node.js 24) has no such path.
const v8v14InsertionSortMax = 7

// sortSlice is SortSlice where lengths up to insertionMax use a plain binary insertion sort.
func sortSlice[T any](s []T, cmp func(a, b T) float64, insertionMax int) {
	if len(s) < 2 {
		return
	}
	st := &sorter[T]{a: s, cmp: cmp, minGallop: minGallopWins}
	if len(s) <= insertionMax {
		st.binaryInsertionSort(0, 0, len(s))
		return
	}
	st.sort()
}

// SortStrings sorts s in place as Array.prototype.sort() with no comparator does for strings:
// ascending UTF-16 code unit order, stable.
func SortStrings(s []string) {
	slices.SortStableFunc(s, Compare)
}

// SortNumbersDefault sorts s in place as Array.prototype.sort() with no comparator does for
// numbers: by their String() forms in code unit order, stable.
func SortNumbersDefault(s []float64) {
	keys := make(map[float64]string, len(s))
	str := func(f float64) string {
		if math.IsNaN(f) {
			return "NaN"
		}
		k, ok := keys[f]
		if !ok {
			k = NumberToString(f)
			keys[f] = k
		}
		return k
	}
	slices.SortStableFunc(s, func(a, b float64) int { return Compare(str(a), str(b)) })
}

func (st *sorter[T]) compare(x, y T) float64 {
	v := st.cmp(x, y)
	if math.IsNaN(v) {
		return 0
	}
	return v
}

func (st *sorter[T]) sort() {
	n := len(st.a)
	remaining := n
	low := 0
	minRun := minRunLength(remaining)
	for remaining != 0 {
		cur := st.countAndMakeRun(low, low+remaining)
		if cur < minRun {
			forced := min(minRun, remaining)
			st.binaryInsertionSort(low, low+cur, low+forced)
			cur = forced
		}
		st.runs = append(st.runs, sortRun{low, cur})
		st.mergeCollapse()
		low += cur
		remaining -= cur
	}
	st.mergeForceCollapse()
}

func minRunLength(n int) int {
	r := 0
	for n >= 64 {
		r |= n & 1
		n >>= 1
	}
	return n + r
}

func (st *sorter[T]) binaryInsertionSort(low, startArg, high int) {
	a := st.a
	start := startArg
	if low == startArg {
		start++
	}
	for ; start < high; start++ {
		left, right := low, start
		pivot := a[start]
		for left < right {
			mid := left + ((right - left) >> 1)
			if st.compare(pivot, a[mid]) < 0 {
				right = mid
			} else {
				left = mid + 1
			}
		}
		copy(a[left+1:start+1], a[left:start])
		a[left] = pivot
	}
}

func (st *sorter[T]) countAndMakeRun(lowArg, high int) int {
	a := st.a
	low := lowArg + 1
	if low == high {
		return 1
	}
	runLength := 2
	elementLow, elementLowPre := a[low], a[low-1]
	order := st.compare(elementLow, elementLowPre)
	descending := order < 0
	previous := elementLow
	for idx := low + 1; idx < high; idx++ {
		current := a[idx]
		order = st.compare(current, previous)
		if descending {
			if order >= 0 {
				break
			}
		} else if order < 0 {
			break
		}
		previous = current
		runLength++
	}
	if descending {
		slices.Reverse(a[lowArg : lowArg+runLength])
	}
	return runLength
}

func (st *sorter[T]) runInvariantEstablished(n int) bool {
	if n < 2 {
		return true
	}
	return st.runs[n-2].length > st.runs[n-1].length+st.runs[n].length
}

func (st *sorter[T]) mergeCollapse() {
	for len(st.runs) > 1 {
		n := len(st.runs) - 2
		if !st.runInvariantEstablished(n+1) || !st.runInvariantEstablished(n) {
			if st.runs[n-1].length < st.runs[n+1].length {
				n--
			}
			st.mergeAt(n)
		} else if st.runs[n].length <= st.runs[n+1].length {
			st.mergeAt(n)
		} else {
			break
		}
	}
}

func (st *sorter[T]) mergeForceCollapse() {
	for len(st.runs) > 1 {
		n := len(st.runs) - 2
		if n > 0 && st.runs[n-1].length < st.runs[n+1].length {
			n--
		}
		st.mergeAt(n)
	}
}

func (st *sorter[T]) mergeAt(i int) {
	a := st.a
	baseA, lengthA := st.runs[i].base, st.runs[i].length
	baseB, lengthB := st.runs[i+1].base, st.runs[i+1].length
	st.runs[i].length = lengthA + lengthB
	if i == len(st.runs)-3 {
		st.runs[i+1] = st.runs[i+2]
	}
	st.runs = st.runs[:len(st.runs)-1]

	k := st.gallopRight(a, a[baseB], baseA, lengthA, 0)
	baseA += k
	lengthA -= k
	if lengthA == 0 {
		return
	}
	lengthB = st.gallopLeft(a, a[baseA+lengthA-1], baseB, lengthB, lengthB-1)
	if lengthB == 0 {
		return
	}
	if lengthA <= lengthB {
		st.mergeLow(baseA, lengthA, baseB, lengthB)
	} else {
		st.mergeHigh(baseA, lengthA, baseB, lengthB)
	}
}

// gallopLeft returns the offset in 0..length at which key belongs in arr[base:base+length],
// left of any equal elements, searching from hint.
func (st *sorter[T]) gallopLeft(arr []T, key T, base, length, hint int) int {
	lastOfs, offset := 0, 1
	if st.compare(arr[base+hint], key) < 0 {
		maxOfs := length - hint
		for offset < maxOfs {
			if st.compare(arr[base+hint+offset], key) >= 0 {
				break
			}
			lastOfs = offset
			offset = (offset << 1) + 1
			if offset <= 0 {
				offset = maxOfs
			}
		}
		if offset > maxOfs {
			offset = maxOfs
		}
		lastOfs += hint
		offset += hint
	} else {
		maxOfs := hint + 1
		for offset < maxOfs {
			if st.compare(arr[base+hint-offset], key) < 0 {
				break
			}
			lastOfs = offset
			offset = (offset << 1) + 1
			if offset <= 0 {
				offset = maxOfs
			}
		}
		if offset > maxOfs {
			offset = maxOfs
		}
		lastOfs, offset = hint-offset, hint-lastOfs
	}
	lastOfs++
	for lastOfs < offset {
		m := lastOfs + ((offset - lastOfs) >> 1)
		if st.compare(arr[base+m], key) < 0 {
			lastOfs = m + 1
		} else {
			offset = m
		}
	}
	return offset
}

// gallopRight is gallopLeft placing key right of any equal elements.
func (st *sorter[T]) gallopRight(arr []T, key T, base, length, hint int) int {
	lastOfs, offset := 0, 1
	if st.compare(key, arr[base+hint]) < 0 {
		maxOfs := hint + 1
		for offset < maxOfs {
			if st.compare(key, arr[base+hint-offset]) >= 0 {
				break
			}
			lastOfs = offset
			offset = (offset << 1) + 1
			if offset <= 0 {
				offset = maxOfs
			}
		}
		if offset > maxOfs {
			offset = maxOfs
		}
		lastOfs, offset = hint-offset, hint-lastOfs
	} else {
		maxOfs := length - hint
		for offset < maxOfs {
			if st.compare(key, arr[base+hint+offset]) < 0 {
				break
			}
			lastOfs = offset
			offset = (offset << 1) + 1
			if offset <= 0 {
				offset = maxOfs
			}
		}
		if offset > maxOfs {
			offset = maxOfs
		}
		lastOfs += hint
		offset += hint
	}
	lastOfs++
	for lastOfs < offset {
		m := lastOfs + ((offset - lastOfs) >> 1)
		if st.compare(key, arr[base+m]) < 0 {
			offset = m
		} else {
			lastOfs = m + 1
		}
	}
	return offset
}

func (st *sorter[T]) tempArray(n int) []T {
	if cap(st.tmp) < n {
		st.tmp = make([]T, n)
	}
	return st.tmp[:n]
}

func (st *sorter[T]) mergeLow(baseA, lengthA, baseB, lengthB int) {
	a := st.a
	tmp := st.tempArray(lengthA)
	copy(tmp, a[baseA:baseA+lengthA])
	dest, cursorTemp, cursorB := baseA, 0, baseB

	a[dest] = a[cursorB]
	dest++
	cursorB++

	succeed := func() {
		if lengthA > 0 {
			copy(a[dest:dest+lengthA], tmp[cursorTemp:cursorTemp+lengthA])
		}
	}
	copyB := func() {
		copy(a[dest:dest+lengthB], a[cursorB:cursorB+lengthB])
		a[dest+lengthB] = tmp[cursorTemp]
	}

	lengthB--
	if lengthB == 0 {
		succeed()
		return
	}
	if lengthA == 1 {
		copyB()
		return
	}
	minGallop := st.minGallop
	for {
		winsA, winsB := 0, 0
		for {
			if st.compare(a[cursorB], tmp[cursorTemp]) < 0 {
				a[dest] = a[cursorB]
				dest++
				cursorB++
				winsB++
				lengthB--
				winsA = 0
				if lengthB == 0 {
					succeed()
					return
				}
				if winsB >= minGallop {
					break
				}
			} else {
				a[dest] = tmp[cursorTemp]
				dest++
				cursorTemp++
				winsA++
				lengthA--
				winsB = 0
				if lengthA == 1 {
					copyB()
					return
				}
				if winsA >= minGallop {
					break
				}
			}
		}
		minGallop++
		for first := true; first || winsA >= minGallopWins || winsB >= minGallopWins; first = false {
			minGallop = max(1, minGallop-1)
			st.minGallop = minGallop

			winsA = st.gallopRight(tmp, a[cursorB], cursorTemp, lengthA, 0)
			if winsA > 0 {
				copy(a[dest:dest+winsA], tmp[cursorTemp:cursorTemp+winsA])
				dest += winsA
				cursorTemp += winsA
				lengthA -= winsA
				if lengthA == 1 {
					copyB()
					return
				}
				if lengthA == 0 {
					succeed()
					return
				}
			}
			a[dest] = a[cursorB]
			dest++
			cursorB++
			lengthB--
			if lengthB == 0 {
				succeed()
				return
			}

			winsB = st.gallopLeft(a, tmp[cursorTemp], cursorB, lengthB, 0)
			if winsB > 0 {
				copy(a[dest:dest+winsB], a[cursorB:cursorB+winsB])
				dest += winsB
				cursorB += winsB
				lengthB -= winsB
				if lengthB == 0 {
					succeed()
					return
				}
			}
			a[dest] = tmp[cursorTemp]
			dest++
			cursorTemp++
			lengthA--
			if lengthA == 1 {
				copyB()
				return
			}
		}
		minGallop++
		st.minGallop = minGallop
	}
}

func (st *sorter[T]) mergeHigh(baseA, lengthA, baseB, lengthB int) {
	a := st.a
	tmp := st.tempArray(lengthB)
	copy(tmp, a[baseB:baseB+lengthB])
	dest, cursorTemp, cursorA := baseB+lengthB-1, lengthB-1, baseA+lengthA-1

	a[dest] = a[cursorA]
	dest--
	cursorA--

	succeed := func() {
		if lengthB > 0 {
			copy(a[dest-(lengthB-1):dest+1], tmp[:lengthB])
		}
	}
	copyA := func() {
		dest -= lengthA
		cursorA -= lengthA
		copy(a[dest+1:dest+1+lengthA], a[cursorA+1:cursorA+1+lengthA])
		a[dest] = tmp[cursorTemp]
	}

	lengthA--
	if lengthA == 0 {
		succeed()
		return
	}
	if lengthB == 1 {
		copyA()
		return
	}
	minGallop := st.minGallop
	for {
		winsA, winsB := 0, 0
		for {
			if st.compare(tmp[cursorTemp], a[cursorA]) < 0 {
				a[dest] = a[cursorA]
				dest--
				cursorA--
				winsA++
				lengthA--
				winsB = 0
				if lengthA == 0 {
					succeed()
					return
				}
				if winsA >= minGallop {
					break
				}
			} else {
				a[dest] = tmp[cursorTemp]
				dest--
				cursorTemp--
				winsB++
				lengthB--
				winsA = 0
				if lengthB == 1 {
					copyA()
					return
				}
				if winsB >= minGallop {
					break
				}
			}
		}
		minGallop++
		for first := true; first || winsA >= minGallopWins || winsB >= minGallopWins; first = false {
			minGallop = max(1, minGallop-1)
			st.minGallop = minGallop

			k := st.gallopRight(a, tmp[cursorTemp], baseA, lengthA, lengthA-1)
			winsA = lengthA - k
			if winsA > 0 {
				dest -= winsA
				cursorA -= winsA
				copy(a[dest+1:dest+1+winsA], a[cursorA+1:cursorA+1+winsA])
				lengthA -= winsA
				if lengthA == 0 {
					succeed()
					return
				}
			}
			a[dest] = tmp[cursorTemp]
			dest--
			cursorTemp--
			lengthB--
			if lengthB == 1 {
				copyA()
				return
			}

			k = st.gallopLeft(tmp, a[cursorA], 0, lengthB, lengthB-1)
			winsB = lengthB - k
			if winsB > 0 {
				dest -= winsB
				cursorTemp -= winsB
				copy(a[dest+1:dest+1+winsB], tmp[cursorTemp+1:cursorTemp+1+winsB])
				lengthB -= winsB
				if lengthB == 1 {
					copyA()
					return
				}
				if lengthB == 0 {
					succeed()
					return
				}
			}
			a[dest] = a[cursorA]
			dest--
			cursorA--
			lengthA--
			if lengthA == 0 {
				succeed()
				return
			}
		}
		minGallop++
		st.minGallop = minGallop
	}
}

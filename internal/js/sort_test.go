package js

import (
	"math"
	"slices"
	"testing"
)

// genSortValues mirrors genValues in dev/js/gen_sort.mjs.
func genSortValues(seed uint32, n, pat, dup int) []float64 {
	r := &mulberry32{a: seed}
	v := make([]float64, n)
	switch pat {
	case 0:
		for i := range v {
			v[i] = float64(r.next() % uint32(dup))
		}
	case 1:
		for i := range v {
			v[i] = float64(i / dup)
		}
	case 2:
		for i := range v {
			v[i] = float64((n - i) / dup)
		}
	case 3:
		for i := 0; i < n; {
			length := 1 + int(r.next()%40)
			base := float64(r.next() % 1000)
			up := r.next()&1 != 0
			for j := 0; j < length && i < n; j, i = j+1, i+1 {
				if up {
					v[i] = base + float64(j)
				} else {
					v[i] = base - float64(j)
				}
			}
		}
	case 4:
		for i := range v {
			v[i] = float64(i % (dup + 1))
		}
	case 5:
		for i := range v {
			v[i] = 7
		}
	case 6:
		for i := range v {
			v[i] = float64(i)
		}
		swaps := 1 + int(r.next()%8)
		for k := 0; k < swaps && n > 1; k++ {
			a, b := r.next()%uint32(n), r.next()%uint32(n)
			v[a], v[b] = v[b], v[a]
		}
	default:
		for i := 0; i < n; {
			length := 1 + int(r.next()%200)
			base := float64(r.next() % 100000)
			for j := 0; j < length && i < n; j, i = j+1, i+1 {
				v[i] = base - 3*float64(j)
			}
		}
	}
	return v
}

// sortCmp mirrors makeCmp in dev/js/gen_sort.mjs.
func sortCmp(kind int, vals []float64, r *mulberry32, calls *int) func(a, b int) float64 {
	f := []func(a, b int) float64{
		func(a, b int) float64 { return vals[a] - vals[b] },
		func(a, b int) float64 {
			if vals[a] > vals[b] {
				return 1
			}
			return -1
		},
		func(a, b int) float64 {
			if vals[a] < vals[b] {
				return 1
			}
			return -1
		},
		func(a, b int) float64 { return float64(int(r.next()%3) - 1) },
		func(a, b int) float64 { return float64(r.next())/4294967296 - 0.5 },
		func(a, b int) float64 {
			if vals[a] == vals[b] {
				return math.NaN()
			}
			return vals[a] - vals[b]
		},
		func(a, b int) float64 {
			h := uint32(int32(vals[a]))*0x9e3779b1 ^ uint32(int32(vals[b]))
			return float64(int(h%3) - 1)
		},
		func(a, b int) float64 { return vals[b] - vals[a] },
		func(a, b int) float64 { return math.Mod(vals[a], 10) - math.Mod(vals[b], 10) },
		func(a, b int) float64 {
			if r.next()&1 != 0 {
				return vals[a] - vals[b]
			}
			return vals[b] - vals[a]
		},
		func(a, b int) float64 {
			if r.next()%50 == 0 {
				return vals[b] - vals[a]
			}
			return vals[a] - vals[b]
		},
		func(a, b int) float64 {
			if vals[a] >= vals[b] {
				return 0.5
			}
			return math.Inf(-1)
		},
	}[kind]
	return func(a, b int) float64 {
		*calls++
		return f(a, b)
	}
}

func TestSortVectors(t *testing.T) {
	checkSortVectors(t, "sort.jsonl.gz", "SortSlice", 0)
}

func TestSortVectorsNode26(t *testing.T) {
	checkSortVectors(t, "sort_node26.jsonl.gz", "sortSlice (V8 14 short-array path)", v8v14InsertionSortMax)
}

func checkSortVectors(t *testing.T, file, what string, insertionMax int) {
	recs := readVectors(t, file)
	m := &mismatches{t: t, what: what + " permutations"}
	cm := &mismatches{t: t, what: what + " comparator call counts"}
	perKind := map[int]int{}
	maxN := 0
	for _, r := range recs {
		if str(r, "k") == "meta" {
			t.Logf("vectors from %s", str(r, "runtime"))
			continue
		}
		seed := uint32(num(r, "seed"))
		n, pat, dup, kind := int(num(r, "n")), int(num(r, "pat")), int(num(r, "dup")), int(num(r, "kind"))
		vals := genSortValues(seed, n, pat, dup)
		calls := 0
		cmp := sortCmp(kind, vals, &mulberry32{a: seed ^ 0x5bd1e995}, &calls)
		ids := make([]int, n)
		for i := range ids {
			ids[i] = i
		}
		sortSlice(ids, cmp, insertionMax)
		h := uint32(fnvBasis)
		for _, id := range ids {
			h = (h ^ uint32(id)) * 16777619
		}
		ok := h == uint32(num(r, "h"))
		if want, has := r.Get("ids"); has {
			var w []int
			for _, x := range want.(Array) {
				w = append(w, int(x.(float64)))
			}
			ok = ok && slices.Equal(ids, w)
		}
		m.check(ok, "seed=%d n=%d pat=%d dup=%d kind=%d: permutation differs", seed, n, pat, dup, kind)
		cm.check(calls == int(num(r, "calls")), "seed=%d n=%d kind=%d: %d comparator calls want %v", seed, n, kind, calls, num(r, "calls"))
		perKind[kind]++
		maxN = max(maxN, n)
	}
	t.Logf("arrays per comparator kind: %v, largest n=%d", perKind, maxN)
	m.done()
	cm.done()
}

func TestSortStrings(t *testing.T) {
	s := []string{"b", "\U0001F600", "\uff5e", "a", "", "B", "aa", "\u00e9"}
	SortStrings(s)
	want := []string{"", "B", "a", "aa", "b", "\u00e9", "\U0001F600", "\uff5e"}
	if !slices.Equal(s, want) {
		t.Fatalf("got %q want %q", s, want)
	}
}

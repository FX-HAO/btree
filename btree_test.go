package btree

import (
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func init() {
	seed := time.Now().Unix()
	fmt.Println(seed)
	rand.Seed(seed)
}

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []Item) {
	for _, v := range rand.Perm(n) {
		out = append(out, Int(v))
	}
	return
}

// rang returns an ordered list of Int items in the range [0, n).
func rang(n int) (out []Item) {
	for i := 0; i < n; i++ {
		out = append(out, Int(i))
	}
	return
}

// rangerev returns a reversed ordered list of Int items in the range [0, n).
func rangrev(n int) (out []Item) {
	for i := n - 1; i >= 0; i-- {
		out = append(out, Int(i))
	}
	return
}

// // all extracts all items from a tree in order as a slice.
func all(t *BTree) (out []Item) {
	t.Ascend(func(a Item) bool {
		out = append(out, a)
		return true
	})
	return
}

// // rangerev returns a reversed ordered list of Int items in the range [0, n).
// func rangrev(n int) (out []Item) {
// 	for i := n - 1; i >= 0; i-- {
// 		out = append(out, Int(i))
// 	}
// 	return
// }

// // allrev extracts all items from a tree in reverse order as a slice.
// func allrev(t *BTree) (out []Item) {
// 	t.Descend(func(a Item) bool {
// 		out = append(out, a)
// 		return true
// 	})
// 	return
// }

var btreeDegree = flag.Int("degree", 2, "B-Tree degree")

func TestBTree(t *testing.T) {
	tr := New(*btreeDegree)
	const treeSize = 1000
	for i := 0; i < 1; i++ {
		if min := tr.Min(); min != nil {
			t.Fatalf("empty min, got %+v", min)
		}
		if max := tr.Max(); max != nil {
			t.Fatalf("empty max, got %+v", max)
		}
		a := perm(treeSize)
		// a := []Int{3, 8, 1, 6, 2, 9, 5, 4, 7, 0}
		for _, item := range a {
			if x := tr.ReplaceOrInsert(item); x != nil {
				t.Fatal("insert found item", item)
			}
		}
		for _, item := range a {
			if x := tr.ReplaceOrInsert(item); x == nil {
				t.Fatal("insert didn't find item", item)
			}
		}
		if min, want := tr.Min(), Item(Int(0)); min != want {
			t.Fatalf("min: want %+v, got %+v", want, min)
		}
		if max, want := tr.Max(), Item(Int(treeSize-1)); max != want {
			t.Fatalf("max: want %+v, got %+v", want, max)
		}
		got := all(tr)
		want := rang(treeSize)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch:\n got: %v\nwant: %v", got, want)
		}

		// 	// 	// gotrev := allrev(tr)
		// 	// 	// wantrev := rangrev(treeSize)
		// 	// 	// if !reflect.DeepEqual(gotrev, wantrev) {
		// 	// 	// 	t.Fatalf("mismatch:\n got: %v\nwant: %v", got, want)
		// 	// 	// }

		// fmt.Println("Deleting...")
		// printNode(tr.root)
		// fmt.Println()
		for _, item := range perm(treeSize) {
			x := tr.Delete(item)
			// fmt.Println("")
			// fmt.Printf("%v - ", item)
			// printNode(tr.root)
			if x == nil {
				t.Fatalf("didn't find %v", item)
			}
			// if x := tr.Delete(item); x == nil {
			// 	t.Fatalf("didn't find %v", item)
			// }
		}
		if got = all(tr); len(got) > 0 {
			t.Fatalf("some left!: %v", got)
		}
	}
}

func TestAscendRange(t *testing.T) {
	tree := New(*btreeDegree)
	treeSize := 1000
	order := perm(treeSize)
	for _, i := range order {
		tree.ReplaceOrInsert(i)
	}
	k := 0
	tree.AscendRange(Int(2), Int(5), func(item Item) bool {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := Int(2 + k)
		i2 := item.(Int)
		if i1 != i2 {
			t.Errorf("expecting %v, got %v", i1, i2)
		}
		k++
		return true
	})
	if k != 3 {
		t.Errorf("expecting %v, got %v,", 3, k)
	}
}

func TestDescendRange(t *testing.T) {
	tr := New(2)
	for _, v := range perm(100) {
		tr.ReplaceOrInsert(v)
	}
	var got []Item
	tr.DescendRange(Int(60), Int(40), func(a Item) bool {
		got = append(got, a)
		return true
	})
	if want := rangrev(100)[39:59]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendrange:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	tr.DescendRange(Int(60), Int(40), func(a Item) bool {
		if a.(Int) < 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := rangrev(100)[39:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendGreaterThan(t *testing.T) {
	tr := New(*btreeDegree)
	for _, v := range perm(100) {
		tr.ReplaceOrInsert(v)
	}
	var got []Item
	tr.DescendGreaterThan(Int(40), func(a Item) bool {
		got = append(got, a)
		return true
	})
	if want := rangrev(100)[:59]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendgreaterthan:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	tr.DescendGreaterThan(Int(40), func(a Item) bool {
		if a.(Int) < 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := rangrev(100)[:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendgreaterthan:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendLessOrEqual(t *testing.T) {
	tr := New(*btreeDegree)
	for _, v := range perm(100) {
		tr.ReplaceOrInsert(v)
	}
	var got []Item
	tr.DescendLessOrEqual(Int(40), func(a Item) bool {
		got = append(got, a)
		return true
	})
	if want := rangrev(100)[59:]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendlessorequal:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	tr.DescendLessOrEqual(Int(60), func(a Item) bool {
		if a.(Int) < 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := rangrev(100)[39:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendlessorequal:\n got: %v\nwant: %v", got, want)
	}
}

func printNode(n *node) {
	if n == nil {
		fmt.Print(make(children, 0))
		return
	}

	fmt.Printf("%v-", n.items)
	for i := range n.children {
		fmt.Print(n.children[i].items)
	}
	fmt.Print("-")

	for i := range n.children {
		for j := range n.children[i].children {
			fmt.Print(n.children[i].children[j].items)
		}
	}
}

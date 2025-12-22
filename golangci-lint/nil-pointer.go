package golangcilint

import "fmt"

// A few different nil-pointer patterns. These are intentionally wrong to see what
// staticcheck/typecheck report. They compile, but would panic if executed.

func nilDerefDirect() int {
	var p *int
	return *p // nil dereference
}

func foo(x int) {

}

func fn(x *int) {
    fmt.Println(*x)

    // This nil check is equally important for the previous dereference
    if x != nil {
        foo(*x)
    }
}


func nilDeref_AfterNilCheck() int {
	var p *int
	if p == nil {
		return *p // definitely nil; staticcheck usually flags this (SA5011)
	}
	return 0
}

type node struct {
	next *node
	val  int
}

func nilDeref_ThroughField() int {
	var n *node
	return n.val // nil dereference through struct pointer
}

func nilDeref_ChainedField() int {
	var n *node
	if n == nil {
		// Still dereferencing (chained) even though we just proved it's nil.
		return n.next.val
	}
	return 0
}

func nilDeref_MapValuePointer() int {
	// Indexing a nil map is OK (returns zero value), but dereferencing the
	// resulting *int (which is nil) will panic if executed.
	var m map[string]*int
	return *m["missing"]
}

func nilDeref_SliceOfPointers() int {
	// Slice indexing is fine; element is nil pointer, deref panics if executed.
	s := make([]*int, 1)
	return *s[0]
}

// Returns a *node, but may return nil.
func maybeNode() *node {
    return nil
}

// Caller does not check for nil before dereferencing.
func nilDerefReturnedPointerNoCheck() int {
    n := maybeNode()
    return n.val // nil dereference if executed
}

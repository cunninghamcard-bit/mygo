package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestKeyAfterChildrenOrStatePanics checks that calling Element.Key after
// building children or reading the element's state panics with a clear
// message directing callers to Context.Key, rather than silently dropping
// input, closing locals or leaving children parented to the unkeyed ID.
func TestKeyAfterChildrenOrStatePanics(t *testing.T) {
	cases := []struct {
		name string
		want string
		fn   func(c *Context)
	}{
		{
			"Children",
			"Key after Children",
			func(c *Context) {
				Row(c).Children(func() { Text(c, "inside") }).Key("k")
			},
		},
		{
			"Clicked",
			"Key after the element's state was read",
			func(c *Context) {
				b := ButtonBase(c).Size(80, 32)
				_ = b.Clicked()
				b.Key("k")
			},
		},
		{
			"Focused",
			"Key after the element's state was read",
			func(c *Context) {
				b := ButtonBase(c).Size(80, 32)
				_ = b.Focused()
				b.Key("k")
			},
		},
		{
			"OnUnmount",
			"Key after the element's state was read",
			func(c *Context) {
				Box(c).OnUnmount(func() {}).Key("k")
			},
		},
		{
			"Hovered",
			"Key after the element's state was read",
			func(c *Context) {
				b := Row(c)
				_ = b.Hovered()
				b.Key("k")
			},
		},
		{
			"AnimateTo",
			"Key after the element's state was read",
			func(c *Context) {
				b := Box(c)
				_ = b.AnimateTo("x", 10, Animation{Duration: 100 * time.Millisecond})
				b.Key("k")
			},
		},
		{
			"FocusBind",
			"Key after the element's state was read",
			func(c *Context) {
				want := 1
				Box(c).FocusBind(&want, 1).Key("k")
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil || !strings.Contains(fmt.Sprint(r), tc.want) {
					t.Errorf("recovered %v, want substring %q", r, tc.want)
				}
			}()
			NewTester(tc.fn, 200, 100)
		})
	}
}

// TestKeyAfterStyleAndOnButtonWorks checks that pure layout/style chaining
// before Key, Context.Key before a custom control, and Button(c, label).Key
// continue to work and leave no orphaned positional states behind.
func TestKeyAfterStyleAndOnButtonWorks(t *testing.T) {
	clicks := 0
	custom := func(c *Context, label string) Element {
		b := ButtonBase(c).Padding(6, 12).Radius(6)
		n := Local(b, "n", func() int { return 1 })
		*n++
		if b.Clicked() {
			clicks++
		}
		b.Children(func() { Text(c, label) })
		return b
	}
	tt := NewTester(func(c *Context) {
		Column(c).Padding(8).Gap(4).Key("col").Children(func() {
			custom(c.Key("run"), "Run")
			for i := range 20 {
				Button(c, fmt.Sprintf("Item %d", i)).Padding(4).Key(i)
			}
		})
	}, 300, 400)
	if err := tt.Click("Run"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 {
		t.Fatalf("clicks = %d, want 1", clicks)
	}
	// 1 root + 1 col + 1 custom button + 1 custom text + 20 buttons + 20 button texts = 44 states at most.
	if got := len(tt.rt.states); got > 46 {
		t.Errorf("len(states) = %d, leaked positional states on Button.Key", got)
	}
}

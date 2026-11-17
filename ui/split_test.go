package ui

import (
	"image/color"
	"testing"
)

// SplitGap spaces the panes by gap DIPs from the first frame, paints no
// divider line over the background, and resizes across the gap by pointer
// drag, keyboard arrows and SplitReset, in DIPs and percentages.
func TestSplitGap(t *testing.T) {
	bg := RGB(250, 248, 242)
	left := float32(100)
	var a, b Rect
	tt := NewTester(func(c *Context) {
		c.Root().Background(bg)
		Split(c, &left, func() {
			a = Box(c).Fill().Background(RGB(220, 230, 245)).Bounds()
		}, func() {
			b = Box(c).Fill().Background(RGB(230, 245, 225)).Bounds()
		}).SplitGap(8).SplitReset(120).Fill()
	}, 300, 100)
	tt.Frame()
	if a != (Rect{0, 0, 100, 100}) || b != (Rect{108, 0, 192, 100}) {
		t.Fatalf("panes at %v and %v, want {0 0 100 100} and {108 0 192 100}", a, b)
	}
	// Hovering the gap shows the resize cursor without painting a line.
	tt.Move(104, 50)
	if tt.Cursor() != CursorResizeEW {
		t.Fatalf("cursor over gap = %v, want CursorResizeEW", tt.Cursor())
	}
	if got := tt.Image().RGBAAt(104, 50); got != (color.RGBA{bg.R, bg.G, bg.B, 255}) {
		t.Fatalf("gap pixel = %v, want background %v", got, bg)
	}
	// Dragging anywhere across the 8 DIP gap moves the split.
	tt.Press(107, 50)
	tt.Move(137, 50)
	tt.Release(137, 50)
	tt.Frame()
	if left != 130 || b.X != 138 || b.W != 162 {
		t.Fatalf("after drag: left=%v, second=%v", left, b)
	}
	// Keyboard arrows and double-click SplitReset respect the gap.
	tt.Key(0, KeyRight)
	tt.Frame()
	if left != 140 || b.X != 148 {
		t.Fatalf("after KeyRight: left=%v, second=%v", left, b)
	}
	tt.ClickAt(144, 50)
	tt.ClickAt(144, 50)
	tt.Frame()
	if left != 120 || b.X != 128 {
		t.Fatalf("after SplitReset double-click: left=%v, second=%v", left, b)
	}
}

// SplitGap(0) puts the panes flush with no 1 DIP line, and SplitPercent
// and SplitVertical share the room left after the gap.
func TestSplitGapZeroAndPercent(t *testing.T) {
	bg := RGB(255, 255, 255)
	pct := float32(50)
	top := float32(40)
	var a, b, va, vb Rect
	tt := NewTester(func(c *Context) {
		t := LightTheme()
		t.Background = bg
		c.SetTheme(t)
		Column(c).Fill().Children(func() {
			SplitPercent(c, &pct, func() {
				a = Box(c).Fill().Bounds()
			}, func() {
				b = Box(c).Fill().Bounds()
			}).SplitGap(0).Height(40).FillWidth()
			SplitVertical(c, &top, func() {
				va = Box(c).Fill().Bounds()
			}, func() {
				vb = Box(c).Fill().Bounds()
			}).SplitGap(10).Grow(1).FillWidth()
		})
	}, 200, 160)
	tt.Frame()
	if a != (Rect{0, 0, 100, 40}) || b != (Rect{100, 0, 100, 40}) {
		t.Fatalf("SplitGap(0) percent panes: %v and %v", a, b)
	}
	if got := tt.Image().RGBAAt(100, 20); got != (color.RGBA{255, 255, 255, 255}) {
		t.Fatalf("SplitGap(0) painted %v at seam x=100", got)
	}
	if va != (Rect{0, 40, 200, 40}) || vb != (Rect{0, 90, 200, 70}) {
		t.Fatalf("SplitVertical SplitGap(10) panes: %v and %v", va, vb)
	}
}

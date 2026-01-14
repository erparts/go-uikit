package uikit

import "image"

// shiftSubtree shifts the Rect of w (and all descendants if w is a Layout) by (dx,dy).
// It returns a restore function that must be called to restore all original rects.
func shiftSubtree(w Widget, dx, dy int) func() {
	type entry struct {
		b   *Base
		rec image.Rectangle
	}
	var stack []Widget
	var saved []entry

	stack = append(stack, w)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		b := cur.Base()
		saved = append(saved, entry{b: b, rec: b.Rect})
		b.Rect = image.Rect(b.Rect.Min.X+dx, b.Rect.Min.Y+dy, b.Rect.Max.X+dx, b.Rect.Max.Y+dy)

		if l, ok := any(cur).(Layout); ok {
			chs := l.Children()
			for i := len(chs) - 1; i >= 0; i-- {
				stack = append(stack, chs[i])
			}
			continue
		}
		if hw, ok := any(cur).(interface{ Children() []Widget }); ok {
			chs := hw.Children()
			for i := len(chs) - 1; i >= 0; i-- {
				stack = append(stack, chs[i])
			}
		}
	}

	return func() {
		for i := len(saved) - 1; i >= 0; i-- {
			saved[i].b.Rect = saved[i].rec
		}
	}
}

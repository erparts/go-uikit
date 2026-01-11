package ui

// Rect is an integer rectangle in screen pixels.
type Rect struct {
	X, Y, W, H int
}

func (r Rect) Right() int  { return r.X + r.W }
func (r Rect) Bottom() int { return r.Y + r.H }

func (r Rect) Inset(pxX, pxY int) Rect {
	w := r.W - pxX*2
	h := r.H - pxY*2
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return Rect{X: r.X + pxX, Y: r.Y + pxY, W: w, H: h}
}

func (r Rect) Contains(pxX, pxY int) bool {
	return pxX >= r.X && pxX < r.Right() && pxY >= r.Y && pxY < r.Bottom()
}

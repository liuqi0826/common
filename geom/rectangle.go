package geom

type Rectangle struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func (this *Rectangle) Rectangle(x, y, width, height float32) {
	this.X = x
	this.Y = y
	this.Width = width
	this.Height = height
}
func (this *Rectangle) IntersectPoint(p *Vector2) bool {
	if p.X >= this.X && p.X <= this.X+this.Width && p.Y >= this.Y && p.Y <= this.Y+this.Height {
		return true
	}
	return false
}
func (this *Rectangle) Clone() Rectangle {
	return Rectangle{X: this.X, Y: this.Y, Width: this.Width, Height: this.Height}
}

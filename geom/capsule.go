package geom

type Capsule struct {
	Radius float32
	Height float32
}

func (this *Capsule) Capsule(radius, height float32) {
	this.Radius = radius
	this.Height = height
}

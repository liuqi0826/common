package geom

//++++++++++++++++++++ Capsule ++++++++++++++++++++
type Capsule struct {
	Radius float32
	Height float32
}

func (this *Capsule) Constructor(radius, height float32) {
	this.Radius = radius
	this.Height = height
}

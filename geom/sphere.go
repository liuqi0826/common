package geom

//++++++++++++++++++++ Sphere ++++++++++++++++++++
type Sphere struct {
	Center Vector4
	Radius float32
}

func (this *Sphere) Constructor(center *Vector4, radius float32) {
	this.Center = center.Clone()
	this.Radius = radius
}
func (this *Sphere) Clone() (s Sphere) {
	var c = this.Center.Clone()
	s = Sphere{}
	s.Constructor(&c, this.Radius)
	return s
}

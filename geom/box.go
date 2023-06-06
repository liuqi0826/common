package geom

//  4----5
// /|   /|
//0----1 |
//| 7--|-6
//3----2

//++++++++++++++++++++ Box ++++++++++++++++++++
type Box struct {
	Mix Vector4
	Max Vector4
}

func (this *Box) Constructor(mix, max *Vector4) {
	this.Mix = mix.Clone()
	this.Max = max.Clone()
}
func (this *Box) GetCenter() (center *Vector4) {
	center = &Vector4{
		X: this.Max.X - this.Mix.X,
		Y: this.Max.Y - this.Mix.Y,
		Z: this.Max.Z - this.Mix.Z,
		W: 1.0,
	}
	return
}
func (this *Box) GetRadius() float32 {
	return Vector4Distance(&this.Max, &this.Mix)
}
func (this *Box) GetVertex() (list [8]*Vector4) {
	list = [8]*Vector4{}
	list[0] = &Vector4{X: this.Mix.X, Y: this.Max.Y, Z: this.Mix.Z}
	list[1] = &Vector4{X: this.Max.X, Y: this.Max.Y, Z: this.Mix.Z}
	list[2] = &Vector4{X: this.Max.X, Y: this.Mix.Y, Z: this.Mix.Z}
	list[3] = &Vector4{X: this.Mix.X, Y: this.Mix.Y, Z: this.Mix.Z}
	list[4] = &Vector4{X: this.Mix.X, Y: this.Max.Y, Z: this.Max.Z}
	list[5] = &Vector4{X: this.Max.X, Y: this.Max.Y, Z: this.Max.Z}
	list[6] = &Vector4{X: this.Max.X, Y: this.Mix.Y, Z: this.Max.Z}
	list[7] = &Vector4{X: this.Mix.X, Y: this.Mix.Y, Z: this.Max.Z}
	return
}

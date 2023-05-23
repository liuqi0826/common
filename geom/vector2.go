package geom

import (
	"math"
)

type Vector2 struct {
	X float32
	Y float32
}

func (this *Vector2) Vector2() {
	this.X = 0.0
	this.Y = 0.0
}
func (this *Vector2) Length() float32 {
	return float32(math.Sqrt(float64(this.X*this.X + this.Y*this.Y)))
}
func (this *Vector2) Normalize() {
	magSQ := this.X*this.X + this.Y*this.Y
	if magSQ > 0.0 {
		oneOverMag := 1.0 / float32(math.Sqrt(float64(magSQ)))
		this.X = this.X * oneOverMag
		this.Y = this.Y * oneOverMag
	}
}
func (this *Vector2) Add(v *Vector2) {
	this.X += v.X
	this.Y += v.Y
}
func (this *Vector2) Sub(v *Vector2) {
	this.X -= v.X
	this.Y -= v.Y
}

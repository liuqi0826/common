package geom

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// ++++++++++++++++++++ Matrix3x3 ++++++++++++++++++++
type Matrix3x3 struct {
	A  float32
	B  float32
	C  float32
	D  float32
	TX float32
	TY float32
}

func (this *Matrix3x3) Constructor() {
	this.Identity()
}
func (this *Matrix3x3) Identity() {
	this.A = 1.0
	this.B = 0.0
	this.C = 0.0
	this.D = 1.0
	this.TX = 0.0
	this.TY = 0.0
}
func (this *Matrix3x3) Invert() {
	var norm = this.A*this.D - this.B*this.C
	if norm == 0 {
		this.A = 0
		this.B = 0
		this.C = 0
		this.D = 0
		this.TX = -this.TX
		this.TY = -this.TY
	} else {
		norm = 1.0 / norm
		var a1 = this.D * norm
		this.D = this.A * norm
		this.A = a1
		this.B *= -norm
		this.C *= -norm
		var tx1 = -this.A*this.TX - this.C*this.TY
		this.TY = -this.B*this.TX - this.D*this.TY
		this.TX = tx1
	}
}
func (this *Matrix3x3) Rotate(theta float32) {
	var cos = float32(math.Cos(float64(theta)))
	var sin = float32(math.Sin(float64(theta)))

	var a1 = this.A*cos - this.B*sin
	this.B = this.A*sin + this.B*cos
	this.A = a1

	var c1 = this.C*cos - this.D*sin
	this.D = this.C*sin + this.D*cos
	this.C = c1

	var tx1 = this.TX*cos - this.TY*sin
	this.TY = this.TX*sin + this.TY*cos
	this.TX = tx1
}
func (this *Matrix3x3) Scale(sx, sy float32) {
	this.A *= sx
	this.B *= sy
	this.C *= sx
	this.D *= sy
	this.TX *= sx
	this.TY *= sy
}
func (this *Matrix3x3) Translate(dx, dy float32) {
	this.TX += dx
	this.TY += dy
}
func (this *Matrix3x3) TransformPoint(point *Vector2) *Vector2 {
	return &Vector2{X: point.X*this.A + point.Y*this.C + this.TX, Y: point.X*this.B + point.Y*this.D + this.TY}
}
func (this *Matrix3x3) DeltaTransformPoint(point *Vector2) *Vector2 {
	return &Vector2{X: point.X*this.A + point.Y*this.C, Y: point.X*this.B + point.Y*this.D}
}
func (this *Matrix3x3) Append(lhs *Matrix3x3) {
	var a1 = this.A*lhs.A + this.B*lhs.C
	this.B = this.A*lhs.B + this.B*lhs.D
	this.A = a1

	var c1 = this.C*lhs.A + this.D*lhs.C
	this.D = this.C*lhs.B + this.D*lhs.D
	this.C = c1

	var tx1 = this.TX*lhs.A + this.TY*lhs.C + lhs.TX
	this.TY = this.TX*lhs.B + this.TY*lhs.D + lhs.TY
	this.TX = tx1
}
func (this *Matrix3x3) Clone() (mtx Matrix3x3) {
	mtx = Matrix3x3{}
	mtx.A = this.A
	mtx.B = this.B
	mtx.C = this.C
	mtx.D = this.D
	mtx.TX = this.TX
	mtx.TY = this.TY
	return
}
func (this *Matrix3x3) ToBinary() (byteArray []byte) {
	var buff = bytes.NewBuffer(byteArray)
	binary.Write(buff, binary.BigEndian, this.A)
	binary.Write(buff, binary.BigEndian, this.B)
	binary.Write(buff, binary.BigEndian, this.C)
	binary.Write(buff, binary.BigEndian, this.D)
	binary.Write(buff, binary.BigEndian, this.TX)
	binary.Write(buff, binary.BigEndian, this.TY)
	return
}
func (this *Matrix3x3) ToString() string {
	return "Matrix3x3[A:" + fmt.Sprint(this.A) + ",B:" + fmt.Sprint(this.B) + ",C:" + fmt.Sprint(this.C) + ",D:" + fmt.Sprint(this.D) + ",TX:" + fmt.Sprint(this.TX) + ",TY:" + fmt.Sprint(this.TY) + "]"
}

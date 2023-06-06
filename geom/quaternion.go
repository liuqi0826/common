package geom

import (
	"math"
)

// ++++++++++++++++++++ Quaternion ++++++++++++++++++++
type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

func (this *Quaternion) Constructor() {
	this.X = 0.0
	this.Y = 0.0
	this.Z = 0.0
	this.W = 1.0
}
func (this *Quaternion) Magnitude() float32 {
	return float32(math.Sqrt(float64(this.X*this.X + this.Y*this.Y + this.Z*this.Z + this.W*this.W)))
}
func (this *Quaternion) Normalize() {
	var mag = 1.0 / float32(math.Sqrt(float64(this.X*this.X+this.Y*this.Y+this.Z*this.Z+this.W*this.W)))
	this.X *= mag
	this.Y *= mag
	this.Z *= mag
	this.W *= mag
}
func (this *Quaternion) Append(q *Quaternion) {
	var x1 = this.X
	var y1 = this.Y
	var z1 = this.Z
	var w1 = this.W
	var x2 = q.X
	var y2 = q.Y
	var z2 = q.Z
	var w2 = q.W
	this.W = w1*w2 - x1*x2 - y1*y2 - z1*z2
	this.X = w1*x2 + x1*w2 + y1*z2 - z1*y2
	this.Y = w1*y2 - x1*z2 + y1*w2 + z1*x2
	this.Z = w1*z2 + x1*y2 - y1*x2 + z1*w2
}
func (this *Quaternion) MultiplyVector(v *Vector4) (target Quaternion) {
	var vx, vy, vz = v.X, v.Y, v.Z
	target = Quaternion{
		W: -this.X*vx - this.Y*vy - this.Z*vz,
		X: this.W*vx + this.Y*vz - this.Z*vy,
		Y: this.W*vy - this.X*vz + this.Z*vx,
		Z: this.W*vz + this.X*vy - this.Y*vx,
	}
	return
}
func (this *Quaternion) FromAxisAngle(axis *Vector4, angle float32) {
	var sin_a = float32(math.Sin(float64(angle * .5)))
	var cos_a = float32(math.Cos(float64(angle * .5)))
	this.X = axis.X * sin_a
	this.Y = axis.Y * sin_a
	this.Z = axis.Z * sin_a
	this.W = cos_a
	this.Normalize()
}
func (this *Quaternion) ToAxisAngle(axis *Vector4) (angle float32) {
	var sqrLength = this.X*this.X + this.Y*this.Y + this.Z*this.Z
	if sqrLength > 0.0 {
		angle = float32(2.0 * math.Acos(float64(this.W)))
		sqrLength = float32(1.0 / math.Sqrt(float64(sqrLength)))
		axis.X = this.X * sqrLength
		axis.Y = this.Y * sqrLength
		axis.Z = this.Z * sqrLength
	} else {
		angle = float32(0)
		axis.X = 1.0
		axis.Y = 0.0
		axis.Z = 0.0
	}
	angle = angle / CONST_PI_Over_180
	return
}
func (this *Quaternion) FromEulerAngles(ea *Vector4) {
	var eulerAngles = Vector4{X: ea.X, Y: ea.Y, Z: ea.Z, W: ea.W}
	eulerAngles.X *= DEGREES_TO_RADIANS
	eulerAngles.Y *= DEGREES_TO_RADIANS
	eulerAngles.Z *= DEGREES_TO_RADIANS
	var halfX, halfY, halfZ = eulerAngles.X * .5, eulerAngles.Y * .5, eulerAngles.Z * .5
	var cosX, sinX = float32(math.Cos(float64(halfX))), float32(math.Sin(float64(halfX)))
	var cosY, sinY = float32(math.Cos(float64(halfY))), float32(math.Sin(float64(halfY)))
	var cosZ, sinZ = float32(math.Cos(float64(halfZ))), float32(math.Sin(float64(halfZ)))
	this.W = cosX*cosY*cosZ + sinX*sinY*sinZ
	this.X = sinX*cosY*cosZ - cosX*sinY*sinZ
	this.Y = cosX*sinY*cosZ + sinX*cosY*sinZ
	this.Z = cosX*cosY*sinZ - sinX*sinY*cosZ
}
func (this *Quaternion) ToEulerAngles() (target Vector4) {
	target = Vector4{}
	target.X = float32(math.Atan2(float64(2.0*(this.W*this.X+this.Y*this.Z)), float64(1.0-2.0*(this.X*this.X+this.Y*this.Y))))
	var temp = 2.0 * (this.W*this.Y - this.Z*this.X)
	temp = Clampf(temp, -1.0, 1.0)
	target.Y = float32(math.Asin(float64(temp)))
	target.Z = float32(math.Atan2(float64(2.0*(this.W*this.Z+this.X*this.Y)), float64(1.0-2.0*(this.Y*this.Y+this.Z*this.Z))))
	target.X *= RADIANS_TO_DEGREES
	target.Y *= RADIANS_TO_DEGREES
	target.Z *= RADIANS_TO_DEGREES
	return
}
func (this *Quaternion) FromMatrix4x4(mtx *Matrix4x4) {
	var v = mtx.Decompose(QUATERNION)[1]
	this.X = v.X
	this.Y = v.Y
	this.Z = v.Z
	this.W = v.W
}
func (this *Quaternion) ToMatrix4x4() (target Matrix4x4) {
	var rawData = [16]float32{}
	var xy2, xz2, xw2 = 2.0 * this.X * this.Y, 2.0 * this.X * this.Z, 2.0 * this.X * this.W
	var yz2, yw2, zw2 = 2.0 * this.Y * this.Z, 2.0 * this.Y * this.W, 2.0 * this.Z * this.W
	var xx, yy, zz, ww = this.X * this.X, this.Y * this.Y, this.Z * this.Z, this.W * this.W
	rawData[0] = xx - yy - zz + ww
	rawData[4] = xy2 - zw2
	rawData[8] = xz2 + yw2
	rawData[12] = 0.0
	rawData[1] = xy2 + zw2
	rawData[5] = -xx + yy - zz + ww
	rawData[9] = yz2 - xw2
	rawData[13] = 0.0
	rawData[2] = xz2 - yw2
	rawData[6] = yz2 + xw2
	rawData[10] = -xx - yy + zz + ww
	rawData[14] = 0.0
	rawData[3] = 0.0
	rawData[7] = 0.0
	rawData[11] = 0.0
	rawData[15] = 1
	target = Matrix4x4{}
	target.Constructor(&rawData)
	return
}
func (this *Quaternion) Slerp(qa *Quaternion, qb *Quaternion, t float32) {
	var w1, x1, y1, z1 = qa.W, qa.X, qa.Y, qa.Z
	var w2, x2, y2, z2 = qb.W, qb.X, qb.Y, qb.Z
	var dot = w1*w2 + x1*x2 + y1*y2 + z1*z2
	if dot < 0 {
		dot = -dot
		w2 = -w2
		x2 = -x2
		y2 = -y2
		z2 = -z2
	}
	if dot < .95 {
		var angle = float32(math.Acos(float64(dot)))
		var s = float32(1.0 / math.Sin(float64(angle)))
		var s1 = float32(math.Sin(float64(angle*(1-t)))) * s
		var s2 = float32(math.Sin(float64(angle*t))) * s
		this.W = w1*s1 + w2*s2
		this.X = x1*s1 + x2*s2
		this.Y = y1*s1 + y2*s2
		this.Z = z1*s1 + z2*s2
	} else {
		this.W = w1 + t*(w2-w1)
		this.X = x1 + t*(x2-x1)
		this.Y = y1 + t*(y2-y1)
		this.Z = z1 + t*(z2-z1)
		var len = float32(1.0 / math.Sqrt(float64(this.W*this.W+this.X*this.X+this.Y*this.Y+this.Z*this.Z)))
		this.W *= len
		this.X *= len
		this.Y *= len
		this.Z *= len
	}
}
func (this *Quaternion) Lerp(qa *Quaternion, qb *Quaternion, t float32) {
	var w1, x1, y1, z1 = qa.W, qa.X, qa.Y, qa.Z
	var w2, x2, y2, z2 = qb.W, qb.X, qb.Y, qb.Z
	if w1*w2+x1*x2+y1*y2+z1*z2 < 0 {
		w2 = -w2
		x2 = -x2
		y2 = -y2
		z2 = -z2
	}
	this.W = w1 + t*(w2-w1)
	this.X = x1 + t*(x2-x1)
	this.Y = y1 + t*(y2-y1)
	this.Z = z1 + t*(z2-z1)
	var len = float32(1.0 / math.Sqrt(float64(this.W*this.W+this.X*this.X+this.Y*this.Y+this.Z*this.Z)))
	this.W *= len
	this.X *= len
	this.Y *= len
	this.Z *= len
}
func (this *Quaternion) TransformVector(v *Vector4) (target Vector4) {
	var x2, y2, z2 = v.X, v.Y, v.Z
	var w1 = -this.X*x2 - this.Y*y2 - this.Z*z2
	var x1 = this.W*x2 + this.Y*z2 - this.Z*y2
	var y1 = this.W*y2 - this.X*z2 + this.Z*x2
	var z1 = this.W*z2 + this.X*y2 - this.Y*x2
	target = Vector4{
		X: -w1*this.X + x1*this.W - y1*this.Z + z1*this.Y,
		Y: -w1*this.Y + x1*this.Z + y1*this.W - z1*this.X,
		Z: -w1*this.Z - x1*this.Y + y1*this.X + z1*this.W,
		W: 1.0,
	}
	return
}

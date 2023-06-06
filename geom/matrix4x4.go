package geom

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func CreatePerspective(width float32, height float32, zNear float32, zFar float32) (mtx Matrix4x4) {
	mtx.Constructor(&[16]float32{2.0 * zNear / width, 0, 0, 0, 0, 2.0 * zNear / height, 0, 0, 0, 0, zFar / (zFar - zNear), zNear * zFar / (zNear - zFar), 0, 0, 1, 0})
	return
}
func CreateOrtho(width float32, height float32, zNear float32, zFar float32) (mtx Matrix4x4) {
	mtx.Constructor(&[16]float32{2.0 / width, 0, 0, 0, 0, 2.0 / height, 0, 0, 0, 0, 1 / (zFar - zNear), zNear / (zNear - zFar), 0, 0, 0, 1})
	return
}
func InterpolateMatrix4x4(m1 *Matrix4x4, m2 *Matrix4x4, percent float32) (mtx *Matrix4x4) {
	var arr = [16]float32{}
	for i := 0; i < 16; i++ {
		arr[i] = m1.Raw[i] + (m2.Raw[i]-m1.Raw[i])*percent
	}
	mtx = &Matrix4x4{}
	mtx.Constructor(&arr)
	return
}

// ++++++++++++++++++++ Matrix4x4 ++++++++++++++++++++
type Matrix4x4 struct {
	Raw [16]float32
}

func (this *Matrix4x4) Constructor(raw *[16]float32) {
	if raw != nil {
		for idx, v := range raw {
			this.Raw[idx] = v
		}
	} else {
		this.Identity()
	}
}
func (this *Matrix4x4) Identity() {
	this.Raw = [16]float32{1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0}
}
func (this *Matrix4x4) Transpose() {
	var raw = [16]float32{this.Raw[0], this.Raw[4], this.Raw[8], this.Raw[12], this.Raw[1], this.Raw[5], this.Raw[9], this.Raw[13], this.Raw[2], this.Raw[6], this.Raw[10], this.Raw[14], this.Raw[3], this.Raw[7], this.Raw[11], this.Raw[15]}
	this.Constructor(&raw)
}
func (this *Matrix4x4) Invert() bool {
	var d = this.Determinant()
	if math.Abs(float64(d)) > 0.00000000001 {
		d = 1 / d
		var data = this.GetRaw()
		this.Raw[0] = d * (data[5]*(data[10]*data[15]-data[14]*data[11]) - data[9]*(data[6]*data[15]-data[14]*data[7]) + data[13]*(data[6]*data[11]-data[10]*data[7]))
		this.Raw[1] = -d * (data[1]*(data[10]*data[15]-data[14]*data[11]) - data[9]*(data[2]*data[15]-data[14]*data[3]) + data[13]*(data[2]*data[11]-data[10]*data[3]))
		this.Raw[2] = d * (data[1]*(data[6]*data[15]-data[14]*data[7]) - data[5]*(data[2]*data[15]-data[14]*data[3]) + data[13]*(data[2]*data[7]-data[6]*data[3]))
		this.Raw[3] = -d * (data[1]*(data[6]*data[11]-data[10]*data[7]) - data[5]*(data[2]*data[11]-data[10]*data[3]) + data[9]*(data[2]*data[7]-data[6]*data[3]))
		this.Raw[4] = -d * (data[4]*(data[10]*data[15]-data[14]*data[11]) - data[8]*(data[6]*data[15]-data[14]*data[7]) + data[12]*(data[6]*data[11]-data[10]*data[7]))
		this.Raw[5] = d * (data[0]*(data[10]*data[15]-data[14]*data[11]) - data[8]*(data[2]*data[15]-data[14]*data[3]) + data[12]*(data[2]*data[11]-data[10]*data[3]))
		this.Raw[6] = -d * (data[0]*(data[6]*data[15]-data[14]*data[7]) - data[4]*(data[2]*data[15]-data[14]*data[3]) + data[12]*(data[2]*data[7]-data[6]*data[3]))
		this.Raw[7] = d * (data[0]*(data[6]*data[11]-data[10]*data[7]) - data[4]*(data[2]*data[11]-data[10]*data[3]) + data[8]*(data[2]*data[7]-data[6]*data[3]))
		this.Raw[8] = d * (data[4]*(data[9]*data[15]-data[13]*data[11]) - data[8]*(data[5]*data[15]-data[13]*data[7]) + data[12]*(data[5]*data[11]-data[9]*data[7]))
		this.Raw[9] = -d * (data[0]*(data[9]*data[15]-data[13]*data[11]) - data[8]*(data[1]*data[15]-data[13]*data[3]) + data[12]*(data[1]*data[11]-data[9]*data[3]))
		this.Raw[10] = d * (data[0]*(data[5]*data[15]-data[13]*data[7]) - data[4]*(data[1]*data[15]-data[13]*data[3]) + data[12]*(data[1]*data[7]-data[5]*data[3]))
		this.Raw[11] = -d * (data[0]*(data[5]*data[11]-data[9]*data[7]) - data[4]*(data[1]*data[11]-data[9]*data[3]) + data[8]*(data[1]*data[7]-data[5]*data[3]))
		this.Raw[12] = -d * (data[4]*(data[9]*data[14]-data[13]*data[10]) - data[8]*(data[5]*data[14]-data[13]*data[6]) + data[12]*(data[5]*data[10]-data[9]*data[6]))
		this.Raw[13] = d * (data[0]*(data[9]*data[14]-data[13]*data[10]) - data[8]*(data[1]*data[14]-data[13]*data[2]) + data[12]*(data[1]*data[10]-data[9]*data[2]))
		this.Raw[14] = -d * (data[0]*(data[5]*data[14]-data[13]*data[6]) - data[4]*(data[1]*data[14]-data[13]*data[2]) + data[12]*(data[1]*data[6]-data[5]*data[2]))
		this.Raw[15] = d * (data[0]*(data[5]*data[10]-data[9]*data[6]) - data[4]*(data[1]*data[10]-data[9]*data[2]) + data[8]*(data[1]*data[6]-data[5]*data[2]))
		return true
	}
	return false
}
func (this *Matrix4x4) Append(lhs *Matrix4x4) {
	if lhs == nil {
		return
	}
	var data [16]float32 = [16]float32{this.Raw[0], this.Raw[1], this.Raw[2], this.Raw[3], this.Raw[4], this.Raw[5], this.Raw[6], this.Raw[7], this.Raw[8], this.Raw[9], this.Raw[10], this.Raw[11], this.Raw[12], this.Raw[13], this.Raw[14], this.Raw[15]}
	this.Raw[0] = data[0]*lhs.Raw[0] + data[1]*lhs.Raw[4] + data[2]*lhs.Raw[8] + data[3]*lhs.Raw[12]
	this.Raw[1] = data[0]*lhs.Raw[1] + data[1]*lhs.Raw[5] + data[2]*lhs.Raw[9] + data[3]*lhs.Raw[13]
	this.Raw[2] = data[0]*lhs.Raw[2] + data[1]*lhs.Raw[6] + data[2]*lhs.Raw[10] + data[3]*lhs.Raw[14]
	this.Raw[3] = data[0]*lhs.Raw[3] + data[1]*lhs.Raw[7] + data[2]*lhs.Raw[11] + data[3]*lhs.Raw[15]
	this.Raw[4] = data[4]*lhs.Raw[0] + data[5]*lhs.Raw[4] + data[6]*lhs.Raw[8] + data[7]*lhs.Raw[12]
	this.Raw[5] = data[4]*lhs.Raw[1] + data[5]*lhs.Raw[5] + data[6]*lhs.Raw[9] + data[7]*lhs.Raw[13]
	this.Raw[6] = data[4]*lhs.Raw[2] + data[5]*lhs.Raw[6] + data[6]*lhs.Raw[10] + data[7]*lhs.Raw[14]
	this.Raw[7] = data[4]*lhs.Raw[3] + data[5]*lhs.Raw[7] + data[6]*lhs.Raw[11] + data[7]*lhs.Raw[15]
	this.Raw[8] = data[8]*lhs.Raw[0] + data[9]*lhs.Raw[4] + data[10]*lhs.Raw[8] + data[11]*lhs.Raw[12]
	this.Raw[9] = data[8]*lhs.Raw[1] + data[9]*lhs.Raw[5] + data[10]*lhs.Raw[9] + data[11]*lhs.Raw[13]
	this.Raw[10] = data[8]*lhs.Raw[2] + data[9]*lhs.Raw[6] + data[10]*lhs.Raw[10] + data[11]*lhs.Raw[14]
	this.Raw[11] = data[8]*lhs.Raw[3] + data[9]*lhs.Raw[7] + data[10]*lhs.Raw[11] + data[11]*lhs.Raw[15]
	this.Raw[12] = data[12]*lhs.Raw[0] + data[13]*lhs.Raw[4] + data[14]*lhs.Raw[8] + data[15]*lhs.Raw[12]
	this.Raw[13] = data[12]*lhs.Raw[1] + data[13]*lhs.Raw[5] + data[14]*lhs.Raw[9] + data[15]*lhs.Raw[13]
	this.Raw[14] = data[12]*lhs.Raw[2] + data[13]*lhs.Raw[6] + data[14]*lhs.Raw[10] + data[15]*lhs.Raw[14]
	this.Raw[15] = data[12]*lhs.Raw[3] + data[13]*lhs.Raw[7] + data[14]*lhs.Raw[11] + data[15]*lhs.Raw[15]
}
func (this *Matrix4x4) AppendRotation(degrees float32, axis *Vector4, pivotPoint *Vector4) {
	var tx, ty, tz = float32(0), float32(0), float32(0)
	if pivotPoint != nil {
		tx, ty, tz = pivotPoint.X, pivotPoint.Y, pivotPoint.Z
	}
	var radian = degrees * RADIANS_TO_DEGREES
	var cos = float32(math.Cos(float64(radian)))
	var sin = float32(math.Sin(float64(radian)))
	var x, y, z = axis.X, axis.Y, axis.Z
	var x2, y2, z2 = x * x, y * y, z * z
	var ls = x2 + y2 + z2
	if ls != 0 {
		var l = float32(math.Sqrt(float64(ls)))
		x, y, z, x2, y2, z2 = x/l, y/l, z/l, x2/ls, y2/ls, z2/ls
	}
	var ccos = 1 - cos
	var mtx = new(Matrix4x4)
	mtx.Raw[0] = x2 + (y2+z2)*cos
	mtx.Raw[1] = x*y*ccos + z*sin
	mtx.Raw[2] = x*z*ccos - y*sin
	mtx.Raw[4] = x*y*ccos - z*sin
	mtx.Raw[5] = y2 + (x2+z2)*cos
	mtx.Raw[6] = y*z*ccos + x*sin
	mtx.Raw[8] = x*z*ccos + y*sin
	mtx.Raw[9] = y*z*ccos - x*sin
	mtx.Raw[10] = z2 + (x2+y2)*cos
	mtx.Raw[12] = (tx*(y2+z2)-x*(ty*y+tz*z))*ccos + (ty*z-tz*y)*sin
	mtx.Raw[13] = (ty*(x2+z2)-y*(tx*x+tz*z))*ccos + (tz*x-tx*z)*sin
	mtx.Raw[14] = (tz*(x2+y2)-z*(tx*x+ty*y))*ccos + (tx*y-ty*x)*sin
	this.Append(mtx)
}
func (this *Matrix4x4) AppendScale(xScale float32, yScale float32, zScale float32) {
	var mtx = new(Matrix4x4)
	mtx.Constructor(nil)
	mtx.Raw[0] = xScale
	mtx.Raw[5] = yScale
	mtx.Raw[10] = zScale
	this.Append(mtx)
}
func (this *Matrix4x4) AppendTranslation(x float32, y float32, z float32) {
	this.Raw[12] += x
	this.Raw[13] += y
	this.Raw[14] += z
}
func (this *Matrix4x4) Prepend(rhs *Matrix4x4) {
	if rhs == nil {
		return
	}
	var data [16]float32 = [16]float32{this.Raw[0], this.Raw[1], this.Raw[2], this.Raw[3], this.Raw[4], this.Raw[5], this.Raw[6], this.Raw[7], this.Raw[8], this.Raw[9], this.Raw[10], this.Raw[11], this.Raw[12], this.Raw[13], this.Raw[14], this.Raw[15]}
	this.Raw[0] = rhs.Raw[0]*data[0] + rhs.Raw[1]*data[4] + rhs.Raw[2]*data[8] + rhs.Raw[3]*data[12]
	this.Raw[1] = rhs.Raw[0]*data[1] + rhs.Raw[1]*data[5] + rhs.Raw[2]*data[9] + rhs.Raw[3]*data[13]
	this.Raw[2] = rhs.Raw[0]*data[2] + rhs.Raw[1]*data[6] + rhs.Raw[2]*data[10] + rhs.Raw[3]*data[14]
	this.Raw[3] = rhs.Raw[0]*data[3] + rhs.Raw[1]*data[7] + rhs.Raw[2]*data[11] + rhs.Raw[3]*data[15]
	this.Raw[4] = rhs.Raw[4]*data[0] + rhs.Raw[5]*data[4] + rhs.Raw[6]*data[8] + rhs.Raw[7]*data[12]
	this.Raw[5] = rhs.Raw[4]*data[1] + rhs.Raw[5]*data[5] + rhs.Raw[6]*data[9] + rhs.Raw[7]*data[13]
	this.Raw[6] = rhs.Raw[4]*data[2] + rhs.Raw[5]*data[6] + rhs.Raw[6]*data[10] + rhs.Raw[7]*data[14]
	this.Raw[7] = rhs.Raw[4]*data[3] + rhs.Raw[5]*data[7] + rhs.Raw[6]*data[11] + rhs.Raw[7]*data[15]
	this.Raw[8] = rhs.Raw[8]*data[0] + rhs.Raw[9]*data[4] + rhs.Raw[10]*data[8] + rhs.Raw[11]*data[12]
	this.Raw[9] = rhs.Raw[8]*data[1] + rhs.Raw[9]*data[5] + rhs.Raw[10]*data[9] + rhs.Raw[11]*data[13]
	this.Raw[10] = rhs.Raw[8]*data[2] + rhs.Raw[9]*data[6] + rhs.Raw[10]*data[10] + rhs.Raw[11]*data[14]
	this.Raw[11] = rhs.Raw[8]*data[3] + rhs.Raw[9]*data[7] + rhs.Raw[10]*data[11] + rhs.Raw[11]*data[15]
	this.Raw[12] = rhs.Raw[12]*data[0] + rhs.Raw[13]*data[4] + rhs.Raw[14]*data[8] + rhs.Raw[15]*data[12]
	this.Raw[13] = rhs.Raw[12]*data[1] + rhs.Raw[13]*data[5] + rhs.Raw[14]*data[9] + rhs.Raw[15]*data[13]
	this.Raw[14] = rhs.Raw[12]*data[2] + rhs.Raw[13]*data[6] + rhs.Raw[14]*data[10] + rhs.Raw[15]*data[14]
	this.Raw[15] = rhs.Raw[12]*data[3] + rhs.Raw[13]*data[7] + rhs.Raw[14]*data[11] + rhs.Raw[15]*data[15]
}
func (this *Matrix4x4) PrependRotation(degrees float32, axis *Vector4, pivotPoint *Vector4) {
	var tx, ty, tz = float32(0), float32(0), float32(0)
	if pivotPoint != nil {
		tx, ty, tz = pivotPoint.X, pivotPoint.Y, pivotPoint.Z
	}
	var radian = degrees * RADIANS_TO_DEGREES
	var cos = float32(math.Cos(float64(radian)))
	var sin = float32(math.Sin(float64(radian)))
	var x, y, z = axis.X, axis.Y, axis.Z
	var x2, y2, z2 = x * x, y * y, z * z
	var ls = x2 + y2 + z2
	if ls != 0 {
		var l = float32(math.Sqrt(float64(ls)))
		x, y, z, x2, y2, z2 = x/l, y/l, z/l, x2/ls, y2/ls, z2/ls
	}
	var ccos = 1 - cos
	var mtx = new(Matrix4x4)
	mtx.Raw[0] = x2 + (y2+z2)*cos
	mtx.Raw[1] = x*y*ccos + z*sin
	mtx.Raw[2] = x*z*ccos - y*sin
	mtx.Raw[4] = x*y*ccos - z*sin
	mtx.Raw[5] = y2 + (x2+z2)*cos
	mtx.Raw[6] = y*z*ccos + x*sin
	mtx.Raw[8] = x*z*ccos + y*sin
	mtx.Raw[9] = y*z*ccos - x*sin
	mtx.Raw[10] = z2 + (x2+y2)*cos
	mtx.Raw[12] = (tx*(y2+z2)-x*(ty*y+tz*z))*ccos + (ty*z-tz*y)*sin
	mtx.Raw[13] = (ty*(x2+z2)-y*(tx*x+tz*z))*ccos + (tz*x-tx*z)*sin
	mtx.Raw[14] = (tz*(x2+y2)-z*(tx*x+ty*y))*ccos + (tx*y-ty*x)*sin
	this.Prepend(mtx)
}
func (this *Matrix4x4) PrependScale(xScale float32, yScale float32, zScale float32) {
	var mtx = new(Matrix4x4)
	mtx.Constructor(nil)
	mtx.Raw[0] = xScale
	mtx.Raw[5] = yScale
	mtx.Raw[10] = zScale
	this.Prepend(mtx)
}
func (this *Matrix4x4) PrependTranslation(x float32, y float32, z float32) {
	var mtx = new(Matrix4x4)
	mtx.Constructor(nil)
	mtx.Raw[12] = x
	mtx.Raw[13] = y
	mtx.Raw[14] = z
	this.Prepend(mtx)
}
func (this *Matrix4x4) DeltaTransformVector(v *Vector4) (target *Vector4) {
	target = &Vector4{}
	target.X = v.X*this.Raw[0] + v.Y*this.Raw[4] + v.Z*this.Raw[8] + this.Raw[3]
	target.Y = v.X*this.Raw[1] + v.Y*this.Raw[5] + v.Z*this.Raw[9] + this.Raw[7]
	target.Z = v.X*this.Raw[2] + v.Y*this.Raw[6] + v.Z*this.Raw[10] + this.Raw[11]
	target.W = v.X*this.Raw[12] + v.Y*this.Raw[13] + v.Z*this.Raw[14] + this.Raw[15]
	if target.W != 1 {
		if target.W != 0 {
			var v = 1.0 / target.W
			target.X = target.X * v
			target.Y = target.Y * v
			target.Z = target.Z * v
			target.W = 1.0
		} else {
			target.W = 1.0
		}
	}
	return
}
func (this *Matrix4x4) TransformVector(v *Vector4) (target *Vector4) {
	target = &Vector4{}
	target.X = v.X*this.Raw[0] + v.Y*this.Raw[4] + v.Z*this.Raw[8] + this.Raw[12]
	target.Y = v.X*this.Raw[1] + v.Y*this.Raw[5] + v.Z*this.Raw[9] + this.Raw[13]
	target.Z = v.X*this.Raw[2] + v.Y*this.Raw[6] + v.Z*this.Raw[10] + this.Raw[14]
	target.W = v.X*this.Raw[3] + v.Y*this.Raw[7] + v.Z*this.Raw[11] + this.Raw[15]
	if target.W != 1 {
		if target.W != 0 {
			var v = 1.0 / target.W
			target.X = target.X * v
			target.Y = target.Y * v
			target.Z = target.Z * v
			target.W = 1.0
		} else {
			target.W = 1.0
		}
	}
	return
}
func (this *Matrix4x4) TransformVectorList(vl []*Vector4) (list []*Vector4) {
	for _, v := range vl {
		var tg = this.TransformVector(v)
		list = append(list, tg)
	}
	return
}
func (this *Matrix4x4) TransformLinerVector(v *Vector4) (target *Vector4) {
	target = &Vector4{}
	target.X = v.X*this.Raw[0] + v.Y*this.Raw[4] + v.Z*this.Raw[8]
	target.Y = v.X*this.Raw[1] + v.Y*this.Raw[5] + v.Z*this.Raw[9]
	target.Z = v.X*this.Raw[2] + v.Y*this.Raw[6] + v.Z*this.Raw[10]
	target.W = v.X*this.Raw[3] + v.Y*this.Raw[7] + v.Z*this.Raw[11]
	if target.W != 1 {
		if target.W != 0 {
			var v = 1.0 / target.W
			target.X = target.X * v
			target.Y = target.Y * v
			target.Z = target.Z * v
			target.W = 1.0
		} else {
			target.W = 1.0
		}
	}
	return
}
func (this *Matrix4x4) TransformLinerVectorList(vl []*Vector4) (list []*Vector4) {
	for _, v := range vl {
		var tg = this.TransformLinerVector(v)
		list = append(list, tg)
	}
	return
}
func (this *Matrix4x4) Decompose(orientationStyle string) (list [3]*Vector4) {
	var mr = this.GetRaw()                                                                                     //var mr = m.rawData.copy();
	var pos = Vector4{X: mr[12], Y: mr[13], Z: mr[14], W: 1.0}                                                 //var pos = new Vector3D(mr[12], mr[13], mr[14]);
	mr[12], mr[13], mr[14] = 0, 0, 0                                                                           //mr[12] = 0;mr[13] = 0;mr[14] = 0;
	var scale = Vector4{}                                                                                      //var scale = new Vector3D();
	scale.X = float32(math.Sqrt(float64(mr[0]*mr[0] + mr[1]*mr[1] + mr[2]*mr[2])))                             //scale.x = Math.sqrt(mr[0] * mr[0] + mr[1] * mr[1] + mr[2] * mr[2]);
	scale.Y = float32(math.Sqrt(float64(mr[4]*mr[4] + mr[5]*mr[5] + mr[6]*mr[6])))                             //scale.y = Math.sqrt(mr[4] * mr[4] + mr[5] * mr[5] + mr[6] * mr[6]);
	scale.Z = float32(math.Sqrt(float64(mr[8]*mr[8] + mr[9]*mr[9] + mr[10]*mr[10])))                           //scale.z = Math.sqrt(mr[8] * mr[8] + mr[9] * mr[9] + mr[10] * mr[10]);
	if mr[0]*(mr[5]*mr[10]-mr[6]*mr[9])-mr[1]*(mr[4]*mr[10]-mr[6]*mr[8])+mr[2]*(mr[4]*mr[9]-mr[5]*mr[8]) < 0 { //if (mr[0] * (mr[5] * mr[10] - mr[6] * mr[9]) - mr[1] * (mr[4] * mr[10] - mr[6] * mr[8]) + mr[2] * (mr[4] * mr[9] - mr[5] * mr[8]) < 0)
		scale.Z = -scale.Z //scale.z = -scale.z;
	}
	var ox = 1 / scale.X //
	mr[0] *= ox          //mr[0] /= scale.x;
	mr[1] *= ox          //mr[1] /= scale.x;
	mr[2] *= ox          //mr[2] /= scale.x;
	var oy = 1 / scale.Y //
	mr[4] *= oy          //mr[4] /= scale.y;
	mr[5] *= oy          //mr[5] /= scale.y;
	mr[6] *= oy          //mr[6] /= scale.y;
	var oz = 1 / scale.Z //
	mr[8] *= oz          //mr[8] /= scale.z;
	mr[9] *= oz          //mr[9] /= scale.z;
	mr[10] *= oz         //mr[10] /= scale.z;
	var rot = Vector4{}  //
	if orientationStyle != AXIS_ANGLE && orientationStyle != QUATERNION && orientationStyle != EULER_ANGLES {
		orientationStyle = EULER_ANGLES
	}
	switch orientationStyle {
	case AXIS_ANGLE:
		rot.W = float32(math.Acos(float64((mr[0] + mr[5] + mr[10] - 1) * .5)))                                                         //rot.w = Math.acos((mr[0] + mr[5] + mr[10] - 1) / 2);
		var len = float32(math.Sqrt(float64((mr[6]-mr[9])*(mr[6]-mr[9]) + (mr[8]-mr[2])*(mr[8]-mr[2]) + (mr[1]-mr[4])*(mr[1]-mr[4])))) //var len = Math.sqrt((mr[6] - mr[9]) * (mr[6] - mr[9]) + (mr[8] - mr[2]) * (mr[8] - mr[2]) + (mr[1] - mr[4]) * (mr[1] - mr[4]));
		if len != 0 {                                                                                                                  //if (len != 0)
			var ol = 1 / len             //
			rot.X = (mr[6] - mr[9]) * ol //rot.x = (mr[6] - mr[9]) / len;
			rot.Y = (mr[8] - mr[2]) * ol //rot.y = (mr[8] - mr[2]) / len;
			rot.Z = (mr[1] - mr[4]) * ol //rot.z = (mr[1] - mr[4]) / len;
		} else {
			rot.X, rot.Y, rot.Z = 0, 0, 0 //rot.x = rot.y = rot.z = 0;
		}
	case QUATERNION:
		var tr = mr[0] + mr[5] + mr[10] //var tr = mr[0] + mr[5] + mr[10];
		if tr > 0 {                     //if (tr > 0)
			rot.W = float32(math.Sqrt(float64(1+tr) * .5)) //rot.w = Math.sqrt(1 + tr) / 2;
			rot.X = (mr[6] - mr[9]) / (4 * rot.W)          //rot.x = (mr[6] - mr[9]) / (4 * rot.w);
			rot.Y = (mr[8] - mr[2]) / (4 * rot.W)          //rot.y = (mr[8] - mr[2]) / (4 * rot.w);
			rot.Z = (mr[1] - mr[4]) / (4 * rot.W)          //rot.z = (mr[1] - mr[4]) / (4 * rot.w);
		} else if (mr[0] > mr[5]) && (mr[0] > mr[10]) { //else if ((mr[0] > mr[5]) && (mr[0] > mr[10]))
			rot.X = float32(math.Sqrt(float64(1+mr[0]-mr[5]-mr[10])) * .5) //rot.x = Math.sqrt(1 + mr[0] - mr[5] - mr[10]) / 2;
			rot.W = (mr[6] - mr[9]) / (4 * rot.X)                          //rot.w = (mr[6] - mr[9]) / (4 * rot.x);
			rot.Y = (mr[1] + mr[4]) / (4 * rot.X)                          //rot.y = (mr[1] + mr[4]) / (4 * rot.x);
			rot.Z = (mr[8] + mr[2]) / (4 * rot.X)                          //rot.z = (mr[8] + mr[2]) / (4 * rot.x);
		} else if mr[5] > mr[10] { //else if (mr[5] > mr[10])
			rot.Y = float32(math.Sqrt(float64(1+mr[5]-mr[0]-mr[10])) * .5) //rot.y = Math.sqrt(1 + mr[5] - mr[0] - mr[10]) / 2;
			rot.X = (mr[1] + mr[4]) / (4 * rot.Y)                          //rot.x = (mr[1] + mr[4]) / (4 * rot.y);
			rot.W = (mr[8] - mr[2]) / (4 * rot.Y)                          //rot.w = (mr[8] - mr[2]) / (4 * rot.y);
			rot.Z = (mr[6] + mr[9]) / (4 * rot.Y)                          //rot.z = (mr[6] + mr[9]) / (4 * rot.y);
		} else {
			rot.Z = float32(math.Sqrt(float64(1+mr[10]-mr[0]-mr[5])) * .5) //rot.z = Math.sqrt(1 + mr[10] - mr[0] - mr[5]) / 2;
			rot.X = (mr[8] + mr[2]) / (4 * rot.Z)                          //rot.x = (mr[8] + mr[2]) / (4 * rot.z);
			rot.Y = (mr[6] + mr[9]) / (4 * rot.Z)                          //rot.y = (mr[6] + mr[9]) / (4 * rot.z);
			rot.W = (mr[1] - mr[4]) / (4 * rot.Z)                          //rot.w = (mr[1] - mr[4]) / (4 * rot.z);
		}
	case EULER_ANGLES:
		rot.Y = float32(math.Asin(float64(-mr[2]))) //rot.y = Math.asin(-mr[2]);
		if mr[2] != 1 && mr[2] != -1 {              //if (mr[2] != 1 && mr[2] != -1)
			rot.X = float32(math.Atan2(float64(mr[6]), float64(mr[10]))) //rot.x = Math.atan2(mr[6], mr[10]);
			rot.Z = float32(math.Atan2(float64(mr[1]), float64(mr[0])))  //rot.z = Math.atan2(mr[1], mr[0]);
		} else {
			rot.Z = 0                                                   //rot.z = 0;
			rot.X = float32(math.Atan2(float64(mr[4]), float64(mr[5]))) //rot.x = Math.atan2(mr[4], mr[5]);
		}
	}
	list = [3]*Vector4{}
	list[0] = &pos   //vec.push(pos);
	list[1] = &rot   //vec.push(rot);
	list[2] = &scale //vec.push(scale);
	return
}
func (this *Matrix4x4) Recompose(components [3]*Vector4, orientationStyle string) bool {
	if components[2].X == 0 || components[2].Y == 0 || components[2].Z == 0 {
		return false
	}
	if orientationStyle != AXIS_ANGLE && orientationStyle != QUATERNION && orientationStyle != EULER_ANGLES {
		orientationStyle = EULER_ANGLES
	}
	this.Identity()
	var scale = [12]float32{
		components[2].X, components[2].X, components[2].X, 0, //scale[0] = scale[1] = scale[2] = components[2].x;
		components[2].Y, components[2].Y, components[2].Y, 0, //scale[4] = scale[5] = scale[6] = components[2].y;
		components[2].Z, components[2].Z, components[2].Z, 0} //scale[8] = scale[9] = scale[10] = components[2].z;
	switch orientationStyle {
	case EULER_ANGLES:
		var cx = float32(math.Cos(float64(components[1].X))) //var cx = Math.cos(components[1].x);
		var cy = float32(math.Cos(float64(components[1].Y))) //var cy = Math.cos(components[1].y);
		var cz = float32(math.Cos(float64(components[1].Z))) //var cz = Math.cos(components[1].z);
		var sx = float32(math.Sin(float64(components[1].X))) //var sx = Math.sin(components[1].x);
		var sy = float32(math.Sin(float64(components[1].Y))) //var sy = Math.sin(components[1].y);
		var sz = float32(math.Sin(float64(components[1].Z))) //var sz = Math.sin(components[1].z);
		this.Raw[0] = cy * cz * scale[0]                     //rawData[0] = cy * cz * scale[0];
		this.Raw[1] = cy * sz * scale[1]                     //rawData[1] = cy * sz * scale[1];
		this.Raw[2] = -sy * scale[2]                         //rawData[2] = -sy * scale[2];
		this.Raw[3] = 0                                      //rawData[3] = 0;
		this.Raw[4] = (sx*sy*cz - cx*sz) * scale[4]          //rawData[4] = (sx * sy * cz - cx * sz) * scale[4];
		this.Raw[5] = (sx*sy*sz + cx*cz) * scale[5]          //rawData[5] = (sx * sy * sz + cx * cz) * scale[5];
		this.Raw[6] = sx * cy * scale[6]                     //rawData[6] = sx * cy * scale[6];
		this.Raw[7] = 0                                      //rawData[7] = 0;
		this.Raw[8] = (cx*sy*cz + sx*sz) * scale[8]          //rawData[8] = (cx * sy * cz + sx * sz) * scale[8];
		this.Raw[9] = (cx*sy*sz - sx*cz) * scale[9]          //rawData[9] = (cx * sy * sz - sx * cz) * scale[9];
		this.Raw[10] = cx * cy * scale[10]                   //rawData[10] = cx * cy * scale[10];
		this.Raw[11] = 0                                     //rawData[11] = 0;
		this.Raw[12] = components[0].X                       //rawData[12] = components[0].x;
		this.Raw[13] = components[0].Y                       //rawData[13] = components[0].y;
		this.Raw[14] = components[0].Z                       //rawData[14] = components[0].z;
		this.Raw[15] = 1                                     //rawData[15] = 1;
	default:
		var x = components[1].X //var x = components[1].x;
		var y = components[1].Y //var y = components[1].y;
		var z = components[1].Z //var z = components[1].z;
		var w = components[1].W //var w = components[1].w;
		if orientationStyle == AXIS_ANGLE {
			x *= float32(math.Sin(float64(w * .5))) //x *= Math.sin(w / 2);
			y *= float32(math.Sin(float64(w * .5))) //y *= Math.sin(w / 2);
			z *= float32(math.Sin(float64(w * .5))) //z *= Math.sin(w / 2);
			w = float32(math.Cos(float64(w * .5)))  //w = Math.cos(w / 2);
		}
		this.Raw[0] = (1 - 2*y*y - 2*z*z) * scale[0]   //rawData[0] = (1 - 2 * y * y - 2 * z * z) * scale[0];
		this.Raw[1] = (2*x*y + 2*w*z) * scale[1]       //rawData[1] = (2 * x * y + 2 * w * z) * scale[1];
		this.Raw[2] = (2*x*z - 2*w*y) * scale[2]       //rawData[2] = (2 * x * z - 2 * w * y) * scale[2];
		this.Raw[3] = 0                                //rawData[3] = 0;
		this.Raw[4] = (2*x*y - 2*w*z) * scale[4]       //rawData[4] = (2 * x * y - 2 * w * z) * scale[4];
		this.Raw[5] = (1 - 2*x*x - 2*z*z) * scale[5]   //rawData[5] = (1 - 2 * x * x - 2 * z * z) * scale[5];
		this.Raw[6] = (2*y*z + 2*w*x) * scale[6]       //rawData[6] = (2 * y * z + 2 * w * x) * scale[6];
		this.Raw[7] = 0                                //rawData[7] = 0;
		this.Raw[8] = (2*x*z + 2*w*y) * scale[8]       //rawData[8] = (2 * x * z + 2 * w * y) * scale[8];
		this.Raw[9] = (2*y*z - 2*w*x) * scale[9]       //rawData[9] = (2 * y * z - 2 * w * x) * scale[9];
		this.Raw[10] = (1 - 2*x*x - 2*y*y) * scale[10] //rawData[10] = (1 - 2 * x * x - 2 * y * y) * scale[10];
		this.Raw[11] = 0                               //rawData[11] = 0;
		this.Raw[12] = components[0].X                 //rawData[12] = components[0].x;
		this.Raw[13] = components[0].Y                 //rawData[13] = components[0].y;
		this.Raw[14] = components[0].Z                 //rawData[14] = components[0].z;
		this.Raw[15] = 1                               //rawData[15] = 1;
	}
	if components[2].X == 0 {
		this.Raw[0] = 1e-15
	}
	if components[2].Y == 0 {
		this.Raw[5] = 1e-15
	}
	if components[2].Z == 0 {
		this.Raw[10] = 1e-15
	}
	return !(components[2].X == 0 || components[2].Y == 0 || components[2].Y == 0)
}
func (this *Matrix4x4) Determinant() float32 {
	return ((this.Raw[0]*this.Raw[5]-this.Raw[4]*this.Raw[1])*(this.Raw[10]*this.Raw[15]-this.Raw[14]*this.Raw[11]) - (this.Raw[0]*this.Raw[9]-this.Raw[8]*this.Raw[1])*(this.Raw[6]*this.Raw[15]-this.Raw[14]*this.Raw[7]) + (this.Raw[0]*this.Raw[13]-this.Raw[12]*this.Raw[1])*(this.Raw[6]*this.Raw[11]-this.Raw[10]*this.Raw[7]) + (this.Raw[4]*this.Raw[9]-this.Raw[8]*this.Raw[5])*(this.Raw[2]*this.Raw[15]-this.Raw[14]*this.Raw[3]) - (this.Raw[4]*this.Raw[13]-this.Raw[12]*this.Raw[5])*(this.Raw[2]*this.Raw[11]-this.Raw[10]*this.Raw[3]) + (this.Raw[8]*this.Raw[13]-this.Raw[12]*this.Raw[9])*(this.Raw[2]*this.Raw[7]-this.Raw[6]*this.Raw[3]))
}
func (this *Matrix4x4) Position() *Vector4 {
	return &Vector4{X: this.Raw[12], Y: this.Raw[13], Z: this.Raw[14], W: 1.0}
}
func (this *Matrix4x4) GetRaw() (raw [16]float32) {
	raw = [16]float32{this.Raw[0], this.Raw[1], this.Raw[2], this.Raw[3], this.Raw[4], this.Raw[5], this.Raw[6], this.Raw[7], this.Raw[8], this.Raw[9], this.Raw[10], this.Raw[11], this.Raw[12], this.Raw[13], this.Raw[14], this.Raw[15]}
	return
}
func (this *Matrix4x4) GetRawSlice() (raw []float32) {
	raw = []float32{this.Raw[0], this.Raw[1], this.Raw[2], this.Raw[3], this.Raw[4], this.Raw[5], this.Raw[6], this.Raw[7], this.Raw[8], this.Raw[9], this.Raw[10], this.Raw[11], this.Raw[12], this.Raw[13], this.Raw[14], this.Raw[15]}
	return
}
func (this *Matrix4x4) Clone() (mtx Matrix4x4) {
	mtx = Matrix4x4{}
	mtx.Constructor(&this.Raw)
	return
}
func (this *Matrix4x4) CopyTo(mtx *Matrix4x4) {
	if mtx != nil {
		for i := 0; i < 16; i++ {
			mtx.Raw[i] = this.Raw[i]
		}
	}
}
func (this *Matrix4x4) CopyFrom(mtx *Matrix4x4) {
	if mtx != nil {
		for i := 0; i < 16; i++ {
			this.Raw[i] = mtx.Raw[i]
		}
	}
}
func (this *Matrix4x4) ToArray() (array []float32) {
	array = this.Raw[:]
	return
}
func (this *Matrix4x4) ToBinary() (byteArray []byte) {
	var buff = bytes.NewBuffer(byteArray)
	for _, v := range this.Raw {
		binary.Write(buff, binary.BigEndian, v)
	}
	return
}
func (this *Matrix4x4) ToString() (str string) {
	str = "Matrix4x4[ "
	for _, v := range this.Raw {
		str += fmt.Sprintf("%f ", v)
	}
	str += "]"
	return
}

package geom

import (
	"math"
)

// 弧度至角度
const RADIANS_TO_DEGREES float32 = float32(180 / math.Pi)

// 角度至弧度
const DEGREES_TO_RADIANS float32 = float32(math.Pi / 180)

const CONST_2_PI float32 = math.Pi * 2.0
const CONST_PI_Over_2 float32 = math.Pi / 2.0
const CONST_1_Over_PI float32 = 1.0 / math.Pi
const CONST_PI_Over_180 float32 = math.Pi / 180.0
const CONST_180_Over_PI float32 = 180.0 / math.Pi

const AXIS_ANGLE string = "axisAgnle"
const QUATERNION string = "quaternion"
const EULER_ANGLES string = "eulerAngles"

const ALIGN_ANY string = "alignAny"
const ALIGN_XY_AXIS string = "alignXYAxis"
const ALIGN_YZ_AXIS string = "alignYZAxis"
const ALIGN_XZ_AXIS string = "alignXZAxis"

const FRONT string = "front"
const BACK string = "back"
const INTERSECT string = "intersect"

func FloatEqual(f0 float32, f1 float32) bool {
	return math.Abs(float64(f0-f1)) < 0.00000001
}
func Clampf(value float32, min_inclusive float32, max_inclusive float32) float32 {
	if min_inclusive > max_inclusive {
		var temp = min_inclusive
		min_inclusive = max_inclusive
		max_inclusive = temp
	}
	if value < min_inclusive {
		return min_inclusive
	}
	if value > max_inclusive {
		return max_inclusive
	}
	return value
}
func Sqrt(x float64) float64 {
	var z = 1.0
	for {
		var tmp = z - (z*z-x)/(2*z)
		if tmp == z || math.Abs(tmp-z) < 0.000000000001 {
			break
		}
		z = tmp
	}
	return z
}

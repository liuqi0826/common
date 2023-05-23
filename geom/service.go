package geom

var vector2Pool []Vector2Proxy
var vector4Pool []Vector4Proxy
var matrix3x3Pool []Matrix3x3Proxy
var matrix4x4Pool []Matrix4x4Proxy

func init() {
	vector2Pool = make([]Vector2Proxy, 1024)
	vector4Pool = make([]Vector4Proxy, 1024)
	matrix3x3Pool = make([]Matrix3x3Proxy, 1024)
	matrix4x4Pool = make([]Matrix4x4Proxy, 1024)
}
func CreateVector2() (vet *Vector2Proxy) {
	for _, v := range vector2Pool {
		if !v.locked {
			vet = &v
			v.locked = true
			break
		}
	}
	if vet == nil {
		var idx = len(vector2Pool)
		vector2Pool = append(vector2Pool, Vector2Proxy{})
		vector2Pool[idx].locked = true
		vet = &vector2Pool[idx]
	}
	return
}
func CreateVector4() (vet *Vector4Proxy) {
	for _, v := range vector4Pool {
		if !v.locked {
			vet = &v
			v.locked = true
			break
		}
	}
	if vet == nil {
		var idx = len(vector4Pool)
		vector4Pool = append(vector4Pool, Vector4Proxy{})
		vector4Pool[idx].locked = true
		vet = &vector4Pool[idx]
	}
	return
}
func CreateMatrix3x3() (mtx *Matrix3x3Proxy) {
	for _, m := range matrix3x3Pool {
		if !m.locked {
			mtx = &m
			m.locked = true
			break
		}
	}
	if mtx == nil {
		var idx = len(matrix3x3Pool)
		matrix3x3Pool = append(matrix3x3Pool, Matrix3x3Proxy{})
		matrix3x3Pool[idx].locked = true
		mtx = &matrix3x3Pool[idx]
	}
	return
}
func CreateMatrix4x4() (mtx *Matrix4x4Proxy) {
	for _, m := range matrix4x4Pool {
		if !m.locked {
			mtx = &m
			m.locked = true
			break
		}
	}
	if mtx == nil {
		var idx = len(matrix4x4Pool)
		matrix4x4Pool = append(matrix4x4Pool, Matrix4x4Proxy{})
		matrix4x4Pool[idx].locked = true
		mtx = &matrix4x4Pool[idx]
	}
	return
}

type Vector2Proxy struct {
	Vector2 Vector2
	locked  bool
}

func (this *Vector2Proxy) Dispose() {
	this.locked = false
}

type Vector4Proxy struct {
	Vector4 Vector4
	locked  bool
}

func (this *Vector4Proxy) Dispose() {
	this.locked = false
}

type Matrix3x3Proxy struct {
	Matrix3x3 Matrix3x3
	locked    bool
}

func (this *Matrix3x3Proxy) Dispose() {
	this.locked = false
}

type Matrix4x4Proxy struct {
	Matrix4x4 Matrix4x4
	locked    bool
}

func (this *Matrix4x4Proxy) Dispose() {
	this.locked = false
}

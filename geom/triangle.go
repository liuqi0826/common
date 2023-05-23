package geom

type Triangle struct {
	Vector [3]*Vector4
}

func (this *Triangle) Triangle(v0, v1, v2 *Vector4) {
	this.Vector = [3]*Vector4{}
	if v0 != nil {
		vct := v0.Clone()
		this.Vector[0] = &vct
	} else {
		this.Vector[0] = &Vector4{}
		this.Vector[0].W = 1.0
	}
	if v1 != nil {
		vct := v1.Clone()
		this.Vector[1] = &vct
	} else {
		this.Vector[1] = &Vector4{}
		this.Vector[1].W = 1.0
	}
	if v2 != nil {
		vct := v2.Clone()
		this.Vector[2] = &vct
	} else {
		this.Vector[2] = &Vector4{}
		this.Vector[2].W = 1.0
	}
}
func (this *Triangle) GetNormal() *Vector4 {
	e1, e2 := this.Vector[1].Clone(), this.Vector[2].Clone()
	e1.Subtract(this.Vector[0])
	e2.Subtract(this.Vector[1])
	normal := e1.CrossProduct(&e2)
	normal.Normalize()
	return normal
}

func (this *Triangle) PointInTriangle(point *Vector4) bool {
	var dot00, dot01, dot02, dot11, dot12 float32
	var inverDeno, u, v float32
	var vector = [3]*Vector4{}

	vector[0] = new(Vector4)
	vector[0].X = this.Vector[2].X - this.Vector[0].X
	vector[0].Y = this.Vector[2].Y - this.Vector[0].Y
	vector[0].Z = this.Vector[2].Z - this.Vector[0].Z
	vector[1] = new(Vector4)
	vector[1].X = this.Vector[1].X - this.Vector[0].X
	vector[1].Y = this.Vector[1].Y - this.Vector[0].Y
	vector[1].Z = this.Vector[1].Z - this.Vector[0].Z
	vector[2] = new(Vector4)
	vector[2].X = point.X - this.Vector[0].X
	vector[2].Y = point.Y - this.Vector[0].Y
	vector[2].Z = point.Z - this.Vector[0].Z

	dot00 = vector[0].DotProduct(vector[0])
	dot01 = vector[0].DotProduct(vector[1])
	dot02 = vector[0].DotProduct(vector[2])
	dot11 = vector[1].DotProduct(vector[1])
	dot12 = vector[1].DotProduct(vector[2])

	inverDeno = 1 / (dot00*dot11 - dot01*dot01)

	u = (dot11*dot02 - dot01*dot12) * inverDeno
	if u < 0 || u > 1 {
		return false
	}

	v = (dot00*dot12 - dot01*dot02) * inverDeno
	if v < 0 || v > 1 {
		return false
	}

	return u+v <= 1
}

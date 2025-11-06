package llarmvp

var None = Version{"none"}

type Version struct {
	Ver string
}

func (v1 Version) Equal(v2 Version) bool {
	return v1.Ver == v2.Ver
}

func (v1 Version) IsNone() bool {
	return v1.Ver == "none"
}

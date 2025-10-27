package ParametrosStructs

type ParametrosMKDisk struct {
	Size int
	Fit  string
	Unit string
	Path string
}

type ParametrosRMDisk struct {
	Path string
}

type ParametrosFDisk struct {
	Size int
	Name string
	Unit string
	Type string
	Fit  string
	Path string
}

type ParametrosMount struct {
	Path string
	Name string
}

type ParametrosMkfs struct {
	Id   string
	Type string
}

type ParametrosLogin struct {
	User string
	Pass string
	Id   string
}

type ParametrosMkUser struct {
	User string
	Pass string
	Grp  string
}

type ParametrosMkGrp struct {
	Name string
}

type ParametrosRmGrp struct {
	Name string
}

type ParametrosRmUser struct {
	User string
}

type ParametrosMkFile struct {
	Path string
	R    bool
	Cont string
	Size int
}

type ParametrosCat struct {
	ListaPath []string
}

type ParametrosMkDir struct {
	Path string
	P    bool
}

type ParametrosChGrp struct {
	User string
	Grp  string
}

type ParametrosRep struct {
	Name       string
	Path       string
	Id         string
	PathFileLs string
}

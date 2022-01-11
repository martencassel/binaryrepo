package goproxy

type GoProxy struct {
	upstream string
}

func NewGoProxy(url string) *GoProxy {
	g := &GoProxy{
		upstream: url,
	}
	return g
}

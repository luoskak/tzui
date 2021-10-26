package tzui

var (
	BaseTzComponent = func(name string, child ...ITzComponent) *TzComponent {
		return &TzComponent{Name: name, Children: child}
	}
)

type TzComponent struct {
	Name     string         `json:"name"`
	Children []ITzComponent `json:"children"`
}

type ITzComponent interface {
	IReq
	IRes
}

type IReq interface {
	// ReqStr return request struct with default value
	ReqStr() interface{}
}

type IRes interface {
	// ResStr return response struct with default value
	ResStr() interface{}
}

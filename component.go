package tzui

var (
	BaseTzComponent = func(name string, child ...ITzComponent) *TzComponent {
		return &TzComponent{Name: name, Children: child}
	}
)

type TzComponent struct {
	// 使用snacked name命名
	Name     string         `json:"name"`
	Children []ITzComponent `json:"children"`
}

func (c TzComponent) ComponentName() string {
	return c.Name
}

type ITzComponent interface {
	IReq
	IRes
	IName
}

type IName interface {
	ComponentName() string
}

type HasSourceComponent interface {
	SetSource(source string)
}

type IReq interface {
	// ReqStr return request struct with default value
	ReqStr() interface{}
}

type IRes interface {
	// ResStr return response struct with default value
	ResStr() interface{}
}

type IParseTag interface {
	ParseTag([]*Field)
}

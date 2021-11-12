package tzui

import "context"

var (
	BaseTzComponent = func(name string, child ...ITzComponent) *TzComponent {
		return &TzComponent{Name: name, Children: child}
	}
)

type TzComponent struct {
	// 使用snacked name命名
	Name     string         `json:"name"`
	Children []ITzComponent `json:"children"`
	Sources  []*TzSource    `json:"sources"`
}

type TzSource struct {
	Name        string
	URL         string
	initialized bool
	Binder      Binder  `json:"-"`
	Handler     Handler `json:"-"`
}

func (c TzComponent) ComponentName() string {
	return c.Name
}

type ITzComponent interface {
	IName
}

type IName interface {
	ComponentName() string
}

type Binder func(body []byte) (interface{}, error)

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type IParseTag interface {
	ParseTag([]*Field, *TagManager)
}

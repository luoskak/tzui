package devui

type LinkType string

const (
	_         LinkType = ""
	LT_ROUTER LinkType = "routerLink"
	LT_HREF   LinkType = "hrefLink"
)

type MenuItem struct {
	ID       int         `json:"id"`
	Title    string      `json:"title"`
	Link     string      `json:"link,omitempty"`
	LinkType LinkType    `json:"linkType,omitempty"`
	Disabled bool        `json:"disabled,omitempty"`
	MenuIcon Icon        `json:"menuIcon,omitempty"`
	Children []*MenuItem `json:"children,omitempty"`
}

func NewMenuItem() *MenuItem {
	return &MenuItem{
		LinkType: LT_ROUTER,
	}
}

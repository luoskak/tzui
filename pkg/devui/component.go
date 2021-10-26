package devui

import "gitlab.com/tz/tzui/pkg/tzui"

var (
	_ tzui.ITzComponent = DataTableTzComponent{}
)

// DataTable
var (
	DataTable = func(source string) *DataTableTzComponent {
		return &DataTableTzComponent{
			DataSourceURL: source,
		}
	}
	HeaderFixedDataTable = func(source string) *DataTableTzComponent {
		cmp := DataTable(source)
		cmp.FixHeader = true
		cmp.ContainFixHeaderHeight = true
		return cmp
	}
)

package devui

import "gitlab.com/tz/tzui"

const (
	dataTableTzComponentName = "data_table"
)

var (
	_ tzui.ITzComponent = DataTableTzComponent{}
)

// DataTable
var (
	DataTable = func() *DataTableTzComponent {
		return &DataTableTzComponent{
			TzComponent: tzui.TzComponent{Name: dataTableTzComponentName},
		}
	}
	HeaderFixedDataTable = func() *DataTableTzComponent {
		cmp := DataTable()
		cmp.FixHeader = true
		cmp.ContainFixHeaderHeight = true
		return cmp
	}
)

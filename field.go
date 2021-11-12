package tzui

import "gitlab.com/tz/tzui/pkg/utils"

type Field struct {
	fieldName string
	Tags      map[string]string
}

func (f Field) CasedName() string {
	return utils.CasedName(f.fieldName)
}

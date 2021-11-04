package tzui

type Field struct {
	fieldName string
	Tags      map[string]string
}

func (f Field) CasedName() string {
	return f.fieldName
}

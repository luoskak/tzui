package tzui

import "fmt"

type TagManager struct {
	cbuild *controllerBuilder
	tags   map[string]Tag
}

func (tm *TagManager) GetTag(name string) Tag {
	tag, ok := tm.tags[name]
	if !ok {
		return emptyTag{}
	}
	return tag
}

func (tm *TagManager) AddTag(name string, tag Tag) {
	_, duplicated := tm.tags[name]
	if duplicated {
		tm.cbuild.addError(fmt.Errorf("duplicated add tag %s", name))
		return
	}
	tm.tags[name] = tag
}

type Tag interface {
	// Check tag is valid
	IsValid(name string)
}

var (
	_ Tag = DictionaryTag{}
)

type emptyTag struct {
}

func (t emptyTag) IsValid(name string) {

}

type DictionaryTag struct {
	modifiedTags map[string]bool
	tm           *TagManager
}

func (t DictionaryTag) IsValid(name string) {
	_, ok := t.modifiedTags[name]
	if !ok {
		t.tm.cbuild.addError(fmt.Errorf("dictionary %s not modify", name))
	}
}

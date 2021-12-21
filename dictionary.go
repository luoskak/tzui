package tzui

import "gitlab.com/tz/tzui/pkg/utils"

var json = utils.JsonAPI

type TzDictionary struct {
	Name string
	URL  string
}

type TzDictionaryRequest struct {
	Current int64
}

func (*TzDictionaryRequest) Bind(body []byte) (interface{}, error) {
	req := new(TzDictionaryRequest)
	err := json.Unmarshal(body, req)
	return req, err
}

package tzui

type ITzComponentBuilder interface {
	Model() interface{}
	Build() ITzComponent
	ComponentName() string
	Sources() []*TzSource
}

// type HasSourceComponentBuilder interface {
// 	SourceURL() string
// 	Handler(ctx context.Context, req interface{}) (interface{}, error)
// }

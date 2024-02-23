package client

type optionData struct {
	ttl           uint32
	context       IContext
	correlationId string
}

type IOption func(data interface{})

func WithTTL(ttl uint32) IOption {
	return func(i interface{}) {
		if data, ok := i.(*optionData); ok {
			data.ttl = ttl
		}
	}
}

func WithCond(context IContext) IOption {
	return func(i interface{}) {
		if data, ok := i.(*optionData); ok {
			data.context = context
		}
	}
}

func WithCorrelationId(id string) IOption {
	return func(i interface{}) {
		if data, ok := i.(*optionData); ok {
			data.correlationId = id
		}
	}
}

func newOptionData(opts ...IOption) *optionData {
	data := &optionData{}
	for _, op := range opts {
		op(data)
	}
	return data
}

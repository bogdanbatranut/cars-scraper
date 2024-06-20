package autoklass

type RequestBuilder struct {
	namingMapper *AutoklassRONamingMapper
}

func NewRequestBuilder() *RequestBuilder {
	namingMapper := NewAutoklassRoNamingMapper()
	return &RequestBuilder{namingMapper: namingMapper}
}

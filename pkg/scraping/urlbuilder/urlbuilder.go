package urlbuilder

import (
	"carscraper/pkg/jobs"
)

type IURLBuilder interface {
	GetURL(job jobs.SessionJob) *string
}

type URLBuilderMapper struct {
	builders map[string]IURLBuilder
}

func NewURLBuildMapper() *URLBuilderMapper {
	return &URLBuilderMapper{builders: make(map[string]IURLBuilder)}
}

func (mapper URLBuilderMapper) AddBuilder(market string, builder IURLBuilder) {
	mapper.builders[market] = builder
}

func (mapper URLBuilderMapper) GetURLBuilder(market string) IURLBuilder {
	return mapper.builders[market]
}

package id_generator

import (
	"context"
	"id-generator/id_generate"
	"id-generator/no"
)

type generate interface {
	Acquire(ctx context.Context) (int64, error)
}

const (
	RandType     int32 = 0
	IncreaseType int32 = 1
)

func NewGenerate(persistence Persistence) (generate, error){
	t, err := id_generate.GetGenerateType()
	if err != nil {
		return nil, err
	}
	if t == RandType {
		return no.NewNoGenerate(persistence)
	}
	//todo 递增
	return nil, err
}



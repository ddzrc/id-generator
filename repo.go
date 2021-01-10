package id_generator

import "time"

type Persistence interface {
	//4个返回值，分别， 1. 从数据库中获取的数据， 这个数组长度越大，对资源层的压力越小，2. 刷新间隔， 3. 服务器中内存存放个数, 这个数越大并发越高
	GetNextNums() ([]int64, time.Duration, int64, error)
	GetGenerateType() (int32, error)
}


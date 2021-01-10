package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var IdBizMap = sync.Map{}

type BizMap struct {
	c chan int64
	lock sync.Mutex
}

var ServiceDisruptionBizName = "service_disruption"

func GetNoBiz(ctx context.Context, Url string, bizName string, capCount ...int64) (int64, error) {
	if Url == "" || bizName == "" {
		return 0, errors.New(fmt.Sprintf("url:%v or bizName:%v is null", Url, bizName))
	}
	var cBiz *BizMap
	if v, ok := IdBizMap.Load(bizName); ok {
		cBiz = v.(*BizMap)
	}
	var count int64
	if len(capCount) > 0 {
		count = capCount[0]
	}
	if count > 1000 {
		count = 1000
	}
	if count == 0 {
		count = 10
	}
	if cBiz == nil {
		cBiz = &BizMap{c: make(chan  int64, count)}
		IdBizMap.Store(bizName, cBiz)
	}

	errChan := make(chan error, 1)

	remoteGetIdFunc := func(ctx context.Context) {
		cBiz.lock.Lock()
		defer cBiz.lock.Unlock()
		if len(cBiz.c) > 0 {
			return
		}
		reqURL := fmt.Sprintf("%v/v1/orderidserv/acquire/biz?biz_name=%v&count=%v", Url, bizName, count)
		result, err := rest.NewTraceClient(ctx, reqURL, rest.ClientConfig{
			EndPointURLToken:     "",                              //URL Token
			Timeout:              10,                              //10s
			NumMaxRetries:        0,                               //最大重试次数
			RetryMaxTimeDuration: time.Duration(60) * time.Second, //总计重试时间
			TraceOption:          &rest.TraceOption{RequestHeader: true, RequestBody: true, RespBody: true},
		}).GetAndParseResult()
		if err != nil {
			errChan <- err
			return
		}
		if !result.Success {
			errChan <- fmt.Errorf("get id failed, err:%v", result.Message)
		}
		var idList []int64
		err = json.Unmarshal(result.Result, &idList)
		if err != nil {
			errChan <- err
			return
		}
		for _, v := range idList {
			cBiz.c <- v
		}

	}
	defer close(errChan)
	for {
		select {
		case id := <- cBiz.c:
			return id, nil
		case err := <- errChan:
			return 0, err
		default:
			remoteGetIdFunc(ctx)
		}
	}

}

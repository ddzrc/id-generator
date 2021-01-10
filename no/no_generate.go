package no

import (
	"context"
	"errors"
	"id-generator"
	"math/rand"
	"sync"
	"time"
)

/*
noGenerate 生成的id可以向外保留，具有可防被预测的功能， 除开网络性损耗，单核qps 50万
 */
type noGenerate struct {
	//数据库中持久化数据
	repoNumList []int64
	//阀值
	valve int
	//内存中存放的个数
	memoCount int64
	//内存中的数据
	currentChan chan int64
	//持久化
	persistence id_generator.Persistence
	lock        sync.Mutex
	forceFresh  chan int64

	ticker         *time.Ticker
	ticketInterval time.Duration
}

func NewNoGenerate(persistence id_generator.Persistence) (*noGenerate, error) {

	array, tick, memoryCount, err := persistence.GetNextNums()
	if err != nil {
		return nil, err
	}

	if memoryCount < 100 || memoryCount > 10000000 {
		memoryCount = 10000
	}

	c := make(chan int64, memoryCount*2)
	g := &noGenerate{valve: len(array), repoNumList: array, memoCount: memoryCount, persistence: persistence, currentChan: c}
	g.ticker = time.NewTicker(tick)
	go func() {
		//定时刷新当前内存的数据，增加不可预测性
		for {
			select {
			case <-g.ticker.C:
				g.fresh()
			}
		}
	}()

	return g, nil
}

func (g *noGenerate) Acquire(ctx context.Context) (int64, error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	for {
		//备用的持久化数据小于阀值，进行添加
		if len(g.repoNumList) < g.valve {
			go func() {
				err := g.addRepoNum()
				if err != nil {
					errChan <- err
				}
			}()
		}

		select {
		case num := <-g.currentChan:
			if num != 0 {
				return num, nil
			}
		case <-ctx.Done():
			return 0, errors.New("ctx done or time-out")
		case err := <-errChan:
			return 0, err
		default:
			//currentChan中没有，进行chan增加
			err := g.addMemoryNum()
			if err != nil {
				return 0, err
			}
		}
	}
}

//添加数据库中的号
func (g *noGenerate) addRepoNum() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(g.repoNumList) >= g.valve {
		return nil
	}

	array, ticketDuration, memoryCount, err := g.persistence.GetNextNums()
	if g.ticketInterval != ticketDuration {
		g.ticker = time.NewTicker(ticketDuration)
	}

	if err != nil {
		return err
	}

	if memoryCount < 100 || memoryCount > 10000000 {
		memoryCount = 10000
	}
	g.memoCount = memoryCount
	g.valve = len(array)
	g.repoNumList = append(g.repoNumList, array...)
	return nil
}

//刷新内存中的id，防止当前可预见性
func (g *noGenerate) fresh() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(g.repoNumList) == 0 {
		return errors.New("repo num is null")
	}
	current := g.repoNumList[0]
	g.repoNumList = g.repoNumList[1:]

	array := make([]int64, g.memoCount)
	for i, _ := range array {
		array[i] = int64(i)
	}
	rand.Seed(time.Now().UnixNano())

	for i, _ := range array {
		j := rand.Int63n(g.memoCount)
		array[i], array[j] = array[j], array[i]
	}
	c := make(chan int64, g.memoCount*2)
	for i, _ := range array {
		c <- current + array[i]
	}
	old := g.currentChan
	g.currentChan = c
	close(old)
	return nil
}

func (g *noGenerate) addMemoryNum() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(g.currentChan) > int(g.memoCount/5) {
		return nil
	}
	if len(g.repoNumList) == 0 {
		return errors.New("repo num is null")
	}
	current := g.repoNumList[0]
	g.repoNumList = g.repoNumList[1:]

	array := make([]int64, g.memoCount)
	for i, _ := range array {
		array[i] = int64(i)
	}
	rand.Seed(time.Now().UnixNano())
	for i, _ := range array {
		j := rand.Int63n(g.memoCount)
		array[i], array[j] = array[j], array[i]
	}

	for i, _ := range array {
		g.currentChan <- current + array[i]
	}
	return nil
}

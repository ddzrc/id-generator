package no

import (
	"context"
	"database/sql"
	"fmt"
	"id-generator/repoimpl"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestGetNum(t *testing.T) {
	cpu := runtime.NumCPU()
	fmt.Println(cpu)
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}
	p := repoimpl.NewJDBCPersistence("test", db)
	g, err := NewNoGenerate(p)
	if err != nil {
		t.Fatal(err)
	}
	getNumFunc := func() {
		result := make([]int64, 0)
		for i := 0; i <10000; i++ {
			id, err := g.Acquire(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if id == 0 {
				fmt.Println("==============")
			}
			result = append(result, id)
		}
	}
	w := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			getNumFunc()
			w.Done()
		}()
	}
	w.Wait()

}

func TestGetNum1(t *testing.T) {
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}
	p := repoimpl.NewJDBCPersistence("test", db)
	g, err := NewNoGenerate(p)
	if err != nil {
		t.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	getNumFunc := func() {
		for {
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
			id, err := g.Acquire(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(id)
		}
	}

	w := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			getNumFunc()
			w.Done()
		}()
	}
	w.Wait()

}


func TestGetNum2(t *testing.T) {
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}
	p := repoimpl.NewJDBCPersistence("test1", db)
	g, err := NewNoGenerate(p)
	if err != nil {
		t.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	getNumFunc := func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))
			id, err := g.Acquire(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(id)
		}
	}

	w := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			getNumFunc()
			w.Done()
		}()
	}
	w.Wait()

}

func TestMine(t *testing.T){

	fmt.Println(1 << 63 -1)
}
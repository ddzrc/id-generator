package id_generator

import (
	"database/sql"
	"fmt"
	"id-generator/no"
	"id-generator/repoimpl"
	"sync"
)

type IDGenerateService struct {
	IDGenerateMap  map[string]generate
	M sync.RWMutex
}


func (igf *IDGenerateService) GetIDGenerate(businessName string, db  *sql.DB) (generate, error){
	if db == nil {
		return nil, fmt.Errorf("db is null:%v", db)
	}
	var g generate
	igf.M.RLock()
	if v, ok := igf.IDGenerateMap[businessName]; ok {
		g = v
		igf.M.RUnlock()
		return v, nil
	}
	igf.M.RUnlock()
	igf.M.Lock()
	defer igf.M.Unlock()
	p := repoimpl.NewJDBCPersistence(businessName, db)
	g, err := no.NewNoGenerate(p)
	if igf.IDGenerateMap == nil {
		igf.IDGenerateMap = make(map[string]generate, 0)
	}
	igf.IDGenerateMap[businessName] = g
	if err != nil {
		return nil, err
	}
	return g, nil
}




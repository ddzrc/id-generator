package repoimpl

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type JDBCPersistence struct {
	db           *sql.DB
	biz          string
	count        int64
}

func NewJDBCPersistence(biz string, db *sql.DB) *JDBCPersistence {
	return &JDBCPersistence{db: db, biz: biz}
}

var querySql = `
	select ifnull(val, 0), ifnull(biz, null), ifnull(tick_interval, 60), ifnull(count, 1), ifnull(memory_count, 10000)
	from id_generate where biz = ? for update
`

var updateSql = `
	update id_generate set val = ?, update_time_utc = ? where biz = ?
`

/*
CREATE TABLE id_generate
(
id INTEGER AUTO_INCREMENT,
biz VARCHAR(255) UNIQUE COMMENT '业务号',
val BIGINT COMMENT '当前使用到的值',
count INTEGER COMMENT '一次从数据库中读取几个数据',
tick_interval INTEGER COMMENT '刷新间隔',
create_time_utc DATETIME ,
update_time_utc DATETIME,
memory_count INTEGER COMMENT '内存中存放的个数',
type tinyint COMMENT '0：no非递增， 1：id递增',
PRIMARY KEY (id)
);

CREATE UNIQUE INDEX biz_idx ON id_generate (biz);

insert into id_generate(biz, val, count, tick_interval, memory_count, create_time_utc) values("test1", 100000000, 10, 10000000, 10000, "2020-12-28 06:30:00")


*/

func (jdbc *JDBCPersistence) GetGenerateType() (int32, error) {

	sql := `
select 
	type
from id_generate
where 
	biz = ?
`
	var t int32
	err := jdbc.db.QueryRow(sql, jdbc.biz).Scan(&t, jdbc.biz)
	if err != nil {
		return 0, err
	}

	return t, nil

}

func (jdbc *JDBCPersistence) GetNextNums() ([]int64, time.Duration, int64, error) {
	var val int64
	var count int
	var tick int64
	var memoryCount int64
	tx, err := jdbc.db.Begin()
	if err != nil {
		return nil, 0, 0, err
	}
	defer tx.Rollback()
	err = tx.QueryRow(querySql, jdbc.biz).Scan(&val, &jdbc.biz, &tick, &count, &memoryCount)
	if err != nil {
		return nil, 0, 0, err
	}
	_, err = tx.Exec(updateSql, val + int64(count) * memoryCount, time.Now(), jdbc.biz)
	if err != nil {
		return nil, 0, 0, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, 0, 0, err
	}
	result := make([]int64, 0)
	for i := 0; i < count; i++ {
		result = append(result, val + int64(i) * memoryCount)
	}
	if tick < 100 {
		tick = 100
	}
	return result, time.Duration(tick) * time.Millisecond, memoryCount, nil
}

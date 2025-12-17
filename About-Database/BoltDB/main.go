package main

import (
	"BoltDB/src/bolt"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var db_name = "xray.db"
var bucket = []byte("bucket")
var key = []byte("foo")
var wg sync.WaitGroup

func main() {

	// 起始时间
	t_start := time.Now()

	// 获取数据库对象
	db, err := bolt.GetDB(db_name)
	if err != nil {
		fmt.Errorf("failed to get db %w", err)
	}
	defer db.Close()
	// 插入记录
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			value := []byte(strconv.Itoa(i))
			bolt.SetKV(db, bucket, key, value)
			wg.Done()
		}(i)

		// value := []byte(strconv.Itoa(i))
		// bolt.UpdateKV(db, bucket, key, value)
	}
	wg.Wait()

	// 删除前值
	val_befor, _ := bolt.GetKV(db, bucket, key)
	// 删除记录
	bolt.DelKV(db, bucket, key)
	// 删除后值
	val_after, _ := bolt.GetKV(db, bucket, key)
	// 结束时间
	t_end := time.Now()
	// 时间差
	diff := t_end.Sub(t_start)

	fmt.Printf("Total is {%v} \n", diff)
	fmt.Printf("Before delete value is {%s} \n", val_befor)
	fmt.Printf("After delete value is {%s} \n", val_after)

}

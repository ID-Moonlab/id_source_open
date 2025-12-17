package bolt

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	dbFileMode = 0600            // 文件权限
	dbTimeOut  = 3 * time.Second // 数据库超时
)

func GetDB(db_name string) (*bolt.DB, error) {
	db, err := bolt.Open(db_name, dbFileMode, &bolt.Options{Timeout: dbTimeOut})
	if err != nil {
		fmt.Errorf("failed to get db %w", err)
		return nil, err
	}
	return db, nil
}

func CloseDB(db *bolt.DB) error {
	return db.Close()
}

func GetKV(db *bolt.DB, bucket []byte, key []byte) ([]byte, error) {
	value := []byte("")
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		value = bucket.Get(key)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("get key:%v from db:%v and bucket:%v failed: %v", key, db, bucket, err)
	}

	return value, nil
}

func SetKV(db *bolt.DB, bucket []byte, key []byte, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("create bucket:%v to db:%v failed: %v", bucket, db, err)
		}

		err = bucket.Put(key, []byte(value))
		if err != nil {
			return fmt.Errorf("put key:%v,value:%v to db:%v and bucket:%v failed: %v", key, value, db, bucket, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("update db:%v failed: %v", db, err)
	}

	return nil
}

func DelDB(db_name string, db *bolt.DB) error {
	os.Remove(db_name)
	CloseDB(db)
	return nil
}

func DelBucket(db *bolt.DB, bucket []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			return fmt.Errorf("delete db:%v,bucket:%v failed: %v", db, bucket, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("delete db:%v,bucket:%v failed: %v", db, bucket, err)
	}

	return nil
}

func DelKV(db *bolt.DB, bucket []byte, key []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		err := bucket.Delete(key)
		if err != nil {
			return fmt.Errorf("delete key:%v from db:%v,bucket:%v failed: %v", key, db, bucket, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("delete key:%v from db:%v,bucket:%v failed: %v", key, db, bucket, err)
	}

	return nil
}

func BeginROTx(db *bolt.DB) (*bolt.Tx, error) {
	return db.Begin(false)
}

func BeginRWTx(db *bolt.DB) (*bolt.Tx, error) {
	return db.Begin(true)
}

func CommitTx(Tx *bolt.Tx) error {
	err := Tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx:%v failed: %v", Tx, err)
	}
	return nil
}

func RollbackTx(Tx *bolt.Tx) error {
	err := Tx.Rollback()
	if err != nil {
		return fmt.Errorf("rollback tx:%v failed: %v", Tx, err)
	}
	return nil
}

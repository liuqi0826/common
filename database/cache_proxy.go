package database

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/ssdb/gossdb/ssdb"
	"github.com/syndtr/goleveldb/leveldb"
)

type CacheProxy struct {
	sync.RWMutex

	databaseRedis   *redis.Client
	databaseSSDB    *ssdb.Client
	databaseLevelDB *leveldb.DB

	dbType string
	state  int
}

func (this *CacheProxy) Constructor() {
}
func (this *CacheProxy) ConnectRedis(address, passwrod string, database int) (err error) {
	this.databaseRedis = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: passwrod,
		DB:       database,
		PoolSize: 100,
	})
	_, err = this.databaseRedis.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("创建Redis数据库连接：", address, "::", database)
	this.dbType = DATABASE_REDIS
	this.state = STATE_CONNECTED
	return
}
func (this *CacheProxy) ConnectSSDB(address string, port int) (err error) {
	this.databaseSSDB, err = ssdb.Connect(address, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("创建SSDB数据库连接：", address, "::", port)
	this.dbType = DATABASE_SSDB
	this.state = STATE_CONNECTED
	return
}
func (this *CacheProxy) ConnectLevelDB(path string) (err error) {
	this.databaseLevelDB, err = leveldb.OpenFile(path, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("创建LevelDB数据库连接：", path)
	this.dbType = DATABASE_LEVELDB
	this.state = STATE_CONNECTED
	return
}
func (this *CacheProxy) Set(key, value string) (string, error) {
	if this.state == STATE_CONNECTED {
		switch this.dbType {
		case DATABASE_REDIS:
			return this.databaseRedis.Set(key, value, time.Millisecond*100).Result()
		case DATABASE_SSDB:
			var res, err = this.databaseSSDB.Set(key, value)
			if err != nil {
				return "", err
			}
			if result, ok := res.(string); ok {
				return result, nil
			} else {
				return "", errors.New("数据转换失败")
			}
		case DATABASE_LEVELDB:
			var err = this.databaseLevelDB.Put([]byte(key), []byte(value), nil)
			return "", err
		}
	}
	return "", errors.New("数据库连接尚未建立...")
}
func (this *CacheProxy) Get(key string) (string, error) {
	if this.state == STATE_CONNECTED {
		switch this.dbType {
		case DATABASE_REDIS:
			return this.databaseRedis.Get(key).Result()
		case DATABASE_SSDB:
			var res, err = this.databaseSSDB.Get(key)
			if err != nil {
				return "", err
			}
			if result, ok := res.(string); ok {
				return result, nil
			} else {
				return "", errors.New("数据转换失败")
			}
		case DATABASE_LEVELDB:
			var res, err = this.databaseLevelDB.Get([]byte(key), nil)
			return string(res), err
		}
	}
	return "", errors.New("数据库连接尚未建立...")
}
func (this *CacheProxy) Del(key string) (any, error) {
	if this.state == STATE_CONNECTED {
		switch this.dbType {
		case DATABASE_REDIS:
			return this.databaseRedis.Del(key).Result()
		case DATABASE_SSDB:
			return this.databaseSSDB.Del(key)
		case DATABASE_LEVELDB:
			var err = this.databaseLevelDB.Delete([]byte(key), nil)
			return nil, err
		}
	}
	return "", errors.New("数据库连接尚未建立...")
}
func (this *CacheProxy) Close() error {
	if this.state == STATE_CONNECTED {
		switch this.dbType {
		case DATABASE_REDIS:
			return this.databaseRedis.Close()
		case DATABASE_SSDB:
			return this.databaseSSDB.Close()
		case DATABASE_LEVELDB:
			return this.databaseLevelDB.Close()
		}
	}
	return errors.New("数据库连接尚未建立...")
}
func (this *CacheProxy) GetRedisContext() *redis.Client {
	if this.dbType == DATABASE_REDIS {
		return this.databaseRedis
	}
	return nil
}
func (this *CacheProxy) GetSSDBContext() *ssdb.Client {
	if this.dbType == DATABASE_SSDB {
		return this.databaseSSDB
	}
	return nil
}
func (this *CacheProxy) GetLevelDBContext() *leveldb.DB {
	if this.dbType == DATABASE_LEVELDB {
		return this.databaseLevelDB
	}
	return nil
}

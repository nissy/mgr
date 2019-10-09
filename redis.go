package mgr

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/nissy/mgr/decoder"
)

type ToRedis struct {
	SourceFile    string          `yaml:"source_file"`
	Address       string          `yaml:"address"`
	Migrates      []*MigrateRedis `yaml:"migrates"`
	conn          redis.Conn
	toDB          int
	isTo          bool
	decoderExpiry int64
	decoderDB     int
	decoderList   int
	decoder.Nop
}

type MigrateRedis struct {
	SourceDB int `yaml:"source_db"`
	ToDB     int `yaml:"to_db"`
}

func (t *ToRedis) Do() (err error) {
	f, err := os.Open(t.SourceFile)
	if err != nil {
		return err
	}
	if t.conn, err = redis.Dial("tcp", t.Address); err != nil {
		return err
	}

	defer t.conn.Close()
	return decoder.Decode(f, t)
}

func (t *ToRedis) StartDatabase(n int, offset int) error {
	t.decoderDB = n
	t.isTo = false
	for _, v := range t.Migrates {
		if t.decoderDB == v.SourceDB {
			if _, err := t.conn.Do("SELECT", v.ToDB); err != nil {
				return err
			}
			t.toDB = v.ToDB
			t.isTo = true
			return nil
		}
	}
	return nil
}

func (t *ToRedis) exec(cmd string, key []byte, args ...interface{}) (err error) {
	if t.isTo {
		if _, err = t.conn.Do(cmd, append([]interface{}{key}, args...)...); err == nil {
			fmt.Printf("SOURCE=%v TO=%v %s %s\n", t.decoderDB, t.toDB, cmd, key)
		}
	}
	return err
}

func (t *ToRedis) execExpire(key []byte) error {
	if t.decoderExpiry > 0 {
		if _, err := t.conn.Do("EXPIRE", key, t.decoderExpiry); err != nil {
			return err
		}
		t.decoderExpiry = 0
	}
	return nil
}

func (t *ToRedis) Set(key, value []byte, expiry int64) error {
	if err := t.exec("SET", key, value); err != nil {
		return err
	}
	t.decoderExpiry = expiry
	return t.execExpire(key)
}

func (t *ToRedis) StartHash(key []byte, length, expiry int64) (err error) {
	t.decoderExpiry = expiry
	return err
}

func (t *ToRedis) Hset(key, field, value []byte) error {
	if err := t.exec("HSET", key, field, value); err != nil {
		return err
	}
	return t.execExpire(key)
}

func (t *ToRedis) StartSet(key []byte, cardinality, expiry int64) (err error) {
	t.decoderExpiry = expiry
	return err
}

func (t *ToRedis) Sadd(key, member []byte) error {
	if err := t.exec("SADD", key, member); err != nil {
		return err
	}
	return t.execExpire(key)
}

func (t *ToRedis) StartList(key []byte, length, expiry int64) error {
	t.decoderList = 0
	t.decoderExpiry = expiry
	return nil
}

func (t *ToRedis) Rpush(key, value []byte) error {
	t.decoderList++
	if err := t.exec("RPUSH", key, value); err != nil {
		return err
	}
	return t.execExpire(key)
}

func (t *ToRedis) StartZSet(key []byte, cardinality, expiry int64) error {
	t.decoderList = 0
	t.decoderExpiry = expiry
	return nil
}

func (t *ToRedis) Zadd(key []byte, score float64, member []byte) error {
	t.decoderList++
	if err := t.exec("ZADD", key, score, member); err != nil {
		return err
	}
	return t.execExpire(key)
}

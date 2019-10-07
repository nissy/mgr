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
	decoderExpire int64
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

func (t *ToRedis) redisDo(cmd string, key []byte, args ...interface{}) (err error) {
	if t.isTo {
		if _, err = t.conn.Do(cmd, append([]interface{}{key}, args...)...); err == nil {
			fmt.Printf("SOURCE=%v TO=%v %s %s\n", t.decoderDB, t.toDB, cmd, key)
		}
	}
	return err
}

func (t *ToRedis) Set(key, value []byte, expire int64) error {
	if err := t.redisDo("SET", key, value); err != nil {
		return err
	}
	if expire > 0 {
		return t.redisDo("EXPIRE", key, expire)
	}
	return nil
}

func (t *ToRedis) Hset(key, field, value []byte) error {
	return t.redisDo("HSET", key, field, value)
}

func (t *ToRedis) Sadd(key, member []byte) error {
	return t.redisDo("SADD", key, member)
}

func (t *ToRedis) StartList(key []byte, length, expire int64) error {
	t.decoderList = 0
	t.decoderExpire = expire
	return nil
}

func (t *ToRedis) Rpush(key, value []byte) error {
	t.decoderList++
	if err := t.redisDo("RPUSH", key, value); err != nil {
		return err
	}
	if t.decoderExpire > 0 {
		if err := t.redisDo("EXPIRE", key, t.decoderExpire); err != nil {
			return err
		}
		t.decoderExpire = 0
	}
	return nil
}

func (t *ToRedis) StartZSet(key []byte, cardinality, expire int64) error {
	t.decoderList = 0
	t.decoderExpire = expire
	return nil
}

func (t *ToRedis) Zadd(key []byte, score float64, member []byte) error {
	t.decoderList++
	if err := t.redisDo("ZADD", key, score, member); err != nil {
		return err
	}
	if t.decoderExpire > 0 {
		if err := t.redisDo("EXPIRE", key, t.decoderExpire); err != nil {
			return err
		}
		t.decoderExpire = 0
	}
	return nil
}

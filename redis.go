package mgr

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/nissy/mgr/decoder"
)

type ToRedis struct {
	SourceFile string          `yaml:"source_file"`
	Address    string          `yaml:"address"`
	Migrates   []*MigrateRedis `yaml:"migrates"`

	isTo        bool       `yaml:"-"`
	conn        redis.Conn `yaml:"-"`
	decoderDB   int        `yaml:"-"`
	decoderList int        `yaml:"-"`
	decoder.Nop `yaml:"-"`
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
			t.isTo = true
			return nil
		}
	}

	return nil
}

func (t *ToRedis) Aux(auxkey, auxval []byte) error {
	fmt.Printf("decoderDB=%d %q -> %q\n", t.decoderDB, auxkey, auxval)
	return nil
}

func (t *ToRedis) ModuleAux(modName []byte) error {
	fmt.Printf("decoderDB=%d %q \n", t.decoderDB, modName)
	return nil
}

func (t *ToRedis) Set(key, value []byte, expire int64) error {
	if t.isTo {
		if _, err := t.conn.Do("MULTI"); err != nil {
			return err
		}
		if err := t.conn.Send("SET", key, value); err != nil {
			return err
		}
		if expire > 0 {
			if err := t.conn.Send("EXPIRE", key, expire); err != nil {
				return err
			}
		}
		if _, err := t.conn.Do("EXEC"); err != nil {
			return err
		}
	}

	return nil
}

func (t *ToRedis) Hset(key, field, value []byte) error {
	if t.isTo {
		if _, err := t.conn.Do("HSET", key, field, value); err != nil {
			return err
		}
	}

	return nil
}

func (t *ToRedis) Sadd(key, member []byte) error {
	if t.isTo {
		if _, err := t.conn.Do("SADD", key, member); err != nil {
			return err
		}
	}

	return nil
}

func (t *ToRedis) StartList(key []byte, length, expire int64) error {
	t.decoderList = 0
	return nil
}

func (t *ToRedis) Rpush(key, value []byte) error {
	if t.isTo {
		if _, err := t.conn.Do("RPUSH", key, value); err != nil {
			return err
		}
		t.decoderList++
	}

	return nil
}

func (t *ToRedis) StartZSet(key []byte, cardinality, expire int64) error {
	t.decoderList = 0
	return nil
}

func (t *ToRedis) Zadd(key []byte, score float64, member []byte) error {
	if t.isTo {
		if _, err := t.conn.Do("ZADD", key, score, member); err != nil {
			return err
		}
		t.decoderList++
	}

	return nil
}

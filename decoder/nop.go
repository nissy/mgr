package decoder

type Nop struct {
}

func (d Nop) StartRDB() (err error) {
	return err
}

func (d Nop) StartDatabase(n int, offset int) (err error) {
	return err
}

func (d Nop) Aux(key, value []byte) (err error) {
	return err
}

func (d Nop) ModuleAux(modName []byte) (err error) {
	return err
}

func (d Nop) ResizeDatabase(dbSize, expiresSize uint32) (err error) {
	return err
}

func (d Nop) EndDatabase(n int, offset int) (err error) {
	return err
}

func (d Nop) EndRDB(offset int) (err error) {
	return err
}

func (d Nop) Set(key, value []byte, expiry int64) (err error) {
	return err
}

func (d Nop) StartHash(key []byte, length, expiry int64) (err error) {
	return err
}

func (d Nop) Hset(key, field, value []byte) (err error) {
	return err
}

func (d Nop) EndHash(key []byte) (err error) {
	return err
}

func (d Nop) StartSet(key []byte, cardinality, expiry int64) (err error) {
	return err
}

func (d Nop) Sadd(key, member []byte) (err error) {
	return err
}

func (d Nop) EndSet(key []byte) (err error) {
	return err
}

func (d Nop) StartList(key []byte, length, expiry int64) (err error) {
	return err
}

func (d Nop) Rpush(key, value []byte) (err error) {
	return err
}

func (d Nop) EndList(key []byte) (err error) {
	return err
}

func (d Nop) StartZSet(key []byte, cardinality, expiry int64) (err error) {
	return err
}

func (d Nop) Zadd(key []byte, score float64, member []byte) (err error) {
	return err
}

func (d Nop) EndZSet(key []byte) (err error) {
	return err
}

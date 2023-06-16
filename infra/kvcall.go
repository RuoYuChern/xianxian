package infra

import (
	"sync"

	"github.com/dgraph-io/badger/v4"
	"xiyu.com/common"
)

type kvDb struct {
	db     *badger.DB
	mu     sync.Mutex
	seqMap map[string]*badger.Sequence
}

var kvdb *kvDb

func (kvdb *kvDb) Close() {
	if kvdb.db != nil {
		for _, v := range kvdb.seqMap {
			v.Release()
		}
		kvdb.db.Close()
	}
}

func KvOpen() {
	db, err := badger.Open(badger.DefaultOptions(common.GlbBaInfa.Conf.Infra.KvDir))
	if err != nil {
		common.Logger.Errorf("kvopen failed:%s", err.Error())
		panic(err.Error())
	}
	kvdb = &kvDb{db: db}
	common.TaddItem(kvdb)
}

func KvSeq(key string) (uint64, error) {
	kvdb.mu.Lock()
	defer kvdb.mu.Unlock()
	seq, ok := kvdb.seqMap[key]
	if ok {
		return seq.Next()
	}
	seq, err := kvdb.db.GetSequence([]byte(key), 50)
	if err != nil {
		return 0, err
	}
	kvdb.seqMap[key] = seq
	return seq.Next()
}

func KvSet(key string, value string) error {
	return KvBSet([]byte(key), []byte(value))
}

func KvBSet(key []byte, value []byte) error {
	return kvdb.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		err := txn.SetEntry(e)
		if err != nil {
			common.Logger.Infof("Set key [%s] error:%s", key, err.Error())
		}
		return err
	})
}

func KvGet(key string) ([]byte, error) {
	var value []byte
	err := kvdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		item.Value(func(val []byte) error {
			value = append(value, val...)
			return nil
		})

		return nil
	})

	if err != nil {
		common.Logger.Infof("Get key [%s] error:%s", key, err.Error())
		return nil, err
	}
	return value, nil
}

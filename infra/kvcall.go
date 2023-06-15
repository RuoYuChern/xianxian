package infra

import (
	"github.com/dgraph-io/badger/v4"
	"xiyu.com/common"
)

type kvDb struct {
	db *badger.DB
}

var kvdb *kvDb

func (k *kvDb) Close() {
	if k.db != nil {
		k.db.Close()
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

func KvSet(key string, value string) error {
	return kvdb.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value))
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

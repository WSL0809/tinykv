package standalone_storage

import (
	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct {
	db *badger.DB
}
type StandAloneReader struct {
	db *badger.DB
}
type standAloneIter struct {
    txn  *badger.Txn
    iter engine_util.DBIterator
}
func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	db := engine_util.CreateDB(conf.DBPath, false)
	return &StandAloneStorage{db: db}
}

func (s *StandAloneStorage) Start() error {
	// Your Code Here (1).
	return nil
}

func (s *StandAloneStorage) Stop() error {
	return s.db.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	return &StandAloneReader{db: s.db}, nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	wb := engine_util.WriteBatch{}
	for _, m := range batch{
		switch v := m.Data.(type){
			case storage.Put:
				wb.SetCF(v.Cf, v.Key, v.Value)
			case storage.Delete:
				wb.DeleteCF(v.Cf, v.Key)
		}
	}
	return wb.WriteToDB(s.db)
}

func (r *StandAloneReader) GetCF(cf string, key []byte) ([]byte, error) {
	txn := r.db.NewTransaction(false)
	defer txn.Discard()
	
	i, err := txn.Get(engine_util.KeyWithCF(cf, key))
	if err == badger.ErrKeyNotFound {return nil, nil}
	if err != nil {return nil, err}
	
	v, err := i.ValueCopy(nil)
	if err != nil {return nil, err}
	
	return v, nil
}

func (r *StandAloneReader) IterCF(cf string) engine_util.DBIterator {
	txn := r.db.NewTransaction(false)
	inner := engine_util.NewCFIterator(cf, txn)
	
	return &standAloneIter{txn: txn, iter: inner}
}

func (r *StandAloneReader) Close(){}
func (it *standAloneIter) Item() engine_util.DBItem { return it.iter.Item() }
func (it *standAloneIter) Valid() bool              { return it.iter.Valid() }
func (it *standAloneIter) Next()                    { it.iter.Next() }
func (it *standAloneIter) Seek(key []byte)          { it.iter.Seek(key) }
func (it *standAloneIter) Close() {
    it.iter.Close()
    it.txn.Discard()
}
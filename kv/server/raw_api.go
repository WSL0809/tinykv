package server

import (
	"context"

	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// The functions below are Server's Raw API. (implements TinyKvServer).
// Some helper methods can be found in sever.go in the current directory

// RawGet return the corresponding Get response based on RawGetRequest's CF and Key fields
func (server *Server) RawGet(_ context.Context, req *kvrpcpb.RawGetRequest) (*kvrpcpb.RawGetResponse, error) {
	reader, err := server.storage.Reader(req.Context)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	
	value, err := reader.GetCF(req.Cf, req.Key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return &kvrpcpb.RawGetResponse{NotFound: true}, nil
	}
	
	return &kvrpcpb.RawGetResponse{Value: value}, nil
}

// RawPut puts the target data into storage and returns the corresponding response
func (server *Server) RawPut(_ context.Context, req *kvrpcpb.RawPutRequest) (*kvrpcpb.RawPutResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be modified
	m := storage.Modify{Data: storage.Put{Key: req.Key, Value: req.Value, Cf: req.Cf}}
	err := server.storage.Write(req.Context, []storage.Modify{m})
	if err != nil {
		return nil, err
	}
	return &kvrpcpb.RawPutResponse{}, nil
}

// RawDelete delete the target data from storage and returns the corresponding response
func (server *Server) RawDelete(_ context.Context, req *kvrpcpb.RawDeleteRequest) (*kvrpcpb.RawDeleteResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be deleted
	m := storage.Modify{Data: storage.Delete{Key: req.Key, Cf: req.Cf}}
	err := server.storage.Write(req.Context, []storage.Modify{m})
	if err != nil {
		return nil, err
	}
	return &kvrpcpb.RawDeleteResponse{}, nil
}

// RawScan scan the data starting from the start key up to limit. and return the corresponding result
func (server *Server) RawScan(_ context.Context, req *kvrpcpb.RawScanRequest) (*kvrpcpb.RawScanResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using reader.IterCF
	resp := &kvrpcpb.RawScanResponse{}
	reader, err := server.storage.Reader(req.Context)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	
	it := reader.IterCF(req.Cf)
	defer it.Close()
	
	if req.Limit == 0 {return nil, nil}
	
	var count uint32
	for it.Seek(req.StartKey); it.Valid() && count < req.Limit; it.Next() {
		item := it.Item()
		k := item.Key()
		v, err := item.Value()
		if err != nil{
			return nil, err
		}
		
		resp.Kvs = append(resp.Kvs, &kvrpcpb.KvPair{Key: k, Value: v})
		count++
	}
	return resp, nil
}

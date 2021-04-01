package daemon

import (
	"context"
	"time"

	"github.com/powerpuffpenguin/xormcache/cache"
	"github.com/powerpuffpenguin/xormcache/module"
	grpc_cacher "github.com/powerpuffpenguin/xormcache/protocol/cacher"
	"google.golang.org/grpc/codes"
)

var startAt = time.Now()

type Cacher struct {
	module.Service
}

var emptyGetIdsResponse grpc_cacher.GetIdsResponse

func (s Cacher) GetIds(ctx context.Context, request *grpc_cacher.GetIdsRequest) (response *grpc_cacher.GetIdsResponse, e error) {
	cacher := cache.DefaultCacher()
	val := cacher.GetIds(request.TableName, request.Sql)
	if val == nil {
		e = s.Error(codes.NotFound, `not found`)
		return
	}
	ele := val.(cache.Element)
	var nothit bool
	nothit, e = s.ServeMessage(ctx, ele.Modtime, func(nobody bool) error {
		if nobody {
			response = &emptyGetIdsResponse
		} else {
			response = &grpc_cacher.GetIdsResponse{
				Data: ele.Data,
			}
		}
		return nil
	})
	if e == nil && nothit {
		s.SetHTTPCacheMaxAge(ctx, int(cacher.Expired/time.Second))
	}
	return
}

var emptyGetBeanResponse grpc_cacher.GetBeanResponse

func (s Cacher) GetBean(ctx context.Context, request *grpc_cacher.GetBeanRequest) (response *grpc_cacher.GetBeanResponse, e error) {
	cacher := cache.DefaultCacher()
	val := cacher.GetBean(request.TableName, request.Id)
	if val == nil {
		e = s.Error(codes.NotFound, `not found`)
		return
	}
	ele := val.(cache.Element)
	var nothit bool
	nothit, e = s.ServeMessage(ctx, ele.Modtime, func(nobody bool) error {
		if nobody {
			response = &emptyGetBeanResponse
		} else {
			response = &grpc_cacher.GetBeanResponse{
				Data: ele.Data,
			}
		}
		return nil
	})
	if e == nil && nothit {
		s.SetHTTPCacheMaxAge(ctx, int(cacher.Expired/time.Second))
	}
	return
}

var emptyPutIdsResponse grpc_cacher.PutIdsResponse

func (Cacher) PutIds(ctx context.Context, request *grpc_cacher.PutIdsRequest) (response *grpc_cacher.PutIdsResponse, e error) {
	cache.DefaultCacher().PutIds(request.TableName, request.Sql, request.Ids)
	response = &emptyPutIdsResponse
	return
}

var emptyPutBeanResponse grpc_cacher.PutBeanResponse

func (Cacher) PutBean(ctx context.Context, request *grpc_cacher.PutBeanRequest) (response *grpc_cacher.PutBeanResponse, e error) {
	cache.DefaultCacher().PutBean(request.TableName, request.Id, request.Obj)
	response = &emptyPutBeanResponse
	return
}

var emptyDelIdsResponse grpc_cacher.DelIdsResponse

func (Cacher) DelIds(ctx context.Context, request *grpc_cacher.DelIdsRequest) (response *grpc_cacher.DelIdsResponse, e error) {
	cache.DefaultCacher().DelIds(request.TableName, request.Sql)
	response = &emptyDelIdsResponse
	return
}

var emptyDelBeanResponse grpc_cacher.DelBeanResponse

func (Cacher) DelBean(ctx context.Context, request *grpc_cacher.DelBeanRequest) (response *grpc_cacher.DelBeanResponse, e error) {
	cache.DefaultCacher().DelBean(request.TableName, request.Id)
	response = &emptyDelBeanResponse
	return
}

var emptyClearIdsResponse grpc_cacher.ClearIdsResponse

func (Cacher) ClearIds(ctx context.Context, request *grpc_cacher.ClearIdsRequest) (response *grpc_cacher.ClearIdsResponse, e error) {
	cache.DefaultCacher().ClearIds(request.TableName)
	response = &emptyClearIdsResponse
	return
}

var emptyClearBeansResponse grpc_cacher.ClearBeansResponse

func (Cacher) ClearBeans(ctx context.Context, request *grpc_cacher.ClearBeansRequest) (response *grpc_cacher.ClearBeansResponse, e error) {
	cache.DefaultCacher().ClearBeans(request.TableName)
	response = &emptyClearBeansResponse
	return
}

var emptyDetailResponse grpc_cacher.DetailResponse

func (s Cacher) Detail(ctx context.Context, request *grpc_cacher.DetailRequest) (response *grpc_cacher.DetailResponse, e error) {
	cacher := cache.DefaultCacher()
	nothit, e := s.ServeMessage(ctx, startAt, func(nobody bool) error {
		if nobody {
			response = &emptyDetailResponse
		} else {
			response = &grpc_cacher.DetailResponse{
				MaxAge:         uint32(cacher.Expired / time.Second),
				MaxElementSize: uint32(cacher.MaxElementSize),
			}
		}
		return nil
	})
	if e == nil && nothit {
		s.SetHTTPCacheMaxAge(ctx, int(cacher.Expired/time.Second))
	}
	return
}

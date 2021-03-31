package module

import (
	"context"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var unixEpochTime = time.Unix(0, 0)

func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

// CheckNotModified if NotModified return true, only check on method == "GET" OR "HEAD" .
func (Service) CheckNotModified(method, ims string, modtime time.Time) bool {
	return checkIfModifiedSince(method, ims, modtime) == condFalse
}

// SetHTTPCacheMaxAge .
func (Service) SetHTTPCacheMaxAge(ctx context.Context, maxAge int) error {
	return grpc.SetHeader(ctx, metadata.Pairs(`Cache-Control`, `max-age=`+strconv.Itoa(maxAge)))
}

// SetHTTPCode .
func (Service) SetHTTPCode(ctx context.Context, code int) error {
	return grpc.SetHeader(ctx, metadata.Pairs(`x-http-code`, strconv.Itoa(code)))
}

// ToHTTPError .
func (s Service) ToHTTPError(ctx context.Context, id string, e error) error {
	if os.IsNotExist(e) {
		return s.Error(codes.NotFound, `not exists : `+id)
	}
	if os.IsExist(e) {
		return s.Error(codes.PermissionDenied, `already exists : `+id)
	}
	if os.IsPermission(e) {
		return s.Error(codes.PermissionDenied, `forbidden : `+id)
	}
	return s.Error(codes.Unknown, e.Error())
}

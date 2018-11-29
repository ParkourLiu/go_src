package retry

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/go-kit/kit/endpoint"
)

var r *retrier.Retrier

func init() {
	r = retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
}

func Retry() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			err = r.Run(func() error {
				resp, err = next(ctx, request)
				if err != nil {
					return err
				}
				if reflect.TypeOf(resp).Kind() == reflect.Struct {
					s := reflect.ValueOf(resp)
					errStr := s.FieldByName("Err").String()
					if errStr != "" {
						return errors.New(errStr)
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}

package datasource

import "context"

type Datasource interface {
	Get(ctx context.Context, key, field string) (string, error)
	Set(ctx context.Context, key, field string, value string) error
	Delete(ctx context.Context, key, field string) error
	DeleteAll(ctx context.Context, key string) error
	Select(ctx context.Context, key string, options *SelectOptions) ([]string, error)
	SelectAll(ctx context.Context, key string) (map[string]string, error)
	Count(ctx context.Context, key string) (int64, error)
	Incr(ctx context.Context, key, field string, value int64) error
}

type SelectOptions struct {
	Count int64  `json:"count,omitempty"`
	Filed string `json:"filed,omitempty"`
}

func NewDatasource(datasourceType string) (Datasource, error) {
	switch datasourceType {
	case "mysql":
		return NewMySQLDatasource()
	case "redis":
		return NewRedisDatasource()
	}
	panic("not implement")
}

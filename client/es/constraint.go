package es

import (
	"context"
	"encoding/json"
	"io"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/7/27
    @desc:

***************************/

type Client interface {
	IndexName(name string) Client
	Index() Index

	Not(filters ...Filter) Client
	Where(filters ...Filter) Client
	Or(filters ...Filter) Client
	OrderBy(field string, isDesc bool) Client
	Offset(uint64) Client
	Agg(aggs ...Agg) Client
	Start(uint64) Client
	Limit(uint64, uint64) Client
	Fields(...string) Client
	Search(ctx context.Context, result interface{}) (uint64, error)
	RawSQL(ctx context.Context, closer io.ReadCloser, result interface{}) (uint64, error)
}

type Filter interface {
	Term(field string, value interface{}) Filter
	Terms(field string, values ...interface{}) Filter
	Between(field string, start, end int64) Filter
	Gt(field string, value int64) Filter
	Gte(field string, value int64) Filter
	Lt(field string, value int64) Filter
	Lte(field string, value int64) Filter
	Result() []interface{}
}

type Index interface {
	Exists(ctx context.Context) (bool, error)
	Create(ctx context.Context, mapping map[string]interface{}) error
	List(ctx context.Context) ([]error, error)
	Mapping(ctx context.Context) (map[string]interface{}, error)
}

type Agg interface {
	Name(string) Agg
	//Filter(f Filter) Agg
	DateHistogram(field, interval, format, offset, timeZone string) Agg
	Distinct(field string, number int64) Agg
	//Metric(...AggDateHistogramMetric) Agg
	Result() (string, interface{})
}

type AggDateHistogramMetric interface {
	Add(name, operator string)
	Result() interface{}
}

// SearchResult index 返回数据，直接解析到对应的结构体
type SearchResult struct {
	Took    uint64      `json:"took"`
	TimeOut bool        `json:"time_out"`
	Error   interface{} `json:"error"`
	Shards  struct {
		Total   uint64 `json:"total"`
		Success uint64 `json:"success"`
		Skipped uint64 `json:"skipped"`
		Failed  uint64 `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    uint64 `json:"value"`
			Relation string `json:"relation"`
		}
		MaxScore  float64                 `json:"max_score"`
		IndexHits []SearchResultHitResult `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations"`
}

// SearchResultHitResult index 返回数据，直接解析到对应的结构体
type SearchResultHitResult struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	Id     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

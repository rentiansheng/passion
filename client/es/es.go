package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	ges "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/rentiansheng/mapper"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/11/3
    @desc:

***************************/

var (
	esClient *ges.Client
)

type es struct {
	isAgg bool
	// 	  field:desc
	sorts  []string
	fields []string
	index  string
	from   uint64
	size   uint64
	cond   cond
	agg    map[string]interface{}
}

type cond struct {
	must   []interface{}
	filter []interface{}
	should []interface{}
	not    []interface{}
	exists []interface{}
}

func SetESClient(c *ges.Client) {
	esClient = c
}

func (e es) IndexName(name string) Client {
	e.index = name
	return e
}

func (e es) Index() Index {
	//TODO implement me
	panic("implement me")
}

func (e es) Not(filters ...Filter) Client {
	for _, filter := range filters {
		e.cond.not = append(e.cond.not, filter.Result()...)
	}
	return e
}

func (e es) Where(filters ...Filter) Client {
	for _, filter := range filters {
		e.cond.must = append(e.cond.must, filter.Result()...)
	}
	return e
}

func (e es) Or(filters ...Filter) Client {
	for _, filter := range filters {
		e.cond.should = append(e.cond.should, filter.Result()...)
	}
	return e
}

func (e es) OrderBy(field string, isDesc bool) Client {
	if isDesc {
		e.sorts = append(e.sorts, field+":desc")
	} else {
		e.sorts = append(e.sorts, field+":asc")
	}
	return e
}

func (e es) Agg(aggs ...Agg) Client {
	e.isAgg = true
	for _, agg := range aggs {
		name, value := agg.Result()
		e.agg[name] = value
	}
	return e
}

func (e es) Offset(u uint64) Client {
	e.size = u
	return e
}

func (e es) Start(u uint64) Client {
	e.from = u
	return e
}

func (e es) Limit(form uint64, size uint64) Client {
	e.from, e.size = form, size
	return e
}

func (e es) Find(ctx context.Context, result interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (e es) Search(ctx context.Context, result interface{}) (uint64, error) {

	cond := esCondition{
		Query: esConditionQuery{Bool: esQueryBool{
			Must:   e.cond.must,
			Not:    e.cond.not,
			Should: e.cond.should,
		}},
		Agg: e.agg,
	}
	queryBody := &bytes.Buffer{}
	if err := json.NewEncoder(queryBody).Encode(cond); err != nil {
		return 0, fmt.Errorf("search condition build error. %s", err.Error())
	}

	searchOpts := []func(*esapi.SearchRequest){
		esClient.Search.WithContext(ctx),
		esClient.Search.WithIndex(e.index),
		esClient.Search.WithBody(queryBody),
		esClient.Search.WithSort(e.sorts...),
	}
	if e.isAgg {
		searchOpts = append(searchOpts, esClient.Search.WithSize(0))
	} else {
		searchOpts = append(searchOpts, esClient.Search.WithTrackTotalHits(true))
		searchOpts = append(searchOpts, esClient.Search.WithSourceIncludes(e.fields...))
		searchOpts = append(searchOpts, esClient.Search.WithFrom(int(e.from)))
		searchOpts = append(searchOpts, esClient.Search.WithSize(int(e.size)))
	}

	res, err := esClient.Search(
		searchOpts...,
	)

	if err != nil {
		return 0, fmt.Errorf("unexpected error when get: %s", err)
	}

	return e.parseSearchRespResult(ctx, res.Body, result)
}

func (e es) TranslateSQL(ctx context.Context, sql string) (bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (e es) RawSQL(ctx context.Context, closer io.ReadCloser, result interface{}) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (e es) GetByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (e es) UpdateByID(ctx context.Context, id string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (e es) DeleteByID(ctx context.Context, ids ...string) error {
	//TODO implement me
	panic("implement me")
}

func (e es) Count(ctx context.Context) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (e es) Query(ctx context.Context) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e es) Fields(fields ...string) Client {
	e.fields = fields
	return e
}

func ES() Client {
	return &es{
		isAgg:  false,
		sorts:  []string{},
		fields: nil,
		index:  "",
		from:   0,
		size:   0,
		cond:   cond{},
		agg:    make(map[string]interface{}, 0),
	}
}

func (e es) parseSearchRespResult(ctx context.Context, respBody io.ReadCloser, results interface{}) (uint64, error) {

	resultV := reflect.ValueOf(results)
	if resultV.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("results argument must be pointer")
	}
	resp, err := e.parseSearchRespDefaultDecode(ctx, respBody)
	if err != nil {
		return 0, err
	}

	total := uint64(0)
	if e.isAgg {
		if err := mapper.Mapper(ctx, resp.Aggregations, results); err != nil {
			return 0, err
		}

	} else {
		if resultV.Elem().Kind() != reflect.Slice {
			return 0, fmt.Errorf("results argument must be a slice address")
		}
		total = resp.Hits.Total.Value

		elemt := resultV.Elem().Type().Elem()
		slice := reflect.MakeSlice(resultV.Elem().Type(), 0, 10)
		for _, indexHit := range resp.Hits.IndexHits {
			elem := reflect.New(elemt)
			err := e.parseSearchResultIndexHit(ctx, indexHit, elem)
			if err != nil {
				return total, err
			}
			slice = reflect.Append(slice, elem.Elem())
		}
		resultV.Elem().Set(slice)
	}
	return total, nil
}

func (e es) parseSearchRespDefaultDecode(ctx context.Context, respBody io.ReadCloser) (SearchResult, error) {
	var resp SearchResult

	d := json.NewDecoder(respBody)
	d.UseNumber()
	err := d.Decode(&resp)
	if err != nil {
		return resp, err
	}
	if resp.Error != nil {
		return resp, fmt.Errorf("%s", resp.Error)
	}
	if resp.TimeOut {
		return resp, fmt.Errorf(" time_out, took: %v", resp.Took)
	}
	return resp, nil
}

func (e es) parseSearchResultIndexHit(ctx context.Context, indexHit SearchResultHitResult, elemp reflect.Value) error {

	if err := json.Unmarshal(indexHit.Source, elemp.Interface()); nil != err {
		return err
	}
	elemt := elemp.Elem().Type()
	// add _id
	//if searchOpt.id != nil {
	if elemt.Kind() == reflect.Map {
		elemp.Elem().SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(indexHit.Id))
	} else if elemt.Kind() == reflect.Struct {
		for i := 0; i < elemt.NumField(); i++ {
			if !elemp.IsValid() {
				return fmt.Errorf("struct IsValid false")
			}
			if !elemp.Elem().CanSet() {
				return fmt.Errorf("struct not allow change")
			}
			field := elemt.Field(i)
			tags := strings.Split(field.Tag.Get("json"), ",")
			for _, tag := range tags {
				if tag == "_id" || tag == "es_id" {
					elemp.Elem().Field(i).Set(reflect.ValueOf(indexHit.Id))
				}
			}

		}
	}

	return nil
}

var _ Client = (*es)(nil)

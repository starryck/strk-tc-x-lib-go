package xbgorm

import (
	"context"
	"database/sql"

	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xbctnr"
)

type ModelDAO[T any] struct {
	client *Client
	entity *T
}

func (dao *ModelDAO[T]) Model(entity *T) *ModelDAO[T] {
	dao.entity = entity
	dao.client = dao.client.Model(entity)
	return dao
}

func (dao *ModelDAO[T]) Find(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Find(data, args...)
	return dao
}

func (dao *ModelDAO[T]) FindInBatches(data any, size int, batcher ClientBatcher) *ModelDAO[T] {
	dao.client = dao.client.FindInBatches(data, size, batcher)
	return dao
}

func (dao *ModelDAO[T]) First(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.First(data, args...)
	return dao
}

func (dao *ModelDAO[T]) Last(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Last(data, args...)
	return dao
}

func (dao *ModelDAO[T]) Take(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Take(data, args...)
	return dao
}

func (dao *ModelDAO[T]) Count(number *int64) *ModelDAO[T] {
	dao.client = dao.client.Count(number)
	return dao
}

func (dao *ModelDAO[T]) Pluck(field string, data any) *ModelDAO[T] {
	dao.client = dao.client.Pluck(field, data)
	return dao
}

func (dao *ModelDAO[T]) Where(query any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Where(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Or(query any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Or(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Not(query any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Not(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Order(column any) *ModelDAO[T] {
	dao.client = dao.client.Order(column)
	return dao
}

func (dao *ModelDAO[T]) Limit(number int) *ModelDAO[T] {
	dao.client = dao.client.Limit(number)
	return dao
}

func (dao *ModelDAO[T]) Offset(offset int) *ModelDAO[T] {
	dao.client = dao.client.Offset(offset)
	return dao
}

func (dao *ModelDAO[T]) Select(query any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Select(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Omit(fields ...string) *ModelDAO[T] {
	dao.client = dao.client.Omit(fields...)
	return dao
}

func (dao *ModelDAO[T]) Distinct(args ...any) *ModelDAO[T] {
	dao.client = dao.client.Distinct(args...)
	return dao
}

func (dao *ModelDAO[T]) Joins(query string, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Joins(query, args...)
	return dao
}

func (dao *ModelDAO[T]) InnerJoins(query string, args ...any) *ModelDAO[T] {
	dao.client = dao.client.InnerJoins(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Preload(query string, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Preload(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Group(name string) *ModelDAO[T] {
	dao.client = dao.client.Group(name)
	return dao
}

func (dao *ModelDAO[T]) Having(query any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Having(query, args...)
	return dao
}

func (dao *ModelDAO[T]) Scopes(modifiers ...ClientModifier) *ModelDAO[T] {
	dao.client = dao.client.Scopes(modifiers...)
	return dao
}

func (dao *ModelDAO[T]) Clauses(expressions ...Expression) *ModelDAO[T] {
	dao.client = dao.client.Clauses(expressions...)
	return dao
}

func (dao *ModelDAO[T]) MapColumns(mappings map[string]string) *ModelDAO[T] {
	dao.client = dao.client.MapColumns(mappings)
	return dao
}

func (dao *ModelDAO[T]) Create(data any) *ModelDAO[T] {
	dao.client = dao.client.Create(data)
	return dao
}

func (dao *ModelDAO[T]) CreateInBatches(data any, size int) *ModelDAO[T] {
	dao.client = dao.client.CreateInBatches(data, size)
	return dao
}

func (dao *ModelDAO[T]) FirstOrInit(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.FirstOrInit(data, args...)
	return dao
}

func (dao *ModelDAO[T]) FirstOrCreate(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.FirstOrCreate(data, args...)
	return dao
}

func (dao *ModelDAO[T]) Save(data any) *ModelDAO[T] {
	dao.client = dao.client.Save(data)
	return dao
}

func (dao *ModelDAO[T]) Attrs(dataset ...any) *ModelDAO[T] {
	dao.client = dao.client.Attrs(dataset...)
	return dao
}

func (dao *ModelDAO[T]) Assign(dataset ...any) *ModelDAO[T] {
	dao.client = dao.client.Assign(dataset...)
	return dao
}

func (dao *ModelDAO[T]) Update(field string, value any) *ModelDAO[T] {
	dao.client = dao.client.Update(field, value)
	return dao
}

func (dao *ModelDAO[T]) Updates(data any) *ModelDAO[T] {
	dao.client = dao.client.Updates(data)
	return dao
}

func (dao *ModelDAO[T]) UpdateColumn(field string, value any) *ModelDAO[T] {
	dao.client = dao.client.UpdateColumn(field, value)
	return dao
}

func (dao *ModelDAO[T]) UpdateColumns(data any) *ModelDAO[T] {
	dao.client = dao.client.UpdateColumns(data)
	return dao
}

func (dao *ModelDAO[T]) Delete(data any, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Delete(data, args...)
	return dao
}

func (dao *ModelDAO[T]) Table(name string, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Table(name, args...)
	return dao
}

func (dao *ModelDAO[T]) Raw(sql string, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Raw(sql, args...)
	return dao
}

func (dao *ModelDAO[T]) Exec(sql string, args ...any) *ModelDAO[T] {
	dao.client = dao.client.Exec(sql, args...)
	return dao
}

func (dao *ModelDAO[T]) Row() *sql.Row {
	return dao.client.Row()
}

func (dao *ModelDAO[T]) Rows() (*sql.Rows, error) {
	return dao.client.Rows()
}

func (dao *ModelDAO[T]) Scan(data any) *ModelDAO[T] {
	dao.client = dao.client.Scan(data)
	return dao
}

func (dao *ModelDAO[T]) ScanRows(rows *sql.Rows, data any) error {
	return dao.client.ScanRows(rows, data)
}

func (dao *ModelDAO[T]) Debug() *ModelDAO[T] {
	dao.client = dao.client.Debug()
	return dao
}

func (dao *ModelDAO[T]) Unscoped() *ModelDAO[T] {
	dao.client = dao.client.Unscoped()
	return dao
}

func (dao *ModelDAO[T]) WithContext(ctx context.Context) *ModelDAO[T] {
	dao.client = dao.client.WithContext(ctx)
	return dao
}

func (dao *ModelDAO[T]) Session(config *Session) *ModelDAO[T] {
	dao.client = dao.client.Session(config)
	return dao
}

func (dao *ModelDAO[T]) Association(field string) *Association {
	return dao.client.Association(field)
}

func (dao *ModelDAO[T]) Get(key string) (any, bool) {
	return dao.client.Get(key)
}

func (dao *ModelDAO[T]) Set(key string, value any) *ModelDAO[T] {
	dao.client = dao.client.Set(key, value)
	return dao
}

func (dao *ModelDAO[T]) InstanceGet(key string) (any, bool) {
	return dao.client.InstanceGet(key)
}

func (dao *ModelDAO[T]) InstanceSet(key string, value any) *ModelDAO[T] {
	dao.client = dao.client.InstanceSet(key, value)
	return dao
}

func (dao *ModelDAO[T]) GetClient() *Client {
	return dao.client
}

func (dao *ModelDAO[T]) SetClient(client *Client) {
	dao.client = client
}

func (dao *ModelDAO[T]) GetEntity() *T {
	return dao.entity
}

func (dao *ModelDAO[T]) SetEntity(entity *T) {
	dao.Model(entity)
}

func (dao *ModelDAO[T]) GetError() error {
	return dao.client.Error
}

func (dao *ModelDAO[T]) SetError(err error) {
	dao.client.AddError(err)
}

type (
	ClientBatcher  = func(client *Client, size int) error
	ClientModifier = func(client *Client) *Client
)

type ModelRepository struct {
	client *Client
}

func (repository *ModelRepository) GetClient() *Client {
	return repository.client
}

func (repository *ModelRepository) SetClient(client *Client) {
	repository.client = client
}

type ModelService struct {
	client   *Client
	sessions *xbctnr.Deque[*Client]
}

func (service *ModelService) Initialize() {
	service.sessions = &xbctnr.Deque[*Client]{}
}

func (service *ModelService) GetClient() *Client {
	return service.client
}

func (service *ModelService) SetClient(client *Client) {
	service.client = client
}

func (service *ModelService) WithConnection(operator ClientOperator, args ...any) error {
	return service.client.Connection(func(client *Client) error {
		service.enterSession(client)
		defer func() {
			service.leaveSession()
		}()
		return operator(args...)
	})
}

func (service *ModelService) WithTransaction(operator ClientOperator, args ...any) error {
	return service.client.Transaction(func(client *Client) error {
		service.enterSession(client)
		defer func() {
			service.leaveSession()
		}()
		return operator(args...)
	})
}

func (service *ModelService) enterSession(client *Client) {
	service.sessions.Push(service.client)
	service.client = client
}

func (service *ModelService) leaveSession() {
	if client, ok := service.sessions.RPull(); ok {
		service.client = client
	}
}

type ClientOperator = func(args ...any) error

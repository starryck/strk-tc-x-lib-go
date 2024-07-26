package xbgorm

import (
	"context"
	"database/sql"
)

type DAO[T any] struct {
	client *Client
	entity *T
}

func (dao *DAO[T]) Model(entity *T) *DAO[T] {
	dao.entity = entity
	dao.client = dao.client.Model(entity)
	return dao
}

func (dao *DAO[T]) Find(data any, args ...any) *DAO[T] {
	dao.client = dao.client.Find(data, args...)
	return dao
}

func (dao *DAO[T]) FindInBatches(data any, size int, batcher ClientBatcher) *DAO[T] {
	dao.client = dao.client.FindInBatches(data, size, batcher)
	return dao
}

func (dao *DAO[T]) First(data any, args ...any) *DAO[T] {
	dao.client = dao.client.First(data, args...)
	return dao
}

func (dao *DAO[T]) Last(data any, args ...any) *DAO[T] {
	dao.client = dao.client.Last(data, args...)
	return dao
}

func (dao *DAO[T]) Take(data any, args ...any) *DAO[T] {
	dao.client = dao.client.Take(data, args...)
	return dao
}

func (dao *DAO[T]) Count(number *int64) *DAO[T] {
	dao.client = dao.client.Count(number)
	return dao
}

func (dao *DAO[T]) Pluck(field string, data any) *DAO[T] {
	dao.client = dao.client.Pluck(field, data)
	return dao
}

func (dao *DAO[T]) Where(query any, args ...any) *DAO[T] {
	dao.client = dao.client.Where(query, args...)
	return dao
}

func (dao *DAO[T]) Or(query any, args ...any) *DAO[T] {
	dao.client = dao.client.Or(query, args...)
	return dao
}

func (dao *DAO[T]) Not(query any, args ...any) *DAO[T] {
	dao.client = dao.client.Not(query, args...)
	return dao
}

func (dao *DAO[T]) Order(column any) *DAO[T] {
	dao.client = dao.client.Order(column)
	return dao
}

func (dao *DAO[T]) Limit(number int) *DAO[T] {
	dao.client = dao.client.Limit(number)
	return dao
}

func (dao *DAO[T]) Offset(offset int) *DAO[T] {
	dao.client = dao.client.Offset(offset)
	return dao
}

func (dao *DAO[T]) Select(query any, args ...any) *DAO[T] {
	dao.client = dao.client.Select(query, args...)
	return dao
}

func (dao *DAO[T]) Omit(fields ...string) *DAO[T] {
	dao.client = dao.client.Omit(fields...)
	return dao
}

func (dao *DAO[T]) Distinct(args ...any) *DAO[T] {
	dao.client = dao.client.Distinct(args...)
	return dao
}

func (dao *DAO[T]) Joins(query string, args ...any) *DAO[T] {
	dao.client = dao.client.Joins(query, args...)
	return dao
}

func (dao *DAO[T]) InnerJoins(query string, args ...any) *DAO[T] {
	dao.client = dao.client.InnerJoins(query, args...)
	return dao
}

func (dao *DAO[T]) Preload(query string, args ...any) *DAO[T] {
	dao.client = dao.client.Preload(query, args...)
	return dao
}

func (dao *DAO[T]) Group(name string) *DAO[T] {
	dao.client = dao.client.Group(name)
	return dao
}

func (dao *DAO[T]) Having(query any, args ...any) *DAO[T] {
	dao.client = dao.client.Having(query, args...)
	return dao
}

func (dao *DAO[T]) Scopes(modifiers ...ClientModifier) *DAO[T] {
	dao.client = dao.client.Scopes(modifiers...)
	return dao
}

func (dao *DAO[T]) Clauses(expressions ...Expression) *DAO[T] {
	dao.client = dao.client.Clauses(expressions...)
	return dao
}

func (dao *DAO[T]) MapColumns(mappings map[string]string) *DAO[T] {
	dao.client = dao.client.MapColumns(mappings)
	return dao
}

func (dao *DAO[T]) Create(data any) *DAO[T] {
	dao.client = dao.client.Create(data)
	return dao
}

func (dao *DAO[T]) CreateInBatches(data any, size int) *DAO[T] {
	dao.client = dao.client.CreateInBatches(data, size)
	return dao
}

func (dao *DAO[T]) FirstOrInit(data any, args ...any) *DAO[T] {
	dao.client = dao.client.FirstOrInit(data, args...)
	return dao
}

func (dao *DAO[T]) FirstOrCreate(data any, args ...any) *DAO[T] {
	dao.client = dao.client.FirstOrCreate(data, args...)
	return dao
}

func (dao *DAO[T]) Save(data any) *DAO[T] {
	dao.client = dao.client.Save(data)
	return dao
}

func (dao *DAO[T]) Attrs(dataset ...any) *DAO[T] {
	dao.client = dao.client.Attrs(dataset...)
	return dao
}

func (dao *DAO[T]) Assign(dataset ...any) *DAO[T] {
	dao.client = dao.client.Assign(dataset...)
	return dao
}

func (dao *DAO[T]) Update(field string, value any) *DAO[T] {
	dao.client = dao.client.Update(field, value)
	return dao
}

func (dao *DAO[T]) Updates(data any) *DAO[T] {
	dao.client = dao.client.Updates(data)
	return dao
}

func (dao *DAO[T]) UpdateColumn(field string, value any) *DAO[T] {
	dao.client = dao.client.UpdateColumn(field, value)
	return dao
}

func (dao *DAO[T]) UpdateColumns(data any) *DAO[T] {
	dao.client = dao.client.UpdateColumns(data)
	return dao
}

func (dao *DAO[T]) Delete(data any, args ...any) *DAO[T] {
	dao.client = dao.client.Delete(data, args...)
	return dao
}

func (dao *DAO[T]) Table(name string, args ...any) *DAO[T] {
	dao.client = dao.client.Table(name, args...)
	return dao
}

func (dao *DAO[T]) Raw(sql string, args ...any) *DAO[T] {
	dao.client = dao.client.Raw(sql, args...)
	return dao
}

func (dao *DAO[T]) Exec(sql string, args ...any) *DAO[T] {
	dao.client = dao.client.Exec(sql, args...)
	return dao
}

func (dao *DAO[T]) Row() *sql.Row {
	return dao.client.Row()
}

func (dao *DAO[T]) Rows() (*sql.Rows, error) {
	return dao.client.Rows()
}

func (dao *DAO[T]) Scan(data any) *DAO[T] {
	dao.client = dao.client.Scan(data)
	return dao
}

func (dao *DAO[T]) ScanRows(rows *sql.Rows, data any) error {
	return dao.client.ScanRows(rows, data)
}

func (dao *DAO[T]) Debug() *DAO[T] {
	dao.client = dao.client.Debug()
	return dao
}

func (dao *DAO[T]) WithContext(ctx context.Context) *DAO[T] {
	dao.client = dao.client.WithContext(ctx)
	return dao
}

func (dao *DAO[T]) Session(config *Session) *DAO[T] {
	dao.client = dao.client.Session(config)
	return dao
}

func (dao *DAO[T]) Association(field string) *Association {
	return dao.client.Association(field)
}

func (dao *DAO[T]) Get(key string) (any, bool) {
	return dao.client.Get(key)
}

func (dao *DAO[T]) Set(key string, value any) *DAO[T] {
	dao.client = dao.client.Set(key, value)
	return dao
}

func (dao *DAO[T]) InstanceGet(key string) (any, bool) {
	return dao.client.InstanceGet(key)
}

func (dao *DAO[T]) InstanceSet(key string, value any) *DAO[T] {
	dao.client = dao.client.InstanceSet(key, value)
	return dao
}

func (dao *DAO[T]) GetClient() *Client {
	return dao.client
}

func (dao *DAO[T]) SetClient(client *Client) {
	dao.client = client
}

func (dao *DAO[T]) GetEntity() *T {
	return dao.entity
}

func (dao *DAO[T]) SetEntity(entity *T) {
	dao.Model(entity)
}

func (dao *DAO[T]) GetError() error {
	return dao.client.Error
}

func (dao *DAO[T]) SetError(err error) {
	dao.client.AddError(err)
}

type (
	ClientBatcher  = func(client *Client, size int) error
	ClientModifier = func(client *Client) *Client
)

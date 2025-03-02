package xbcache

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"

	arc "github.com/hashicorp/golang-lru/arc/v2"

	"github.com/starryck/x-lib-go/source/core/utility/xbctnr"
)

var sequenceKindSet = xbctnr.NewSet(reflect.Array, reflect.Slice)

func MakeCacheKey(prefix string, keysep string, keyfrags ...any) string {
	keyParts := &bytes.Buffer{}
	keyParts.WriteString(makePrefixKeyPart(prefix, keysep))
	for _, keyfrag := range keyfrags {
		if isSequenceKind(keyfrag) {
			keyParts.WriteString(makeSequenceKeyPart(keyfrag, keysep))
		} else {
			keyParts.WriteString(makeNonsequenceKeyPart(keyfrag, keysep))
		}
	}
	return keyParts.String()
}

func isSequenceKind(value any) bool {
	if valueType := reflect.TypeOf(value); valueType != nil {
		return sequenceKindSet.Has(valueType.Kind())
	}
	return false
}

func makePrefixKeyPart(prefix string, keysep string) string {
	if prefix == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", keysep, prefix)
}

func makeSequenceKeyPart(keyfrag any, keysep string) string {
	sequence := reflect.ValueOf(keyfrag)
	elements := []string{}
	elementSet := xbctnr.NewSet[string]()
	for i := range sequence.Len() {
		element := fmt.Sprintf("%v", sequence.Index(i))
		if ok := elementSet.Has(element); ok {
			continue
		}
		elements = append(elements, element)
		elementSet.Add(element)
	}
	slices.Sort(elements)
	return fmt.Sprintf("%s%v", keysep, elements)
}

func makeNonsequenceKeyPart(keyfrag any, keysep string) string {
	return fmt.Sprintf("%s%v", keysep, keyfrag)
}

const (
	defaultKeyPrefix    = ""
	defaultKeySeparator = ":"
	defaultProxySize    = 2 << 7
)

func NewARCCache[T any](options *ARCCacheOptions) *ARCCache[T] {
	cache := (&arcCacheBuilder[T]{options: options}).
		initialize().
		setPrefix().
		setKeysep().
		setSize().
		setProxy().
		build()
	return cache
}

type ARCCache[T any] struct {
	prefix string
	keysep string
	size   int
	proxy  *arc.ARCCache[string, T]
}

func (cache *ARCCache[T]) Has(keyfrags []any) bool {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, keyfrags...)
	ok := cache.proxy.Contains(cacheKey)
	return ok
}

func (cache *ARCCache[T]) Get(keyfrags []any) (T, bool) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, keyfrags...)
	cacheVal, ok := cache.proxy.Get(cacheKey)
	return cacheVal, ok
}

func (cache *ARCCache[T]) Set(keyfrags []any, value T) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, keyfrags...)
	cache.proxy.Add(cacheKey, value)
	return
}

func (cache *ARCCache[T]) Peek(keyfrags []any) (T, bool) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, keyfrags...)
	cacheVal, ok := cache.proxy.Peek(cacheKey)
	return cacheVal, ok
}

func (cache *ARCCache[T]) Delete(keyfrags []any) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, keyfrags...)
	cache.proxy.Remove(cacheKey)
	return
}

func (cache *ARCCache[T]) Length() int {
	length := cache.proxy.Len()
	return length
}

func (cache *ARCCache[T]) Keys() []string {
	cacheKeys := cache.proxy.Keys()
	return cacheKeys
}

func (cache *ARCCache[T]) Clear() {
	cache.proxy.Purge()
	return
}

func (cache *ARCCache[T]) String() string {
	return fmt.Sprintf("<ARCCache| prefix: `%s`, keysep: `%s`, size: `%v`>",
		cache.prefix, cache.keysep, cache.size)
}

type arcCacheBuilder[T any] struct {
	cache   *ARCCache[T]
	options *ARCCacheOptions
}

type ARCCacheOptions struct {
	Prefix *string
	Keysep *string
	Size   *int
}

func (builder *arcCacheBuilder[T]) build() *ARCCache[T] {
	return builder.cache
}

func (builder *arcCacheBuilder[T]) initialize() *arcCacheBuilder[T] {
	builder.cache = &ARCCache[T]{}
	if builder.options == nil {
		builder.options = &ARCCacheOptions{}
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setPrefix() *arcCacheBuilder[T] {
	prefix := builder.options.Prefix
	if prefix != nil {
		builder.cache.prefix = *prefix
	} else {
		builder.cache.prefix = defaultKeyPrefix
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setKeysep() *arcCacheBuilder[T] {
	keysep := builder.options.Keysep
	if keysep != nil {
		builder.cache.keysep = *keysep
	} else {
		builder.cache.keysep = defaultKeySeparator
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setSize() *arcCacheBuilder[T] {
	size := builder.options.Size
	if size != nil {
		builder.cache.size = *size
	} else {
		builder.cache.size = defaultProxySize
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setProxy() *arcCacheBuilder[T] {
	builder.cache.proxy, _ = arc.NewARC[string, T](builder.cache.size)
	return builder
}

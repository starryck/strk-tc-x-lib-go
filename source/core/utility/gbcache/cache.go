package gbcache

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"

	lru "github.com/hashicorp/golang-lru/v2"
)

var sequenceKindSet = map[reflect.Kind]bool{
	reflect.Array: true,
	reflect.Slice: true,
}

func MakeCacheKey(prefix string, keysep string, parts ...any) string {
	keyParts := &bytes.Buffer{}
	keyParts.WriteString(makePrefixKeyPart(prefix, keysep))
	for _, part := range parts {
		if isPartSequenceKind(part) {
			keyParts.WriteString(makeSequenceKeyPart(part, keysep))
		} else {
			keyParts.WriteString(makeNonsequenceKeyPart(part, keysep))
		}
	}
	return keyParts.String()
}

func isPartSequenceKind(part any) bool {
	var ok bool
	if partType := reflect.TypeOf(part); partType != nil {
		_, ok = sequenceKindSet[partType.Kind()]
	}
	return ok
}

func makePrefixKeyPart(prefix string, keysep string) string {
	if prefix == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", keysep, prefix)
}

func makeSequenceKeyPart(part any, keysep string) string {
	sequence := reflect.ValueOf(part)
	elements := []string{}
	elementSet := map[string]bool{}
	for i := 0; i < sequence.Len(); i++ {
		element := fmt.Sprintf("%v", sequence.Index(i))
		if _, ok := elementSet[element]; ok {
			continue
		}
		elements = append(elements, element)
		elementSet[element] = true
	}
	sort.Strings(elements)
	return fmt.Sprintf("%s%v", keysep, elements)
}

func makeNonsequenceKeyPart(part any, keysep string) string {
	return fmt.Sprintf("%s%v", keysep, part)
}

const (
	DefaultKeyPrefix    = ""
	DefaultKeySeparator = ":"
	DefaultProxySize    = 2 << 7
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
	proxy  *lru.ARCCache[string, T]
}

func (cache *ARCCache[T]) Has(parts []any) bool {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, parts...)
	ok := cache.proxy.Contains(cacheKey)
	return ok
}

func (cache *ARCCache[T]) Get(parts []any) (T, bool) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, parts...)
	cacheVal, ok := cache.proxy.Get(cacheKey)
	return cacheVal, ok
}

func (cache *ARCCache[T]) Set(parts []any, value T) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, parts...)
	cache.proxy.Add(cacheKey, value)
	return
}

func (cache *ARCCache[T]) Peek(parts []any) (T, bool) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, parts...)
	cacheVal, ok := cache.proxy.Peek(cacheKey)
	return cacheVal, ok
}

func (cache *ARCCache[T]) Delete(parts []any) {
	cacheKey := MakeCacheKey(cache.prefix, cache.keysep, parts...)
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
		builder.cache.prefix = DefaultKeyPrefix
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setKeysep() *arcCacheBuilder[T] {
	keysep := builder.options.Keysep
	if keysep != nil {
		builder.cache.keysep = *keysep
	} else {
		builder.cache.keysep = DefaultKeySeparator
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setSize() *arcCacheBuilder[T] {
	size := builder.options.Size
	if size != nil {
		builder.cache.size = *size
	} else {
		builder.cache.size = DefaultProxySize
	}
	return builder
}

func (builder *arcCacheBuilder[T]) setProxy() *arcCacheBuilder[T] {
	builder.cache.proxy, _ = lru.NewARC[string, T](builder.cache.size)
	return builder
}

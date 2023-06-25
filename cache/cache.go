package cache

import (
	"fmt"
	"github.com/klauspost/compress/s2"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"time"
)

const (
	noCompression = 0x0
	s2Compression = 0x1
)

type MarshalFunction func(value interface{}) ([]byte, error)
type UnmarshalFunction func(b []byte, value interface{}) error

type Item[T any] struct {
	Key   string
	Value T
	TTL   time.Duration
}

func (item *Item[T]) ttl() time.Duration {
	const defaultTTL = time.Hour

	if item.TTL < 0 {
		return 0
	}

	if item.TTL != 0 {
		if item.TTL < time.Second {
			log.Printf("too short TTL for key=%q: %s", item.Key, item.TTL)
			return defaultTTL
		}
		return item.TTL
	}

	return defaultTTL
}

type Interface[T any] interface {
	Set(key string, value T)
	Get(key string) T
	_marshal(value interface{}) ([]byte, error)
	_unmarshal(b []byte, value interface{}) error
}

type Cache[T any] struct {
	Interface[T]
	options   *Options
	marshal   MarshalFunction
	unmarshal UnmarshalFunction
}

type Options struct {
	Marshal   MarshalFunction
	Unmarshal UnmarshalFunction
	TTL       time.Duration
}

const (
	compressionThreshold = 64
	timeLen              = 4
)

func compress(data []byte) []byte {
	if len(data) < compressionThreshold {
		n := len(data) + 1
		b := make([]byte, n, n+timeLen)
		copy(b, data)
		b[len(b)-1] = noCompression
		return b
	}

	n := s2.MaxEncodedLen(len(data)) + 1
	b := make([]byte, n, n+timeLen)
	b = s2.Encode(b, data)
	b = append(b, s2Compression)
	return b
}

func New[T any](opts Options) *Cache[T] {
	ce := &Cache[T]{}

	if opts.Unmarshal != nil {
		ce.unmarshal = opts.Unmarshal
	} else {
		ce.unmarshal = ce._unmarshal
	}

	if opts.Marshal != nil {
		ce.marshal = opts.Marshal
	} else {
		ce.marshal = ce._marshal
	}

	if opts.TTL != 0 {
		opts.TTL = time.Minute
	}

	return ce
}

func (c Cache[T]) _marshal(value interface{}) ([]byte, error) {
	switch value := value.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
	}

	b, err := msgpack.Marshal(value)
	if err != nil {
		return nil, err
	}

	return compress(b), nil
}

func (c Cache[T]) _unmarshal(bytes []byte, value interface{}) error {
	if len(bytes) == 0 {
		return nil
	}

	switch value := value.(type) {
	case nil:
		return nil
	case *[]byte:
		clone := make([]byte, len(bytes))
		copy(clone, bytes)
		*value = clone
		return nil
	case *string:
		*value = string(bytes)
		return nil
	}

	switch c := bytes[len(bytes)-1]; c {
	case noCompression:
		bytes = bytes[:len(bytes)-1]
	case s2Compression:
		bytes = bytes[:len(bytes)-1]

		var err error
		bytes, err = s2.Decode(nil, bytes)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown compression method: %x", c)
	}

	var result interface{} = nil
	return msgpack.Unmarshal(bytes, result)
}

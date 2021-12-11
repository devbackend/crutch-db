package storage_test

import (
	"sort"
	"testing"

	"github.com/devbackend/crutch-db/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetFromEmpty(t *testing.T) {
	st := storage.New()

	_, ok := st.Get("any")

	assert.Equal(t, false, ok)
}

func TestStorage_ReadSetValue(t *testing.T) {
	st := storage.New()

	err := st.Set("test", 123)

	assert.Equal(t, nil, err)

	val, ok := st.Get("test")

	assert.Equal(t, 123, val)
	assert.Equal(t, true, ok)
}

func TestStorage_DeleteValue(t *testing.T) {
	st := storage.New()

	_ = st.Set("test", 123)
	err := st.Delete("test")

	assert.Equal(t, nil, err)

	val, ok := st.Get("test")

	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
}

func TestStorage_KeysOnEmpty(t *testing.T) {
	st := storage.New()

	assert.Equal(t, []string{}, st.Keys())
}

func TestStorage_Keys(t *testing.T) {
	st := storage.New()

	_ = st.Set("test1", 123)
	_ = st.Set("test2", 456)
	_ = st.Set("test3", 789)

	actual := st.Keys()

	sort.Strings(actual)

	assert.Equal(t, []string{"test1", "test2", "test3"}, actual)
}

func BenchmarkStorageSet_RewriteKey(b *testing.B) {
	benchmarks := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "small string",
			value: string(make([]byte, 64)),
		},
		{
			name:  "big string",
			value: string(make([]byte, 1024*1024)),
		},
		{
			name:  "small int slice",
			value: make([]int, 64),
		},
		{
			name:  "big int slice",
			value: make([]int, 1024*1024),
		},
		{
			name:  "small string slice",
			value: make([]string, 64),
		},
		{
			name:  "big string slice",
			value: make([]string, 1024*1024),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			st := storage.New()

			for i := 0; i < b.N; i++ {
				_ = st.Set("test-key", bm.value)
			}
		})
	}
}

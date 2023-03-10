package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCacheSetGet(t *testing.T) {
	c := NewCache(1)

	// Кладем ключ, которого еще нет
	first := c.Set("first", 100)
	value, isExists := c.Get("first")

	require.Equal(t, false, first)
	require.Equal(t, value, 100)
	require.True(t, isExists)

	// Перезаписываем ключ
	second := c.Set("first", 200)
	value, isExists = c.Get("first")

	require.Equal(t, true, second)
	require.Equal(t, value, 200)
	require.True(t, isExists)

	// Превышаем капасити и проверяем, что старый ключ удален
	trhird := c.Set("second", 400)
	value, isExists = c.Get("first")

	require.Equal(t, false, trhird)
	require.Nil(t, value)
	require.False(t, isExists)
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(10)

		// Наполняем кэш данными
		c.Set("first", 10)
		c.Set("second", 20)

		// Очищаем и добавляем ключ после очистки
		c.Clear()
		c.Set("third", 30)

		// Проверяем, что старые ключи удалены
		value, isExists := c.Get("first")
		require.Nil(t, value)
		require.False(t, isExists)

		// Проверяем, что новый ключ добавлен
		value, isExists = c.Get("third")
		require.Equal(t, value, 30)
		require.True(t, isExists)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Clear()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

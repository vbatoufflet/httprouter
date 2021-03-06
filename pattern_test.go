package httprouter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatternMatch(t *testing.T) {
	var (
		p   *pattern
		ctx context.Context
		ok  bool
	)

	p = newPattern("/")
	ctx, ok = p.match(context.TODO(), "/")
	assert.Equal(t, context.TODO(), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo")
	assert.Nil(t, ctx)
	assert.False(t, ok)

	p = newPattern("/*")
	ctx, ok = p.match(context.TODO(), "/")
	assert.Equal(t, context.TODO(), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo")
	assert.Equal(t, context.TODO(), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/:key")
	ctx, ok = p.match(context.TODO(), "/foo/a")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, "a"), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/")
	assert.Nil(t, ctx)
	assert.False(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a/b")
	assert.Nil(t, ctx)
	assert.False(t, ok)

	p = newPattern("/foo/:key_a")
	ctx, ok = p.match(context.TODO(), "/foo/a")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key_a"}, "a"), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/:key/bar")
	ctx, ok = p.match(context.TODO(), "/foo/a/bar")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, "a"), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a")
	assert.Nil(t, ctx)
	assert.False(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a/")
	assert.Nil(t, ctx)
	assert.False(t, ok)

	p = newPattern("/foo/:key/*")
	ctx, ok = p.match(context.TODO(), "/foo/a/bar")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, "a"), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/:key1/bar/:key2")
	ctx, ok = p.match(context.TODO(), "/foo/a/bar/b")
	assert.Equal(t, context.WithValue(
		context.WithValue(context.TODO(), contextKey{"key1"}, "a"),
		contextKey{"key2"},
		"b",
	), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a/bar")
	assert.Nil(t, ctx)
	assert.False(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a/bar/")
	assert.Nil(t, ctx)
	assert.False(t, ok)

	p = newPattern("/foo/:key.ext")
	ctx, ok = p.match(context.TODO(), "/foo/.ext")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, ""), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a.ext")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, "a"), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/:key.")
	ctx, ok = p.match(context.TODO(), "/foo/.")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, ""), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a.")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, "a"), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/:key1:key2")
	ctx, ok = p.match(context.TODO(), "/foo/a")
	assert.Equal(t, context.WithValue(
		context.WithValue(context.TODO(), contextKey{"key1"}, "a"),
		contextKey{"key2"},
		"",
	), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/ab")
	assert.Equal(t, context.WithValue(
		context.WithValue(context.TODO(), contextKey{"key1"}, "ab"),
		contextKey{"key2"},
		"",
	), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/:key1.:key2")
	ctx, ok = p.match(context.TODO(), "/foo/a.b")
	assert.Equal(t, context.WithValue(
		context.WithValue(context.TODO(), contextKey{"key1"}, "a"),
		contextKey{"key2"},
		"b",
	), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/a.")
	assert.Equal(t, context.WithValue(
		context.WithValue(context.TODO(), contextKey{"key1"}, "a"),
		contextKey{"key2"},
		"",
	), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/.b")
	assert.Equal(t, context.WithValue(
		context.WithValue(context.TODO(), contextKey{"key1"}, ""),
		contextKey{"key2"},
		"b",
	), ctx)
	assert.True(t, ok)

	p = newPattern("/foo/bar:key")
	ctx, ok = p.match(context.TODO(), "/foo/bar")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, ""), ctx)
	assert.True(t, ok)
	ctx, ok = p.match(context.TODO(), "/foo/bara")
	assert.Equal(t, context.WithValue(context.TODO(), contextKey{"key"}, "a"), ctx)
	assert.True(t, ok)
}

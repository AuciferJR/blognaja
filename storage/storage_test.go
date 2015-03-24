package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {

	cb := &CouchbaseAdaptor{}
	op, err := cb.InsertBlog("", "", "", "")

	assert.EqualError(t, err, "title or content cannot be nil", "Should be error if title or content is nil")
	assert.Equal(t, "", op, "Should be blank")

	op, err = cb.InsertBlog("test", "content", "blog", "")
	assert.NotEmpty(t, op, "should be inserted")

	blog, err := cb.GetBlog(op)
	assert.EqualValues(t, "test", blog.Value.Title, "title should be test")
	assert.EqualValues(t, "content", blog.Value.Content, "title should be test")
	assert.EqualValues(t, op, blog.Id, "title should be test")

	err = cb.DeleteBlog("")
	assert.EqualError(t, err, "key cannot be nil")

	err = cb.DeleteBlog(op)
	assert.Nil(t, err, "Error should be nil")

}

func TestConnect(t *testing.T) {
	ct, err := connect()

	assert.Nil(t, err, "Error Should be nil")

	assert.NotNil(t, ct, "CT Should not be nil")
}

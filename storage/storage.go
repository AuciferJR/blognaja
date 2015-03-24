package storage

import (
	"blognaja/blog/define"
	"fmt"
	"github.com/couchbaselabs/gocb"
	"github.com/satori/go.uuid"
	"time"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func DBError(text string) error {
	return &errorString{text}
}

type CouchbaseAdaptor struct {
}

var couchbaseDB map[string]string = map[string]string{
	"Url":      "http://52.74.90.239:8091/",
	"Bucket":   "blognaja",
	"Username": "Administrator",
	"Password": "Noentry987",
}

var conn *gocb.Cluster
var bucket *gocb.Bucket

func timeNow() string {
	t := time.Now()
	return t.Format("Mon, 02 Jan 2006 15:04:05")
}

func connect() (bg *gocb.Bucket, err error) {

	if conn == nil {
		conn = &gocb.Cluster{}
	}

	if bucket == nil {
		bucket = &gocb.Bucket{}
	}

	conn, err = gocb.Connect(couchbaseDB["Url"])
	bucket, err := conn.OpenBucket(couchbaseDB["Bucket"], couchbaseDB["Password"])
	if err != nil {
		return nil, DBError(fmt.Sprintf("Cannot OpenBucket name : '%s' | Reason : '%s'", couchbaseDB["Bucket"], err.Error()))
	}

	bg = bucket

	return bg, err
}

func (c *CouchbaseAdaptor) InsertBlog(title string, content string, btype string, category string) (string, error) {

	var key string = ""

	if title == "" || content == "" {
		return key, DBError("title or content cannot be nil")
	}

	bg, err := connect()
	if err != nil {
		return key, err
	}

	key = uuid.NewV4().String()
	_, err = bg.Insert(key, map[string]string{"title": title, "content": content, "type": btype, "category": category, "createon": timeNow(), "updateon": ""}, 0)
	if err != nil {
		return key, DBError(fmt.Sprintf("Cannot Insert Document to Bucket : '%s' | Reason : '%s'", couchbaseDB["Bucket"], err.Error()))
	}

	return key, nil
}

func (c *CouchbaseAdaptor) UpdateBlog(key string, title string, content string, category string) error {

	if key == "" || title == "" || content == "" {
		return DBError("key or title or content cannot be nil")
	}

	bg, err := connect()
	if err != nil {
		return err
	}

	jsonData, err := c.GetBlog(key)
	if err != nil {
		return err
	}

	jsonData.Value.Title = title
	jsonData.Value.Content = content
	jsonData.Value.Updateon = timeNow()
	jsonData.Value.CategoryKey = category
	_, err = bg.Replace(key, jsonData.Value, 0, 0)
	if err != nil {
		return DBError(fmt.Sprintf("Cannot Update Document to Bucket : '%s' , Key : '%s' | Reason : '%s'", couchbaseDB["Bucket"], key, err.Error()))
	}

	return nil
}

func (c *CouchbaseAdaptor) DeleteBlog(key string) error {

	if key == "" {
		return DBError("key cannot be nil")
	}

	bg, err := connect()
	if err != nil {
		return err
	}

	_, err = bg.Remove(key, 0)
	if err != nil {
		return DBError(fmt.Sprintf("Cannot Delete Document from Bucket : '%s' , Key : '%s' | Reason : '%s'", couchbaseDB["Bucket"], key, err.Error()))
	}

	return nil
}

func (c *CouchbaseAdaptor) GetAllBlog() (*[]define.BlogData, error) {

	bg, err := connect()
	if err != nil {
		return nil, err
	}

	viewQuery := gocb.NewViewQuery("blogview", "getallblog")
	viewQuery.Reduce(false)
	viewQuery.Order(gocb.Descending)
	viewQuery.Stale(gocb.Before)
	result := bg.ExecuteViewQuery(viewQuery)

	var allBlog []define.BlogData
	jsonData := define.BlogData{}
	for result.Next(&jsonData) {
		allBlog = append(allBlog, jsonData)
	}

	// bt, err := json.Marshal(allBlog)
	// if err != nil {
	// 	return "", DBError(fmt.Sprintf("JSON Encoding Error : %s", err.Error()))
	// }

	// jsonText := string(bt[:])

	return &allBlog, nil
}

func (c *CouchbaseAdaptor) GetAllCategory() (*[]define.BlogData, error) {

	bg, err := connect()
	if err != nil {
		return nil, err
	}

	viewQuery := gocb.NewViewQuery("blogview", "getallcat")
	viewQuery.Reduce(false)
	viewQuery.Order(gocb.Descending)
	viewQuery.Stale(gocb.Before)
	result := bg.ExecuteViewQuery(viewQuery)

	var allBlog []define.BlogData
	jsonData := define.BlogData{}
	for result.Next(&jsonData) {
		allBlog = append(allBlog, jsonData)
	}

	return &allBlog, nil
}

func (c *CouchbaseAdaptor) GetBlog(key string) (*define.BlogData, error) {

	bg, err := connect()
	if err != nil {
		return nil, err
	}

	blogReturn := &define.BlogData{}
	jsonData := &define.BlogContent{}
	bg.Get(key, jsonData)

	blogReturn.Id = key
	blogReturn.Key = jsonData.Createon
	blogReturn.Value = *jsonData

	return blogReturn, nil
}

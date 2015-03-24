package blog

import (
	"blognaja/blog/define"
	"blognaja/storage"
)

var dbcon define.DBCommand

type BlogCategory struct {
	Key         string
	Name        string
	Description string
}

type Blog struct {
	Key         string
	Title       string
	Content     string
	Type        string
	CategoryKey string
}

type BlogManager struct {
}

func Init() *BlogManager {

	dbcon = &storage.CouchbaseAdaptor{}

	return &BlogManager{}
}

func (bm *BlogManager) AddNewBlog(b *Blog) (string, error) {

	key, err := dbcon.InsertBlog(b.Title, b.Content, b.Type, b.CategoryKey)

	if err != nil {
		return "", err
	}

	return key, err
}

func (bm *BlogManager) EditBlog(b *Blog) error {

	err := dbcon.UpdateBlog(b.Key, b.Title, b.Content, b.CategoryKey)

	if err != nil {
		return err
	}

	return nil
}

func (bm *BlogManager) DeleteBlog(key string) error {

	err := dbcon.DeleteBlog(key)

	if err != nil {
		return err
	}

	return nil
}

func (bm *BlogManager) GetAllBlog() (objResponse *[]define.BlogData, err error) {

	objResponse, err = dbcon.GetAllBlog()

	if err != nil {
		return nil, err
	}

	return objResponse, nil
}

func (bm *BlogManager) GetAllCategory() (objResponse *[]define.BlogData, err error) {

	objResponse, err = dbcon.GetAllCategory()

	if err != nil {
		return nil, err
	}

	return objResponse, nil
}

func (bm *BlogManager) GetBlogByKey(key string) (objResponse *define.BlogData, err error) {

	objResponse, err = dbcon.GetBlog(key)

	if err != nil {
		return nil, err
	}

	return objResponse, nil
}

func (bm *BlogManager) AddNewCategory(c *BlogCategory) (string, error) {

	// key, err := dbcon.

	// if err != nil {
	// 	return "", err
	// }
	return "", nil
	// return key, err
}

func (bm *BlogManager) EditCategory(c *BlogCategory) error {
	return nil
}

func (bm *BlogManager) DeleteCategory(key string) error {
	return nil
}

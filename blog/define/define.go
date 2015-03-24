package define

type BlogData struct {
	Id    string      `json:"id"`
	Key   string      `json:"key"`
	Value BlogContent `json:"value"`
}

type BlogContent struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	CategoryKey string `json:"category"`
	Createon    string `json:"createon"`
	Updateon    string `json:"updateon"`
}

type DBCommand interface {
	InsertBlog(title string, content string, btype string, category string) (string, error)
	UpdateBlog(key string, title string, content string, category string) error
	DeleteBlog(key string) error
	GetAllBlog() (*[]BlogData, error)
	GetAllCategory() (*[]BlogData, error)
	GetBlog(key string) (*BlogData, error)
}

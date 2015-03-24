package main

import (
	"blognaja/blog"
	"blognaja/blog/define"
	"github.com/plimble/ace"
	"github.com/plimble/ace-contrib/pongo2"
	"log"
	"strings"
)

type jsonResponseData struct {
	Status string             `json:"status"`
	Detail string             `json:"detail"`
	Data   *[]define.BlogData `json:"data"`
}

func (j *jsonResponseData) resFmt(s string, d string, dt *[]define.BlogData) *jsonResponseData {
	j.Status = s
	j.Detail = d
	j.Data = dt
	return j
}

const listen string = ":3001"

func wslogger(c *ace.C) {
	log.Println("----> Request Webservice :", c.Request.RequestURI)
	c.Next()
}

func addblog(c *ace.C) {
	bm := blog.Init()
	res := &jsonResponseData{}

	title := strings.TrimSpace(c.MustPostString("title", ""))
	content := strings.TrimSpace(c.MustPostString("content", ""))
	btype := strings.TrimSpace(c.MustPostString("btype", ""))
	category := strings.TrimSpace(c.MustPostString("category", ""))

	if title == "" || content == "" || btype == "" {
		log.Println("      ERROR : title or content or btype is blank value")
		res.resFmt("-1", "title or content or btype shoud have any value", nil)
	} else {
		b := &blog.Blog{
			Title:       title,
			Content:     content,
			Type:        btype,
			CategoryKey: category,
		}

		key, err := bm.AddNewBlog(b)
		if err != nil {
			log.Printf("      ERROR : %s\n", err.Error())
			res.resFmt("-2", err.Error(), nil)
		}
		log.Printf("      SUCCESS : Add Blog Complete with KEY = %s\n", key)
		res.resFmt("0", "success:"+key, nil)
	}

	c.JSON(200, *res)
}

func editblog(c *ace.C) {
	bm := blog.Init()
	res := &jsonResponseData{}

	key := strings.TrimSpace(c.MustPostString("key", ""))
	title := strings.TrimSpace(c.MustPostString("title", ""))
	content := strings.TrimSpace(c.MustPostString("content", ""))
	category := strings.TrimSpace(c.MustPostString("category", ""))

	if title == "" || content == "" || key == "" {
		log.Println("      ERROR : title or content or key is blank value")
		res.resFmt("-1", "title or content or key shoud have any value", nil)
	} else {
		b := &blog.Blog{
			Key:         key,
			Title:       title,
			Content:     content,
			CategoryKey: category,
		}

		err := bm.EditBlog(b)
		if err != nil {
			log.Printf("      ERROR : %s\n", err.Error())
			res.resFmt("-2", err.Error(), nil)
		} else {
			log.Printf("      SUCCESS : Edit Blog Complete with KEY = %s\n", key)
			res.resFmt("0", "success:"+key, nil)
		}
	}

	c.JSON(200, *res)
}

func deleteblog(c *ace.C) {
	bm := blog.Init()
	res := &jsonResponseData{}

	key := strings.TrimSpace(c.MustQueryString("key", ""))
	if key == "" {
		log.Println("      ERROR : key is blank value")
		res.resFmt("-1", "key shoud have any value", nil)
	} else {
		err := bm.DeleteBlog(key)
		if err != nil {
			res.resFmt("-2", err.Error(), nil)
		} else {
			res.resFmt("0", "success:"+key, nil)
		}
	}

	c.JSON(200, *res)
}

func getallblog(c *ace.C) {
	bm := blog.Init()
	res := &jsonResponseData{}

	objRes, err := bm.GetAllBlog()
	if err != nil {
		res.resFmt("-2", err.Error(), nil)
	} else {
		res.resFmt("0", "success", objRes)
	}
	c.JSON(200, *res)
}

func getallcategory(c *ace.C) {
	bm := blog.Init()
	res := &jsonResponseData{}

	objRes, err := bm.GetAllCategory()
	if err != nil {
		res.resFmt("-2", err.Error(), nil)
	} else {
		res.resFmt("0", "success", objRes)
	}
	c.JSON(200, *res)
}

func getblog(c *ace.C) {
	bm := blog.Init()
	res := &jsonResponseData{}

	key := strings.TrimSpace(c.MustQueryString("key", ""))
	if key == "" {
		log.Println("      ERROR : key is blank value")
		res.resFmt("-1", "key shoud have any value", nil)
	} else {
		var objResMap []define.BlogData
		objRes, err := bm.GetBlogByKey(key)
		objResMap = append(objResMap, *objRes)
		if err != nil {
			res.resFmt("-2", err.Error(), nil)
		} else {
			res.resFmt("0", "success", &objResMap)
		}
	}

	c.JSON(200, *res)
}

func main() {

	a := ace.New()

	// Webservice group
	ws := a.Group("/webservice", wslogger)

	ws.POST("/addblog", addblog)
	ws.POST("/editblog", editblog)
	ws.GET("/deleteblog", deleteblog)
	ws.GET("/getallblog", getallblog)
	ws.GET("/getallcategory", getallcategory)
	ws.GET("/getblog", getblog)

	// Main Web
	tmo := &pongo2.TemplateOptions{
		Directory:     "./templates",
		IsDevelopment: true,
	}
	render := pongo2.Pongo2(tmo)

	a.Static("/static", "./static")
	a.HtmlTemplate(render)
	a.GET("/", func(c *ace.C) {
		c.HTML("index.html", nil)
	})

	a.GET("/category", func(c *ace.C) {
		c.HTML("category.html", nil)
	})

	log.Println("BlogNaja running on", listen)
	a.Run(listen)
}

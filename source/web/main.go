package main

import (
	"flag"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"

	"saynice.xyz/src/util"
)

// APIMessage API 消息体
type APIMessage struct {
	Code int         `json:"code"`
	Erro string      `json:"erro"`
	Data interface{} `json:"data"`
}

// NewPostData new/post.html 页面片段
type NewPostData struct {
	Title                           string
	Description                     string
	PostTextareaPlaceholder         string
	PostsAButtonText                string
	PreviewButtonText               string
	CommunityComplianceCheckBoxText string
	APIDomain                       string
	SubmitButtonText                string
}

// VoteData 随机匿名空间表决数据对象
type VoteData struct {
	Vote        interface{}
	VoteURL     string
	ArticlesURL string
}

// IndexData index.html 页面片段
type IndexData struct {
	Title            string
	Description      string
	ReportButtonText string
	APIDomain        string
	Posts            interface{}
	HasVote          bool
	VoteData         VoteData
}

// PostData 主题单页数据
type PostData struct {
	Title       string
	Description string
	Code        int
	Erro        string
	Post        interface{}
}

var (
	// Domain 请求域
	Domain string
)

func uri(path string) string {
	return Domain + path
}

func init() {
	flag.StringVar(&Domain, "D", "http://127.0.0.1:18823/v1", "请求域")
}

func initTmpl(tmpl *template.Template) {
	box := packr.NewBox("./res/layouts")

	htmls := []string{
		"fragmentHead.html",
		"fragmentFooter.html",
		"fragmentToolbar.html",
		"index.html",
		"newPost.html",
		"post.html",
		"fragmentVote.html",
		"fragmentPreview.html",
		"fragmentWelcome.html",
		"fragmentAbout.html",
	}

	for _, v := range htmls {
		indexTmpl := tmpl.New(v)
		data, _ := box.FindString(v)
		indexTmpl.Parse(data)
	}
}

func main() {
	flag.Parse()
	flag.Usage()

	tmpl := template.New("")

	initTmpl(tmpl)

	router := gin.Default()

	router.GET("/", indexHTML)
	router.GET("/posts/:offset", indexHTML)
	router.GET("/new/post", newPostHTML)
	router.GET("/post/:id", postHTML)

	router.SetHTMLTemplate(tmpl)
	router.StaticFS("/static", packr.NewBox("./res/static"))
	router.Run(":19548")
}

func indexHTML(c *gin.Context) {
	data := &IndexData{
		Title:            "SayNice - 匿名情感倾诉社区、完美树洞、你的 OK 工具人",
		Description:      "",
		ReportButtonText: "举报",
		APIDomain:        Domain,
	}

	offset := c.Param("offset")

	var msg APIMessage

	e := util.GetJSON(uri("/posts?offset="+offset), &msg)

	if nil != e {
		c.JSON(http.StatusOK, e.Error())
		return
	} else if 0 != msg.Code {
		c.JSON(http.StatusOK, msg)
		return
	}

	data.Posts = msg.Data

	e = util.GetJSON(uri("/space/"), &msg)

	if nil == e && 0 == msg.Code {
		data.HasVote = true
		data.VoteData = VoteData{
			Vote:        msg.Data,
			VoteURL:     uri("/vote"),
			ArticlesURL: uri("/articles"),
		}
	}

	c.HTML(http.StatusOK, "index.html", data)
}

func newPostHTML(c *gin.Context) {
	data := &NewPostData{
		Title:                           "SayNice - 匿名情感倾诉社区、完美树洞、你的 OK 工具人",
		Description:                     "",
		PostTextareaPlaceholder:         "写点什么呢",
		PostsAButtonText:                "Say Nice 社区",
		PreviewButtonText:               "预览",
		CommunityComplianceCheckBoxText: "是否同意 SayNice.xyz 社区守约",
		APIDomain:                       Domain,
		SubmitButtonText:                "发布",
	}

	c.HTML(http.StatusOK, "newPost.html", data)
}

func postHTML(c *gin.Context) {
	var msg APIMessage

	id := c.Param("id")
	e := util.GetJSON(uri("/post/"+id), &msg)

	if nil != e {
		c.JSON(http.StatusOK, e.Error())
		return
	}

	post := &PostData{
		Title:       "SayNice - 匿名情感倾诉社区、完美树洞、你的 OK 工具人",
		Description: "",
		Code:        msg.Code,
		Erro:        msg.Erro,
		Post:        msg.Data,
	}

	c.HTML(http.StatusOK, "post.html", post)
}

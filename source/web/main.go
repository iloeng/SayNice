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

// Post 主题
// Text comment:'主题内容';
// Feelings comment:'作者的感受';
// Attitudes comment:'主题状态, 0: IDLE, 1: ALLOWED, 2: UNALLOWED, 3: ABSTAIN, 4: REPORTED';
// Remark comment:'评语，备注';
type Post struct {
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	Remark    string `json:"remark" gorm:"type:text;default:'';"`
	CreatedAt string `json:"createdAt,omitempty"`
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

// IndexData index.html 页面片段
type IndexData struct {
	Title            string
	Description      string
	ReportButtonText string
	APIDomain        string
	Posts            interface{}
	HasVote          bool
	Vote             interface{}
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

func main() {
	flag.Parse()
	flag.Usage()

	box := packr.NewBox("./res/layouts")

	tmpl := template.New("")

	indexTmpl := tmpl.New("index")
	data, _ := box.FindString("index.html")
	indexTmpl.Parse(data)

	indexTmpl = tmpl.New("newPost")
	data, _ = box.FindString("newPost.html")
	indexTmpl.Parse(data)

	indexTmpl = tmpl.New("result")
	data, _ = box.FindString("result.html")
	indexTmpl.Parse(data)

	indexTmpl = tmpl.New("fragmentVote")
	data, _ = box.FindString("fragmentVote.html")
	indexTmpl.Parse(data)

	staticBox := packr.NewBox("./res/static")

	router := gin.Default()

	router.GET("/", indexHTML)
	router.GET("/new/post", newPostHTML)
	router.POST("/result", resultHTML)

	router.SetHTMLTemplate(tmpl)
	router.StaticFS("/static", staticBox)
	router.Run(":19548")
}

func indexHTML(c *gin.Context) {
	data := &IndexData{
		Title:            "SayNice - 匿名情感倾诉社区、完美树洞、你的 OK 工具人",
		Description:      "",
		ReportButtonText: "举报",
		APIDomain:        Domain,
	}

	var msg APIMessage

	e := util.GetJSON(uri("/posts"), &msg)

	if nil != e {
		c.JSON(http.StatusOK, e.Error())
		return
	} else if 0 != msg.Code && 10020 != msg.Code {
		c.JSON(http.StatusOK, msg)
		return
	}

	data.Posts = msg.Data

	e = util.GetJSON(uri("/space/"), &msg)

	if nil == e && 0 == msg.Code {
		data.HasVote = true
		data.Vote = msg.Data
	}

	c.HTML(http.StatusOK, "index", data)
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

	c.HTML(http.StatusOK, "newPost", data)
}

func resultHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "result", nil)
}

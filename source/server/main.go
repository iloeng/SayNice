package main

import (
	"database/sql/driver"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// TODO 标签功能、添加用户感受功能、随机匿名空间拆分为独立服务

/////////////////////////////////////////////// CONST DEFINE START

// 主题、随机匿名空间审查和用户表决状态
const (
	// StatusIdle 初始状态或进行中
	StatusIdle = 0
	// StatusAllowed 主题允许发布或表决结果为允许
	StatusAllowed = 1
	// StatusUnallowed 主题不允许发布或表决结果为不允许
	StatusUnallowed = 2
	// StatusAbstain 用户放弃表决或主题存在争议，需要开启空间发起表决
	StatusAbstain = 3
	// StatusReported 主题被人举报了
	StatusReported = 4
)

// APIMessage code 值
const (
	// CodeSuccess 服务执行成功
	CodeSuccess = 0
)

const (
	// MaxWordCount 最大字数
	MaxWordCount = 5000
)

// 缓存 Key 的前缀
const (
	PrefixRASPID = "RASPID"
	PrefixTOKEN  = "TOKEN"
)

// 缓存超时时间
const (
	Expire3Days = 259200 // 3 days‬
)

const (
	// MaxVoteSize 最大表决人数
	MaxVoteSize = 9

	// _3thVoteNum 十分之三的表决数
	_3thVoteNum = MaxVoteSize / 3
)

/////////////////////////////////////////////// CONST DEFINE END
//
/////////////////////////////////////////////// STRUCT DEFINE START

// APIMessage API 消息体
type APIMessage struct {
	Code int         `json:"code"`
	Erro string      `json:"erro"`
	Data interface{} `json:"data"`
}

func (msg *APIMessage) Error() string {
	return msg.Erro
}

// JSONTime Json 序列化时用的时间格式
type JSONTime struct {
	time.Time
}

// MarshalJSON 序列化
func (t *JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))), nil
}

// UnmarshalJSON 反序列化
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	var err error

	loc, _ := time.LoadLocation("Asia/Shanghai")
	t.Time, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), loc)
	if err != nil {
		return err
	}

	return nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Feeling 作者的感觉、感受和心情
// Name comment:'感受、心情';
type Feeling struct {
	ID        uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Name      string    `json:"name" gorm:"type:text;default:'';not null;"`
	CreatedAt *JSONTime `json:"createdAt,omitempty"`
}

// Emoji 表情
// Name comment:'表情、态度';
// Feelings comment:'这种表情或态度所代表的感受或心情'
type Emoji struct {
	ID       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Name     string    `json:"name" gorm:"type:text;default:'';not null;"`
	Feelings []Feeling `json:"feelings,omitempty" gorm:"many2many:emoji_feeling;"`
}

// Attitude 他人的意见、态度
// Post comment:'评价的主题';
// Emoji comment:'他们对此主题的表情';
type Attitude struct {
	ID        uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Post      Post      `json:"post,omitempty" gorm:"foreignkey:PostID;"`
	PostID    uint      `json:"postId"`
	Emoji     *Emoji    `json:"emoji,omitempty" gorm:"foreignkey:EmojiID;"`
	EmojiID   uint      `json:"emojiId"`
	CreatedAt *JSONTime `json:"createdAt,omitempty"`
}

// Post 主题
// Text comment:'主题内容';
// Feelings comment:'作者的感受';
// Attitudes comment:'主题状态, 0: IDLE, 1: ALLOWED, 2: UNALLOWED, 3: ABSTAIN, 4: REPORTED';
// Remark comment:'评语，备注';
type Post struct {
	ID        uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Text      string     `json:"text" gorm:"type:text;default:'';not null;"`
	Feelings  []Feeling  `json:"feelings,omitempty" gorm:"many2many:post_feeling;"`
	Attitudes []Attitude `json:"attitudes,omitempty"`
	Token     string     `json:"-" gorm:"type:text;default:'';not null;"`
	Status    int        `json:"status" gorm:"type:integer;default:0;not null;"`
	Remark    string     `json:"remark" gorm:"type:text;default:'';"`
	CreatedAt *JSONTime  `json:"createdAt,omitempty"`
}

// RASpace 随机匿名空间
// Title comment:'空间标题';
// Post comment:'随机匿名空间审查的主题';
// MaxVotedCount comment:'最大表决数量';
// Status comment:'表决状态, 0: IDLE, 1: ALLOWED, 2: UNALLOWED, 3: ABSTAIN, 4: REPORTED';
type RASpace struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Title         string    `json:"subject" gorm:"type:text;default:'';not null;"`
	Post          *Post     `json:"post,omitempty" gorm:"foreignkey:PostID;"`
	PostID        uint      `json:"postId"`
	Votes         []Vote    `json:"votes,omitempty"`
	MaxVotedCount int       `json:"-" gorm:"type:integer;default:0;not null;"`
	Status        int       `json:"status" gorm:"type:integer;default:0;not null;"`
	CreatedAt     *JSONTime `json:"createdAt,omitempty"`
}

// Vote 表决
// RASpace comment:'该表决在哪个随机匿名空间中进行的';
// Token comment:'表决者token';
// Status comment:'表决状态, 0: IDLE, 1: ALLOWED, 2: UNALLOWED, 3: ABSTAIN, 4: REPORTED';
// Remark comment:'评语，备注';
type Vote struct {
	RASpace   *RASpace  `json:"raSpace,omitempty" gorm:"foreignkey:RASpaceID;"`
	RASpaceID uint      `json:"raSpaceId"`
	Token     string    `json:"token" gorm:"type:text;default:'';not null;"`
	Status    int       `json:"status" gorm:"type:integer;default:0;not null;"`
	Remark    string    `json:"remark" gorm:"type:text;default:'';"`
	CreatedAt *JSONTime `json:"createdAt,omitempty"`
}

// SyncSlice 当前存在的 RASpace
type SyncSlice struct {
	lock sync.Mutex
	ids  []uint
}

// Add 添加一个 RASpace ID
func (box *SyncSlice) Add(id uint) {
	box.lock.Lock()
	box.ids = append(box.ids, id)
	box.lock.Unlock()
}

// Del 删除指定的 RASpace ID
func (box *SyncSlice) Del(id uint) {
	box.lock.Lock()
	for i := 0; i < len(box.ids); i++ {
		if box.ids[i] == id {
			box.ids = append(box.ids[:i], box.ids[i+1:]...)
		}
	}
	box.lock.Unlock()
}

// Remove 移除指定位置的 RASpace ID
func (box *SyncSlice) Remove(index int) {
	box.lock.Lock()
	for i := 0; i < len(box.ids); i++ {
		if i == index {
			box.ids = append(box.ids[:i], box.ids[i+1:]...)
		}
	}
	box.lock.Unlock()
}

// Get 获取指定位置的 RASpace ID
func (box *SyncSlice) Get(index int) uint {
	return box.ids[index]
}

// Last 获取最后一个元素
func (box *SyncSlice) Last() uint {
	return box.ids[len(box.ids)-1]
}

// Has 匹配指定 id 是否存在
func (box *SyncSlice) Has(id uint) bool {
	ok := false

	box.lock.Lock()
	for _, value := range box.ids {
		if value == id {
			ok = true
			break
		}
	}
	box.lock.Unlock()

	return ok
}

// Count 获取 ID 总数
func (box *SyncSlice) Count() int {
	return len(box.ids)
}

/////////////////////////////////////////////// STRUCT DEFINE END

var (
	// Debug 是否查看日志
	Debug bool
	// DBPath SQLite3 数据库文件地址
	DBPath string

	db    *gorm.DB
	cache *freecache.Cache
	box   *SyncSlice
	npid  *SyncSlice // new post id
)

func checkError(err error) {
	if nil != err {
		panic(err)
	}
}

func init() {
	flag.BoolVar(&Debug, "debug", false, "是否查看日志")
	flag.StringVar(&DBPath, "db", "saynice.db", "SQLite3 数据库文件地址")
}

func initDB() {
	var err error
	db, err = gorm.Open("sqlite3", DBPath)
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Feeling{}).Error
	checkError(err)
	err = db.AutoMigrate(&Emoji{}).Error
	checkError(err)
	err = db.AutoMigrate(&Attitude{}).Error
	checkError(err)
	err = db.AutoMigrate(&Post{}).Error
	checkError(err)
	err = db.AutoMigrate(&RASpace{}).Error
	checkError(err)
	err = db.AutoMigrate(&Vote{}).Error
	checkError(err)
}

func initCache() {
	box = &SyncSlice{}
	npid = &SyncSlice{}
	cacheSize := 64 * 1024 * 1024 // N兆
	cache = freecache.NewCache(cacheSize)
}

func main() {
	flag.Parse()
	flag.Usage()

	rand.Seed(time.Now().UnixNano())

	initDB()
	initCache()

	if Debug {
		db.LogMode(true)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	defer db.Close()

	startRASpace()

	router := gin.Default()
	v1 := router.Group("/v1")

	v1.GET("/posts", output(listPost))
	v1.GET("/emoji", output(listEmoji))

	v1.GET("/post/:id", output(getPost))

	v1.GET("/new/post", output(newPost))
	v1.POST("/post/:id", output(submitPost))
	v1.POST("/say", output(submitAttitude))

	v1.GET("/space/*token", output(getRASpace))
	v1.POST("/vote/:token", output(votePost))
	v1.POST("/report/:id", output(reportPost))

	router.Run(":18823")
}

func startRASpace() {
	var raSpaces []RASpace

	e := db.Where("status=?", StatusIdle).Find(&raSpaces).Error

	if nil != e {
		return
	}

	for _, space := range raSpaces {
		setRASpaceCache(space)
	}
}

func output(fn func(*gin.Context) (int, string, interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, erro, data := fn(c)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.JSON(http.StatusOK, APIMessage{
			Code: code,
			Erro: erro,
			Data: data,
		})
	}
}

// 显示所有主题
func listPost(c *gin.Context) (int, string, interface{}) {
	offset, e := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if nil != e {
		offset = 0
	}

	limit, e := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if nil != e {
		limit = 10
	}

	order, e := strconv.Atoi(c.DefaultQuery("order", "0"))

	if nil != e {
		order = 0
	}

	var postList []Post

	tx := db.Where("status=?", StatusAllowed).Limit(limit).Offset(offset)

	if 0 == order {
		e = tx.Order("created_at DESC").Find(&postList).Error
	} else {
		e = tx.Find(&postList).Error
	}

	if nil != e {
		return 10010, e.Error(), nil
	} else if 0 == len(postList) {
		return 10020, "No posts", nil
	}

	return CodeSuccess, "", postList
}

// 显示所有 emoji 语言
func listEmoji(c *gin.Context) (int, string, interface{}) {
	feelingStr := c.DefaultQuery("feelings", "all")

	var e error
	var emojiList []Emoji

	if "all" == feelingStr {
		e = db.Find(&emojiList).Error
	} else {
		feelings := strings.Split(feelingStr, ",")
		e = db.Where("feeling IN (?)", feelings).Find(&emojiList).Error
	}

	if nil != e {
		return 20001, e.Error(), nil
	}

	return CodeSuccess, "", emojiList
}

// 获取指定的主题内容
func getPost(c *gin.Context) (int, string, interface{}) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)

	if 0 == id || nil != e {
		return 30010, "Parameter id value is wrong", nil
	}

	var post Post

	e = db.Where("id=?", id).Find(&post).Error

	if nil != e {
		return 30020, e.Error(), nil
	}

	switch post.Status {
	case StatusAllowed:
		return CodeSuccess, "", post
	case StatusUnallowed:
		return 30030, "Review failed, it's " + post.Remark, nil
	case StatusAbstain:
		return 30040, "Content is too controversial, it's " + post.Remark, nil
	case StatusReported:
		return 30050, "Post reported, it's " + post.Remark, nil
	default:
		return 30060, "Pending review", nil
	}
}

func newPost(c *gin.Context) (int, string, interface{}) {
	var id uint

	if 0 < npid.Count() {
		id = npid.Last() + 1
	} else {
		var post Post
		e := db.Last(&post).Error

		if nil != e {
			var count int
			db.Table("posts").Count(&count)

			if 0 < count {
				return 41010, e.Error(), nil
			}

			// 主题 ID 的起始值
			id = 10000
		} else {
			id = post.ID + 1
		}
	}

	npid.Add(id)

	return CodeSuccess, "", id
}

// 提交一个主题
func submitPost(c *gin.Context) (int, string, interface{}) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)

	if nil != e || id <= 0 || !npid.Has(uint(id)) {
		return 40001, "Cannot find the post with the specified id.", nil
	}

	var data Post

	e = c.ShouldBindJSON(&data)

	if nil != e {
		return 40010, e.Error(), nil
	} else if MaxWordCount < len(data.Text) {
		msg := strconv.Itoa(len(data.Text)) +
			" words are too many, please limit to " +
			strconv.Itoa(MaxWordCount) +
			" words"
		return 40020, msg, nil
	}

	post := Post{ID: uint(id), Text: data.Text, Token: uuid.New().String()}

	e = db.Set("gorm:association_autoupdate", false).Create(&post).Error

	if nil != e {
		return 40030, e.Error(), nil
	}

	npid.Del(post.ID)

	go openRASpace(post.ID, "请检查新主题中是否含有违反守约的内容，谢谢。")

	return CodeSuccess, "", post.ID
}

// 表达对指定主题的看法及态度
func submitAttitude(c *gin.Context) (int, string, interface{}) {
	var data Attitude

	e := c.ShouldBindJSON(&data)

	if nil != e {
		return 50001, e.Error(), nil
	}

	attitude := Attitude{
		PostID:  data.PostID,
		EmojiID: data.EmojiID,
	}

	e = db.Set("gorm:association_autoupdate", false).Create(&attitude).Error

	if nil != e {
		return 50002, e.Error(), nil
	}

	return CodeSuccess, "", attitude.ID
}

// 进入一个随机匿名空间
func getRASpace(c *gin.Context) (int, string, interface{}) {
	if 0 == box.Count() {
		return 60000, "Welcome to create a new Post", nil
	}

	var id uint

	token := c.Param("token")[1:]

	// 如果通过 token 进来的, 则可以直接访问
	if "" != token {
		idByte, err := cache.Get(key(PrefixTOKEN, token))

		if nil != err {
			return 60010, "Parameter token is wrong or invalid.", nil
		}

		id64, err := strconv.ParseUint(string(idByte), 10, 64)

		if nil != err {
			return 60011, err.Error(), nil
		}

		id = uint(id64)
	}

	if 0 == id {
		if 30 < rand.Intn(100) { // 没有用户系统的替代方案
			return 60020, "See you next time.", nil
		}

		idIndex := rand.Intn(box.Count())
		id = box.Get(idIndex)
	}

	var space RASpace

	err := db.Transaction(func(tx *gorm.DB) error {
		spaceByte, e := cache.Get(key(PrefixRASPID, id))

		if nil != e {
			e = tx.Where("id=?", id).Find(&space).Error

			if nil != e {
				return msg(60030, e.Error(), nil)
			}

			box.Del(space.ID)
			setRASpaceCache(space)
		} else {
			e = json.Unmarshal(spaceByte, &space)

			if nil != e {
				return msg(60031, e.Error(), nil)
			}
		}

		var post Post

		e = tx.Model(&space).Association("Post").Find(&post).Error

		if nil != e {
			return msg(60032, e.Error(), nil)
		}

		space.Post = &post

		return nil
	})

	if nil == err {
		if "" == token {
			token = uuid.New().String()
		}

		vote := Vote{
			RASpace:   &space,
			RASpaceID: space.ID,
			Token:     token,
		}

		cache.Set(key(PrefixTOKEN, token), val(space.ID), Expire3Days)
		return CodeSuccess, "", vote
	} else if msg, ok := err.(*APIMessage); ok {
		return msg.Code, msg.Erro, msg.Data
	}

	return 60040, "Unknown", err.Error()
}

// 对一个主题进行表决
func votePost(c *gin.Context) (int, string, interface{}) {
	var data Vote

	token := c.Param("token")

	if "" == token {
		return 70001, "No token.", nil
	}

	raSpaceIDByte, e := cache.Get(key(PrefixTOKEN, token))

	if nil != e {
		return 70002, "Parameter token is wrong or invalid.", nil
	}

	raSpaceID64, e := strconv.ParseUint(string(raSpaceIDByte), 10, 64)

	if nil != e {
		return 70003, e.Error(), nil
	}

	raSpaceID := uint(raSpaceID64)

	if raSpaceID <= 0 {
		return 70010, "Internal error and id <= 0.", nil
	} else if !box.Has(raSpaceID) {
		return 70011, "This RASpace is closed.", nil
	}

	e = c.ShouldBindJSON(&data)

	if nil != e {
		return 70012, e.Error(), nil
	} else if !(StatusAllowed == data.Status || StatusUnallowed == data.Status || StatusAbstain == data.Status) {
		return 70013, "Status value must be between {1, 2, 3}, 1: allow, 2: unallow, 3: abstain.", nil
	}

	vote := Vote{
		RASpaceID: raSpaceID,
		Token:     data.Token, // vote token 与 url token 并非同一个 token, url token 是访问 rasp 空间的令牌，而 vote token 是表决者身份令牌, vote token 可以设为空, 由系统生成。
		Status:    data.Status,
		Remark:    data.Remark,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if "" == vote.Token {
			vote.Token = uuid.New().String()
		} else {
			var count int
			tx.Model(&Vote{}).Where(&Vote{RASpaceID: vote.RASpaceID, Token: vote.Token}).Count(&count)

			if 0 < count {
				return msg(70014, "Repeated voting.", nil)
			}
		}

		var space RASpace

		spaceByte, e := cache.Get(key(PrefixRASPID, vote.RASpaceID))

		if nil != e {
			e = tx.Where("id=?", vote.RASpaceID).Find(&space).Error

			if nil != e {
				return msg(70015, e.Error(), nil)
			}

			box.Del(space.ID)
			setRASpaceCache(space)
		} else {
			e = json.Unmarshal(spaceByte, &space)

			if nil != e {
				return msg(70016, e.Error(), nil)
			}
		}

		e = tx.Set("gorm:association_autoupdate", false).Create(&vote).Error

		if nil != e {
			return msg(70017, e.Error(), nil)
		}

		var votedCount int
		tx.Model(&Vote{}).Where(&Vote{RASpaceID: space.ID}).Count(&votedCount)

		if votedCount < MaxVoteSize {
			// 当表决人数未达到最大表决人数时，暂不处理
		} else {
			// 表决结束，统计结果,
			// 同步 Post, RASpace 两张表,
			// 更新缓存

			var votes []Vote

			e = tx.Where(&Vote{RASpaceID: space.ID}).Find(&votes).Error

			if nil != e {
				return msg(70018, e.Error(), nil)
			}

			var allowCount int
			var unallowCount int
			var unallowRemark string

			for _, vote := range votes {
				if StatusAllowed == vote.Status {
					allowCount++
				} else if StatusUnallowed == vote.Status {
					unallowCount++
					unallowRemark += vote.Remark + ";"
				}
			}

			status := StatusIdle
			post := &Post{}

			if _3thVoteNum <= allowCount && _3thVoteNum <= unallowCount {
				if unallowCount <= allowCount {
					// 表示主题审查通过
					status = StatusAllowed
				} else {
					// 表示主题审查未通过
					status = StatusUnallowed
					post.Remark = unallowRemark
				}
			} else if _3thVoteNum <= allowCount {
				// 表示主题审查通过
				status = StatusAllowed
			} else if _3thVoteNum <= unallowCount {
				// 表示主题审查未通过
				status = StatusUnallowed
				post.Remark = unallowRemark
			} else {
				// 表示大多数人放弃了表决
				status = StatusAbstain

				// 有两种可能,
				// 一种是主题内容有太多争议之处，无法判断,
				// 另一种是表决者懒惰，不想审查.
				// 两种情况无法通过程序处理,
				// 故而重置表决结果，重新表决.
				// 此处随机匿名空间开启理由采用第一种情况,
				// 因为我相信我们不是懒惰.
				go openRASpace(space.PostID, "该主题存在诸多争议，请认真审查再做判断，谢谢")
			}

			post.Status = status
			e = tx.Model(&Post{}).Where("id=?", space.PostID).Update(post).Error

			if nil != e {
				return msg(70020, e.Error(), nil)
			}

			e = tx.Model(&RASpace{}).Where("id=?", space.ID).Update(&RASpace{Status: status}).Error

			if nil != e {
				return msg(70021, e.Error(), nil)
			}

			box.Del(space.ID)
			cache.Del(key(PrefixRASPID, space.ID))
		}

		return nil
	})

	if nil == err {
		cache.Del(key(PrefixTOKEN, token))

		return CodeSuccess, "", nil
	} else if msg, ok := err.(*APIMessage); ok {
		return msg.Code, msg.Erro, msg.Data
	}

	return 70030, "Unknown", err.Error()
}

// 举报一个主题
func reportPost(c *gin.Context) (int, string, interface{}) {
	postID, e := strconv.ParseUint(c.Param("id"), 10, 64)

	if postID <= 0 || nil != e {
		return 80001, "Parameter ID type error", nil
	}

	remark := c.PostForm("remark")

	updateMap := map[string]interface{}{"status": 3, "remark": remark}

	e = db.Model(&Post{}).Where("id=?", postID).Update(updateMap).Error

	if nil != e {
		return 80002, e.Error(), nil
	}

	go openRASpace(uint(postID), "有人举报了该主题，请检查举报是否属实，谢谢。")

	return CodeSuccess, "", nil
}

////////////////////////////////////////////////////// TOOLS FUNC START

func openRASpace(postID uint, title string) {
	space := RASpace{
		Title:         title,
		PostID:        postID,
		MaxVotedCount: MaxVoteSize,
	}

	e := db.Set("gorm:association_autoupdate", false).Create(&space).Error

	if nil != e {
		return
	}

	setRASpaceCache(space)
}

func setRASpaceCache(space RASpace) error {
	spaceByte, _ := json.Marshal(space)

	box.Add(space.ID)

	return cache.Set(key(PrefixRASPID, space.ID), spaceByte, Expire3Days)
}

func key(prefix string, t interface{}) []byte {
	if str, ok := t.(string); ok {
		return []byte(prefix + str)
	}

	return []byte(prefix + fmt.Sprint(t))
}

func val(t interface{}) []byte {
	if str, ok := t.(string); ok {
		return []byte(str)
	}

	return []byte(fmt.Sprint(t))
}

func msg(code int, erro string, data interface{}) *APIMessage {
	return &APIMessage{
		Code: code,
		Erro: erro,
		Data: data,
	}
}

////////////////////////////////////////////////////// TOOLS FUNC END

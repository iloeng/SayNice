package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

// Feeling 作者的感觉、感受和心情
// Name comment:'感受、心情';
type Feeling struct {
	ID   uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Name string `json:"name" gorm:"type:text;default:'';not null;"`
}

// Emoji 表情
// Code comment:'emoji Unicode 编码';
// Name comment:'表情、态度';
// Feelings comment:'这种表情或态度所代表的感受或心情'
type Emoji struct {
	ID       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT:0"`
	Code     string    `json:"code" gorm:"type:string;"`
	Show     string    `json:"show" gorm:"type:string;"`
	Name     string    `json:"name" gorm:"type:text;default:'';not null;"`
	Feelings []Feeling `json:"feelings,omitempty" gorm:"many2many:emoji_feeling;"`
}

func main() {
	// smiling := Feeling {ID: 10000, Name: "smiling"}
	// affection := Feeling {ID: 10002, Name: "affection"}
	// tongue := Feeling {ID: 10003, Name: "tongue"}
	// skeptical := Feeling {ID: 10004, Name: "skeptical"}
	// sleepy := Feeling {ID: 10005, Name: "sleepy"}
	// unwell := Feeling {ID: 10006, Name: "unwell"}
	// hat := Feeling {ID: 10007, Name: "hat"}
	// glasses := Feeling {ID: 10008, Name: "glasses"}
	// concerned := Feeling {ID: 10009, Name: "concerned"}
	// negative := Feeling {ID: 10010, Name: "negative"}
	// costume := Feeling {ID: 10011, Name: "costume"}
	// cat := Feeling {ID: 10012, Name: "cat"}
	// monkey := Feeling {ID: 10013, Name: "monkey"}
	// emotion := Feeling {ID: 10014, Name: "emotion"}
	// weather := Feeling {ID: 10015, Name: "weather"}

	fi, e := os.Open("./emoji.txt")
	if nil != e {
		panic(e)
	}
	defer fi.Close()

	feelings := make(map[string]Feeling)
	var feeling Feeling
	var emojiList []Emoji

	buf := bufio.NewReader(fi)
	skip := false

	feelingID := uint(10000)
	emojiID := uint(10000)

	for {
		a, _, e := buf.ReadLine()
		if io.EOF == e {
			break
		}

		if skip {
			skip = false
			continue
		}

		// fmt.Println(len(a))
		line := string(a)
		re := regexp.MustCompile(" +")

		if "\n" == line || 0 == len(line) {
			continue
		} else if "#" == line[0:1] {
			feeling = addFeelingIfNeed(feelingID, line[12:], feelings)
			feelingID++

			continue
		}

		line = re.ReplaceAllString(line, " ")
		fileds := strings.Split(line, " ")

		var emoji Emoji

		emoji.ID = emojiID
		emoji.Code = fileds[0]

		if ";" == fileds[1] {
			emoji.Show = fileds[4]
			emoji.Name = strings.Join(fileds[6:], " ")
		} else {
			emoji.Show = fileds[5]
			emoji.Name = strings.Join(fileds[7:], " ")
			skip = true
		}

		emoji.Feelings = append(emoji.Feelings, feeling)

		emojiList = append(emojiList, emoji)

		emojiID++
	}

	emojiListByte, e := json.Marshal(emojiList)

	if nil != e {
		panic(e)
	}

	e = ioutil.WriteFile(time.Now().Local().Format("emoji_20060102150405.json"), emojiListByte, os.ModeAppend)

	if nil != e {
		panic(e)
	}
}

func addFeelingIfNeed(id uint, name string, feelings map[string]Feeling) Feeling {
	fields := strings.Split(name, "-")

	for i := 0; i < len(fields); i++ {
		if fields[i] == "face" {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}

	feelingName := strings.Join(fields, " ")

	feeling, ok := feelings[feelingName]

	if ok {
		return feeling
	}

	feeling = Feeling{
		ID:   id,
		Name: feelingName,
	}

	feelings[feelingName] = feeling

	return feeling
}

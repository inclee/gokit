package pagecache

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Feed struct {
	Score int `json:"score"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func (f Feed)Description ()string{
	return fmt.Sprintf("score: %d, title :%s ,content: %s",f.Score,f.Title,f.Content)
}

func TestNewPageCache(t *testing.T) {
	conf := Config{}
	note := NoteConfig{}
	note.Host = "localhost"
	note.Port = 6379
	note.DbNum = 1
	note.Key = "redis_cluster_1"

	conf.Notes = append(conf.Notes,note)
	rc,err := NewRedisPage(conf)
	if err != nil {
		t.Errorf("new redis page %s ",err)
	}
	Register("redis",rc)
	c := NewPageCache("redis")
	for i := 100; i < 200 ; i++  {
		fs := Feed{
			Score:i,
			Title:fmt.Sprintf("title_%d",i),
			Content:fmt.Sprintf("content_%d",i),
		}
		da, err := json.Marshal(fs)
		if err != nil {
			t.Errorf("json marshal %s",err)
		}
		err = c.Append("feed",da,i)
		if err !=  nil {
			t.Errorf("append %s",err)
		}
	}
	items ,err := c.RePage("feed",0,10)
	if err != nil {
		t.Errorf("repage %s",err.Error())
	}
	for _,item := range items  {
		da := item.([]byte)
		f := Feed{}
		err := json.Unmarshal(da,&f)
		if err != nil {
			t.Errorf("json unmarshal %s",err)
		}else{
			t.Logf("repage %s",f.Description())
		}
	}
}

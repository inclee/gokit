package pagecache

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"stathat.com/c/consistent"
	"errors"
	"time"
)

type Config struct {
	Notes []NoteConfig `json:"notes"`
}
type NoteConfig struct {
	Host  string `json:"host"`
	Port  int  `json:"port"`
	DbNum int `json:"dbNum"`
	Key string `json:"key"`
	Password string `json:"password"`
	MaxIdel int `json:"max_idel"`
}


type RedisPage struct {
	cluster *consistent.Consistent
	notes map[string]*note
}


func NewRedisPage(cf Config)(*RedisPage,error){
	ins := &RedisPage{}
	err := ins.init(cf)
	return ins,err
}

func (rp *RedisPage)init(conf Config) error{
	rp.cluster = consistent.New()
	rp.notes = make(map[string]*note,0)
	for _,ci := range  conf.Notes {
		rp.cluster.Add(ci.Key)
		n := &note{}
		err := n.init(ci)
		if err != nil {
			return  err
		}
		rp.notes[ci.Key] = n
	}
	return nil
}
func  (rp *RedisPage) getNote(key string)(*note,error){
	noteK,err := rp.cluster.Get(key)
	if err != nil {
		return nil,err
	}
	logs.Debug("Find NoteK %s",noteK)
	for nk,_ := range rp.notes{
		logs.Debug("Note Key %s",nk)
	}
	if note,ok := rp.notes[noteK] ;ok{
		return note,nil
	}
	return nil,errors.New("can't find cache note for key:"+key)
}

func (rp *RedisPage)Append(key string,data interface{},score int)error{
	if note, err := rp.getNote(key);err == nil {
		return note.append(key,data,score)
	}else {
		return err
	}
}
func (rp *RedisPage)Remove(key string,data interface{})error{
	if note, err := rp.getNote(key);err == nil {
		return note.remove(key,data)
	}else {
		return err
	}
}
func (rp *RedisPage)RePage(key string,start int ,stop int)([]interface{},error){
	if note, err := rp.getNote(key);err == nil {
		return note.rePage(key,start,stop)
	}else {
		return nil,err
	}
}

type note struct {
	p        *redis.Pool // redis connection pool
	host     string
	port     int
	conninfo int
	dbNum    int
	key      string
	password string
	maxIdle  int
}
func (n *note)init(cf NoteConfig)error{
	if len(cf.Key) == 0{
		cf.Key = "Redis"
	}
	if cf.MaxIdel == 0 {
		cf.MaxIdel = 3
	}
	n.key = cf.Key
	n.host = cf.Host
	n.port = cf.Port
	n.dbNum = cf.DbNum
	n.password = cf.Password
	n.maxIdle = cf.MaxIdel

	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", fmt.Sprintf("%s:%d",n.host,n.port))
		if err != nil {
			return nil, err
		}

		if n.password != "" {
			if _, err := c.Do("AUTH", n.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", n.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	n.p = &redis.Pool{
		MaxIdle:     n.maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}

	c := n.p.Get()
	defer c.Close()
	return c.Err()
}


func (n *note)append(key string,data interface{},score int)error{
	c := n.p.Get()
	defer c.Close()
	args := make([]interface{},0,0)
	args =  append(args, key)
	args =  append(args, score)
	args =  append(args, data)
	_,err := c.Do("ZADD",args...)
	return err
}

func (n *note)remove(key string,data interface{})error{
	c := n.p.Get()
	defer c.Close()
	args := make([]interface{},0,0)
	args =  append(args, key)
	args =  append(args, data)
	_,err := c.Do("ZREM",args...)
	return err
}
func (n *note)rePage(key string,start int ,stop int)([]interface{},error){
	c := n.p.Get()
	defer c.Close()
	args := make([]interface{},0,0)
	args =  append(args, key)
	args =  append(args,start)
	args =  append(args,stop)
	ret,err := c.Do("ZREVRANGE",args...)
	return ret.([]interface{}),err
}
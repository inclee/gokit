package pagecache


type PageCache interface {
	Append(key string,data interface{},score int)error
	Remove(key string,data interface{})error
	RePage(key string,start int ,stop int)([]interface{},error)
}
var (
	caches = make(map[string]PageCache,0)
)

func Register(adapter string,cache PageCache)  {
	caches[adapter] = cache
}

func NewPageCache(adapter string)PageCache{
	if  cache ,ok := caches[adapter];ok{
		return cache
	}
	return nil
}

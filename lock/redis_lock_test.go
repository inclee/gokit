package lock

import (
	"fmt"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func TestRedisLock(t *testing.T) {
	lock := NewRedisLock(redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	}), "my_idx_lock")
	ok, err := lock.TryLock()
	fmt.Println("-", ok, err)
	alock := NewRedisLock(redis.NewClient(&redis.Options{}), "my_idx_lock")
	i := 0
	for {
		ok, err = alock.TryLock()
		fmt.Println("二 lock ", ok, err)
		time.Sleep(time.Second * 10)
		i += 1
		if i > 5 {
			break
		}
	}
	fmt.Println("一 unlock", lock.UnLock())
	ok, err = alock.TryLock()
	fmt.Println("二 lock ", ok, err)
	fmt.Println(" 二 unlock", alock.UnLock())
}

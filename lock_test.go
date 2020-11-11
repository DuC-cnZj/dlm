package lock

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
	"time"
)

var rdx *redis.Client

func TestMain(t *testing.M) {
	rdx = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	t.Run()
}

func TestLock(t *testing.T) {
	var (
		num int
		wg  = &sync.WaitGroup{}
		n   = 1000
	)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			lockOne := NewRpcLock(rdx, "lock", WithEX(100), WithOwner("duc"))
			if lockOne.Acquire() {
				defer lockOne.Release()
				time.Sleep(1 * time.Second)
				num++
			}
		}()
	}
	wg.Wait()
	fmt.Println(num)
}

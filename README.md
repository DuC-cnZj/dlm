# DLM golang 分布式锁

## Usage

```go
lock := NewLock(rdx, "lock", WithEX(100), WithOwner("duc"))
if lock.Acquire() {
    // do sth.
    lock.Release()
}
```
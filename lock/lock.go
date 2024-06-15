package lock

type Lock interface {
	TryLock() (ok bool, err error)
	UnLock() error
}

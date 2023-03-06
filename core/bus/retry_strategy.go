package bus

type Retriable func() (any, error)

// RetryStrategy allows to define a retry strategy to be used in EventBus.
// For some noncritical operations, it is better to fail as soon as possible rather than retry a coupe of times.
// For example, it is better to fail right after a smaller number of retries with only a short delay between retries,
// and display a message to the user.
type RetryStrategy interface {
	Acquire(task Retriable) (any, error)
}

package connection

type Factory interface {
	Create(options *Options) (Connection, error)
}

package credentials

// A Provider allows you to provide credentials and load it in runtime.
// The configuration will invoke one provider at a time and only continue to the next
// if no credentials have been located. For example, if the process finds values defined via the 'BRIGHT_TOKEN'
// environment variables, the file at '.sectesterrc' will not be read.
type Provider interface {
	Get() *Credentials
}

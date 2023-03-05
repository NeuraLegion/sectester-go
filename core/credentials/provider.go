package credentials

type Provider interface {
	Get() *Credentials
}

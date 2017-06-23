package gateways

type Gateway interface {
	Name() (name string)
	Start() (err error)
	Stop()
}

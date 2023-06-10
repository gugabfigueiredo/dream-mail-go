package service

type IProvider interface {
	SendMail() error
}

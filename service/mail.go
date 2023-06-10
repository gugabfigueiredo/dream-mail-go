package service

import "github.com/gugabfigueiredo/dream-mail-go/model"

type IProvider interface {
	Send(model.Mail) error
}

type Provider struct {
	Client any
}

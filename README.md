# dream-mail-go

A coding challenge generic email service written in [Golang](https://go.dev/) that can be spun up to send emails 
through HTTP API calls or by implementing directly into apps.

This is a work in progress and is not ready for production use.

## Running

clone the project locally and build the docker image
```bash
git clone
cd dream-mail-go && make docker-build
```

update your .env file with desired provider credentials

run the docker container:
```bash
make docker-run
```

## Testing

to run the tests:
```bash
go test ./...
```

### API

check the API documentation [here](.swagger/swagger.yaml)

## Usage

You can import the package and use it in your code:

```go
package main

import (
	"github.com/gugabfigueiredo/dream-mail-go/service"
	"github.com/gugabfigueiredo/dream-mail-go/models"
)

//use the builtin providers or implement your own provider with IProvider interface
type MyProvider struct {
	//...
}

func NewMyProvider() *MyProvider {
	return &MyProvider{}
}

func (p *MyProvider) SendMail(mail *models.Mail) error {
	//...
}

//...

func main() {

	mailService := service.NewService([]service.IProvider{
		NewMyProvider(),
		service.NewSMTPProvider("smtp.domain.com", 587, "username", "password", Logger),
	}, Logger)

	mailService.SendMail(&models.Mail{
		From:    "sender@domain.com",
		To:      []string{"recipient@domain.com"},
		Subject: "Hello World",
		Text:    "Hello World",
	})
}
```

## Providers

The service supports the following providers:
- [SendGrid](https://sendgrid.com/)
- [Sparkpost](https://www.sparkpost.com/)
- [Amazon SES](https://aws.amazon.com/ses/)

it defaults to sending pure SMTP messages if no provider is available

## License

[MIT](https://choosealicense.com/licenses/mit/)
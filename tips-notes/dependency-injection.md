# Dependency Injection

## 1. Overview

- [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection): A dependency can be anything that effects the behavior or outcome of your logic.
  - Other services: Making your code more modular, less duplicate code and more testable.
  - Configuration: Such as a database passwords, URL endpoints.
  - System or environment state: Such as the clock or file system.
  - _Stubs or external APIs. So that API requests can be mocked within the system during tests to keep things stable and quick._
- Terminology:
  - _Service_: an instance of a class.
  - _Container_: a collection of services. Services are lazy-loaded and only initialized when they are requested from the container.
  - _Singleton_: an instance that is initialised once, but can be reused many times.

## 2. Example

- Source: https://elliotchance.medium.com/a-new-simpler-way-to-do-dependency-injection-in-go-9e191bef50d5

```go
// service sends an email
type SendEmail struct {
    From string
}

func (sender *SendEmail) Send(to, subject, body string) error {
    // It sends an email here, and perhaps returns an error.
}

// service welcomes new customers
type CustomerWelcome struct{}func (welcomer *CustomerWelcome) Welcome(name, email string) error {
    body := fmt.Sprintf("Hi, %s!", name)
    subject := "Welcome"    emailer := &SendEmail{
        From: "hi@welcome.com",
    }
    return emailer.Send(email, subject, body)
}

// main.go
welcomer := &CustomerWelcome{}
err := welcomer.Welcome("Bob", "bob@smith.com")
// check error...

// ----------------------------------------------------------------
// Major problem: difficult to unit test. Can't actually send email,
// verify that correct customer receives the correctly formatted
// email message.
// ----------------------------------------------------------------
// DI
// ----------------------------------------------------------------
// EmailSender provides an interface so we can swap out the
// implementation of SendEmail under tests.
type EmailSender interface {
    Send(to, subject, body string) error
}

type CustomerWelcome struct{
    Emailer EmailSender
}

func (welcomer *CustomerWelcome) Welcome(name, email string) error {
    body := fmt.Sprintf("Hi, %s!", name)
    subject := "Welcome"

    return welcomer.Emailer.Send(email, subject, body)
}
// main.go
emailer := &SendEmail{
    From: "hi@welcome.com",
}
welcomer := &CustomerWelcome{
    Emailer: emailer,
}
err := welcomer.Welcome("Bob", "bob@smith.com")
// check error...
// ----------------------------------------------------------------
// write unit test
// ----------------------------------------------------------------
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type FakeEmailSender struct {
    mock.Mock
}

func (mock *FakeEmailSender) Send(to, subject, body string) error {
    args := mock.Called(to, subject, body)
    return args.Error(0)
}

func TestCustomerWelcome_Welcome(t *testing.T) {
    emailer := &FakeEmailSender{}
    emailer.On("Send",
        "bob@smith.com", "Welcome", "Hi, Bob!").Return(nil)

    welcomer := &CustomerWelcome{
        Emailer: emailer,
    }
    err := welcomer.Welcome("Bob", "bob@smith.com")
    assert.NoError(t, err)
    emailer.AssertExpectations(t)
}
```

- Trade off:
  - **Duplicate code**: Imagine needing to use CustomerWelcome in more than one place (several, or even dozens). Now we have duplicate code that initialises the SendEmail. Especially repeating the From.
  - **Complexity and misunderstanding**: To use a service we now have to know how to setup all of its dependencies. Each of its dependencies may, in turn, have their own. Also, there may be several ways to satisfy dependencies that compile correctly but provide the wrong runtime logic. For example, if we has a EmailCustomer and EmailSupplier that both implemented EmailSender . We might provide the wrong service and a customer receives a message that should have been sent to a supplier.
  - **Maintainability rapidly decreases**: If a service initialization needs to change, or it needs new dependencies you now have to refactor all cases where it is used.
- Building the Services with Functions

```go
func CreateSendEmail() *SendEmail {
    return &SendEmail{
        From: "hi@welcome.com",
    }
}

func CreateCustomerWelcome() *CustomerWelcome {
    return &CustomerWelcome{
        Emailer: CreateSendEmail(),
    }
}

welcomer := CreateCustomerWelcome()
err := welcomer.Welcome("Bob", "bob@smith.com")
// check error...

// unit test
func TestCustomerWelcome_Welcome(t *testing.T) {
    emailer := &FakeEmailSender{}
    emailer.On("Send",
        "bob@smith.com", "Welcome", "Hi, Bob!").Return(nil)

    welcomer := CreateCustomerWelcome()
    welcomer.Emailer = emailer
    err := welcomer.Welcome("Bob", "bob@smith.com")
    assert.NoError(t, err)

    emailer.AssertExpectations(t)
}
```

- Singleton.

```go
type Container struct {
    CustomerWelcome *CustomerWelcome
    SendEmail       EmailSender
)

func (container *Container) GetSendEmail() EmailSender {
    if container.SendEmail == nil {
        container.SendEmail = &SendEmail{
            From: "hi@welcome.com",
        }
    }

    return container.SendEmail
}

func (container *Container) GetCustomerWelcome() *CustomerWelcome {
    if container.CustomerWelcome == nil {
        container.CustomerWelcome = &CustomerWelcome{
            Emailer: container.GetSendEmail(),
        }
    }

    return container.CustomerWelcome
}

// unit test
func TestCustomerWelcome_Welcome(t *testing.T) {
    emailer := &FakeEmailSender{}
    emailer.On("Send",
        "bob@smith.com", "Welcome", "Hi, Bob!").Return(nil)

    container := &Container{}
    container.SendEmail = emailer

    welcomer := container.GetCustomerWelcome()
    err := welcomer.Welcome("Bob", "bob@smith.com")
    assert.NoError(t, err)
    emailer.AssertExpectations(t)
}
```

## 3. DI Framework

### 3.1. Google wire

- [Wire](https://github.com/google/wire) is a code generation tool that automates connecting components using dependency injection

### 3.2. Uber-go fx

- [Fx](https://github.com/uber-go/fx) is an application framework for Go that:
  - Make dependency injection easy.
  - Eliminate the need for global state and `func init()`.

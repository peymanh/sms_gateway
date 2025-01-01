## Introduction

This is SMS Gateway Golang project, done as the code challenge of ArvanCloud interview and hiring process.

## Technology Stack 
| Name                  | Version |
|-----------------------|---------|
| Go                    | 1.22.0  |
| Postgre DB            | 4.2.17  |
| Docker                | 23.0.5  |


## First Time?

1. Clone the project:
```bash
git clone https://github.com/peymanh/sms_gateway
```

2. change directory to the project folder
```bash
cd sms_gateway
```

3. Install the dependencies:
```bash
go get -u
```

4. Build and start the database container with docker compose:
```bash
docker compose up --build
```
5.Run the migrations
```bash
go run migrate/migrate.go
```

6. Run the project
```bash
go run main.go
```


Now, the APIs are available on: **localhost:8000**

## APIs Descriptions

| route           | Description                             | Method |
|-----------------|-----------------------------------------|--------|
| /api/healthchecker        | healt check the server is running       | GET    |
| /api/auth/register        | Register new user                       | POST   |
| /api/auth/login           | Login an existing user                  | POST   |
| /api/auth/refresh         | Refresh access token                    | POST   |
| /api/auth/logout          | Logout an existing user                 | POST   |
| /api/users/me             | Get information of the logged in user   | POST   |
| /api/users/balance/update | Update user balance with specfic amount | POST   |
| /api/sms/send             | Send SMS                                | POST   |
| /api/sms/logs             | Get SMS Logs sent by the logged in user | POST   |

## Implementation Description
__________________
### User model and Authorization

> /register Method: POST

> /login Method: POST

In this project, We have used the Access token and Refresh JWT token for the authentication. Example of register request:

```json
{
	"name": "peyman",
    "email": "hassanabadi.peyman@gmail.com",
    "password": "asdasdasd",
    "passwordConfirm": "asdasdasd"
}
```
and in the response of login service an `access_token` is returned that should be in header of all other services.

```json
{
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzU3Mjk4NzksImlhdCI6MTczNTcyODk3OSwibmJmIjoxNzM1NzI4OTc5LCJzdWIiOiIxZWU3MGY3ZS0xOWFmLTRmNzctOTJiMS00MjE4MGE4ZmY0MDMifQ.cGU5maN27YbDyDIoZPXwUDvqtAIdKmIfT4WsNf8AJNXF1Vd_SvsP6uFFfUKposFbm-54FAemDj95LaZ6YI29jg",
    "status": "success"
}
```

The User model is defined this way:

```go
const (
	UserTypeNormal  UserType = 0
	UserTypePremium UserType = 1
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Balance   int       `gorm:"not null;default:0"`
	Class     UserType  `gorm:"not null;default:0"`
	SMSLogs   []SMSLog  `json: "logs"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
```
We have defined two type of users: `Normal` and `Premium`



## Send SMS Service

![SendSMSDiagram.jpg](SendSMSDiagram.jpg)


This is the most important service of the project.
When a User calls `/api/sms/send`, we follow these steps:
1. Check if the user is authenticated.
2. Check if he has sufficient balance.
3. If True, A log is inserted for the sms with `Pending` status.
4. A goroutine is dispatch to send sms. It checks the user type and according to that, sends the work to different channels
5. Receive the sms status (`Success` of `Failed`) and update it in DB.

>Note: the size of 2 channels are different to ensure that the procedure is not blocked for the Premium users.

>Note: Status of the SMS is mocked randomly. SMS gets successful with 0.9 probability.


## Stay in touch

- Author - Peyman Hassanabadi



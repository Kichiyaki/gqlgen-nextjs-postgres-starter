# Backend

Remember to create config.json file!

```
{
  "application": {
    "name": "gqlgen-nextjs-postgres-starter",
    "address": ":1234",
    "frontend": "http://localhost:3000",
    "debug": false,
    "intervalBetweenTokensGeneration": 5,
    "resetPasswordTokenExpiresIn": 5,
    "registrationDisabled": false,
    "cors": {
      "allowOrigins": ["*"],
      "allowCredentials": true
    },
    "defaultLanguage": "en",
    "bodyLimit": "12M"
  },
  "db": {
    "user": "postgres",
    "password": "",
    "addr": "localhost:5432",
    "name": "gqlgen_nextjs_postgres_starter"
  },
  "session": {
    "secret": "sessionSecret",
    "cookie": {
      "sessionName": "starter.sess",
      "secure": false,
      "httpOnly": true,
      "domain": "localhost",
      "sameSite": "lax or strict",
      "maxAge": 86400
    }
  },
  "email": {
    "host": "emailHost",
    "port": 587,
    "username": "emailUsername",
    "password": "emailPassword",
    "address": "its for 'From' email header"
  }
}

```

## Development

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Golang

### Installing

1. Clone this repository.
2. Navigate to the right directory.
3. Create config.json file (example config file is above).
4. Type "go run main.go" in your command prompt/terminal or whatever.
5. App should start.

## Tech/framework used

<b>Built with</b>

- [Echo](https://echo.labstack.com/)
- [gqlgen](https://github.com/99designs/gqlgen)
- [go-pg](https://github.com/go-pg/pg)

## Tests

Remember to run PostgreSQL locally with default credentials or to set these env variables:

1. POSTGRE_USER
2. POSTGRE_PASSWORD
3. POSTGRE_TEST_DATABASE
4. POSTGRE_ADDR

# GraphQL Starter - Backend

## Example config file (cmd/serve/config.json)

```
{
  "application": {
    "port": 8080,
    "frontend": "http://localhost:3000",
    "name": "graphql-starter"
  },
  "database": {
    "uri": "127.0.0.1",
    "username": "postgres",
    "password": "",
    "name": "graphql-starter",
    "port": "5432"
  },
  "email": {
    "uri": "smtp.gmail.com",
    "port": 587,
    "username": "example@gmail.com",
    "password": "passwordf123",
    "address": "example@gmail.com"
  },
  "session": {
    "name": "graphql-starter.sess",
    "secretKey": "key",
    "store": {
      "address": "127.0.0.1:6379",
      "password": "",
      "db": 0
    }
  }
}

```

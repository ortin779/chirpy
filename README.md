# Chirpy Rest api

Chirpy is a rest api that allows users to post something on the internet. This will be visible to other users.

This project is done as part of the [Boot.dev](https://www.boot.dev/) building servers with golang.It has the following features.

1. User creation and updation
2. Authentication and Authorization
3. Posting Chirps and Getting Chirps
4. Supporting a webhook endpoints to thirdparty services

## Setup Locally

- Clone the repo into ur local machine

```bash
git clone https://github.com/ortin779/chirpy.git
```

- Now open this folder in ur preferred code editor.
- Create a new file named `.env` and copy the contents of `.env.example` into it.
- Update the values of `.env` with ur configuration
- now run `go build -o chirpy && ./chirpy`

### API Documentation

/api/users -- [users](./docs/users.md)

/api/chirps -- [chirps](./docs/users.md)

/api/login -- [auth](./docs/auth.md)

# go-auth-service
Authentication service in go

# Prerequisites

Ensure that Docker and Docker Compose are installed on your machine.

# Setting up go-auth-service

Follow these steps to set up the application:

1.Clone the repository:
```
git clone https://github.com/sourav014/go-auth-service.git
```
2.Start the application using Docker Compose:
```
docker-compose up -d
```
Alternatively, if you're using a newer version of Docker:
```
docker compose up -d
```
3.Access the application: After the setup is complete, the application will be running on port 8080.

# API Endpoints

Once the setup is done, feel free to explore the API endpoints listed below.

## User SignUp
```
## Request:
------------------------------------------------------------

curl --location 'http://127.0.0.1:8080/api/v1/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "TestUser",
    "email": "test.user@gmail.com",
    "password": "mypassword",
    "is_admin": true
}'

## Response
{
    "name": "TestUser",
    "email": "test.user@gmail.com",
    "is_admin": false
}

```

## User Login
```
## Request:
------------------------------------------------------------

curl --location 'http://127.0.0.1:8080/api/v1/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test.user@gmail.com",
    "password": "mypassword"
}'

## Response
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiaXNfYWRtaW4iOnRydWUsImV4cCI6MTczNTk5NjYwOCwiaWF0IjoxNzM1OTk1NzA4LCJqdGkiOiJhNWM1MzY4OC0yMGI0LTQ3YmMtOWRlOS0yNjEwODFlNTI4ODMifQ.Ks_roRCgxGzvpgk2c51sV6b8aKb0AnKbfXl14MNeuP0",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiaXNfYWRtaW4iOnRydWUsImV4cCI6MTczNjA4MjEwOCwiaWF0IjoxNzM1OTk1NzA4LCJqdGkiOiI3ZjY2Mjc4Zi1kNTJkLTQyNDMtOWE3NC02ZDI0ZDlhMTdjMzUifQ.TT_DgUucBHDBxTTZvtcdiBFkzb2jg-KTYkqFNwpw4zs",
    "session_id": "7f66278f-d52d-4243-9a74-6d24d9a17c35",
    "user": {
        "name": "TestUser",
        "email": "test.user@gmail.com",
        "is_admin": true
    }
}

```

## User Renew Access Token
```
## Request:
------------------------------------------------------------

curl --location 'http://127.0.0.1:8080/api/v1/auth/renew' \
--header 'Content-Type: application/json' \
--data '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiaXNfYWRtaW4iOnRydWUsImV4cCI6MTczNjA4MjEwOCwiaWF0IjoxNzM1OTk1NzA4LCJqdGkiOiI3ZjY2Mjc4Zi1kNTJkLTQyNDMtOWE3NC02ZDI0ZDlhMTdjMzUifQ.TT_DgUucBHDBxTTZvtcdiBFkzb2jg-KTYkqFNwpw4zs"
}'

## Response
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiaXNfYWRtaW4iOnRydWUsImV4cCI6MTczNTk5NjY3NSwiaWF0IjoxNzM1OTk1Nzc1LCJqdGkiOiJlODRiNTI2ZS05NTZjLTQ2MDMtYThlOC02ZmYxYjdkMGE5ZjAifQ.dzrcrq_Ej07A4c1zFYbyTK7uc4JfuLH4FQ3kRhnlX3k",
    "access_toker_expires_at": "2025-01-04T13:17:55Z"
}

```

## Get User Profile Details
```
## Request:
------------------------------------------------------------

curl --location 'http://127.0.0.1:8080/api/v1/user/profile' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {please use the access_token received from login api response to get the user profie.}'

## Response
{
    "name": "TestUser",
    "email": "test.user@gmail.com",
    "is_admin": false
}

```

## User Revoke Token
please use the session id received from login api response to revoke the token.
```
## Request:
------------------------------------------------------------

curl --location --request POST 'http://127.0.0.1:8080/api/v1/auth/revoke/7f66278f-d52d-4243-9a74-6d24d9a17c35' \
--header 'Content-Type: application/json'

## Response
{
    "message": "token revoked successfully"
}

```

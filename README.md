# UniqueSSO

[![Build Status](https://drone.hustunique.com/api/badges/UniqueStudio/UniqueSSO/status.svg?ref=refs/heads/dev)](https://drone.hustunique.com/UniqueStudio/UniqueSSO)

## Diagram 

### normal login

#### normal

```sequence
client --> sso: normal login 
client ->> sso: post
sso ->> sso: validate
sso ->>  client: success / failure
```

#### with redirect

```sequence
client --> sso: login with redirect
client ->> sso: post with service
sso ->> sso: validate
sso ->> client: 302(reidrect)
```

### third-party login

#### failure

```sequence
client ->> traefik: "request"
traefik ->> sso: ask for validate
sso ->> sso: validate (by session)
sso ->> traefik: 302(false)
traefik ->> client: 302 redirect
client ->> sso: /login
```

#### success

```sequence
client ->> traefik: request
traefik ->> sso: ask for validate
sso ->> sso: validate (by session)
sso ->> traefik: 200(success) \n Append `X-UID` header
traefik ->> app: request with `X-UID` header
app -->> app: 
app ->> traefik: resp
traefik ->> traefik: delete some header
traefik ->> client: resp
```

## Big picture

1. login at `POST /v1/login?type=${loginType}&service=${redirectURI}` with body 

for login, there are four ways to login:

1. phone number with password

2. phone sms

3. email address with password

4. lark oauth

store state as session, which persisted by redis.

## Deployment

1. edit the backend config file 


## TODO list

- [ ] Access APM systems

## Uniform Response

```json
{
  "serviceResponse": {
    "authenticationFailure": {
      "code": "",
      "description": ""
    },
    "authenticationSuccess": {
      "user": "${UID}",
      "attributes": {
        "uid": "",
        "name": "",
        "phone": "",
        "email": ""
      }
    }
  }
}
```
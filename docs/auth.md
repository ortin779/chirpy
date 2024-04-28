# Auth

### Login as an user

```
POST /api/login
```

This endpoint takes the user credentials, and validates them against stored creds. And then generates an access and refresh tokens.

```json
{
  "email": "abc@email.com",
  "password": "abc@123"
}
```

### Refresh the access-token

```
POST /api/refresh
```

- When the access-token expires, we are allowing user to refresh the token using the Refresh-token.
- This endpoint expects the refresh token to present as an Authorization header.
- If the refresh token is valid and not revoked then we will generate new access-token

### Revoke the access-token

```
POST /api/revoke
```

- As we have the refresh token, if user wants to revoke it in case of some security related issues.
- User can hit this endpoint and pass the refresh token as part of the authorization header.
- It will add the given token into revoked entries in our db.

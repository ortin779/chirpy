# Users

## /api/users

### Create a new user

```
POST /api/users
```

Creating a user. We need to pass the following as the request body.

```json
{
  "email": "abc@email.com",
  "password": "abc@123"
}
```

If user created successfully we will get back the user with id.

### Update a user

```
PUT /api/users/{userId}
```

To update an user, we need to pass the access_token as the Authorization header, and the request body should contain the following values

```json
{
  "email": "abc@email.com",
  "password": "abc@123"
}
```

If user updated successfully we will get back the updated user info.

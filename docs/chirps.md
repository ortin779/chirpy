# Chirps

## /api/chirps

### Create a new chirp

```
POST /api/chirps
```

To create a chirp, the user should pass the access_token as the Authorization header. And the the request body should contain the following info

```json
{
  "body": "iam a chirp"
}
```

If chirp created successfully we will get back the chirp with author info.

### Get Chirps

```
GET /api/chirps?sort=asc&author_id=2
```

The Get chirps is a public endpoint. This supports sorting and filtering. We can sort the chirps by their id and filter them by author. This will return an array of chirps.

### Get Chirp by Id

```
GET /api/chirps/{chirpId}
```

This endpoint is also public, which allows to get a particular chirp by its id.

### Delete a Chirp by Id

```
DELETE /api/chirps/{chirpId}
```

This endpoint is private, and requires access-token. We should pass it through Authorization header. If that chirp belongs to the user then we will delete it otherwise we will throw an authorization(403) Error.

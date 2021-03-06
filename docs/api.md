# API Reference

This API strives to follow the [REST](http://en.wikipedia.org/wiki/Representational_State_Transfer)
architectural style. It has resource-oriented
URLs and uses [JSON](https://www.json.org/) as the representation for resources.
It uses standard HTTP response codes, authentication, and verbs.

Resource URL's are documented using URI templates (when necessary),
as defined following the [RFC 6570](https://tools.ietf.org/html/rfc6570).

JSON request/responses are documented by listing the fields and
the expected type for that field (along with an optional annotation
if the field is expected to be optional). For example, with this specification:

```
{
    "field" : <boolean>,
    "optionalField" : <string>(optional)
}
```

You can expect a JSON like this:

```json
{
    "field": true,
    "optionalField" : "data"
}
```

All JSON fields documented as part of request/response bodies are
to be considered obligatory, unless they are explicitly
documented as optional. Also all fields of type **<string>** are
expected to be non-empty by default, unless stated otherwise.


# Error Handling

When an error occurs you can always expect an HTTP status code indicating the
nature of the failure and also a response body with an error message
giving some more information on what went wrong (when appropriate).

It follows this schema:

```
{
    "error": {
        "message" : <string>
    }
}
```

The **message** is intended for human inspection, no programmatic decision
should be made using their contents (code should handle them as opaque strings).


# Authentication

Authentication is done using bearer tokens transmitted via the HTTP header
**Authorization** as specified by the
[RFC 6750](https://tools.ietf.org/html/rfc6750#section-2.1).

Any time a request is mentioned to be "authenticated" it means it
requires a bearer token to be informed on the header according to the RFC.


## Sign In

To sign in send the following request:

```
POST /v1/auth/signin
```

With the following request body:

```
{
    "email" : <string>,
    "password" : <string>
}
```

Where **email** must be a valid email address from a previously created user.

In case of success you can expect an status code 201 and the following response:

```
{
  "access_token":<string>,
  "token_type":<string>,
  "expires_in":<number>
}
```

Where **token_type** indicates the type of the token (for now the API only
supports tokens of type "bearer") and the **access_token**
contains the token itself that you will use for authentication
on further requests.

The **expires_in** field specifies in seconds how long it will
take for the token to expire.


## Sign Out

To sign out send the following authenticated request:

```
POST /v1/auth/signout
```

In case of a success the bearer token used to authenticate the request
will become invalid and can't be used any further.


# Creating a new user

To create a new user send the request:

```
POST /v1/users
```

With the following request body:

```
{
    "fullname" : <string>,
    "email" : <string>,
    "password" : <string>
}
```

Multiple users with the same **email** is forbidden, so the given
**email** must not belong to an user already registered
in the system. The **fullname** is the full name of the user and
does not need to be unique.

In the case of success you can expect a status code 201 and the response:

```
{
    "id" : <string>
}
```

Where the **id** field is the ID of the created user.
You can use the newly created user to sign in.


# Listing Users

Only administrators are allowed to list all users.

To list all users send the authenticated request:

```
GET /v1/users
```

In the case of success you can expect an status code 200
and a list with all the users:

```
[
    {
        "id" : <string>,
        "fullname" : <string>,
        "email" : <string>
    },
    {
        "id" : <string>,
        "fullname" : <string>,
        "email" : <string>
    }
]
```

Ordered by the field **id** in descending order.

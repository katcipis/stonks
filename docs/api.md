# API Reference

This API strives to follow the [REST](http://en.wikipedia.org/wiki/Representational_State_Transfer)
architectural style. It has resource-oriented
URLs and uses [JSON](https://www.json.org/) as the representation for resources.
It uses standard HTTP response codes, authentication, and verbs.

Resource URL's are documented using URI templates, as defined following
the [RFC 6570](https://tools.ietf.org/html/rfc6570).

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
documented as optional.


# Error Handling

When an error occurs you can expect an HTTP status code indicating the
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

With the **email** and **password** in the request body following the schema:

```
{
    "email" : <string>,
    "password" : <string>
}
```

In case of success you will receive the following response:

```
{
  "access_token":<string>,
  "token_type":<string>,
  "expires_in":<number>
}
```

Where **token_type** is always "bearer" and the **access_token**
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


# Creating User

# Listing Users

# Deleting User

# Ulysses
Codename: Ulysses, a generic webhosting/webservice manager backend built in Go featuring high programmability, high concurrency and low resource consumption. 

## Repository Architecture
```
./conf                  # Configuration File Template
./doc                   # Documentations
./src                   # Source Code for binary executable(s)
    ./src/internal      # Source Code for internal utils indirectly related to program business logic
```

## Design Philosophy 

**Ulysses** is designed as a FaaS (Function-as-a-Service) server with RESTful API capabilities. 

In alignment with `Gin`'s specifications, **Ulysses** defines 2 different types of `gin.HandlerFunc`: `endpoint` and `checkpoint`. While both are but `alias` of `gin.HandlerFunc`, following guidelines should be respected:

- `checkpoint` acts like middleware, checking against a specific criteria of a *HTTP Request* represented by a `*gin.Context` and only makes a *HTTP Response* when the check is *FAILED*.
    - `checkpoint` makes responses with only `*gin.Context.AbortWithStatusJSON()`. And should `return` right after the response is made.
    - `checkpoint` is used to check for general purposes such as: Authentication/Authorization, MFA Challenge Response.  
- `endpoint` provides the actual implementation of *FaaS*, invoking a specific function which is usually* defined by a package under `Ulysses.Lib`.
    - `endpoint` assumes the *HTTP Request*, represented by a `*gin.Context`, is valid in terms of: `Authorization` Header, Proper user access, valid MFA challenge response (when needed).
    - `endpoint` does not assume the *HTTP Request*, represented by a `*gin.Context`, is valid in terms of: function input/parameters, volatile status (product stock, wallet balance, etc.), request integrity.

In terms of returning a HTTP Status Code: 

**401 vs 403**
- Return `401 Unauthorized` when the user is not *properly* logged-in. (i.e., Not including a valid `Authorization` header in the request)
- Return `403 Forbidden` when the user is *properly* logged-in, but lacking permission to the requested resource.

**400 vs 500**
- Return `400 Bad Request` when the user's request could be *trivially* proved to be incorrect, e.g., missing required parameters, length mismatch, requesting something unexist.
- Return `500 Internal Server Error` when the server errors unexpectedly, e.g., MySQL error that is not `sql.ErrNoRow`. This also includes the scenario when the user's request is causing the problem.
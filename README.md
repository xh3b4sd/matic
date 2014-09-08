# client-matic

### spec

#### user workflow

- user needs to write server code
- user needs to compile client
- user needs to compile api blueprint

#### server

- server needs to use middleware-server package
- middleware needs to use context responders

#### compiling workflow

- index registered routes
- index middlewares for each route
- index responders for each middleware
- index expected payload for each middleware
- index expected url query params for each middleware
- create client method for each route

#### client requirements

- client needs to be configured with http client
- client needs to be configured with global headers
- client needs to be configured with global url query params

- client methods need to return error
- client methods need to return http status code
- client methods need to return http response body as string
- client methods need to return http response body as read closer

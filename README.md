## Prerequisites 
To clone the repo on your local run
- Run `git clone git@github.com:syedvasil/blogging_platform.git`

Make sure go and mongoDB are installed locally 
Follow [Here]((https://https://github.com/syedvasil/blogging_platform/blob/main/DatabaseREADME.MD))  
Once all the steps mentioned above are completed




## 🔧 Installation

navigate into the directory and make sure you are in **Main** branch.

make sure u have the latest go, and MongoDB installed as this would be required  
- Run `go mod download`, to download dependencies.
- open the config/config.go and make change to ConstCFG if needed 
- Run `go run ./cmd/server/` to instantiate a local http server for development 


once the server is running use curl or postman to call the APIs  with basic Auth
change the port if required
```
curl --location 'http://localhost:8080/api/v1/posts' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic SmFuZURvZTpwYXNzd29yZDI=' \
--data '{
        "title": "test1",
    "content": "temp content"
}'
```

## Endpoints
 
todo add openAPI docs
Posts

    POST /posts - Create a new post 
    GET /posts - Get a list of posts (with optional filters and pagination)
    GET /posts/:id - Get a single post by ID
    PUT /posts/:id - Update a post by ID
    DELETE /posts/:id - Soft delete a post by ID

Users

    POST /users - Create a new user
    GET /users - Get a list of users (with pagination)
    GET /users/:id - Get a single user by ID
    PUT /users/:id - Update a user by ID
    DELETE /users/:id - Soft delete a user by ID

with basic authorization on Read/Write based on owner and Admin can do every thing
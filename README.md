# Fitness CLub Manager API

Fitness CLub Manager is a back-end api project, which provides
endpoints for CRUD operations on the database. There is an authentication as well.
The API is intended to serve the [front-end](https://github.com/Ilia-tod29/FitnessClubManagerClient) 
of the Fitness Club Manager project.

## Endpoints
- POST(“/users”) - Create user profile
- POST(“/users/login”) - Login in existing user profile
- POST(“/tokens/renew_access”) - Renew access token
- PUT(“/users/:id”) - Update user; Only the field "suspended" can be altered
- GET(“/users/:id”) - Get user by id
- GET(“/allusers”) - Get all users
- GET(“/users”) - Get all users by pages; “page_size” and “page_id” must be passed as params.
- DELETE(“/users/:id”) - Delete a user by id
- POST(“/subscriptions”) - Create subscription
- GET(“/subscriptions/:id”) - Get subscription by id
- GET(“/subscriptions/user/:id”) - Get all subscriptions by user id
- GET(“/allsubscriptions”) - Get all subscriptions
- GET(“/subscriptions”) - Get all subscriptions by pages; “page_size” and “page_id” must be passed as params.
- DELETE(“/subscriptions/:id”) - Delete a subscription
- POST(“/inventoryitems”) - Create inventory item
- PUT(“/inventoryitems/:id”) - Edit inventory item
- GET(“/inventoryitems/:id”) - Get inventory item by id
- GET(“/allinventoryitems”) - Get all inventory item
- GET(“/inventoryitems”) - Get all inventory item by pages; “page_size” and “page_id” must be passed as params.
- DELETE(“/inventoryitems/:id”) - Delete inventory item
- POST(“/gallery”) - Create gallery item 
- GET(“/gallery/:id”) - Get gallery item by id
- GET(“/allgallery”) - Get all gallery items
- GET(“/gallery”) - Get all gallery items  by pages; “page_size” and “page_id” must be passed as params.

For more information about the body parameters that should be passed, please see the [api directory](https://github.com/Ilia-tod29/FitnessClubManagerAPI/tree/main/db). In the beginning of each
file there are the structures with ending - request. They contain the needed information.

## Database structure

Please see the [DBdiagram.pdf](https://github.com/Ilia-tod29/FitnessClubManagerAPI/blob/main/DBdiagram.pdf)

## How to run

- Ensure your local docker machine is working
- Pull the latest docker image:

```bash
docker pull postgres:latest
```

- Run the following commands(make sure you are in the folder of project) in order:
```bash
make postgres
```
```bash
make createdb
```
```bash
make migrateup
```
```bash
make server
```

- Set up the stripes webhook listener:
```bash
make server
```

## Development

- To make env changes ether do them in the [app.env file](https://github.com/Ilia-tod29/FitnessClubManagerAPI/blob/main/app.env)
or export the declared in the file variables in your local machine
- To make changes to the sqlc generated code:
  - update/create files in the [query directory](https://github.com/Ilia-tod29/FitnessClubManagerAPI/tree/main/db/query)
  - run:
  ```bash
    export SQLC_AUTH_TOKEN=sqlc_01HKVZ6M86VR8RQC9RDM6TFHPC
    ```
  - run:
  ```bash
    make sqlc
    ```
- To make db changes:
  - alter the createmigration command in the [Makefile](https://github.com/Ilia-tod29/FitnessClubManagerAPI/blob/main/Makefile)
  by adding the name of your migration at the end of the command
  - while having a created db run:
  ```bash
    make migrateup
    ```
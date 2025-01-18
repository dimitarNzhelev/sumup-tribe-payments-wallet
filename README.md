# tribe-payments-wallet-golang-interview-assignment

## Stack

- `Golang`
- `PostgreSQL`


## About  

This project was created as a take-home test for the Junior Software Engineer (Backend) - Payments Tribe role at SumUp. It implements a simple wallet API that allows users to manage their wallets efficiently.  

### Features  
The Wallet API enables users to:  
- **Create a wallet**  
- **View wallet balance**  
- **Deposit money**  
- **Withdraw money**  

Additionally, the API supports user creation and authentication, ensuring that users can only manage their own wallets. Every wallet transaction is securely logged in a transactions table, providing a clear record of all operations.  

---

## Getting Started  

### Prerequisites  
Before starting, ensure you have a `.env` file configured with the following environment variables:  

```plaintext
POSTGRES_USER=docker
POSTGRES_PASSWORD=docker
POSTGRES_DB=test_database
JWT_SECRET=secret
```
### Running the project
1. Create the `.env` file with the required variables.
2. Start all necessary containers and run database migrations using Docker Compose:
```bash
> docker-compose up
...
postgres-db-1  | 2025-01-18 15:52:29.062 UTC [1] LOG:  database system is ready to accept connections
migrate-1      | 1/u init (32.626855ms)
migrate-1 exited with code 0
go-api-1       | {"level":"info","time":"2025-01-18T15:52:29.816Z","caller":"cmd/api.go:70","msg":"Connecting to database: postgres://docker:docker@postgres-db:5432/test_database?sslmode=disable"}
go-api-1       | {"level":"info","time":"2025-01-18T15:52:29.820Z","caller":"http/server.go:94","msg":"Starting HTTP server","name":"api","address":"0.0.0.0:8080"}
```

### API and Database Details
 - **API**: Runs on http://localhost:8080.
 - **PostgreSQL** Database: Accessible on port 5432.


## ENDPOINTS

### Create User
* `verb` - `POST`
* `endpoint` - `/v1/user`

```bash
curl -X POST -H "Content-Type: application/json" \
-d '{
  "first_name": "John",
  "last_name": "Doe",
  "email": "johndoe@example.com",
  "password": "securepassword123"
}' \
http://localhost:8080/v1/user
```

* Possible responses:

  - 200 - Successfully created a user
  - 400 - Bad arguments/request
  - 422 - Error creating a wallet/Unprocessable entity
---
### Login
* `verb` - `POST`
* `endpoint` - `/v1/login`

```bash
curl -X POST -H "Content-Type: application/json" \
-d '{
  "email": "johndoe@example.com",
  "password": "securepassword123"
}' \
http://localhost:8080/v1/login
```
* Possible responses:

  - 200 - Successfully logged in
  - 400 - Bad arguments/request
  - 401 - Unauthorized (Wrong Password)
  - 422 - Error creating a JWT Token/Unprocessable entity

####  Successful Response Example (HTTP 200):

If the login is successful, the response body includes a JWT token that can be used for authentication in subsequent requests:

```JSON
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5kb2VAZXhhbXBsZS5jb20iLCJleHAiOjE3MzczMDI0OTgsImlhdCI6MTczNzIxNjA5OCwiaWQiOiIxMjZmMDFiMy1mYzgzLTQzOWItOWJkZi0xMTVmOGZmMGEyZmIifQ.dbfLDRdzzIm_v0FZPkxIhv0XPpz9QGBad-6z_D3sIIg"}
```

#### Using the Token:

The token must be included in the Authorization header as a Bearer Token to access wallet-related functionality. For example:

```
Authorization: Bearer <JWT_TOKEN>
```
---
### Create wallet

This endpoint allows authenticated users to create a new wallet. The wallet is associated with the authenticated user and initialized with a balance of 0.

* `verb` - `POST`
* `endpoint` - `/v1/wallet`

```bash
curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Bearer <JWT_TOKEN>" \
http://localhost:8080/v1/wallet
```

* Possible responses:

  - 201 - Successfully created a wallet
  - 401 - Unauthorized (Either wrong JWT token or missing user id from the context)
  - 422 - Error creating a wallet/Unprocessable entity

#### Example Response (201 Created)

If the wallet creation is successful, the server returns a JSON response with the newly created wallet's details:

```JSON
{
    "wallet_id":"f3742598-32ff-47c4-b254-6803ba31c2e6",
    "balance":0,
    "version":1,
    "created_at":"2025-01-18T16:06:28.855165Z",
    "updated_at":"2025-01-18T16:06:28.855165Z"
}
```
---
### View wallet
Get a wallet information.

* `verb` - `GET`
* `endpoint` - `/v1/wallet/{UUID}`


```bash
curl -X GET -H "Content-Type: application/json" \
-H "Authorization: Bearer <JWT_TOKEN>" \
http://localhost:8080/v1/wallet/{UUID}
```
UUID - The wallet ID

* Possible responses:

  - 200 - Successfully found wallet
  - 404 - Can't find wallet
  - 400 - Bad arguments/request

#### Example Response (200 OK)

If the wallet belongs to the user and the user is logged in, the API returns an example response as JSON:

```JSON
{
    "wallet_id":"f3742598-32ff-47c4-b254-6803ba31c2e6",
    "balance":0,
    "version":1,
    "created_at":"2025-01-18T16:06:28.855165Z",
    "updated_at":"2025-01-18T16:06:28.855165Z"
}
```
---
### Make a transaction
Make a deposit or withdraw.

* `verb` - `POST`
* `endpoint` - `/v1/wallet/{UUID}/transaction`

#### Request body
The body must include:
- amount: The amount of money to deposit or withdraw (a positive number).
- transaction_type: Either `"deposit"` or `"withdraw"`.

#### Example request for deposit
```bash
curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Bearer <JWT_TOKEN>" \
-d '{
  "amount": 100.50,
  "transaction_type": "deposit"
}' \
http://localhost:8080/v1/wallet/{UUID}/transaction
```

#### Example request for withdraw
```bash
curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Bearer <JWT_TOKEN>" \
-d '{
  "amount": 50.25,
  "transaction_type": "withdraw"
}' \
http://localhost:8080/v1/wallet/{UUID}/transaction
```

* Possible responses:

  - 200 - Successfully made the transaction
  - 400 - Bad arguments/request - Invalid or missing body.
  - 401 - Unauthorized - The JWT token is missing or invalid.
  - 422 - Unprocessable Entity - An error occurred while processing the transaction (e.g., insufficient balance for a withdrawal).

#### Example Response (200 OK)
If the transaction is completed successfully the client receives the updated wallet as JSON.

```JSON
{
    "wallet_id":"f3742598-32ff-47c4-b254-6803ba31c2e6",
    "balance":100.5,
    "version":2,
    "created_at":"2025-01-18T16:06:28.855165Z",
    "updated_at":"2025-01-18T16:26:01.27295Z"
}
```

## Future Improvements  

1. **JWT Authentication**:  
   I plan to upgrade the JWT authentication algorithm from **HS256** to **RS256** for improved security using asymmetric encryption.  

2. **Token Management**:  
   I aim to implement two types of tokens:  
   - **Access Token**: A short-lived token for authenticating API requests.  
   - **Refresh Token**: A longer-lived token for securely generating new access tokens without requiring the user to log in again.  
   This improves both security and user experience by limiting token exposure and enabling seamless token renewal.  

3. **Soft Delete for Users**:  
   Instead of using `ON DELETE SET NULL` for the `user_id` in the `wallets` table, I plan to implement a **soft delete** mechanism.  
   - Add a `deleted_at` column to the `users` table.  
   - When a user is deleted, instead of removing the record, set the `deleted_at` timestamp.  
   - This allows preserving user data for auditing or restoration while excluding "deleted" users from active operations.  
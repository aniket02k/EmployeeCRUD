# Employee Management API's implementation 

Employee management APIs built using Go, Gorilla Mux for routing, and Protobuf for request and response serialization. The API provides functionality to manage user details, including fetching, creating, and updating users.

## Protobuf Definitions

### File: `proto/employee.proto`

- **UserRequest**: Request for getting user details by `userId`.

  ```protobuf
  message UserRequest {
      string userId = 1;
  }
  
  ```

  

- **UserResponse**: Response containing user details including employee information.

  ```
  message UserResponse {
      string userId = 1;
      string firstname = 2;
      string lastname = 3;
      string email = 4;
      string employeeId = 5;
      string designation = 6;
  }
  ```

  

- **CreateUserRequest**: Request for creating a new user with fields like `firstname`, `lastname`, `email`, and `designation`.

  ```protobuf
  // Request for creating a new user
  message CreateUserRequest {
      string firstname = 1;
      string lastname = 2;
      string email = 3;
      string designation = 4;
  }
  
  ```

  

- **CreateUserResponse**: Response after creating a new user, containing the `userId`.

  ```protobuf
  message CreateUserResponse {
      string userId = 1;
  }
  
  ```

  

- **UpdateUserRequest**: Request for updating a user's email, containing `userId` and the new `email`.

  ```protobuf
  message UpdateUserRequest {
      string userId = 1;
      string email = 2;
  }
  ```

  

## API Endpoints

### 1. Get User Details

- **Endpoint**: `/assignment/user`

- **Method**: `GET`

- **Description**: Fetch user details by `userId`.

- **Request**:
   `{
      "userId": "string"
  }`

- **Respponse**:

  ` {
   "userId": "string",
   "firstname": "string",
   "lastname": "string",
   "email": "string",
   "employeeId": "string",
   "designation": "string"
  }`

### 2. Create New User

- **Endpoint**: `/assignment/user`

- **Method**: `POST`

- **Description**: Create a new user with the provided details.

- **Request**:

  `{
      "firstname": "string",
      "lastname": "string",
      "email": "string",
      "designation": "string"
  }`

- Response: 

  `{
      "userId": "string"
  }`

### 3. Update User Details(Email)

- **Endpoint**: `/assignment/user`

- **Method**: `PATCH`

- **Description**: Update the email address of an existing user.

- **Request**:

  `{
      "userId": "string",
      "email": "string"
  }`
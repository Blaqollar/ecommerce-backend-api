# E-Commerce Backend

This is an e-commerce backend written in Golang using the Gin framework. The backend provides various APIs to manage products, orders, and customers for an e-commerce application.

## Setup

1. Clone the repository:

```
git clone <repository-url>
```

2. Install Golang (if not already installed) from the official Golang website: https://golang.org/

3. Install project dependencies:

```
go mod download
```

4. Configure the environment variables:

   - Create a `.env` file in the project root directory.
   - Add the following environment variables to the `.env` file:

   ```plaintext
   MONGODB_DATABASE_NAME=
   MONGODB_TRASH_DATABASE_NAME=
   MONGODB_URI=
   ```

5. Run the application:

```
go run main.go
```

The application will start running on port 8000 if not specified in the `PORT` environment variable.

## API Documentation

### Products

- **GET /products**

  Fetches all products.

- **GET /products/:id**

  Fetches a single product by ID.

- **POST /products**

  Creates a new product.

- **PUT /products/:id**

  Updates a product by ID.

- **DELETE /products/:id**

  Deletes a product by ID.

### Orders

- **GET /orders**

  Fetches all orders.

- **GET /orders/:id**

  Fetches a single order by ID.

- **POST /orders**

  Creates a new order.

- **PUT /orders/:id**

  Updates an order by ID.

- **DELETE /orders/:id**

  Deletes an order by ID.

### Customers

- **GET /customers**

  Fetches all customers.

- **GET /customers/:id**

  Fetches a single customer by ID.

- **POST /customers**

  Creates a new customer.

- **PUT /customers/:id**

  Updates a customer by ID.

- **DELETE /customers/:id**

  Deletes a customer by ID.

## Database Configuration

This backend uses a MongoDB database. To configure the database connection, set the `MONGODB_URI` environment variable in the `.env` file with the URL of your MongoDB database.

## Contributing

Contributions are welcome! If you find any issues or want to add new features, please create a pull request with a detailed description of the changes.

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute the code as per the terms of the license.

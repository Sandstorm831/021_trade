# Trading Platform API

This is a simple API written in Gin that allows creation of user, stocks, reward and other portfolio stats.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   Go
*   PostgreSQL

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/Sandstorm831/021_trade.git
    ```
2.  Navigate to the project directory:
    ```sh
    cd 021_trade
    ```
3.  Install the dependencies:
    ```sh
    go mod tidy
    ```
4.  Create a `.env` file in the root of the project and add the following environment variables:
    ```env
    DB_DSN=
    # Example : "host=assignment user=assignment password=assignment dbname=assignment port=5432 sslmode=disable TimeZone=Asia/Kolkata"
    ```

## Running the application

> Arrange a PostgreSQL Database and add the credential in the `.env` file in the format mentioned above. Database can be hosted, local or docker based.

To build and run the application, execute the following command from the root of the project:

```sh
$ go build -o 021_trade ./cmd/server/

$ ./021_trade
```

to run the application directly
```sh
$ go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. User API

#### `POST` : `/create-user`

Creates a new user.

**Request Body:**

```json
{
  "Name": "John Doe"
}
```
If the `name` field is empty, a random name will be generated.

**Response:**
```json
{
  "User": {
      "ID": "f6189d03-8ed2-4b69-8ad9-9edbeb003489",
      "Name": "JORadYza",
      "CreatedAt": "2025-12-28T22:25:51.514736597+05:30"
  }
}
```

### 2. Stock API

#### `POST` : `/create-stock`

Creates a new stock and sets its initial price.

**Request Body:**
```json
{
  "Name": "Apple Inc.",
  "Symbol": "AAPL",
  "Price": 150.00
}
```

**Response:**
```json
{
  "Price": "243.241",
  "Stock": {
      "Symbol": "AAPL",
      "Name": "Apple",
      "IsActive": true
  }
}
```

### 3. Reward API

#### `POST` : `/reward`

Records a new reward for a user.

**Request Body:**
```json
{
	"UserID": "9bbf1625-e024-4663-9fc8-47440d1620c0",
	"StockSymbol": "INF",
	"Quantity": "25",
	"IdempotencyKey": "faisass",
	"RewardedAt": "2025-12-25 13:59:04.093438+05:30"
}
```
- `RewardedAt` field can be empty, when empty current timestamp will be inserted
- `IdempotencyKey` is any random string to avoid duplicate reward events

**Response:**
```json
{
    "Reward": {
        "ID": "3253d10e-6559-4662-94b5-ddb45341f603",
        "UserID": "9bbf1625-e024-4663-9fc8-47440d1620c0",
        "StockSymbol": "INF",
        "Quantity": "25",
        "IdempotencyKey": "faisass",
        "RewardedAt": "2025-12-25T13:59:04.093438+05:30"
    },
    "cashLedgerEntryID": 33,
    "feeLedgerEntryID": 34,
    "stockLedgerEntryID": 32
}
```

#### `GET` : `/today-stocks/:userId`

Return all stock rewards for the user for today.

**URL Path Parameter:**
* `userId` : UUID

**Response:**
```json
[
  {
      "ID": "caaa6801-3777-4ef9-a146-62dff3e54b12",
      "UserID": "9bbf1625-e024-4663-9fc8-47440d1620c0",
      "StockSymbol": "TAT",
      "Quantity": "2.23",
      "IdempotencyKey": "sonsvsfaasfaio",
      "RewardedAt": "2025-12-28T13:03:16.010641+05:30"
  },
  {...}
]
```

#### `GET` : `/historical-inr/:userId`

Return the INR value of the user’s stock rewards for all past days (up to yesterday)

**URL Parameter:**
* `userId` : UUID

**Response:**
```json
[
  {
      "Date": "2025-12-25",
      "Value": "13632.8466"
  },
  {
      "Date": "2025-12-26",
      "Value": "14617.2132"
  },
  {
      "Date": "2025-12-27",
      "Value": "15601.5798"
  }
]
```

### 4. Stats API

#### `GET` : `/stats/:userId`

Returns:
- Total shares rewarded today (grouped by stock symbol).  
- Current INR value of the user’s portfolio.

**URL Parameter:**
* `userId` : UUID

**Response:**
```json
{
  "today_rewards": [
    {
        "symbol": "INF",
        "total_quantity": "4"
    },
    {
        "symbol": "RIL",
        "total_quantity": "2.942"
    },
    {
        "symbol": "TAT",
        "total_quantity": "2.23"
    }
  ],
  "total_portfolio_value_inr": "29795.1588728"
}
```

#### `GET` : `/portfolio/:userId`

Show holdings per stock symbol with current INR value

**URL Parameter:**
* `userId` : UUID

**Response:**
```json
{
    "Portfolio": [
        {
            "symbol": "INF",
            "total_quantity": "56",
            "current_price": "460.9906",
            "current_value": "25815.4736"
        },
        {
            "symbol": "RIL",
            "total_quantity": "2.942",
            "current_price": "386.4284",
            "current_value": "1136.8723528"
        },
        {
            "symbol": "TAT",
            "total_quantity": "8.92",
            "current_price": "318.701",
            "current_value": "2842.81292"
        }
    ],
    "TotalValue": "29795.1588728"
}
```

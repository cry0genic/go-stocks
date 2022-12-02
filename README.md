# Go-Stocks

## Requirements

* Create a REST API where clients can retrieve stock data for FAANG companies
  * Optionally, clients can request historical data for one or more stocks
* Service should update stock prices every _N_ minutes where _N_ is configurable
* Service and supporting services (e.g., database, web layer, etc.) should be
  packages with Dockerfiles, and a docker-compose.yml file allowing full 
  application instantiation.
* Use git for version control and frequently commit changes.


## API Resources

The API exposes two endpoints: one for retrieving all stocks and another for
requesting quotes of a specific stock symbol. The API endpoints are versioned
with `v1`.

* GET /v1/stocks
* GET /v1/stock/[symbol]

All timestamps returned by the API are in UTC.

#### Last N Quotes

Each API endpoint allows for an optional parameter `last` that will direct the 
API to return the last _n_ quotes, or the maximum observed quotes, whichever
is less.

### GET /v1/stocks

Example: http://localhost:18081/v1/stocks

Response body:
```json
{
  "aapl": [
    {
      "price": 130.4,
      "symbol": "aapl",
      "time": "2021-05-07T19:31:07.000000272Z"
    }
  ],
  "amzn": [
    {
      "price": 3296.16,
      "symbol": "amzn",
      "time": "2021-05-07T19:31:07.000000623Z"
    }
  ],
  "fb": [
    {
      "price": 320.125,
      "symbol": "fb",
      "time": "2021-05-07T19:31:05.000000929Z"
    }
  ],
  "goog": [
    {
      "price": 2402.14,
      "symbol": "goog",
      "time": "2021-05-07T19:30:21.000000201Z"
    }
  ],
  "nflx": [
    {
      "price": 504.08,
      "symbol": "nflx",
      "time": "2021-05-07T19:31:00.000000053Z"
    }
  ]
}
```

### GET /v1/stock/goog

Example: http://localhost:18081/v1/stock/goog

Response body:
```json
[
  {
    "price": 2403.06,
    "symbol": "goog",
    "time": "2021-05-07T19:32:08.000000511Z"
  }
]
```

### GET /v1/stock/fb?last=3

Example: http://localhost:18081/v1/stock/goog?last=3

Response body:
```json
[
  {
    "price": 320.12,
    "symbol": "fb",
    "time": "2021-05-07T19:36:02.000000631Z"
  },
  {
    "price": 319.92,
    "symbol": "fb",
    "time": "2021-05-07T19:35:09.000000338Z"
  },
  {
    "price": 319.99,
    "symbol": "fb",
    "time": "2021-05-07T19:34:08.00000012Z"
  }
]
```
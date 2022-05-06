# Cryptocurrency backend

## As backend for bot devs:
For bot developers, to use this backend server:
1. Clone this at `main` branch for the most stability
2. `cd` to `crypto-backend` directory
3. Create an `.env` file and add the following fields
```dotenv
APIKEY="<Your Live Coin Watch API Key>"
MONGO_URI="<Your MongoDB URI>"
```
4. Execute the following command
```shell
go run .
```

# Bitcoin
The bitcoin app extends the eliona api to provide a few currency rates for bitcoin. It uses an external provider to retrieve the rates. Therefore, the available rates and updated times depend fully on the third-party. By default the external service is Coindesk.

## Configuration

The app needs environment variables and database tables for configuration.

### Eliona
To start and initialize an app in an Eliona environment, the app have to registered in Eliona. For this, an entry in the database table public.eliona_app is necessary.

To run bitcoin app within a mocked eliona envinroment you should build the proper image:
```shell
$ docker build -t eliona/bitcoin:v1 .
```

And, then add the following service to the docker compose file:
```
  bitcoin-app:
    container_name: bitcoin-app
    image: eliona/bitcoin:v1
    environment: 
      APPNAME: bitcoin
      API_TOKEN: secret
      BITCOIN_INTEGRATION_PORT: 3001
      API_ENDPOINT: "http://api-v2:3000/v2"
      CONNECTION_STRING: postgres://postgres:secret@database-mock:5432
      LOG_LEVEL: debug
      TZ: Europe/Zurich
    networks:
      api-v2-mock-network:
    restart: always
    ports:
      - "3001:3001"
    depends_on:
      - "database"
```

### Environment variables

#### APPNAME

The `APPNAME` MUST be set to `bitcoin`. Some resources use this name to identify the app inside an Eliona environment. For running as a Docker container inside an Eliona environment, the APPNAME have to set in the Dockerfile. If the app runs outside an Eliona environment the APPNAME must be set explicitly.

#### API_ENDPOINT

The `API_ENDPOINT` variable configures the [Eliona Api](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/api). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started.

#### CONNECTION_STRING

The `CONNECTION_STRING` variable configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started.

#### DEBUG_LEVEL (optional)

The `DEBUG_LEVEL` variable defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log).

### Database tables ###

The app requires some configuration data that remains in the database. To do this, the app creates its own database schema weather during initialization. The data in this schema should be made editable by eliona frontend. This allows the app to be configured by the user without direct database access.

A good practice is to initialize the app configuration with default values. This allows the user to see how what needs to be configured.

In detail, you need the following configuration data in table weather.configuration (name, value).

-- bitcoin.configuration (name, value)
('endpoint', 'https://api.coindesk.com/v1/bpi/currentprice.json') -- where is the API located

## API Reference

The bitcoin app grabs currency rates from Coindesk's web service and returns it to eliona api clients without saving them.

There a single endpoint to provide the currency rates:
```
# Request
GET http://localhost:3001/v2/bitcoin/rates
```
```json
# Response
[
    {
        "code": "USD",
        "description": "United States Dollar",
        "rate": 19118.069,
        "updated_time": "2022-09-26T20:44:00Z"
    },
    {
        "code": "GBP",
        "description": "British Pound Sterling",
        "rate": 15974.9055,
        "updated_time": "2022-09-26T20:44:00Z"
    },
    {
        "code": "EUR",
        "description": "Euro",
        "rate": 18623.7905,
        "updated_time": "2022-09-26T20:44:00Z"
    }
]
```
# onepixel (1px.li)

"onepixel" is an API-first URL shortener. 

#### Why this name? 
A URL shortener should have a short domain name, possibly 2 or 3 letter in length.   
"1 pixel" is the smallest unit of a screen.   
So `1px.li` stands for `one pixel links` - i.e. smallest possible links :) 


#### Status Badges 
[![codecov](https://codecov.io/gh/championswimmer/onepixel_backend/graph/badge.svg?token=DL3DR6CS40)](https://codecov.io/gh/championswimmer/onepixel_backend)
[![Build and test](https://github.com/championswimmer/onepixel_backend/actions/workflows/build_test.yaml/badge.svg)](https://github.com/championswimmer/onepixel_backend/actions/workflows/build_test.yaml)

## Hosted Instance (1px.li) 
The latest version of the code is automatically deployed via [Railway](https://railway.app?referralCode=T4g5xz)
to   
<big><https://onepixel.link></big>


You can deploy your own instance too
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/xAJ1-J?referralCode=T4g5xz)

## Databases 

There are two databases used in this project 

1. **Application DB**  
   This stores users, shortlinks and link groups 
2. **Events/Analytics DB**   
   This stores every redirection event used for analytics and stats 


#### Production Configuration

| Purpose        | Database                                 |
|----------------|------------------------------------------|
| Application DB | [PostgreSQL](https://www.postgresql.org) |
| Events DB      | [Clickhouse](https://clickhouse.com/)    |

#### Test Configuration 
For tests on CI it is better to use small embedded databases
instead of spinning up a full database server. 

| Purpose        | Database                                    |
|----------------|---------------------------------------------|
| Application DB | [SQLite](https://www.sqlite.org/index.html) |
| Events DB      | [DuckDB](https://duckdb.org/)               |

> [!NOTE]
> You can also override `.onepixel.local.env` to use
> the embedded databases for local development
> but they are not recommended for production use.

## Development 

### Deploy everything (with Docker)

Simplest way to get it running is 

1. make a `./data` directory where your database will be stored 
2. run `docker-compose up`

### Run with hot-reload for local development 

We will use docker to run an instance of our databases, but we will run the project using [air](https://github.com/cosmtrek/air) locally  

1. make the following directories where your database will be stored
   1. `./data` for postgres
   2. `./clk_data` for clickhouse
   3. `./clk_logs` for clickhouse logs
2. run `docker-compose up -d postgres clickhouse`
3. run `air` in the root directory of the project <sup>1</sup>


> Note[1]: you can also run `go run src/main.go` but it will not reload on changes

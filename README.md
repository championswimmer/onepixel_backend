# onepixel (1px.li)

"onepixel" is an API-first URL shortener. 

#### Why this name? 
A URL shortener should have a short domain name, possibly 2 or 3 letter in length.   
"1 pixel" is the smallest unit of a screen.   
So `1px.li` stands for `one pixel links` - i.e. smallest possible links :) 


#### Status Badges 
[![codecov](https://codecov.io/gh/championswimmer/onepixel_backend/graph/badge.svg?token=DL3DR6CS40)](https://codecov.io/gh/championswimmer/onepixel_backend)
[![Build and test](https://github.com/championswimmer/onepixel_backend/actions/workflows/build_test.yaml/badge.svg)](https://github.com/championswimmer/onepixel_backend/actions/workflows/build_test.yaml)

## Test Instance 
The latest version of the code is automatically deployed via [Railway](https://railway.app)
to   
<big><https://onepixel.link></big>


You can deploy your own instance too
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/xAJ1-J?referralCode=T4g5xz)

## Development 

### Deploy everything (with Docker)

Simplest way to get it running is 

1. make a `./data` directory where your database will be stored 
2. run `docker-compose up`

### Run with hot-reload for local development 

We will use docker to run an instance of database, but we will run the project using [air](https://github.com/cosmtrek/air) locally  

1. make a `./data` directory where your database will be stored
2. run `docker-compose up -d postgres`
3. add `127.0.0.1    postgres` to your `/etc/hosts` file <sup>1</sup>
4. run `air` in the root directory of the project <sup>2</sup>


> Note[1] This is because the server is configured to connect to `host=postgres` for the database.

> Note[2]: you can also run `go run src/main.go` but it will not reload on changes

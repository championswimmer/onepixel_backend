// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/championswimmer/onepixel_backend",
        "contact": {
            "name": "Arnav Gupta",
            "email": "dev@championswimmer.in"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/urls": {
            "get": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Get all urls",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Get all urls",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.UrlResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "something went wrong",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Create random short url",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Create random short url",
                "operationId": "create-random-url",
                "parameters": [
                    {
                        "description": "Url",
                        "name": "url",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.CreateUrlRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.UrlResponse"
                        }
                    },
                    "400": {
                        "description": "The request body is not valid",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "long_url is required to create url",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/urls/{shortcode}": {
            "get": {
                "description": "Get URL info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Get URL info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Shortcode",
                        "name": "shortcode",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UrlInfoResponse"
                        }
                    },
                    "404": {
                        "description": "URL not found",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Create specific short url",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "urls"
                ],
                "summary": "Create specific short url",
                "operationId": "create-specific-url",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Shortcode",
                        "name": "shortcode",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Url",
                        "name": "url",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.CreateUrlRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.UrlResponse"
                        }
                    },
                    "400": {
                        "description": "The request body is not valid",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Shortcode is not allowed",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Shortcode already exists",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "long_url is required to create url",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "security": [
                    {
                        "APIKeyAuth": []
                    }
                ],
                "description": "Register new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Register new user",
                "operationId": "register-user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserResponse"
                        }
                    },
                    "400": {
                        "description": "The request body is not valid",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "User with this email already exists",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "email and password are required to create user",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "Login user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Login user",
                "operationId": "login-user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.LoginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid email or password",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{userid}": {
            "get": {
                "description": "Get user info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user info",
                "operationId": "get-user-info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "patch": {
                "description": "Update user info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update user info",
                "operationId": "update-user-info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "dtos.CreateUrlRequest": {
            "type": "object",
            "properties": {
                "long_url": {
                    "type": "string"
                }
            }
        },
        "dtos.CreateUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dtos.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Something went wrong"
                },
                "status": {
                    "type": "integer",
                    "example": 400
                }
            }
        },
        "dtos.LoginUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dtos.UrlInfoResponse": {
            "type": "object",
            "properties": {
                "hit_count": {
                    "type": "integer"
                },
                "long_url": {
                    "type": "string"
                }
            }
        },
        "dtos.UrlResponse": {
            "type": "object",
            "properties": {
                "creator_id": {
                    "type": "integer",
                    "example": 1
                },
                "long_url": {
                    "type": "string",
                    "example": "https://www.google.com"
                },
                "short_url": {
                    "description": "Example url will pick up host and protocol(http/https) based on the env",
                    "type": "string",
                    "example": "https://1px.li/nhg145"
                }
            }
        },
        "dtos.UserResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@test.com"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "token": {
                    "type": "string",
                    "example": "\u003cJWT_TOKEN\u003e"
                }
            }
        }
    },
    "securityDefinitions": {
        "APIKeyAuth": {
            "type": "apiKey",
            "name": "X-API-Key",
            "in": "header"
        },
        "BearerToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "Operations about users",
            "name": "users"
        },
        {
            "description": "Operations about urls",
            "name": "urls"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "onepixel.link",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "onepixel API",
	Description:      "## About\n\n`1px.li` is an URL shortener created by [Arnav Gupta](https://twitter.com/championswimmer)\n\n- Source Code: <https://github.com/championswimmer/onepixel_backend> \n- Admin API: <https://onepixel.link/api/v1/>\n- API Docs: <https://onepixel.link/docs/> \n\n### Purchase Subscription\n\nUsing `onepixel` requires a subscription. You can purchase a subscription below using either Stripe or RazorPay\n\n[![Purchase](https://img.shields.io/badge/Purchase-slateblue?style=for-the-badge&logo=stripe&logoColor=white)](https://buy.stripe.com/bIY7tIfvucSv94A6oo)\n[![Purchase](https://img.shields.io/badge/Purchase-dodgerblue?style=for-the-badge&logo=razorpay&logoColor=white)](https://rzp.io/l/1pxli_1y)\n\nOnce you purchase a subscription, you'll receive an email with details of your account within 2-3 business days.\n",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

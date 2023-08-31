// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/segments": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "segment"
                ],
                "summary": "Create segment",
                "parameters": [
                    {
                        "description": "slug is a segment name, auto_add_percentage is a percentage of users who will have this segment",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.createSegmentBodyRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "segment"
                ],
                "summary": "Delete segment",
                "parameters": [
                    {
                        "description": "slug-segment name to delete",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.deleteSegmentBodyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Add and delete user segments by his id",
                "parameters": [
                    {
                        "description": "user segments to add and delete and his user id",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.addAndDeleteUserSegmentsBodyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/users/active_segments": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get active user segments",
                "parameters": [
                    {
                        "description": "user id to get his segments",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.getActiveUserSegmentsBodyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.getActiveUserSegmentsBodyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/users/report": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "operation"
                ],
                "summary": "Create a CSV file locally and return url to download a file",
                "parameters": [
                    {
                        "description": "date format year-month",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.createCSVRepostAndURLBodyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.createCSVRepostAndURLBodyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/users/report/{id}": {
            "get": {
                "tags": [
                    "operation"
                ],
                "summary": "Get report CSV file to download",
                "parameters": [
                    {
                        "type": "string",
                        "description": "report id",
                        "name": "input",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.addAndDeleteUserSegmentsBodyRequest": {
            "type": "object",
            "properties": {
                "segments_to_add": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "segments_to_delete": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "v1.createCSVRepostAndURLBodyRequest": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                }
            }
        },
        "v1.createCSVRepostAndURLBodyResponse": {
            "type": "object",
            "properties": {
                "report_url": {
                    "type": "string"
                }
            }
        },
        "v1.createSegmentBodyRequest": {
            "type": "object",
            "properties": {
                "auto_add_percentage": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                }
            }
        },
        "v1.deleteSegmentBodyRequest": {
            "type": "object",
            "properties": {
                "slug": {
                    "type": "string"
                }
            }
        },
        "v1.getActiveUserSegmentsBodyRequest": {
            "type": "object",
            "properties": {
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "v1.getActiveUserSegmentsBodyResponse": {
            "type": "object",
            "properties": {
                "segments": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "v1.response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Dynamic User Segmentation API",
	Description:      "Dynamic User Segmentation API for storing users and their segments",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
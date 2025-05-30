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
        "/auth/login": {
            "post": {
                "description": "Вход в систему с email и паролем",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "parameters": [
                    {
                        "description": "Данные для входа",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неверные учетные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка генерации токена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Создание нового аккаунта с подтверждением по email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "parameters": [
                    {
                        "description": "Данные для регистрации",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.RegisterResponse"
                        }
                    },
                    "400": {
                        "description": "Неверные данные/пользователь существует",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/verify": {
            "post": {
                "description": "Верификация email с кодом подтверждения",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "parameters": [
                    {
                        "description": "Данные для верификации",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.VerifyRequest"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/response.VerifyResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный код/сессия",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка генерации токена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/info": {
            "get": {
                "description": "Возвращает данные песни по названию и группе",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "parameters": [
                    {
                        "description": "Данные для поиска",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.GetSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.GetSongResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/info/group": {
            "get": {
                "description": "Возвращает список песен указанной группы с пагинацией",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "parameters": [
                    {
                        "description": "Название группы",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.GetSongsRequest"
                        }
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "description": "Лимит (по умолчанию 10)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "Смещение (по умолчанию 0)",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Song"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный формат параметров",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Группа не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/song": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Создает новую запись о песне. Требует авторизации.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "parameters": [
                    {
                        "description": "Данные песни",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AddSongRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Созданная песня",
                        "schema": {
                            "$ref": "#/definitions/model.Song"
                        }
                    },
                    "400": {
                        "description": "Неверные данные/ошибка авторизации",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/song/{id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Удаляет песню по ID. Требует прав владельца.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Удаленная песня",
                        "schema": {
                            "$ref": "#/definitions/model.Song"
                        }
                    },
                    "400": {
                        "description": "Неверный формат ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Ошибка авторизации",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Обновляет данные песни по её ID. Требует авторизации.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные для обновления",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно обновленная песня",
                        "schema": {
                            "$ref": "#/definitions/model.Song"
                        }
                    },
                    "400": {
                        "description": "Неверный формат ID/данных, ошибка авторизации",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Song": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "request.AddSongRequest": {
            "type": "object",
            "required": [
                "group",
                "link",
                "release_date",
                "song",
                "text"
            ],
            "properties": {
                "group": {
                    "description": "Название группы (обязательное)",
                    "type": "string",
                    "example": "The Beatles"
                },
                "link": {
                    "description": "Ссылка на песню (обязательное)",
                    "type": "string",
                    "example": "https://example.com/yesterday"
                },
                "release_date": {
                    "description": "Дата релиза в формате dd.mm.yyyy (обязательное)",
                    "type": "string",
                    "example": "06.08.1965"
                },
                "song": {
                    "description": "Название песни (обязательное)",
                    "type": "string",
                    "example": "Yesterday"
                },
                "text": {
                    "description": "Текст песни (обязательное)",
                    "type": "string",
                    "example": "Yesterday, all my troubles seemed so far away..."
                }
            }
        },
        "request.GetSongRequest": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "description": "Название группы (обязательное)",
                    "type": "string",
                    "example": "Queen"
                },
                "song": {
                    "description": "Название песни (обязательное)",
                    "type": "string",
                    "example": "Bohemian Rhapsody"
                }
            }
        },
        "request.GetSongsRequest": {
            "type": "object",
            "required": [
                "group"
            ],
            "properties": {
                "group": {
                    "description": "Название группы (обязательное)",
                    "type": "string",
                    "example": "Queen"
                }
            }
        },
        "request.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "qwerty123"
                }
            }
        },
        "request.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "password": {
                    "type": "string",
                    "example": "strong-password"
                }
            }
        },
        "request.UpdateSongRequest": {
            "type": "object",
            "properties": {
                "group": {
                    "description": "Новое название группы",
                    "type": "string",
                    "example": "The Beatles (Remastered)"
                },
                "link": {
                    "description": "Новая ссылка на песню",
                    "type": "string",
                    "example": "https://example.com/yesterday-remastered"
                },
                "release_date": {
                    "description": "Новая дата релиза в формате dd.mm.yyyy",
                    "type": "string",
                    "example": "01.01.2023"
                },
                "song": {
                    "description": "Новое название песни",
                    "type": "string",
                    "example": "Yesterday (Remastered)"
                },
                "text": {
                    "description": "Обновленный текст песни",
                    "type": "string",
                    "example": "Updated lyrics text..."
                }
            }
        },
        "request.VerifyRequest": {
            "type": "object",
            "required": [
                "code",
                "session_id"
            ],
            "properties": {
                "code": {
                    "type": "string",
                    "example": "1234"
                },
                "session_id": {
                    "type": "string",
                    "example": "UexEJzPJ3M"
                }
            }
        },
        "response.GetSongResponse": {
            "type": "object",
            "properties": {
                "release_date": {
                    "type": "string",
                    "example": "06.08.1965"
                },
                "text": {
                    "type": "string",
                    "example": "Yesterday, all my troubles seemed so far away..."
                }
            }
        },
        "response.LoginResponse": {
            "type": "object",
            "properties": {
                "jwt_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                }
            }
        },
        "response.RegisterResponse": {
            "type": "object",
            "properties": {
                "session_id": {
                    "type": "string",
                    "example": "UexEJzPJ3M"
                }
            }
        },
        "response.VerifyResponse": {
            "type": "object",
            "properties": {
                "jwt_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Music Library API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

// Code generated by go-swagger; DO NOT EDIT.

package spec

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "A payments API",
    "version": "1.0.0"
  },
  "paths": {
    "/payment/{paymentID}": {
      "get": {
        "summary": "Gets a payment by its ID",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the payment to retrieve",
            "name": "paymentID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "#/definitions/payment"
            }
          },
          "404": {
            "description": "Payment not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "put": {
        "summary": "Updates a payment by its ID",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the payment to update",
            "name": "paymentID",
            "in": "path",
            "required": true
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/payment"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "a payment",
            "schema": {
              "$ref": "#/definitions/empty"
            }
          },
          "404": {
            "description": "Payment not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "delete": {
        "summary": "Deletes a payment by its ID",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the payment to delete",
            "name": "paymentID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a payment",
            "schema": {
              "$ref": "#/definitions/empty"
            }
          },
          "404": {
            "description": "Payment not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/payments": {
      "get": {
        "summary": "Lists all the payments",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the organisation the payment belongs to.",
            "name": "organisation-id",
            "in": "query"
          },
          {
            "type": "number",
            "description": "Min amount of the payment",
            "name": "min-amount",
            "in": "query"
          },
          {
            "type": "number",
            "description": "Max amount of the payment",
            "name": "max-amount",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "#/definitions/payments"
            }
          }
        }
      },
      "post": {
        "produces": [
          "application/json"
        ],
        "summary": "Creates a new payment",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/payment"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/payment"
            }
          },
          "422": {
            "description": "Malformed request",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "empty": {
      "type": "object",
      "properties": {
        "_": {
          "type": "string"
        }
      }
    },
    "error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "payment": {
      "$ref": "payment.yml"
    },
    "payments": {
      "type": "object",
      "required": [
        "data",
        "links"
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/payment"
          }
        },
        "links": {
          "type": "object",
          "required": [
            "self"
          ],
          "properties": {
            "self": {
              "type": "string"
            }
          }
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "A payments API",
    "version": "1.0.0"
  },
  "paths": {
    "/payment/{paymentID}": {
      "get": {
        "summary": "Gets a payment by its ID",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the payment to retrieve",
            "name": "paymentID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "#/definitions/payment"
            }
          },
          "404": {
            "description": "Payment not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "put": {
        "summary": "Updates a payment by its ID",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the payment to update",
            "name": "paymentID",
            "in": "path",
            "required": true
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/payment"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "a payment",
            "schema": {
              "$ref": "#/definitions/empty"
            }
          },
          "404": {
            "description": "Payment not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "delete": {
        "summary": "Deletes a payment by its ID",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the payment to delete",
            "name": "paymentID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a payment",
            "schema": {
              "$ref": "#/definitions/empty"
            }
          },
          "404": {
            "description": "Payment not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/payments": {
      "get": {
        "summary": "Lists all the payments",
        "parameters": [
          {
            "type": "string",
            "description": "ID of the organisation the payment belongs to.",
            "name": "organisation-id",
            "in": "query"
          },
          {
            "type": "number",
            "description": "Min amount of the payment",
            "name": "min-amount",
            "in": "query"
          },
          {
            "type": "number",
            "description": "Max amount of the payment",
            "name": "max-amount",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "#/definitions/payments"
            }
          }
        }
      },
      "post": {
        "produces": [
          "application/json"
        ],
        "summary": "Creates a new payment",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/payment"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/payment"
            }
          },
          "422": {
            "description": "Malformed request",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "empty": {
      "type": "object",
      "properties": {
        "_": {
          "type": "string"
        }
      }
    },
    "error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "payment": {
      "type": "object",
      "required": [
        "type",
        "id",
        "organisation_id",
        "attributes"
      ],
      "properties": {
        "attributes": {
          "type": "object",
          "properties": {
            "amount": {
              "type": "string",
              "pattern": "^\\d+.\\d+$"
            },
            "beneficiary_party": {
              "type": "object",
              "required": [
                "account_number",
                "bank_id",
                "bank_id_code"
              ],
              "properties": {
                "account_name": {
                  "type": "string"
                },
                "account_number": {
                  "type": "string"
                },
                "account_number_code": {
                  "type": "string"
                },
                "account_type": {
                  "type": "integer"
                },
                "address": {
                  "type": "string"
                },
                "bank_id": {
                  "type": "string",
                  "pattern": "^\\d+$"
                },
                "bank_id_code": {
                  "type": "string"
                },
                "name": {
                  "type": "string"
                }
              }
            },
            "charges_information": {
              "type": "object",
              "properties": {
                "bearer_code": {
                  "type": "string"
                },
                "receiver_charges_amount": {
                  "type": "string",
                  "pattern": "^\\d+.\\d+$"
                },
                "receiver_charges_currency": {
                  "type": "string"
                },
                "sender_charges": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "amount": {
                        "type": "string",
                        "pattern": "^\\d+.\\d+$"
                      },
                      "currency": {
                        "type": "string"
                      }
                    }
                  }
                }
              }
            },
            "currency": {
              "type": "string"
            },
            "debtor_party": {
              "type": "object",
              "required": [
                "account_number",
                "bank_id",
                "bank_id_code"
              ],
              "properties": {
                "account_name": {
                  "type": "string"
                },
                "account_number": {
                  "type": "string"
                },
                "account_number_code": {
                  "type": "string"
                },
                "account_type": {
                  "type": "integer"
                },
                "address": {
                  "type": "string"
                },
                "bank_id": {
                  "type": "string",
                  "pattern": "^\\d+$"
                },
                "bank_id_code": {
                  "type": "string"
                },
                "name": {
                  "type": "string"
                }
              }
            },
            "end_to_end_reference": {
              "type": "string"
            },
            "fx": {
              "type": "object",
              "properties": {
                "contract_reference": {
                  "type": "string"
                },
                "exchange_rate": {
                  "type": "string",
                  "pattern": "^\\d+.\\d+$"
                },
                "original_amount": {
                  "type": "string",
                  "pattern": "^\\d+.\\d+$"
                },
                "original_currency": {
                  "type": "string"
                }
              }
            },
            "numeric_reference": {
              "type": "string",
              "pattern": "^\\d+$"
            },
            "payment_id": {
              "type": "string",
              "pattern": "^\\d+$"
            },
            "payment_purpose": {
              "type": "string"
            },
            "payment_scheme": {
              "type": "string"
            },
            "payment_type": {
              "type": "string"
            },
            "processing_date": {
              "type": "string",
              "format": "date"
            },
            "reference": {
              "type": "string"
            },
            "scheme_payment_sub_type": {
              "type": "string"
            },
            "sponsor_party": {
              "type": "object",
              "required": [
                "account_number",
                "bank_id",
                "bank_id_code"
              ],
              "properties": {
                "account_name": {
                  "type": "string"
                },
                "account_number": {
                  "type": "string"
                },
                "account_number_code": {
                  "type": "string"
                },
                "account_type": {
                  "type": "integer"
                },
                "address": {
                  "type": "string"
                },
                "bank_id": {
                  "type": "string",
                  "pattern": "^\\d+$"
                },
                "bank_id_code": {
                  "type": "string"
                },
                "name": {
                  "type": "string"
                }
              }
            }
          }
        },
        "id": {
          "type": "string",
          "format": "uuid",
          "readOnly": true
        },
        "organisation_id": {
          "type": "string",
          "format": "uuid",
          "readOnly": true
        },
        "type": {
          "type": "string",
          "minLength": 1,
          "readOnly": true
        },
        "version": {
          "type": "integer"
        }
      }
    },
    "payments": {
      "type": "object",
      "required": [
        "data",
        "links"
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/payment"
          }
        },
        "links": {
          "type": "object",
          "required": [
            "self"
          ],
          "properties": {
            "self": {
              "type": "string"
            }
          }
        }
      }
    }
  }
}`))
}

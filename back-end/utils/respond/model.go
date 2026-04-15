package respond

import "io"

type ApiModel[T any] struct {
	Data   T        `json:"data"`
	Errors []string `json:"errors" example:""`
}

type DataParam struct {
	Code     int
	Filename string
	MimeType string
	Data     []byte
}

// ############################
// ###### Swagger Example #####
// ############################

type File io.Writer

type EmptyBody struct{}

type HealthCheck struct {
	Data   string   `json:"data" example:"server running properly"`
	Errors []string `json:"errors"`
}

type Unauthorized struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Unauthorized!"`
}

type ValidationErr struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Name: This field is required,PhoneNumber: This field must be number"`
}

type CustomerNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected customer ID not found!"`
}

type AddressNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected address ID not found!"`
}

type ProductNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected product ID not found!"`
}

type ProductOtherNameNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected product other name ID not found!"`
}

type ProductVariantNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected product variant ID not found!"`
}

type ProductImageNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected product image ID not found!"`
}

type TransactionNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected transaction ID not found!"`
}

type ExpenseNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected expense ID not found!"`
}

type ExpenseEvidenceNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected expense evidence ID not found!"`
}

type FileNotFound struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"Selected file ID not found!"`
}

type ISE struct {
	Data   string   `json:"data" example:""`
	Errors []string `json:"errors" example:"database closed"`
}

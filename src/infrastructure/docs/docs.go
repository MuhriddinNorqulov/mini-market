package docs

import "github.com/swaggo/swag"

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "Mini Market API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  `{"swagger":"2.0","info":{"title":"{{.Title}}","version":"{{.Version}}"},"host":"{{.Host}}","basePath":"{{.BasePath}}","paths":{}}`,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

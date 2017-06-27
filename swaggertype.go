package swaggering

type (
	// c.f. https://github.com/OAI/OpenAPI-Specification/blob/master/versions/1.2.md

	// Swagger is the top level deserialization target for swagger.
	Swagger struct {
		name                   string
		BasePath, ResourcePath string
		Apis                   []Api
		Models                 map[string]*Model
	}

	// Api is the struct that is deserialized from api_*.json files.
	Api struct {
		Path, Description string
		Operations        []*Operation
	}

	// Operation represents an operation on an API in swagger
	Operation struct {
		Nickname, Method, Path, Deprecated string
		HasBody, DTORequest                bool
		Parameters                         []*Parameter
		ResponseMessages                   []*ResponseMessage
		Collection
	}

	// Model represents a Swagger model
	Model struct {
		Id, Description, Discriminator string
		GoPackage, GoName              string
		GoUses                         bool
		Required, SubTypes             []string
		Properties                     map[string]*Property
		Enums                          []Enum
	}

	// SwaggerType represents an abstract datatype as described by swagger.
	SwaggerType struct {
		Type, Format string
		Ref          string `json:"$ref"`
		Enum         []string
	}

	// Collection represents a list-type from swagger
	Collection struct {
		SwaggerType
		Items       SwaggerType
		UniqueItems bool
	}

	// Property represents a field in a swagger model
	Property struct {
		Collection
		SwaggerName string
	}

	// Parameter is a deserialization target for Swagger JSON files.
	Parameter struct {
		Collection
		ParamType, Name         string
		Required, AllowMultiple bool
	}

	// ResponseMessage is a deserialization target for Swagger JSON files.
	ResponseMessage struct {
		Code                   int
		Message, ResponseModel string
		model                  *Model
	}

	// Enum is a deserialization target for Swagger JSON files.
	Enum struct {
		Name   string
		Values []string
	}
)

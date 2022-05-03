package services

type ConsumerFormValidatorService interface {
	Validate(form map[string]interface{}) error
	FormSchema() map[string]interface{}
}

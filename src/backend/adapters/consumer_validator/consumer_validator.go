package consumer_validator

type ConsumerFormValidatorService struct{}

func NewConsumerFormValidatorService() ConsumerFormValidatorService {
	return ConsumerFormValidatorService{}
}

func (s *ConsumerFormValidatorService) Validate(form map[string]interface{}) error {
	return nil
}

func (s *ConsumerFormValidatorService) FormSchema() map[string]interface{} {
	return nil
}

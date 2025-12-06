package bootstrap

const (
	DatabaseAdapter = "database.adapter"
	SocketAdapter   = "socket.adapter"
	ConfigDefName   = "config.definition"

	ProductServiceName  = "product.service"
	CustomerServiceName = "customer.service"
	MerchantServiceName = "merchant.service"
	OrderService        = "order.service"
	AuthCustomerService = "auth_customer.service"
	AuthMerchantService = "auth_merchant.service"

	ProductHandlerName  = "product.handler"
	CustomerHandlerName = "customer.handler"

	ProductRepositoryName  = "product.repository"
	MerchantRepositoryName = "merchant.repository"

	GeminiLLMName  = "gemini.llm.package"
	KolosalLLMName = "kolosal.llm.package"
	OpenAILLMName  = "openai.llm.package"
	MistralLLMName = "mistal.llm.package"
)

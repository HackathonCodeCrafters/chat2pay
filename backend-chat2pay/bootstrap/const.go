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

	ProductHandlerName      = "product.handler"
	CustomerHandlerName     = "customer.handler"
	CustomerAuthHandlerName = "customer_auth.handler"

	AuthMiddlewareName = "auth.middleware"

	ProductRepositoryName      = "product.repository"
	MerchantRepositoryName     = "merchant.repository"
	CustomerRepositoryName     = "customer.repository"
	MerchantUserRepositoryName = "merchant_user.repository"

	MerchantAuthHandlerName = "merchant_auth.handler"
	ShippingHandlerName     = "shipping.handler"
	OrderHandlerName        = "order.handler"

	OrderServiceName        = "order.service"
	OrderRepositoryName     = "order.repository"
	ChatMessageRepositoryName = "chat_message.repository"
	ChatHandlerName         = "chat.handler"

	RajaOngkirName = "rajaongkir.package"

	GeminiLLMName  = "gemini.llm.package"
	KolosalLLMName = "kolosal.llm.package"
	OpenAILLMName  = "openai.llm.package"
	MistralLLMName = "mistal.llm.package"
)

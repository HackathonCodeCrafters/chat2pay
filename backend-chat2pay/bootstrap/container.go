package bootstrap

import (
	"chat2pay/config/yaml"
	"github.com/sarulabs/di/v2"
)

func NewContainer() (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	defs := []di.Def{
		{
			Name: ConfigDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return config.LoadConfig()
			},
		},
		{
			Name: XenditClientDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				xnd := client.New(cfg.XenditPrivateKey)
				return xnd, nil
			},
		},
		//{
		//	Name: UnitOfWorkDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		uow := repositories.NewUnitOfWork(ctn.Get(DBDefName).(*sqlx.DB))
		//		return uow, nil
		//	},
		//},
		//{
		//	Name: RollbarClientDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		cfg := ctn.Get(ConfigDefName).(config.Config)
		//		rb := rollbarKernel.NewRollbarLever(&cfg)
		//		return rb, nil
		//	},
		//},
		//{
		//	Name: RabbitMQDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		cfg := ctn.Get(ConfigDefName).(config.Config)
		//		log := ctn.Get(LoggerDefName).(logger.Logger)
		//		rb := amqp.NewMessageBroker(&cfg, log)
		//		return rb, nil
		//	},
		//},
		//{
		//	Name: ValidatorDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		return validator.New(), nil
		//	},
		//},
		//{
		//	Name: EchoDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		e := echo.New()
		//		validate := ctn.Get(ValidatorDefName).(*validator.Validate)
		//		e.Validator = &CustomValidator{validator: validate}
		//		return e, nil
		//	},
		//},

		// Handler

		//{
		//	Name: UserHandlerDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		userService := ctn.Get(UserServiceDefName).(userService.UserServiceInterface)
		//		return user.NewUserHandler(userService), nil
		//	},
		//},

		//{
		//	Name: MerchantHandlerDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		merchantService := ctn.Get(MerchantServiceDefName).(merchantService.MerchantServiceInterface)
		//		return merchant.NewMerchantHandler(merchantService), nil
		//	},
		//},

		{
			Name: ProductServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				transactionService := ctn.Get(TransactionServiceDefName).(transaction.TransactionServiceInterface)
				transactionRepo := ctn.Get(TransactionRepositoryDefName).(*repositories.TransactionRepository)
				uow := ctn.Get(UnitOfWorkDefName).(unitofwork.UnitOfWork)
				return transactionHandler.NewTransactionHandler(uow, transactionService, transactionRepo), nil
			},
		},

		// Middleware

		//{
		//	Name: AuthMiddlewareDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		cfg := ctn.Get(ConfigDefName).(config.Config)
		//		userService := ctn.Get(UserServiceDefName).(userService.UserServiceInterface)
		//		return middleware.AuthMiddleware(userService, cfg.JWTSecret), nil
		//	},
		//},
		//{
		//	Name: AdminAuthMiddlewareDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		return middleware.AdminMiddleware(), nil
		//	},
		//},
		//{
		//	Name: MerchantServiceDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		cfg := ctn.Get(ConfigDefName).(config.Config)
		//		db := ctn.Get(DBDefName).(*sql.DB)
		//		log := ctn.Get(LoggerDefName).(logger.Logger)
		//		redisClient := ctn.Get(RedisClientDefName).(*redis.RedisClient)
		//
		//		return merchantService.NewMerchantService(db, cfg, log, redisClient), nil
		//	},
		//},

		// Service

		//{
		//	Name: UserServiceDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		cfg := ctn.Get(ConfigDefName).(config.Config)
		//		db := ctn.Get(DBDefName).(*sql.DB) // langsung pakai *sql.DB
		//		log := ctn.Get(LoggerDefName).(logger.Logger)
		//		redisClient := ctn.Get(RedisClientDefName).(*redis.RedisClient)
		//
		//		return userService.NewUserService(db, cfg, cfg.JWTSecret, log, redisClient), nil
		//	},
		//},
		//
		//{
		//	Name: MerchantServiceDefName,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		cfg := ctn.Get(ConfigDefName).(config.Config)
		//		log := ctn.Get(LoggerDefName).(logger.Logger)
		//		xenditClient := ctn.Get(XenditClientDefName).(*client.API)
		//
		//		return transaction.NewChargeService(cfg, log, xenditClient), nil
		//	},
		//},

		{
			Name: ProductServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(yaml.Config)
				//log := ctn.Get(LoggerDefName).(logger.Logger)
				xenditClient := ctn.Get(XenditClientDefName).(*client.API)
				transactionRepo := ctn.Get(TransactionRepositoryDefName).(*repositories.TransactionRepository)
				transactionLogRepo := ctn.Get(TransactionLogRepositoryDefName).(*repositories.TransactionLogRepository)

				return transaction.NewTransactionService(cfg, rollbar, log, xenditClient, transactionRepo, transactionLogRepo), nil
			},
		},
	}

	if err := builder.Add(defs...); err != nil {
		return di.Container{}, err
	}

	if err := builder.Add(*NewRepository()...); err != nil {
		return di.Container{}, err
	}

	if err := builder.Add(*database.LoadDatabase()...); err != nil {
		return di.Container{}, err
	}

	return builder.Build(), nil
}

// CustomValidator adalah custom validator untuk Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

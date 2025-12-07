package service

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/pkg/llm"
	"chat2pay/internal/pkg/logger"
	"chat2pay/internal/pkg/redis"
	"chat2pay/internal/repositories"
	"context"
	"errors"
	"fmt"
)

type ProductService interface {
	Create(ctx context.Context, req *dto.ProductRequest) *presenter.Response
	CreateMultiple(ctx context.Context, req *[]dto.ProductRequest) *presenter.Response
	GetAll(ctx context.Context, merchantId string, page, limit int) *presenter.Response
	GetById(ctx context.Context, id string) *presenter.Response
	Update(ctx context.Context, id string, req *dto.ProductRequest) *presenter.Response
	Delete(ctx context.Context, id string) *presenter.Response
	AskProduct(ctx context.Context, req *dto.AskProduct) *presenter.Response
}

type productService struct {
	productRepo  repositories.ProductRepository
	merchantRepo repositories.MerchantRepository
	llm          llm.LLM
	redisClient  redis.RedisClient
	cfg          *yaml.Config
}

func NewProductService(
	productRepo repositories.ProductRepository,
	merchantRepo repositories.MerchantRepository,
	llm llm.LLM,
	redisClient redis.RedisClient,
	cfg *yaml.Config,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		merchantRepo: merchantRepo,
		llm:          llm,
		redisClient:  redisClient,
		cfg:          cfg,
	}
}

func (s *productService) Create(ctx context.Context, req *dto.ProductRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_create", s.cfg.Logger.Enable)
	)

	product := &entities.Product{
		MerchantID: req.MerchantID,
		OutletID:   req.OutletID,
		//CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: stringPtr(req.Description),
		SKU:         stringPtr(req.SKU),
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      "active",
	}

	if req.Status != "" {
		product.Status = req.Status
	}

	log.Info("creating product")
	created, err := s.productRepo.Create(ctx, product)
	if err != nil {
		log.Error(fmt.Sprintf("error creating product: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create product"))
	}

	//Embedding product
	emb, err := s.llm.EmbedQuery(ctx, formatProductForEmbedding(*product))

	if err != nil {
		log.Error(fmt.Sprintf("error creating product: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create product"))
	}

	err = s.productRepo.CreateProductEmbedding(ctx, &entities.ProductEmbedding{
		ProductId: created.ID,
		Content:   fmt.Sprintf(`%s - %s`, product.Name, *product.Description),
		Embedding: emb,
	})

	if err != nil {
		log.Error(fmt.Sprintf("error creating product: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create product"))
	}

	data := dto.ToProductResponse(created)
	return response.WithCode(201).WithData(data)
}

func (s *productService) CreateMultiple(ctx context.Context, req *[]dto.ProductRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_create", s.cfg.Logger.Enable)
	)

	for _, productPayload := range *req {
		product := &entities.Product{
			MerchantID: productPayload.MerchantID,
			OutletID:   productPayload.OutletID,
			//CategoryID:  productPayload.CategoryID,
			Name:        productPayload.Name,
			Description: stringPtr(productPayload.Description),
			SKU:         stringPtr(productPayload.SKU),
			Price:       productPayload.Price,
			Stock:       productPayload.Stock,
			Status:      "active",
		}

		log.Info("creating product")
		created, err := s.productRepo.Create(ctx, product)
		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed to create product"))
		}

		//Embedding product
		emb, err := s.llm.EmbedQuery(ctx, formatProductForEmbedding(*product))

		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed to create product"))
		}

		err = s.productRepo.CreateProductEmbedding(ctx, &entities.ProductEmbedding{
			ProductId: created.ID,
			Content:   fmt.Sprintf(`%s - %s`, product.Name, *product.Description),
			Embedding: emb,
		})

		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed to create product"))
		}

	}
	return response.WithCode(201).WithData("ok")
}

func (s *productService) AskProduct(ctx context.Context, req *dto.AskProduct) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_create", s.cfg.Logger.Enable)
	)

	classify, err := s.llm.ClassifyIntent(ctx, req.Prompt)
	if err != nil {
		return response.WithCode(500).WithError(errors.New("failed classify intent"))
	}

	switch classify {
	case "chit_chat":

		answer, err := s.llm.ChatWithHistory(ctx, req.Prompt)
		if err != nil {
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		data := dto.ToLLM(nil, answer)
		return response.WithCode(200).WithData(data)

	case "general_product_request":
		prompt := fmt.Sprintf("User message: '%s'. Ask a clarifying question.", req.Prompt)
		answer, err := s.llm.ChatWithHistory(ctx, prompt)
		if err != nil {
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		data := dto.ToLLM(nil, answer)
		return response.WithCode(200).WithData(data)
	case "specific_product_search":
		//Embedding product
		emb, err := s.llm.EmbedQuery(ctx, req.Prompt)

		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		embedding, err := s.productRepo.GetProductEmbeddingList(ctx, emb)

		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		productIds := []string{}
		for _, productEmbedding := range embedding {
			productIds = append(productIds, productEmbedding.ProductId)
		}

		product, err := s.productRepo.FindByIDs(ctx, productIds)

		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}
		data := dto.ToLLM(&product, "Berikut daftar produk yang mungkin relevan:")
		return response.WithCode(200).WithData(data)

	case "follow_up":

		lastMsg, err := s.llm.GetLastMessageContext(ctx)
		if err != nil {
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		contextualPrompt := fmt.Sprintf("Previous context: %s. User follow-up: %s", lastMsg, req.Prompt)

		answer, err := s.llm.ChatWithHistory(ctx, contextualPrompt)
		if err != nil {
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		data := dto.ToLLM(nil, answer)
		return response.WithCode(200).WithData(data)
	}

	return response.WithCode(200).WithData("ok")
}

func (s *productService) GetAll(ctx context.Context, merchantId string, page, limit int) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_getall", s.cfg.Logger.Enable)
	)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, err := s.productRepo.FindAll(ctx, merchantId, limit, offset)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching products: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to fetch products"))
	}

	total, err := s.productRepo.Count(ctx, merchantId)
	if err != nil {
		log.Error(fmt.Sprintf("error counting products: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to count products"))
	}

	data := dto.ToProductListResponse(products, total, page, limit)
	return response.WithCode(200).WithData(data)
}

func (s *productService) GetById(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_getbyid", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching product with id: %d", id))
	product, err := s.productRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching product: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if product == nil {
		log.Warn("product not found")
		return response.WithCode(404).WithError(errors.New("product not found"))
	}

	data := dto.ToProductResponse(product)
	return response.WithCode(200).WithData(data)
}

func (s *productService) Update(ctx context.Context, id string, req *dto.ProductRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_update", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching product with id: %d", id))
	product, err := s.productRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching product: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if product == nil {
		log.Warn("product not found")
		return response.WithCode(404).WithError(errors.New("product not found"))
	}

	product.Name = req.Name
	product.Description = stringPtr(req.Description)
	product.SKU = stringPtr(req.SKU)
	product.Price = req.Price
	product.Stock = req.Stock
	product.OutletID = req.OutletID
	product.CategoryID = req.CategoryID

	if req.Status != "" {
		product.Status = req.Status
	}

	log.Info("updating product")
	updated, err := s.productRepo.Update(ctx, product)
	if err != nil {
		log.Error(fmt.Sprintf("error updating product: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to update product"))
	}

	data := dto.ToProductResponse(updated)
	return response.WithCode(200).WithData(data)
}

func (s *productService) Delete(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("product_service_delete", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("checking if product exists: %d", id))
	product, err := s.productRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching product: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if product == nil {
		log.Warn("product not found")
		return response.WithCode(404).WithError(errors.New("product not found"))
	}

	log.Info("deleting product")
	err = s.productRepo.Delete(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error deleting product: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to delete product"))
	}

	return response.WithCode(200).WithData(map[string]string{"message": "product deleted successfully"})
}

func formatProductForEmbedding(p entities.Product) string {
	return fmt.Sprintf(`
Nama: %s
Deskripsi: %s
Harga: %.0f
Brand: %s
`,
		p.Name,
		ifnil(p.Description),
		p.Price,
		p.Name,
	)
}

func ifnil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

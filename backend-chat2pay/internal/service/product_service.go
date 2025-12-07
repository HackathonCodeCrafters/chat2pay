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
	"strings"
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
		Image:       stringPtr(req.Image),
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
		// Extract budget from prompt using LLM
		budgetPrompt := fmt.Sprintf(`Dari pesan user: "%s"

Ekstrak budget/harga maksimal yang disebutkan dalam Rupiah.
Jika ada angka seperti "15 juta", "15jt", "15.000.000", konversi ke angka.
Jika tidak ada budget disebutkan, output: 0

Output HANYA angka saja tanpa format (contoh: 15000000), tanpa penjelasan.`, req.Prompt)

		budgetStr, _ := s.llm.Chat(ctx, budgetPrompt)
		maxPrice := extractBudget(budgetStr)

		// Embedding product
		emb, err := s.llm.EmbedQuery(ctx, req.Prompt)
		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		// Use price filter if budget was detected
		var embedding []entities.ProductEmbedding
		if maxPrice > 0 {
			embedding, err = s.productRepo.GetProductEmbeddingListWithPrice(ctx, emb, maxPrice)
		} else {
			embedding, err = s.productRepo.GetProductEmbeddingList(ctx, emb)
		}

		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		productIds := []string{}
		for _, productEmbedding := range embedding {
			productIds = append(productIds, productEmbedding.ProductId)
		}

		products, err := s.productRepo.FindByIDs(ctx, productIds)
		if err != nil {
			log.Error(fmt.Sprintf("error creating product: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		// Generate reasoning/recommendation based on products found
		var message string
		if len(products) == 0 {
			if maxPrice > 0 {
				message = fmt.Sprintf("Maaf, saya tidak menemukan produk yang sesuai dengan budget Rp %s. Coba naikkan budget atau ubah kriteria pencarian.", formatPrice(maxPrice))
			} else {
				message = "Maaf, saya tidak menemukan produk yang sesuai dengan kriteria Anda. Coba dengan kata kunci lain."
			}
		} else {
			// Ask LLM to generate recommendation based on products and user query
			recommendationPrompt := fmt.Sprintf(`Berdasarkan permintaan user: "%s"

Saya menemukan %d produk yang relevan. Berikan rekomendasi singkat (2-3 kalimat) tentang:
1. Kenapa produk-produk ini cocok untuk user
2. Saran budget atau fitur yang perlu dipertimbangkan

Jawab dalam Bahasa Indonesia, singkat dan informatif.`, req.Prompt, len(products))

			recommendation, err := s.llm.Chat(ctx, recommendationPrompt)
			if err != nil {
				message = "Berikut produk yang saya temukan untuk Anda:"
			} else {
				message = recommendation
			}
		}

		// Save search to history for context
		s.llm.ChatWithHistory(ctx, fmt.Sprintf("User mencari: %s. Ditemukan %d produk.", req.Prompt, len(products)))

		data := dto.ToLLM(&products, message)
		return response.WithCode(200).WithData(data)

	case "product_question":
		// User is asking about a product that was already shown
		// Use chat history to answer questions about the product
		lastMsg, err := s.llm.GetLastMessageContext(ctx)
		if err != nil || lastMsg == "" {
			answer, _ := s.llm.ChatWithHistory(ctx, req.Prompt)
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		// Generate answer about the shown product
		questionPrompt := fmt.Sprintf(`Konteks percakapan sebelumnya (produk yang sudah ditampilkan): "%s"

Pertanyaan user: "%s"

Jawab pertanyaan user tentang produk tersebut dengan detail dan informatif.
Jika user bertanya "kenapa menyarankan ini", jelaskan alasan berdasarkan spesifikasi dan kecocokan dengan kebutuhan.
Jika user bertanya "speknya apa", jelaskan spesifikasi utama produk.
Jawab dalam Bahasa Indonesia, ramah dan membantu.`, lastMsg, req.Prompt)

		answer, err := s.llm.Chat(ctx, questionPrompt)
		if err != nil {
			answer, _ = s.llm.ChatWithHistory(ctx, req.Prompt)
		}

		// Save to history
		s.llm.ChatWithHistory(ctx, fmt.Sprintf("User bertanya: %s", req.Prompt))

		data := dto.ToLLM(nil, answer)
		return response.WithCode(200).WithData(data)

	case "product_clarification":
		// User is answering clarifying questions - search products with combined context
		lastMsg, err := s.llm.GetLastMessageContext(ctx)
		if err != nil || lastMsg == "" {
			// No context, treat as specific search
			answer, err := s.llm.ChatWithHistory(ctx, req.Prompt)
			if err != nil {
				return response.WithCode(500).WithError(errors.New("failed get product"))
			}
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		// Combine previous context with user's clarification to create search query
		combinePrompt := fmt.Sprintf(`Berdasarkan konteks percakapan sebelumnya: "%s"
Dan jawaban/preferensi user: "%s"

Buatkan query pencarian produk yang menggabungkan keduanya.
Contoh output: "laptop gaming budget 15 juta untuk sehari-hari"
Hanya output query saja, tanpa penjelasan.`, lastMsg, req.Prompt)

		searchQuery, err := s.llm.Chat(ctx, combinePrompt)
		if err != nil {
			searchQuery = fmt.Sprintf("%s %s", lastMsg, req.Prompt)
		}

		// Search products with combined query
		emb, err := s.llm.EmbedQuery(ctx, searchQuery)
		if err != nil {
			log.Error(fmt.Sprintf("error embedding clarification: %v", err))
			answer, _ := s.llm.ChatWithHistory(ctx, req.Prompt)
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		embedding, err := s.productRepo.GetProductEmbeddingList(ctx, emb)
		if err != nil {
			log.Error(fmt.Sprintf("error getting embeddings: %v", err))
			answer, _ := s.llm.ChatWithHistory(ctx, req.Prompt)
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		if len(embedding) == 0 {
			answer, _ := s.llm.ChatWithHistory(ctx, fmt.Sprintf("User mencari: %s. Tidak ada produk yang cocok, berikan saran alternatif.", searchQuery))
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		productIds := []string{}
		for _, productEmbedding := range embedding {
			productIds = append(productIds, productEmbedding.ProductId)
		}

		products, err := s.productRepo.FindByIDs(ctx, productIds)
		if err != nil {
			log.Error(fmt.Sprintf("error finding products: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		// Generate recommendation with context
		recommendationPrompt := fmt.Sprintf(`User awalnya bertanya: "%s"
Kemudian user memberikan preferensi: "%s"

Saya menemukan %d produk yang relevan. Berikan rekomendasi singkat (2-3 kalimat):
1. Kenapa produk ini cocok dengan kebutuhan dan budget user
2. Saran fitur yang perlu diperhatikan

Jawab dalam Bahasa Indonesia, ramah dan informatif.`, lastMsg, req.Prompt, len(products))

		recommendation, err := s.llm.Chat(ctx, recommendationPrompt)
		if err != nil {
			recommendation = "Berdasarkan preferensi Anda, berikut produk yang saya rekomendasikan:"
		}

		// Save to chat history
		s.llm.ChatWithHistory(ctx, fmt.Sprintf("User: %s. Saya menemukan %d produk yang cocok.", req.Prompt, len(products)))

		data := dto.ToLLM(&products, recommendation)
		return response.WithCode(200).WithData(data)

	case "follow_up":
		// User asks for alternatives or modifications
		lastMsg, err := s.llm.GetLastMessageContext(ctx)
		if err != nil || lastMsg == "" {
			answer, err := s.llm.ChatWithHistory(ctx, req.Prompt)
			if err != nil {
				return response.WithCode(500).WithError(errors.New("failed get product"))
			}
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		// Combine previous context with user's follow-up to search for alternatives
		combinePrompt := fmt.Sprintf(`Berdasarkan konteks percakapan sebelumnya: "%s"
Dan permintaan user: "%s"

Buatkan query pencarian produk alternatif.
Contoh: jika user bilang "yang lebih murah", buat query dengan harga lebih rendah.
Hanya output query saja, tanpa penjelasan.`, lastMsg, req.Prompt)

		searchQuery, err := s.llm.Chat(ctx, combinePrompt)
		if err != nil {
			searchQuery = fmt.Sprintf("%s %s", lastMsg, req.Prompt)
		}

		// Now search products with combined query
		emb, err := s.llm.EmbedQuery(ctx, searchQuery)
		if err != nil {
			log.Error(fmt.Sprintf("error embedding follow-up: %v", err))
			// Fallback to chat response
			answer, _ := s.llm.ChatWithHistory(ctx, req.Prompt)
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		embedding, err := s.productRepo.GetProductEmbeddingList(ctx, emb)
		if err != nil {
			log.Error(fmt.Sprintf("error getting embeddings: %v", err))
			answer, _ := s.llm.ChatWithHistory(ctx, req.Prompt)
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		if len(embedding) == 0 {
			// No products found, give helpful response
			answer, _ := s.llm.ChatWithHistory(ctx, fmt.Sprintf("User mencari: %s. Tidak ada produk yang cocok, berikan saran alternatif.", searchQuery))
			data := dto.ToLLM(nil, answer)
			return response.WithCode(200).WithData(data)
		}

		productIds := []string{}
		for _, productEmbedding := range embedding {
			productIds = append(productIds, productEmbedding.ProductId)
		}

		products, err := s.productRepo.FindByIDs(ctx, productIds)
		if err != nil {
			log.Error(fmt.Sprintf("error finding products: %v", err))
			return response.WithCode(500).WithError(errors.New("failed get product"))
		}

		// Generate recommendation with context
		recommendationPrompt := fmt.Sprintf(`User awalnya bertanya: "%s"
Kemudian user memberikan preferensi: "%s"

Saya menemukan %d produk yang relevan. Berikan rekomendasi singkat (2-3 kalimat):
1. Kenapa produk ini cocok dengan kebutuhan dan budget user
2. Saran fitur yang perlu diperhatikan

Jawab dalam Bahasa Indonesia, ramah dan informatif.`, lastMsg, req.Prompt, len(products))

		recommendation, err := s.llm.Chat(ctx, recommendationPrompt)
		if err != nil {
			recommendation = "Berdasarkan preferensi Anda, berikut produk yang saya rekomendasikan:"
		}

		// Save to chat history
		s.llm.ChatWithHistory(ctx, fmt.Sprintf("User: %s. Saya menemukan %d produk yang cocok.", req.Prompt, len(products)))

		data := dto.ToLLM(&products, recommendation)
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

func extractBudget(s string) float64 {
	// Clean the string and extract number
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "Rp", "")
	s = strings.ReplaceAll(s, "rp", "")
	s = strings.ReplaceAll(s, " ", "")

	// Try to parse as float
	var result float64
	fmt.Sscanf(s, "%f", &result)
	return result
}

func formatPrice(price float64) string {
	// Format price with thousand separator
	p := int64(price)
	if p >= 1000000 {
		return fmt.Sprintf("%d juta", p/1000000)
	}
	return fmt.Sprintf("%d", p)
}

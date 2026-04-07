package service

import (
    "errors"
    "fmt"
    "regexp"
    "strings"
    "unicode"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "classifieds/internal/models"
    "classifieds/internal/repository"
)

// Business errors for ads.
var (
    ErrAdNotFound   = errors.New("ad not found")
    ErrAdForbidden  = errors.New("unauthorized access to this ad")
    ErrShopAdsLimit = errors.New("ad limit reached for this plan")
)

// CreateAdInput groups the data required to create an ad.
type CreateAdInput struct {
    UserID     uint
    CategoryID uint
    ShopID     *uint
    Title      string
    Body       string
    Price      *float64
    Currency   string
    City       string
    District   *string
    Status     string
    Attributes []models.AdAttribute
}

// UpdateAdInput groups the editable fields of an ad.
type UpdateAdInput struct {
    Title      *string
    Body       *string
    Price      *float64
    Currency   *string
    City       *string
    District   *string
    Status     *string
    Attributes []models.AdAttribute
}

// AdService contains the business logic for ads.
type AdService struct {
    adRepo   *repository.AdRepo
    shopRepo *repository.ShopRepo
}

// NewAdService creates a new AdService.
func NewAdService(adRepo *repository.AdRepo, shopRepo *repository.ShopRepo) *AdService {
    return &AdService{adRepo: adRepo, shopRepo: shopRepo}
}

// CreateAd creates a new ad after checking permissions and plan limits.
func (s *AdService) CreateAd(input CreateAdInput) (*models.Ad, error) {
    // Check shop limits if applicable
    if input.ShopID != nil {
        if err := s.checkShopAdLimit(*input.ShopID); err != nil {
            return nil, err
        }
    }

    slug, err := s.generateUniqueSlug(input.Title)
    if err != nil {
        return nil, fmt.Errorf("slug generation: %w", err)
    }

    currency := input.Currency
    if currency == "" {
        currency = "MAD"
    }
    status := input.Status
    if status == "" {
        status = "active"
    }

    ad := &models.Ad{
        UserID:     input.UserID,
        CategoryID: input.CategoryID,
        ShopID:     input.ShopID,
        Slug:       slug,
        Title:      input.Title,
        Body:       input.Body,
        Price:      input.Price,
        Currency:   currency,
        City:       input.City,
        District:   input.District,
        Status:     status,
        Attributes: input.Attributes,
    }

    if err := s.adRepo.Create(ad); err != nil {
        return nil, fmt.Errorf("creation ad: %w", err)
    }

    // Recharger avec toutes les relations
    return s.adRepo.FindByID(ad.Model.ID)
}

// UpdateAd updates an existing ad (owner or admin only).
func (s *AdService) UpdateAd(slug string, requesterID uint, requesterRole string, input UpdateAdInput) (*models.Ad, error) {
    ad, err := s.adRepo.FindBySlug(slug)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrAdNotFound
        }
        return nil, err
    }

    if ad.UserID != requesterID && requesterRole != "admin" {
        return nil, ErrAdForbidden
    }

    // Appliquer les modifications (uniquement les champs fournis)
    if input.Title != nil {
        ad.Title = *input.Title
    }
    if input.Body != nil {
        ad.Body = *input.Body
    }
    if input.Price != nil {
        ad.Price = input.Price
    }
    if input.Currency != nil {
        ad.Currency = *input.Currency
    }
    if input.City != nil {
        ad.City = *input.City
    }
    if input.District != nil {
        ad.District = input.District
    }
    if input.Status != nil {
        ad.Status = *input.Status
    }

    if err := s.adRepo.Update(ad); err != nil {
        return nil, fmt.Errorf("update ad: %w", err)
    }

    // Update attributes if provided
    if input.Attributes != nil {
        if err := s.adRepo.UpdateAttributes(ad.Model.ID, input.Attributes); err != nil {
            return nil, fmt.Errorf("update attributs: %w", err)
        }
    }

    return s.adRepo.FindBySlug(slug)
}

// DeleteAd deletes an ad (owner or admin only).
func (s *AdService) DeleteAd(slug string, requesterID uint, requesterRole string) error {
    ad, err := s.adRepo.FindBySlug(slug)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrAdNotFound
        }
        return err
    }

    if ad.UserID != requesterID && requesterRole != "admin" {
        return ErrAdForbidden
    }

    return s.adRepo.Delete(ad.Model.ID)
}

// GetAd retrieves an ad by slug and increments the view counter.
func (s *AdService) GetAd(slug string) (*models.Ad, error) {
    ad, err := s.adRepo.FindBySlug(slug)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrAdNotFound
        }
        return nil, err
    }

    // Increment views non-blocking
    adID := ad.Model.ID
    go s.adRepo.IncrementViews(adID) //nolint:errcheck

    ad.ViewCount++
    return ad, nil
}

// ListAds returns filtered and paginated ads.
func (s *AdService) ListAds(f repository.AdFilters) (*repository.AdListResult, error) {
    return s.adRepo.List(f)
}

// GetUserAds returns a user's ads.
func (s *AdService) GetUserAds(userID uint, page, limit int) (*repository.AdListResult, error) {
    return s.adRepo.FindByUser(userID, page, limit)
}

// checkShopAdLimit verifies if the shop can still publish ads.
func (s *AdService) checkShopAdLimit(shopID uint) error {
    shop, err := s.shopRepo.FindByID(shopID)
    if err != nil {
        return fmt.Errorf("shop not found: %w", err)
    }

    if !shop.IsSubscriptionValid() {
        return ErrShopAdsLimit
    }

    // Fine limits are handled by ShopService
    return nil
}

// generateUniqueSlug creates a unique slug from the title.
func (s *AdService) generateUniqueSlug(title string) (string, error) {
    base := slugify(title)
    if base == "" {
        base = "ad"
    }

    // Short UUID suffix to guarantee uniqueness without a DB query
    suffix := uuid.NewString()[:8]
    return fmt.Sprintf("%s-%s", base, suffix), nil
}

// slugify converts text into a URL-safe slug (without external dependencies).
func slugify(s string) string {
    // Transliteration table for common characters (fr/ar → ASCII)
    replacer := strings.NewReplacer(
        "à", "a", "â", "a", "ä", "a", "á", "a", "ã", "a",
        "è", "e", "é", "e", "ê", "e", "ë", "e",
        "î", "i", "ï", "i", "í", "i", "ì", "i",
        "ô", "o", "ö", "o", "ó", "o", "ò", "o",
        "ù", "u", "û", "u", "ü", "u", "ú", "u",
        "ç", "c", "ñ", "n",
        "À", "a", "Â", "a", "Ä", "a", "Á", "a",
        "È", "e", "É", "e", "Ê", "e", "Ë", "e",
        "Î", "i", "Ï", "i",
        "Ô", "o", "Ö", "o",
        "Ù", "u", "Û", "u", "Ü", "u",
        "Ç", "c", "Ñ", "n",
        "&", "et", "+", "plus",
    )
    result := replacer.Replace(s)

    // Minuscules
    result = strings.ToLower(result)

    // Remove remaining non-ASCII characters (Arabic, etc.) and replace them with dashes
    var sb strings.Builder
    for _, r := range result {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            if r < 128 { // ASCII uniquement
                sb.WriteRune(r)
            }
        } else {
            sb.WriteRune('-')
        }
    }
    result = sb.String()

    // Collapseur de tirets multiples
    re := regexp.MustCompile(`-+`)
    result = re.ReplaceAllString(result, "-")

    // Remove leading and trailing dashes
    result = strings.Trim(result, "-")

    // Limit length to 80 characters
    if len(result) > 80 {
        result = result[:80]
        result = strings.TrimRight(result, "-")
    }

    return result
}

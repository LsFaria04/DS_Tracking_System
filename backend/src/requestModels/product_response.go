package requestModels

// Category represents a product category
type Category struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    ParentID  int    `json:"parent_id"`
    Permalink string `json:"permalink"`
}

// Image represents a product or variant image
type Image struct {
    ID       int    `json:"id"`
    Position int    `json:"position"`
    URL      string `json:"url"`
}

// VariantOption represents an option inside a product variant
type VariantOption struct {
    Name                 string `json:"name"`
    OptionType           string `json:"option_type"`
    Value                string `json:"value"`
    Custom               string `json:"custom"`
    ProductOptionPosition int    `json:"product_option_position"`
    ProductValuePosition  int    `json:"product_value_position"`
}

// Variant represents a product variant
type Variant struct {
    ID              int             `json:"id"`
    Price           float64         `json:"price"`
    SKU             string          `json:"sku"`
    Barcode         string          `json:"barcode"`
    Stock           int             `json:"stock"`
    StockUnlimited  bool            `json:"stock_unlimited"`
    StockThreshold  int             `json:"stock_threshold"`
    StockNotification bool          `json:"stock_notification"`
    CostPerItem     float64         `json:"cost_per_item"`
    CompareAtPrice  float64         `json:"compare_at_price"`
    Options         []VariantOption `json:"options"`
    Image           Image           `json:"image"`
}

// DigitalProduct represents a digital product attached to a product
type DigitalProduct struct {
    ID                int    `json:"id"`
    URL               string `json:"url"`
    ExpirationSeconds int    `json:"expiration_seconds"`
    External          bool   `json:"external"`
}

// Product represents the full product structure
type Product struct {
    ID                  int             `json:"id"`
    Name                string          `json:"name"`
    PageTitle           string          `json:"page_title"`
    Description         string          `json:"description"`
    Type                string          `json:"type"`
    DaysToExpire        int             `json:"days_to_expire"`
    Price               float64         `json:"price"`
    Discount            float64         `json:"discount"`
    Weight              float64         `json:"weight"`
    Stock               int             `json:"stock"`
    StockUnlimited      bool            `json:"stock_unlimited"`
    StockThreshold      int             `json:"stock_threshold"`
    StockNotification   bool            `json:"stock_notification"`
    BackInStockEnabled  bool            `json:"back_in_stock_enabled"`
    CostPerItem         float64         `json:"cost_per_item"`
    CompareAtPrice      float64         `json:"compare_at_price"`
    SKU                 string          `json:"sku"`
    Brand               string          `json:"brand"`
    Barcode             string          `json:"barcode"`
    GoogleProductCategory string        `json:"google_product_category"`
    Featured            bool            `json:"featured"`
    ReviewsEnabled      bool            `json:"reviews_enabled"`
    Status              string          `json:"status"`
    CreatedAt           string          `json:"created_at"`
    UpdatedAt           string          `json:"updated_at"`
    PackageFormat       string          `json:"package_format"`
    Length              float64         `json:"length"`
    Width               float64         `json:"width"`
}

// ProductWrapper matches the outer array structure [{ "product": { ... } }]
type ProductResponse struct {
    Product Product `json:"product"`
}
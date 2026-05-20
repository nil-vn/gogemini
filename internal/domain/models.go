package domain

// User maps to users table.
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role,omitempty"`
	Email        string `json:"email,omitempty"`
	PasswordHash string `json:"-"`
	Status       string `json:"status,omitempty"`
	CreatedDate  string `json:"created_date,omitempty"`
}

// Config maps to config table.
type Config struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// Car maps to car table.
type Car struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Branch            string     `json:"branch,omitempty"`
	Model             string     `json:"model,omitempty"`
	VIN               string     `json:"vin,omitempty"`
	Color             string     `json:"color,omitempty"`
	TradedCompany     string     `json:"traded_company,omitempty"`
	ImportedDate      string     `json:"imported_date,omitempty"`
	InspectionFrom    string     `json:"inspection_from,omitempty"`
	InspectionTo      string     `json:"inspection_to,omitempty"`
	YearOfManufacture string     `json:"year_of_manufacture,omitempty"`
	PurchasePrice     int64      `json:"purchase_price,omitempty"`
	SellingPrice      int64      `json:"selling_price,omitempty"`
	Status            string     `json:"status,omitempty"`
	Note              string     `json:"note,omitempty"`
	LicensePlateNo    string     `json:"license_plate_no,omitempty"`
	CreatedAt         string     `json:"created_at,omitempty"`
	Situation         string     `json:"car_situation,omitempty"`
	Images            []CarImage `json:"images,omitempty"`
}

// Customer maps to customer table.
type Customer struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Gender     string          `json:"gender,omitempty"`
	BirthDay   string          `json:"birth_day,omitempty"`
	Facebook   string          `json:"facebook,omitempty"`
	Phone      string          `json:"phone,omitempty"`
	Address    string          `json:"address,omitempty"`
	LicenseImg string          `json:"license_img,omitempty"`
	GalleryID  int64           `json:"gallery_id,omitempty"`
	LeadSource string          `json:"lead_source,omitempty"`
	Status     string          `json:"status,omitempty"`
	Note       string          `json:"note,omitempty"`
	CreatedAt  string          `json:"created_at,omitempty"`
	Images     []CustomerImage `json:"images,omitempty"`
}

// Transaction maps to transaction table.
type Transaction struct {
	ID            int64             `json:"id"`
	PurchaseDate  string            `json:"purchase_date,omitempty"`
	SellingPrice  int64             `json:"selling_price,omitempty"`
	Status        string            `json:"status,omitempty"`
	Note          string            `json:"note,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	CustomerID    int64             `json:"customer_id,omitempty"`
	DepositAmount int64             `json:"deposit_amount,omitempty"`
	Cars          []Car             `json:"cars,omitempty"`
	Items         []TransactionItem `json:"items,omitempty"`
}

// TransactionItem maps to transaction_item table.
type TransactionItem struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	Name          string `json:"name"`
	Price         int64  `json:"price,omitempty"`
}

// TransactionCar maps to transaction_car pivot table.
type TransactionCar struct {
	TransactionID int64 `json:"transaction_id"`
	CarID         int64 `json:"car_id"`
}

// CarImage maps to car_image table.
type CarImage struct {
	ID        int64  `json:"id"`
	CarID     int64  `json:"car_id"`
	FilePath  string `json:"file_path"`
	CreatedAt string `json:"created_at,omitempty"`
}

// CustomerImage maps to customer_image table.
type CustomerImage struct {
	ID         int64  `json:"id"`
	CustomerID int64  `json:"customer_id"`
	FilePath   string `json:"file_path"`
	CreatedAt  string `json:"created_at,omitempty"`
}

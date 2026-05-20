package domain

type User struct { ID int64 `json:"id"`; Username string `json:"username"`; Role string `json:"role,omitempty"`; Email string `json:"email,omitempty"`; Status string `json:"status,omitempty"` }
type Car struct { ID int64 `json:"id"`; Name string `json:"name"`; Branch string `json:"branch,omitempty"`; Model string `json:"model,omitempty"`; VIN string `json:"vin,omitempty"`; Status string `json:"status,omitempty"`; Situation string `json:"car_situation,omitempty"`; SellingPrice int64 `json:"selling_price,omitempty"` }
type Customer struct { ID int64 `json:"id"`; Name string `json:"name"`; Phone string `json:"phone,omitempty"`; Address string `json:"address,omitempty"`; Status string `json:"status,omitempty"` }
type Transaction struct { ID int64 `json:"id"`; CustomerID int64 `json:"customer_id,omitempty"`; Status string `json:"status,omitempty"`; SellingPrice int64 `json:"selling_price,omitempty"`; PurchaseDate string `json:"purchase_date,omitempty"` }

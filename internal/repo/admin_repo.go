package repo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gogemini/internal/domain"
)

type AdminRepo struct{ DB *sql.DB }

type ListQuery struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Search   string
	Status   string
}

type ListResult[T any] struct {
	Items []T
	Total int64
}

func normalizeListQuery(q ListQuery, allowedSort map[string]string) (page, pageSize int, sortCol, order string) {
	page = q.Page
	if page < 1 {
		page = 1
	}
	pageSize = q.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	sortCol = allowedSort["id"]
	if v, ok := allowedSort[strings.ToLower(strings.TrimSpace(q.Sort))]; ok {
		sortCol = v
	}
	order = "DESC"
	if strings.EqualFold(q.Order, "asc") {
		order = "ASC"
	}
	return
}

func (r AdminRepo) FindUserByUsernameOrEmail(login string) (domain.User, string, error) {
	var u domain.User
	var hash string
	err := r.DB.QueryRow(`SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,''), password_hash FROM users WHERE username = ? OR email = ? LIMIT 1`, login, login).
		Scan(&u.ID, &u.Username, &u.Role, &u.Email, &u.Status, &hash)
	return u, hash, err
}

func queryList[T any](db *sql.DB, q string, args []any, scan func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []T{}
	for rows.Next() {
		v, e := scan(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r AdminRepo) ListUsers(q ListQuery) (ListResult[domain.User], error) {
	page, pageSize, sortCol, order := normalizeListQuery(q, map[string]string{"id": "id", "username": "username", "email": "email", "status": "status"})
	needle := "%" + strings.TrimSpace(q.Search) + "%"
	args := []any{needle, needle, needle}
	countQ := `SELECT COUNT(*) FROM users WHERE (? = '%%' OR username LIKE ? OR email LIKE ?) AND (? = '' OR status = ?)`
	var total int64
	countArgs := append(args, strings.TrimSpace(q.Status), strings.TrimSpace(q.Status))
	if err := r.DB.QueryRow(countQ, countArgs...).Scan(&total); err != nil {
		return ListResult[domain.User]{}, err
	}
	listQ := fmt.Sprintf(`SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,'') FROM users WHERE (? = '%%' OR username LIKE ? OR email LIKE ?) AND (? = '' OR status = ?) ORDER BY %s %s LIMIT ? OFFSET ?`, sortCol, order)
	items, err := queryList(r.DB, listQ, append(countArgs, pageSize, (page-1)*pageSize), func(rows *sql.Rows) (domain.User, error) {
		var v domain.User
		return v, rows.Scan(&v.ID, &v.Username, &v.Role, &v.Email, &v.Status)
	})
	return ListResult[domain.User]{Items: items, Total: total}, err
}
func (r AdminRepo) ListCars(q ListQuery) (ListResult[domain.Car], error) {
	page, pageSize, sortCol, order := normalizeListQuery(q, map[string]string{"id": "id", "name": "name", "status": "status", "vin": "vin"})
	needle := "%" + strings.TrimSpace(q.Search) + "%"
	args := []any{needle, needle, needle, needle}
	countQ := `SELECT COUNT(*) FROM car WHERE (? = '%%' OR name LIKE ? OR vin LIKE ? OR branch LIKE ?) AND (? = '' OR status = ?)`
	var total int64
	countArgs := append(args, strings.TrimSpace(q.Status), strings.TrimSpace(q.Status))
	if err := r.DB.QueryRow(countQ, countArgs...).Scan(&total); err != nil {
		return ListResult[domain.Car]{}, err
	}
	listQ := fmt.Sprintf(`SELECT id, COALESCE(name,''), COALESCE(branch,''), COALESCE(model,''), COALESCE(vin,''), COALESCE(status,''), COALESCE(car_situation,''), COALESCE(selling_price,0) FROM car WHERE (? = '%%' OR name LIKE ? OR vin LIKE ? OR branch LIKE ?) AND (? = '' OR status = ?) ORDER BY %s %s LIMIT ? OFFSET ?`, sortCol, order)
	items, err := queryList(r.DB, listQ, append(countArgs, pageSize, (page-1)*pageSize), func(rows *sql.Rows) (domain.Car, error) {
		var v domain.Car
		return v, rows.Scan(&v.ID, &v.Name, &v.Branch, &v.Model, &v.VIN, &v.Status, &v.Situation, &v.SellingPrice)
	})
	return ListResult[domain.Car]{Items: items, Total: total}, err
}
func (r AdminRepo) ListCustomers(q ListQuery) (ListResult[domain.Customer], error) {
	page, pageSize, sortCol, order := normalizeListQuery(q, map[string]string{"id": "id", "name": "name", "phone": "phone", "status": "status"})
	needle := "%" + strings.TrimSpace(q.Search) + "%"
	args := []any{needle, needle, needle}
	countQ := `SELECT COUNT(*) FROM customer WHERE (? = '%%' OR name LIKE ? OR phone LIKE ?) AND (? = '' OR status = ?)`
	var total int64
	countArgs := append(args, strings.TrimSpace(q.Status), strings.TrimSpace(q.Status))
	if err := r.DB.QueryRow(countQ, countArgs...).Scan(&total); err != nil {
		return ListResult[domain.Customer]{}, err
	}
	listQ := fmt.Sprintf(`SELECT id, COALESCE(name,''), COALESCE(phone,''), COALESCE(address,''), COALESCE(status,'') FROM customer WHERE (? = '%%' OR name LIKE ? OR phone LIKE ?) AND (? = '' OR status = ?) ORDER BY %s %s LIMIT ? OFFSET ?`, sortCol, order)
	items, err := queryList(r.DB, listQ, append(countArgs, pageSize, (page-1)*pageSize), func(rows *sql.Rows) (domain.Customer, error) {
		var v domain.Customer
		return v, rows.Scan(&v.ID, &v.Name, &v.Phone, &v.Address, &v.Status)
	})
	return ListResult[domain.Customer]{Items: items, Total: total}, err
}
func (r AdminRepo) ListTransactions(q ListQuery) (ListResult[domain.Transaction], error) {
	page, pageSize, sortCol, order := normalizeListQuery(q, map[string]string{"id": "id", "status": "status", "purchase_date": "purchase_date", "selling_price": "selling_price"})
	needle := "%" + strings.TrimSpace(q.Search) + "%"
	args := []any{needle, needle}
	countQ := "SELECT COUNT(*) FROM `transaction` WHERE (? = '%%' OR status LIKE ?) AND (? = '' OR status = ?)"
	var total int64
	countArgs := append(args, strings.TrimSpace(q.Status), strings.TrimSpace(q.Status))
	if err := r.DB.QueryRow(countQ, countArgs...).Scan(&total); err != nil {
		return ListResult[domain.Transaction]{}, err
	}
	listQ := fmt.Sprintf("SELECT id, COALESCE(customer_id,0), COALESCE(status,''), COALESCE(selling_price,0), COALESCE(purchase_date,'') FROM `transaction` WHERE (? = '%%' OR status LIKE ?) AND (? = '' OR status = ?) ORDER BY %s %s LIMIT ? OFFSET ?", sortCol, order)
	items, err := queryList(r.DB, listQ, append(countArgs, pageSize, (page-1)*pageSize), func(rows *sql.Rows) (domain.Transaction, error) {
		var v domain.Transaction
		return v, rows.Scan(&v.ID, &v.CustomerID, &v.Status, &v.SellingPrice, &v.PurchaseDate)
	})
	return ListResult[domain.Transaction]{Items: items, Total: total}, err
}

func (r AdminRepo) GetUser(id int64) (domain.User, error) {
	var v domain.User
	err := r.DB.QueryRow(`SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,'') FROM users WHERE id=?`, id).Scan(&v.ID, &v.Username, &v.Role, &v.Email, &v.Status)
	return v, err
}
func (r AdminRepo) GetCar(id int64) (domain.Car, error) {
	var v domain.Car
	err := r.DB.QueryRow(`SELECT id, COALESCE(name,''), COALESCE(branch,''), COALESCE(model,''), COALESCE(vin,''), COALESCE(status,''), COALESCE(car_situation,''), COALESCE(selling_price,0) FROM car WHERE id=?`, id).Scan(&v.ID, &v.Name, &v.Branch, &v.Model, &v.VIN, &v.Status, &v.Situation, &v.SellingPrice)
	return v, err
}
func (r AdminRepo) GetCustomer(id int64) (domain.Customer, error) {
	var v domain.Customer
	err := r.DB.QueryRow(`SELECT id, COALESCE(name,''), COALESCE(phone,''), COALESCE(address,''), COALESCE(status,'') FROM customer WHERE id=?`, id).Scan(&v.ID, &v.Name, &v.Phone, &v.Address, &v.Status)
	return v, err
}
func (r AdminRepo) GetTransaction(id int64) (domain.Transaction, error) {
	var v domain.Transaction
	err := r.DB.QueryRow("SELECT id, COALESCE(customer_id,0), COALESCE(status,''), COALESCE(selling_price,0), COALESCE(purchase_date,'') FROM `transaction` WHERE id=?", id).Scan(&v.ID, &v.CustomerID, &v.Status, &v.SellingPrice, &v.PurchaseDate)
	return v, err
}

func (r AdminRepo) CreateUser(username, role, email, passwordHash, status string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO users (username, role, email, password_hash, status) VALUES (?, ?, ?, ?, ?)`, username, role, email, passwordHash, status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func (r AdminRepo) UpdateUser(id int64, username, role, email, status string) error {
	_, err := r.DB.Exec(`UPDATE users SET username=?, role=?, email=?, status=? WHERE id=?`, username, role, email, status, id)
	return err
}
func (r AdminRepo) UpdateCar(id int64, in domain.Car) error {
	_, err := r.DB.Exec(`UPDATE car SET name=?, branch=?, model=?, vin=?, status=?, car_situation=?, selling_price=? WHERE id=?`, in.Name, in.Branch, in.Model, in.VIN, in.Status, in.Situation, in.SellingPrice, id)
	return err
}
func (r AdminRepo) UpdateCustomer(id int64, in domain.Customer) error {
	_, err := r.DB.Exec(`UPDATE customer SET name=?, phone=?, address=?, status=? WHERE id=?`, in.Name, in.Phone, in.Address, in.Status, id)
	return err
}
func (r AdminRepo) UpdateTransaction(id int64, in domain.Transaction) error {
	_, err := r.DB.Exec("UPDATE `transaction` SET customer_id=?, status=?, selling_price=?, purchase_date=? WHERE id=?", in.CustomerID, in.Status, in.SellingPrice, in.PurchaseDate, id)
	return err
}
func (r AdminRepo) DeleteByID(table string, id int64) error {
	if table != "users" && table != "car" && table != "customer" && table != "transaction" {
		return fmt.Errorf("invalid table")
	}
	target := table
	if table == "transaction" {
		target = "`transaction`"
	}
	_, err := r.DB.Exec(`DELETE FROM `+target+` WHERE id=?`, id)
	return err
}
func (r AdminRepo) CreateCar(in domain.Car) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO car (name, branch, model, vin, status, car_situation, selling_price) VALUES (?, ?, ?, ?, ?, ?, ?)`, in.Name, in.Branch, in.Model, in.VIN, in.Status, in.Situation, in.SellingPrice)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func (r AdminRepo) CreateCustomer(in domain.Customer) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO customer (name, phone, address, status) VALUES (?, ?, ?, ?)`, in.Name, in.Phone, in.Address, in.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func (r AdminRepo) CreateTransaction(in domain.Transaction) (int64, error) {
	res, err := r.DB.Exec("INSERT INTO `transaction` (customer_id, status, selling_price, purchase_date) VALUES (?, ?, ?, ?)", in.CustomerID, in.Status, in.SellingPrice, in.PurchaseDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r AdminRepo) Dashboard() (map[string]int64, error) {
	counts := map[string]int64{}
	queries := map[string]string{"users": "SELECT COUNT(*) FROM users", "cars": "SELECT COUNT(*) FROM car", "customers": "SELECT COUNT(*) FROM customer", "transactions": "SELECT COUNT(*) FROM `transaction`"}
	for k, q := range queries {
		var c int64
		if err := r.DB.QueryRow(q).Scan(&c); err != nil {
			return nil, fmt.Errorf("%s: %w", k, err)
		}
		counts[k] = c
	}

	var totalRevenue int64
	if err := r.DB.QueryRow("SELECT COALESCE(SUM(selling_price),0) FROM `transaction` WHERE status IN ('PAID', 'DEPOSITED')").Scan(&totalRevenue); err != nil {
		return nil, fmt.Errorf("total_revenue: %w", err)
	}
	counts["total_revenue"] = totalRevenue

	now := time.Now().UTC()
	for i := 5; i >= 0; i-- {
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		nextMonthStart := monthStart.AddDate(0, 1, 0)
		label := monthStart.Format("01/2006")

		var monthRevenue int64
		if err := r.DB.QueryRow(
			"SELECT COALESCE(SUM(selling_price),0) FROM `transaction` WHERE created_at >= ? AND created_at < ?",
			monthStart,
			nextMonthStart,
		).Scan(&monthRevenue); err != nil {
			return nil, fmt.Errorf("revenue_6_months.%s: %w", label, err)
		}
		counts["revenue_"+label] = monthRevenue
	}

	return counts, nil
}
func (r AdminRepo) SearchAll(q string) (map[string]any, error) {
	query := strings.TrimSpace(q)
	if query == "" {
		return map[string]any{
			"users":        []domain.User{},
			"cars":         []domain.Car{},
			"customers":    []domain.Customer{},
			"transactions": []domain.Transaction{},
		}, nil
	}
	needle := "%" + query + "%"
	users, err := queryList(r.DB, `SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,'') FROM users WHERE username LIKE ? OR email LIKE ? OR status LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle, needle}, func(rows *sql.Rows) (domain.User, error) {
		var v domain.User
		return v, rows.Scan(&v.ID, &v.Username, &v.Role, &v.Email, &v.Status)
	})
	if err != nil {
		return nil, err
	}
	cars, err := queryList(r.DB, `SELECT id, COALESCE(name,''), COALESCE(branch,''), COALESCE(model,''), COALESCE(vin,''), COALESCE(status,''), COALESCE(car_situation,''), COALESCE(selling_price,0) FROM car WHERE name LIKE ? OR vin LIKE ? OR branch LIKE ? OR model LIKE ? OR status LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle, needle, needle, needle}, func(rows *sql.Rows) (domain.Car, error) {
		var v domain.Car
		return v, rows.Scan(&v.ID, &v.Name, &v.Branch, &v.Model, &v.VIN, &v.Status, &v.Situation, &v.SellingPrice)
	})
	if err != nil {
		return nil, err
	}
	customers, err := queryList(r.DB, `SELECT id, COALESCE(name,''), COALESCE(phone,''), COALESCE(address,''), COALESCE(status,'') FROM customer WHERE name LIKE ? OR phone LIKE ? OR address LIKE ? OR status LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle, needle, needle}, func(rows *sql.Rows) (domain.Customer, error) {
		var v domain.Customer
		return v, rows.Scan(&v.ID, &v.Name, &v.Phone, &v.Address, &v.Status)
	})
	if err != nil {
		return nil, err
	}
	transactions, err := queryList(r.DB, "SELECT id, COALESCE(customer_id,0), COALESCE(status,''), COALESCE(selling_price,0), COALESCE(purchase_date,'') FROM `transaction` WHERE status LIKE ? OR note LIKE ? ORDER BY id DESC LIMIT 20", []any{needle, needle}, func(rows *sql.Rows) (domain.Transaction, error) {
		var v domain.Transaction
		return v, rows.Scan(&v.ID, &v.CustomerID, &v.Status, &v.SellingPrice, &v.PurchaseDate)
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"users": users, "cars": cars, "customers": customers, "transactions": transactions}, nil
}

func (r AdminRepo) GetSettings() (map[string]string, error) {
	defaults := map[string]string{"currency": "JPY", "theme": "dark", "language": "vi"}
	rows, err := r.DB.Query(`SELECT key, COALESCE(value,'') FROM config WHERE key IN ('currency','theme','language')`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		defaults[k] = v
	}
	return defaults, rows.Err()
}

func (r AdminRepo) UpsertSettings(input map[string]string) error {
	for _, k := range []string{"currency", "theme", "language"} {
		v, ok := input[k]
		if !ok {
			continue
		}
		if _, err := r.DB.Exec(`INSERT INTO config (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, k, v); err != nil {
			return err
		}
	}
	return nil
}

package repo

import (
	"database/sql"
	"fmt"
	"strings"

	"gogemini/internal/domain"
)

type AdminRepo struct{ DB *sql.DB }

func (r AdminRepo) FindUserByUsernameOrEmail(login string) (domain.User, string, error) {
	var u domain.User
	var hash string
	err := r.DB.QueryRow(`SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,''), password_hash FROM users WHERE username = ? OR email = ? LIMIT 1`, login, login).
		Scan(&u.ID, &u.Username, &u.Role, &u.Email, &u.Status, &hash)
	return u, hash, err
}

func queryList[T any](db *sql.DB, q string, args []any, scan func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.Query(q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []T{}
	for rows.Next() { v, e := scan(rows); if e != nil { return nil, e }; out = append(out, v) }
	return out, rows.Err()
}

func (r AdminRepo) ListUsers() ([]domain.User, error) { return queryList(r.DB, `SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,'') FROM users ORDER BY id DESC`, nil, func(rows *sql.Rows) (domain.User, error) { var v domain.User; return v, rows.Scan(&v.ID, &v.Username, &v.Role, &v.Email, &v.Status) }) }
func (r AdminRepo) ListCars() ([]domain.Car, error) { return queryList(r.DB, `SELECT id, COALESCE(name,''), COALESCE(branch,''), COALESCE(model,''), COALESCE(vin,''), COALESCE(status,''), COALESCE(car_situation,''), COALESCE(selling_price,0) FROM car ORDER BY id DESC`, nil, func(rows *sql.Rows) (domain.Car, error) { var v domain.Car; return v, rows.Scan(&v.ID, &v.Name, &v.Branch, &v.Model, &v.VIN, &v.Status, &v.Situation, &v.SellingPrice) }) }
func (r AdminRepo) ListCustomers() ([]domain.Customer, error) { return queryList(r.DB, `SELECT id, COALESCE(name,''), COALESCE(phone,''), COALESCE(address,''), COALESCE(status,'') FROM customer ORDER BY id DESC`, nil, func(rows *sql.Rows) (domain.Customer, error) { var v domain.Customer; return v, rows.Scan(&v.ID, &v.Name, &v.Phone, &v.Address, &v.Status) }) }
func (r AdminRepo) ListTransactions() ([]domain.Transaction, error) { return queryList(r.DB, `SELECT id, COALESCE(customer_id,0), COALESCE(status,''), COALESCE(selling_price,0), COALESCE(purchase_date,'') FROM transaction ORDER BY id DESC`, nil, func(rows *sql.Rows) (domain.Transaction, error) { var v domain.Transaction; return v, rows.Scan(&v.ID, &v.CustomerID, &v.Status, &v.SellingPrice, &v.PurchaseDate) }) }

func (r AdminRepo) CreateUser(username, role, email, passwordHash, status string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO users (username, role, email, password_hash, status) VALUES (?, ?, ?, ?, ?)`, username, role, email, passwordHash, status)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (r AdminRepo) UpdateUser(id int64, username, role, email, status string) error {
	_, err := r.DB.Exec(`UPDATE users SET username=?, role=?, email=?, status=? WHERE id=?`, username, role, email, status, id)
	return err
}
func (r AdminRepo) DeleteByID(table string, id int64) error {
	if table != "users" && table != "car" && table != "customer" && table != "transaction" { return fmt.Errorf("invalid table") }
	_, err := r.DB.Exec(`DELETE FROM `+table+` WHERE id=?`, id)
	return err
}
func (r AdminRepo) CreateCar(in domain.Car) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO car (name, branch, model, vin, status, car_situation, selling_price) VALUES (?, ?, ?, ?, ?, ?, ?)`, in.Name, in.Branch, in.Model, in.VIN, in.Status, in.Situation, in.SellingPrice)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (r AdminRepo) CreateCustomer(in domain.Customer) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO customer (name, phone, address, status) VALUES (?, ?, ?, ?)`, in.Name, in.Phone, in.Address, in.Status)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (r AdminRepo) CreateTransaction(in domain.Transaction) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO transaction (customer_id, status, selling_price, purchase_date) VALUES (?, ?, ?, ?)`, in.CustomerID, in.Status, in.SellingPrice, in.PurchaseDate)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (r AdminRepo) Dashboard() (map[string]int64, error) {
	counts := map[string]int64{}
	queries := map[string]string{"users": "SELECT COUNT(*) FROM users", "cars": "SELECT COUNT(*) FROM car", "customers": "SELECT COUNT(*) FROM customer", "transactions": "SELECT COUNT(*) FROM transaction"}
	for k, q := range queries { var c int64; if err := r.DB.QueryRow(q).Scan(&c); err != nil { return nil, fmt.Errorf("%s: %w", k, err) }; counts[k] = c }
	return counts, nil
}
func (r AdminRepo) SearchAll(q string) (map[string]any, error) {
	needle := "%" + strings.TrimSpace(q) + "%"
	users, err := queryList(r.DB, `SELECT id, username, COALESCE(role,''), COALESCE(email,''), COALESCE(status,'') FROM users WHERE username LIKE ? OR email LIKE ? OR status LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle, needle}, func(rows *sql.Rows) (domain.User, error) { var v domain.User; return v, rows.Scan(&v.ID, &v.Username, &v.Role, &v.Email, &v.Status) })
	if err != nil { return nil, err }
	cars, err := queryList(r.DB, `SELECT id, COALESCE(name,''), COALESCE(branch,''), COALESCE(model,''), COALESCE(vin,''), COALESCE(status,''), COALESCE(car_situation,''), COALESCE(selling_price,0) FROM car WHERE name LIKE ? OR vin LIKE ? OR branch LIKE ? OR model LIKE ? OR status LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle, needle, needle, needle}, func(rows *sql.Rows) (domain.Car, error) { var v domain.Car; return v, rows.Scan(&v.ID, &v.Name, &v.Branch, &v.Model, &v.VIN, &v.Status, &v.Situation, &v.SellingPrice) })
	if err != nil { return nil, err }
	customers, err := queryList(r.DB, `SELECT id, COALESCE(name,''), COALESCE(phone,''), COALESCE(address,''), COALESCE(status,'') FROM customer WHERE name LIKE ? OR phone LIKE ? OR address LIKE ? OR status LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle, needle, needle}, func(rows *sql.Rows) (domain.Customer, error) { var v domain.Customer; return v, rows.Scan(&v.ID, &v.Name, &v.Phone, &v.Address, &v.Status) })
	if err != nil { return nil, err }
	transactions, err := queryList(r.DB, `SELECT id, COALESCE(customer_id,0), COALESCE(status,''), COALESCE(selling_price,0), COALESCE(purchase_date,'') FROM transaction WHERE status LIKE ? OR note LIKE ? ORDER BY id DESC LIMIT 20`, []any{needle, needle}, func(rows *sql.Rows) (domain.Transaction, error) { var v domain.Transaction; return v, rows.Scan(&v.ID, &v.CustomerID, &v.Status, &v.SellingPrice, &v.PurchaseDate) })
	if err != nil { return nil, err }
	return map[string]any{"users": users, "cars": cars, "customers": customers, "transactions": transactions}, nil
}


func (r AdminRepo) GetSettings() (map[string]string, error) {
	defaults := map[string]string{"currency": "JPY", "theme": "dark", "language": "vi"}
	rows, err := r.DB.Query(`SELECT key, COALESCE(value,'') FROM config WHERE key IN ('currency','theme','language')`)
	if err != nil { return nil, err }
	defer rows.Close()
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil { return nil, err }
		defaults[k] = v
	}
	return defaults, rows.Err()
}

func (r AdminRepo) UpsertSettings(input map[string]string) error {
	for _, k := range []string{"currency", "theme", "language"} {
		v, ok := input[k]
		if !ok { continue }
		if _, err := r.DB.Exec(`INSERT INTO config (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, k, v); err != nil { return err }
	}
	return nil
}

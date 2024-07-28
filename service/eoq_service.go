package service

import (
	"aplikasieoq/database"
	"aplikasieoq/models"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// service supplier
func CreateSupplier(ctx *gin.Context) (models.Supplier, error) {

	db := database.GetDB()

	var supplier models.Supplier
	if err := ctx.ShouldBindJSON(&supplier); err != nil {
		return supplier, err
	}

	// Set the CreatedAt and UpdatedAt fields
	now := time.Now()
	supplier.CreatedAt = now
	supplier.UpdatedAt = now

	// Query to insert supplier and return the id
	tsql := fmt.Sprintf(`
		INSERT INTO supplier (nama, perusahaan, kontak, alamat, created_at, updated_at) 
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s') RETURNING id`,
		supplier.Nama, supplier.Perusahaan, supplier.Kontak, supplier.Alamat, supplier.CreatedAt.Format(time.RFC3339), supplier.UpdatedAt.Format(time.RFC3339))

	// Execute query and get the returned ID
	var supplierID int
	err := db.Raw(tsql).Row().Scan(&supplierID)
	if err != nil {
		return supplier, err
	}

	supplier.ID = supplierID

	return supplier, nil

}

// GetSuppliers service to get all suppliers
func GetSuppliers(ctx *gin.Context) ([]models.Supplier, error) {
	db := database.GetDB()
	var suppliers []models.Supplier

	// Query to get all suppliers
	tsql := `SELECT id, nama, perusahaan, alamat, kontak, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM supplier`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var supplier models.Supplier
		if err := rows.Scan(&supplier.ID, &supplier.Nama, &supplier.Perusahaan, &supplier.Alamat, &supplier.Kontak, &supplier.CreatedAt, &supplier.UpdatedAt); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers, nil
}

// GetSupplierByID service to get a supplier by ID
func GetSupplierByID(id int) (models.Supplier, error) {
	db := database.GetDB()

	var supplier models.Supplier

	// Query to get supplier by ID
	tsql := `SELECT id, nama, perusahaan, kontak, alamat, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM supplier WHERE id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(&supplier.ID, &supplier.Nama, &supplier.Perusahaan, &supplier.Kontak, &supplier.Alamat, &supplier.CreatedAt, &supplier.UpdatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return supplier, gorm.ErrRecordNotFound
		}
		return supplier, err
	}

	return supplier, nil
}

// UpdateSupplier service to update supplier with optional fields
func UpdateSupplier(updatedSupplier models.Supplier) (models.Supplier, error) {
	db := database.GetDB()

	// Remove ID from the map to prevent updating it
	updatedFields := map[string]interface{}{}
	if updatedSupplier.Nama != "" {
		updatedFields["nama"] = updatedSupplier.Nama
	}
	if updatedSupplier.Perusahaan != "" {
		updatedFields["perusahaan"] = updatedSupplier.Perusahaan
	}
	if updatedSupplier.Kontak != "" {
		updatedFields["kontak"] = updatedSupplier.Kontak
	}
	if updatedSupplier.Alamat != "" {
		updatedFields["alamat"] = updatedSupplier.Alamat
	}
	updatedFields["updated_at"] = time.Now()

	if len(updatedFields) == 0 {
		return updatedSupplier, errors.New("no fields to update")
	}

	// Query to update supplier by ID
	setClause := ""
	args := []interface{}{}
	for field, value := range updatedFields {
		setClause += fmt.Sprintf("%s = ?, ", field)
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-2] // remove the last comma and space
	args = append(args, updatedSupplier.ID)

	tsql := fmt.Sprintf("UPDATE supplier SET %s WHERE id = ?", setClause)

	// Execute query
	result := db.Exec(tsql, args...)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return updatedSupplier, gorm.ErrRecordNotFound
		}
		return updatedSupplier, result.Error
	}

	// Retrieve the updated supplier
	row := db.Raw("SELECT id, nama, perusahaan, kontak, alamat, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM supplier WHERE id = ?", updatedSupplier.ID).Row()
	if err := row.Scan(&updatedSupplier.ID, &updatedSupplier.Nama, &updatedSupplier.Perusahaan, &updatedSupplier.Kontak, &updatedSupplier.Alamat, &updatedSupplier.CreatedAt, &updatedSupplier.UpdatedAt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return updatedSupplier, gorm.ErrRecordNotFound
		}
		return updatedSupplier, err
	}

	return updatedSupplier, nil
}

func DeleteSupplier(id int) error {
	db := database.GetDB()

	// Query to delete supplier by ID
	tsql := `DELETE FROM supplier WHERE id = ?`

	// Execute query
	result := db.Exec(tsql, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	}

	return nil
}

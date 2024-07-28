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

// service penyimpanan
func CreatePenyimpanan(ctx *gin.Context) (models.Penyimpanan, error) {

	db := database.GetDB()

	var penyimpanan models.Penyimpanan
	if err := ctx.ShouldBindJSON(&penyimpanan); err != nil {
		return penyimpanan, err
	}

	// Set the CreatedAt and UpdatedAt fields
	now := time.Now()
	penyimpanan.CreatedAt = now
	penyimpanan.UpdatedAt = now

	// Query to insert supplier and return the id
	tsql := fmt.Sprintf(`
		INSERT INTO penyimpanan (jenis, biaya_penyimpanan, tanggal_penyimpanan, created_at, updated_at) 
		VALUES ('%s', '%f', '%s', '%s', '%s') RETURNING id`,
		penyimpanan.Jenis, penyimpanan.BiayaPenyimpanan, penyimpanan.TanggalPenyimpanan.Format(time.RFC3339), penyimpanan.CreatedAt.Format(time.RFC3339), penyimpanan.UpdatedAt.Format(time.RFC3339))

	// Execute query and get the returned ID
	var penyimpananID int
	err := db.Raw(tsql).Row().Scan(&penyimpananID)
	if err != nil {
		return penyimpanan, err
	}

	penyimpanan.ID = penyimpananID

	return penyimpanan, nil

}

// GetPenyimpanans service to get all Penyimpanans
func GetPenyimpanans(ctx *gin.Context) ([]models.Penyimpanan, error) {
	db := database.GetDB()
	var penyimpanans []models.Penyimpanan

	// Query to get all penyimpanans
	tsql := `SELECT id, jenis, biaya_penyimpanan, tanggal_penyimpanan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM penyimpanan`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var penyimpanan models.Penyimpanan
		if err := rows.Scan(&penyimpanan.ID, &penyimpanan.Jenis, &penyimpanan.BiayaPenyimpanan, &penyimpanan.TanggalPenyimpanan, &penyimpanan.CreatedAt, &penyimpanan.UpdatedAt); err != nil {
			return nil, err
		}
		penyimpanans = append(penyimpanans, penyimpanan)
	}

	return penyimpanans, nil
}

// GetPenyimpananbyId service to get a penyimpanan by ID
func GetPenyimpananbyId(id int) (models.Penyimpanan, error) {
	db := database.GetDB()

	var penyimpanan models.Penyimpanan

	// Query to get penyimpanan by ID
	tsql := `SELECT id, jenis, biaya_penyimpanan, tanggal_penyimpanan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM penyimpanan WHERE id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(&penyimpanan.ID, &penyimpanan.Jenis, &penyimpanan.BiayaPenyimpanan, &penyimpanan.TanggalPenyimpanan, &penyimpanan.CreatedAt, &penyimpanan.UpdatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return penyimpanan, gorm.ErrRecordNotFound
		}
		return penyimpanan, err
	}

	return penyimpanan, nil
}

// UpdatePenyimpanan service to update penyimpanan with optional fields
func UpdatePenyimpanan(updatePenyimpanan models.Penyimpanan) (models.Penyimpanan, error) {
	db := database.GetDB()

	// Remove ID from the map to prevent updating it
	updatedFields := map[string]interface{}{}
	if updatePenyimpanan.Jenis != "" {
		updatedFields["jenis"] = updatePenyimpanan.Jenis
	}
	if updatePenyimpanan.BiayaPenyimpanan != 0 {
		updatedFields["biaya_penyimpanan"] = updatePenyimpanan.BiayaPenyimpanan
	}
	if !updatePenyimpanan.TanggalPenyimpanan.IsZero() {
		updatedFields["tanggal_penyimpanan"] = updatePenyimpanan.TanggalPenyimpanan
	}
	updatedFields["updated_at"] = time.Now()

	if len(updatedFields) == 0 {
		return updatePenyimpanan, errors.New("no fields to update")
	}

	// Query to update penyimpanan by ID
	setClause := ""
	args := []interface{}{}
	for field, value := range updatedFields {
		setClause += fmt.Sprintf("%s = ?, ", field)
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-2] // remove the last comma and space
	args = append(args, updatePenyimpanan.ID)

	tsql := fmt.Sprintf("UPDATE penyimpanan SET %s WHERE id = ?", setClause)

	// Execute query
	result := db.Exec(tsql, args...)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return updatePenyimpanan, gorm.ErrRecordNotFound
		}
		return updatePenyimpanan, result.Error
	}

	// Retrieve the updated penyimpanan
	row := db.Raw("SELECT id, jenis, biaya_penyimpanan, tanggal_penyimpanan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM penyimpanan WHERE id = ?", updatePenyimpanan.ID).Row()
	if err := row.Scan(&updatePenyimpanan.ID, &updatePenyimpanan.Jenis, &updatePenyimpanan.BiayaPenyimpanan, &updatePenyimpanan.TanggalPenyimpanan, &updatePenyimpanan.CreatedAt, &updatePenyimpanan.UpdatedAt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return updatePenyimpanan, gorm.ErrRecordNotFound
		}
		return updatePenyimpanan, err
	}

	return updatePenyimpanan, nil
}

func DeletePenyimpanan(id int) error {
	db := database.GetDB()

	// Query to delete penyimpanan by ID
	tsql := `DELETE FROM penyimpanan WHERE id = ?`

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

// service barang
func CreateBarang(ctx *gin.Context) (models.Barang, error) {

	db := database.GetDB()

	var barang models.Barang
	if err := ctx.ShouldBindJSON(&barang); err != nil {
		return barang, err
	}

	// Set the CreatedAt and UpdatedAt fields
	now := time.Now()
	barang.CreatedAt = now
	barang.UpdatedAt = now

	// Query to insert supplier and return the id
	tsql := fmt.Sprintf(`
		INSERT INTO barang (id_supplier, nama_barang, created_at, updated_at) 
		VALUES ('%d', '%s', '%s', '%s') RETURNING id`,
		barang.IDSupplier, barang.NamaBarang, barang.CreatedAt.Format(time.RFC3339), barang.UpdatedAt.Format(time.RFC3339))

	// Execute query and get the returned ID
	var barangID int
	err := db.Raw(tsql).Row().Scan(&barangID)
	if err != nil {
		return barang, err
	}

	barang.ID = barangID

	return barang, nil

}

// GetBarangs service to get all Barangs
func GetBarangs(ctx *gin.Context) ([]models.Barang, error) {
	db := database.GetDB()
	var barangs []models.Barang

	// Query to get all suppliers
	tsql := `SELECT id, id_supplier, nama_barang, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM barang`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var barang models.Barang
		if err := rows.Scan(&barang.ID, &barang.IDSupplier, &barang.NamaBarang, &barang.CreatedAt, &barang.UpdatedAt); err != nil {
			return nil, err
		}
		barangs = append(barangs, barang)
	}

	return barangs, nil
}

// GetBarangrByID service to get a barang by ID
func GetBarangByID(id int) (models.Barang, error) {
	db := database.GetDB()

	var barang models.Barang

	// Query to get barang by ID
	tsql := `SELECT id, id_supplier, nama_barang, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM barang WHERE id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(&barang.ID, &barang.IDSupplier, &barang.NamaBarang, &barang.CreatedAt, &barang.UpdatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return barang, gorm.ErrRecordNotFound
		}
		return barang, err
	}

	return barang, nil
}

// UpdateBarang service to update barang with optional fields
func UpdateBarang(updatedBarang models.Barang) (models.Barang, error) {
	db := database.GetDB()

	// Remove ID from the map to prevent updating it
	updatedFields := map[string]interface{}{}
	if updatedBarang.IDSupplier != 0 {
		updatedFields["id_supplier"] = updatedBarang.IDSupplier
	}
	if updatedBarang.NamaBarang != "" {
		updatedFields["nama_barang"] = updatedBarang.NamaBarang
	}
	updatedFields["updated_at"] = time.Now()

	if len(updatedFields) == 0 {
		return updatedBarang, errors.New("no fields to update")
	}

	// Query to update barang by ID
	setClause := ""
	args := []interface{}{}
	for field, value := range updatedFields {
		setClause += fmt.Sprintf("%s = ?, ", field)
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-2] // remove the last comma and space
	args = append(args, updatedBarang.ID)

	tsql := fmt.Sprintf("UPDATE barang SET %s WHERE id = ?", setClause)

	// Execute query
	result := db.Exec(tsql, args...)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return updatedBarang, gorm.ErrRecordNotFound
		}
		return updatedBarang, result.Error
	}

	// Retrieve the updated barang
	row := db.Raw("SELECT id, id_supplier, nama_barang, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM barang WHERE id = ?", updatedBarang.ID).Row()
	if err := row.Scan(&updatedBarang.ID, &updatedBarang.IDSupplier, &updatedBarang.NamaBarang, &updatedBarang.CreatedAt, &updatedBarang.UpdatedAt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return updatedBarang, gorm.ErrRecordNotFound
		}
		return updatedBarang, err
	}

	return updatedBarang, nil
}

func DeleteBarang(id int) error {
	db := database.GetDB()

	// Query to delete barang by ID
	tsql := `DELETE FROM barang WHERE id = ?`

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

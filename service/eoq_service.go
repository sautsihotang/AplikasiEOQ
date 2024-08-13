package service

import (
	"aplikasieoq/database"
	"aplikasieoq/models"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define a struct to hold the joined data from Barang and Supplier
type BarangWithSupplier struct {
	ID                 int       `json:"id"`
	IDSupplier         int       `json:"id_supplier"`
	NamaBarang         string    `json:"nama_barang"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	SupplierNama       string    `json:"supplier_nama"`
	SupplierPerusahaan string    `json:"supplier_perusahaan"`
	SupplierKontak     string    `json:"supplier_kontak"`
	SupplierAlamat     string    `json:"supplier_alamat"`
}

// GetBarangs service to get all Barangs with supplier details
func CalculateStock(ctx *gin.Context) ([]models.StockBarangModel, error) {
	db := database.GetDB()

	var reqStock models.ReqStock
	if err := ctx.ShouldBindJSON(&reqStock); err != nil {
		return nil, err
	}

	// Query to get all Barangs with supplier details
	tsql := `SELECT 
                b.id AS id_barang, 
                b.nama_barang AS nama_barang, 
                COALESCE(p.total_kuantitas_pemesanan, 0) AS total_kuantitas_pemesanan,
                COALESCE(j.total_kuantitas_penjualan, 0) AS total_kuantitas_penjualan,
                COALESCE(p.total_kuantitas_pemesanan, 0) - COALESCE(j.total_kuantitas_penjualan, 0) AS stok_barang
            FROM 
                barang b
            LEFT JOIN (
                SELECT
                    id_barang,
                    SUM(COALESCE(kuantitas, 0)) AS total_kuantitas_pemesanan
                FROM
                    pemesanan
                WHERE
                    tanggal_pemesanan BETWEEN ? AND ?
                GROUP BY
                    id_barang
            ) p ON b.id = p.id_barang
            LEFT JOIN (
                SELECT
                    id_barang,
                    SUM(COALESCE(kuantitas, 0)) AS total_kuantitas_penjualan
                FROM
                    penjualan
                WHERE
                    tanggal_penjualan BETWEEN ? AND ?
                GROUP BY
                    id_barang
            ) j ON b.id = j.id_barang
            LEFT JOIN
                supplier s ON b.id_supplier = s.id
            ORDER BY
                b.id;`

	// Execute query
	rows, err := db.Raw(tsql, reqStock.StartDate, reqStock.EndDate, reqStock.StartDate, reqStock.EndDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockBarangModel
	for rows.Next() {
		var stock models.StockBarangModel
		if err := rows.Scan(
			&stock.IDBarang,
			&stock.NamaBarang,
			&stock.TotalKuantitasPemesanan,
			&stock.TotalKuantitasPenjualan,
			&stock.StockBarang,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// CalculateEOQ menghitung EOQ dan frekuensi pemesanan
func CalculateEOQ(ctx *gin.Context) (models.Eoq, error) {
	db := database.GetDB()
	var eoq models.Eoq
	if err := ctx.ShouldBindJSON(&eoq); err != nil {
		return eoq, err
	}

	// Cetak payload JSON
	fmt.Printf("Payload received: %+v\n", eoq.IDBarang)
	fmt.Printf("Payload received: %+v\n", eoq.Periode)
	fmt.Printf("Payload received: %+v\n", eoq.TanggalPerhitungan)

	// Konversi eoq.Periode (string) ke int
	periodeStr := eoq.Periode
	periodeInt, err := strconv.Atoi(periodeStr)
	if err != nil {
		return eoq, fmt.Errorf("gagal mengonversi periode: %v", err)
	}

	// Mengambil frekuensi pemesanan per tahun
	frekuensiPemesananPerBarangPerTahun, err := TotalFrekuensiPemesananPerBarangPerTahun(db, eoq.IDBarang, periodeInt)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan frekuensi pemesanan per barang per tahun: %v", err)
	}

	frekuensiPemesananPerTahun, err := TotalFrekuensiPemesananPerTahun(db, periodeInt)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan frekuensi pemesanan per tahun: %v", err)
	}

	// D
	quantityBarangPerTahun, err := TotalQuantityBarangPerTahun(db, eoq.IDBarang, periodeInt)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan quantity barang per tahun: %v", err)
	}

	quantityPerTahun, err := TotalQuantityPerTahun(db, periodeInt)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan quantity per tahun: %v", err)
	}

	biayaPemesananPerTahun, err := TotalBiayaPemesananPerTahun(db, periodeInt)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan biaya pemesanan per tahun: %v", err)
	}

	biayaPenyimpananPerTahun, err := TotalBiayaPenyimpananPerTahun(db, periodeInt)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan biaya Penyimpanan per tahun: %v", err)
	}

	// Menampilkan frekuensi pemesanan per tahun
	fmt.Printf("Frekuensi pemesanan per barang per tahun: %d\n", frekuensiPemesananPerBarangPerTahun)
	fmt.Printf("Frekuensi pemesanan per tahun: %d\n", frekuensiPemesananPerTahun)
	fmt.Printf("Quantity barang per tahun: %d\n", quantityBarangPerTahun)
	fmt.Printf("Quantity per tahun: %d\n", quantityPerTahun)
	fmt.Printf("Biaya Pemesanan per tahun: %.f\n", biayaPemesananPerTahun)
	fmt.Printf("Biaya Penyimpanan per tahun: %.f\n", biayaPenyimpananPerTahun)

	// S
	s, err := BiayaPemesananSetiapKaliPesan(biayaPemesananPerTahun, frekuensiPemesananPerTahun)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan biaya pemesanan setiap kali pesan: %v", err)
	}

	// H

	h, err := BiayaPenyimpananPerBarang(quantityBarangPerTahun, quantityPerTahun, biayaPenyimpananPerTahun)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan biaya Penyimpanan per barang: %v", err)
	}
	fmt.Printf("biaya pemesanan setiap kali pesan: %d\n", s)
	fmt.Printf("biaya penyimpanan per barang: %.f\n", h)

	fmt.Printf("D %d\n", quantityBarangPerTahun)
	fmt.Printf("S %d\n", s)
	fmt.Printf("H %.f\n", h)

	nilaiEoq, err := calNilaiEoq(quantityBarangPerTahun, s, h)
	if err != nil {
		return eoq, fmt.Errorf("gagal mendapatkan nilai EOQ per barang: %v", err)
	}

	eoq.NilaiEOQ = math.Round(nilaiEoq)

	fmt.Printf("Nilai EOQ per barang: %.f\n", nilaiEoq)

	now := time.Now()
	eoq.CreatedAt = now
	eoq.UpdatedAt = now
	eoq.TanggalPerhitungan = now
	eoq.NilaiEOQ = nilaiEoq // bagaiaman cara agar data yg masuk ke db %.f

	tsql := fmt.Sprintf(`
		INSERT INTO eoq (id_barang, nilai_eoq, periode, tanggal_perhitungan, created_at, updated_at) 
		VALUES ('%d', '%f', '%s', '%s', '%s', '%s') RETURNING id`,
		eoq.IDBarang, eoq.NilaiEOQ, eoq.Periode, eoq.TanggalPerhitungan.Format(time.RFC3339), eoq.CreatedAt.Format(time.RFC3339), eoq.UpdatedAt.Format(time.RFC3339))

	var eoqID int
	err = db.Raw(tsql).Row().Scan(&eoqID)
	if err != nil {
		return eoq, err
	}

	eoq.ID = eoqID

	return eoq, nil
}

// calNilaiEoq menghitung Economic Order Quantity (EOQ)
func calNilaiEoq(d, s int, h float64) (float64, error) {
	if h <= 0 {
		return 0, fmt.Errorf("biaya penyimpanan per unit per tahun harus lebih besar dari nol")
	}
	if d <= 0 {
		return 0, fmt.Errorf("permintaan tahunan harus lebih besar dari nol")
	}

	// Konversi s ke float64 untuk operasi matematika
	sFloat := float64(s)
	dFloat := float64(d)

	// Hitung EOQ menggunakan rumus √(2DS/H)
	result := math.Sqrt((2 * sFloat * dFloat) / h)

	return result, nil
}

// BiayaPenyimpananPerBarang menghitung biaya penyimpanan per barang berdasarkan kuantitas dan total biaya penyimpanan
func BiayaPenyimpananPerBarang(qtyPerBarang, qtyPerTahun int, biayaPenyimpananPertahun float64) (float64, error) {
	if qtyPerTahun == 0 {
		return 0, fmt.Errorf("belum ada pemesanan barang") // Cek pembagi nol
	}

	// Hitung persentase kuantitas per barang terhadap total kuantitas
	persen := (float64(qtyPerBarang) / float64(qtyPerTahun)) * 100

	// Pembulatan persentase ke bilangan bulat terdekat
	persenBulat := math.Round(persen)

	// Hitung biaya penyimpanan per barang
	result := (persenBulat / 100) * biayaPenyimpananPertahun

	result = result / float64(qtyPerBarang)

	return result, nil
}

// BiayaPemesananSetiapKaliPesan menghitung biaya per pemesanan dan membulatkan hasilnya
func BiayaPemesananSetiapKaliPesan(biaya float64, frekuensi int) (int, error) {
	if frekuensi == 0 {
		return 0, fmt.Errorf("belum ada pemesanan barang")
	}

	// Hitung biaya per pemesanan
	s := biaya / float64(frekuensi)

	// Pembulatan ke bilangan bulat terdekat
	bulatan := int(math.Round(s))

	return bulatan, nil
}

func TotalBiayaPenyimpananPerTahun(db *gorm.DB, periode int) (float64, error) {
	var totalBiaya float64

	// Query SQL langsung
	tsql := `
		SELECT SUM(biaya_penyimpanan) AS total_biaya_penyimpanan
		FROM penyimpanan
		WHERE EXTRACT(YEAR FROM tanggal_penyimpanan) = ?
	`

	// Eksekusi query dan ambil hasilnya
	err := db.Raw(tsql, periode).Scan(&totalBiaya).Error
	if err != nil {
		return 0, err // Kembalikan error jika ada masalah dengan query
	}

	return totalBiaya, nil

}

func TotalBiayaPemesananPerTahun(db *gorm.DB, periode int) (float64, error) {
	var result struct {
		TotalBiayaTelepon      float64 `gorm:"column:total_biaya_telepon"`
		TotalBiayaAdm          float64 `gorm:"column:total_biaya_adm"`
		TotalBiayaTransportasi float64 `gorm:"column:total_biaya_transportasi"`
	}

	var biayaPemesanan float64
	// Query SQL langsung
	tsql := `
		SELECT 
			COALESCE(SUM(biaya_telepon), 0) AS total_biaya_telepon,
			COALESCE(SUM(biaya_adm), 0) AS total_biaya_adm,
			COALESCE(SUM(biaya_transportasi), 0) AS total_biaya_transportasi
		FROM pemesanan
		WHERE EXTRACT(YEAR FROM tanggal_pemesanan) = ?
	`

	// Eksekusi query dan ambil hasilnya
	err := db.Raw(tsql, periode).Scan(&result).Error
	if err != nil {
		return 0, err // Kembalikan error jika ada masalah dengan query
	}

	biayaPemesanan = result.TotalBiayaAdm + result.TotalBiayaTelepon + result.TotalBiayaTransportasi

	return biayaPemesanan, nil

}

// TotalQuantityBarangPerTahun menghitung quantity barang berdasarkan ID barang dan tahun
func TotalQuantityBarangPerTahun(db *gorm.DB, idBarang, periode int) (int, error) {
	var quantity int

	tsql := `
		SELECT SUM(kuantitas::integer) AS total_quantity
		FROM pemesanan
		WHERE id_barang = ?
	  	AND EXTRACT(YEAR FROM tanggal_pemesanan) = ?`
	err := db.Raw(tsql, idBarang, periode).Scan(&quantity).Error
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

// TotalQuantityPerTahun menghitung quantity barang berdasarkan ID barang dan tahun
func TotalQuantityPerTahun(db *gorm.DB, periode int) (int, error) {
	var quantity int

	tsql := `
			SELECT SUM(CAST(kuantitas AS INTEGER)) AS total_kuantitas
			FROM pemesanan
			WHERE EXTRACT(YEAR FROM tanggal_pemesanan) = ?`
	err := db.Raw(tsql, periode).Scan(&quantity).Error
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

// TotalFrekuensiPemesananPerBarangPerTahun menghitung frekuensi pemesanan berdasarkan ID barang dan tahun
func TotalFrekuensiPemesananPerBarangPerTahun(db *gorm.DB, idBarang, periode int) (int, error) {
	var frequency int

	// Query SQL langsung
	tsql := `
		SELECT COUNT(*) AS frequency
		FROM pemesanan
		WHERE id_barang = ? 
		AND EXTRACT(YEAR FROM tanggal_pemesanan) = ?
	`

	// Eksekusi query dan ambil hasilnya
	err := db.Raw(tsql, idBarang, periode).Scan(&frequency).Error
	if err != nil {
		return 0, err // Kembalikan error jika ada masalah dengan query
	}

	return frequency, nil
}

func TotalFrekuensiPemesananPerTahun(db *gorm.DB, periode int) (int, error) {
	var frequency int

	// Query SQL langsung
	tsql := `
		SELECT COUNT(*) AS frekuensi_pemesanan
		FROM pemesanan
		WHERE EXTRACT(YEAR FROM tanggal_pemesanan) = ?
	`

	// Eksekusi query dan ambil hasilnya
	err := db.Raw(tsql, periode).Scan(&frequency).Error
	if err != nil {
		return 0, err // Kembalikan error jika ada masalah dengan query
	}

	return frequency, nil
}

// GetEOQ service to get all perhitungan eoq
func GetEOQ(ctx *gin.Context) ([]models.EoqWithBarang, error) {
	db := database.GetDB()
	var eoqs []models.EoqWithBarang

	// Query to get all eoq
	tsql := `SELECT 
            e.id, 
            e.id_barang, 
			b.nama_barang,
            e.nilai_eoq, 
            e. periode, 
            e.tanggal_perhitungan, 
            COALESCE(e.created_at, NOW()) as created_at, 
            COALESCE(e.updated_at, NOW()) as updated_at 
         FROM 
            eoq  e
		INNER JOIN barang b on e.id_barang = b.id
         ORDER BY 
            created_at DESC`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eoq models.EoqWithBarang
		if err := rows.Scan(&eoq.ID, &eoq.IDBarang, &eoq.NamaBarang, &eoq.NilaiEOQ, &eoq.Periode, &eoq.TanggalPerhitungan, &eoq.CreatedAt, &eoq.UpdatedAt); err != nil {
			return nil, err
		}
		eoqs = append(eoqs, eoq)
	}

	return eoqs, nil
}

func DeleteEoq(id int) error {
	db := database.GetDB()

	// Query to delete eoq by ID
	tsql := `DELETE FROM eoq WHERE id = ?`

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

// service penjualan
func CreatePenjualan(ctx *gin.Context) (models.Penjualan, error) {

	db := database.GetDB()

	var penjualan models.Penjualan
	if err := ctx.ShouldBindJSON(&penjualan); err != nil {
		return penjualan, err
	}

	hargaSatuanFloat := float64(penjualan.HargaSatuan)
	totalHarga := float64(penjualan.Kuantitas) * hargaSatuanFloat
	totalHarga = float64(int(totalHarga))

	penjualan.TotalHarga = models.Float64OrString(totalHarga)

	// Set the CreatedAt and UpdatedAt fields
	now := time.Now()
	penjualan.CreatedAt = now
	penjualan.UpdatedAt = now

	// Query to insert supplier and return the id
	tsql := fmt.Sprintf(`
		INSERT INTO penjualan (id_user, id_barang, kuantitas, harga_satuan, total_harga, tanggal_penjualan, created_at, updated_at) 
		VALUES ('%d', '%d', '%d', '%f', '%f','%s', '%s', '%s') RETURNING id`,
		penjualan.IDUser, penjualan.IDBarang, penjualan.Kuantitas, penjualan.HargaSatuan, penjualan.TotalHarga, penjualan.TanggalPenjualan.Format(time.RFC3339), penjualan.CreatedAt.Format(time.RFC3339), penjualan.UpdatedAt.Format(time.RFC3339))

	// Execute query and get the returned ID
	var penjualanID int
	err := db.Raw(tsql).Row().Scan(&penjualanID)
	if err != nil {
		return penjualan, err
	}

	penjualan.ID = penjualanID

	return penjualan, nil

}

// GetPenjualans service to get all Penjualans
func GetPenjualans(ctx *gin.Context) ([]models.PenjualanWithBarangWithSupplier, error) {
	db := database.GetDB()
	var penjualans []models.PenjualanWithBarangWithSupplier

	// Query to get all penjualan with inner join and aliases
	tsql := `SELECT 
                 p.id AS id, 
                 p.id_user AS id_user, 
                 p.id_barang AS id_barang, 
                 b.nama_barang AS barang_nama, 
                 s.nama AS supplier_nama,
                 p.kuantitas AS kuantitas, 
                 p.harga_satuan AS harga_satuan, 
                 p.total_harga AS total_harga, 
                 p.tanggal_penjualan AS tanggal_penjualan, 
                 COALESCE(p.created_at, NOW()) AS created_at, 
                 COALESCE(p.updated_at, NOW()) AS updated_at
            FROM penjualan p
            INNER JOIN barang b ON p.id_barang = b.id
            INNER JOIN supplier s ON b.id_supplier = s.id`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var penjualan models.PenjualanWithBarangWithSupplier
		if err := rows.Scan(
			&penjualan.ID,
			&penjualan.IDUser,
			&penjualan.IDBarang,
			&penjualan.BarangNama,
			&penjualan.SupplierNama,
			&penjualan.Kuantitas,
			&penjualan.HargaSatuan,
			&penjualan.TotalHarga,
			&penjualan.TanggalPenjualan,
			&penjualan.CreatedAt,
			&penjualan.UpdatedAt,
		); err != nil {
			return nil, err
		}
		penjualans = append(penjualans, penjualan)
	}

	return penjualans, nil
}

// GetPenjualanbyId service to get a penjualan by ID
func GetPenjualanbyId(id int) (models.Penjualan, error) {
	db := database.GetDB()

	var penjualan models.Penjualan

	// Query to get penyimpanan by ID
	tsql := `SELECT id, id_user, id_barang, kuantitas, harga_satuan, total_harga, tanggal_penjualan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM penjualan WHERE id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(&penjualan.ID, &penjualan.IDUser, &penjualan.IDBarang, &penjualan.Kuantitas, &penjualan.HargaSatuan, &penjualan.TotalHarga, &penjualan.TanggalPenjualan, &penjualan.CreatedAt, &penjualan.UpdatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return penjualan, gorm.ErrRecordNotFound
		}
		return penjualan, err
	}

	return penjualan, nil
}

// UpdatePenjualan service to update penjualan with optional fields
func UpdatePenjualan(updatePenjualan models.Penjualan) (models.Penjualan, error) {
	db := database.GetDB()

	// Remove ID from the map to prevent updating it
	updatedFields := map[string]interface{}{}
	if updatePenjualan.IDUser != 0 {
		updatedFields["id_user"] = updatePenjualan.IDUser
	}
	if updatePenjualan.IDBarang != 0 {
		updatedFields["id_barang"] = updatePenjualan.IDBarang
	}
	if updatePenjualan.Kuantitas != 0 {
		updatedFields["kuantitas"] = updatePenjualan.Kuantitas
	}
	if updatePenjualan.HargaSatuan != 0 {
		updatedFields["harga_satuan"] = updatePenjualan.HargaSatuan
	}
	if updatePenjualan.Kuantitas != 0 || updatePenjualan.HargaSatuan != 0 {
		updatedFields["total_harga"] = updatePenjualan.Kuantitas * models.IntOrString(updatePenjualan.HargaSatuan)
	}
	// if updatePenjualan.TotalHarga != 0 {
	// 	updatedFields["total_harga"] = updatePenjualan.TotalHarga
	// }
	if !updatePenjualan.TanggalPenjualan.IsZero() {
		updatedFields["tanggal_penjualan"] = updatePenjualan.TanggalPenjualan
	}
	updatedFields["updated_at"] = time.Now()

	if len(updatedFields) == 0 {
		return updatePenjualan, errors.New("no fields to update")
	}

	// Query to update penyimpanan by ID
	setClause := ""
	args := []interface{}{}
	for field, value := range updatedFields {
		setClause += fmt.Sprintf("%s = ?, ", field)
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-2] // remove the last comma and space
	args = append(args, updatePenjualan.ID)

	tsql := fmt.Sprintf("UPDATE penjualan SET %s WHERE id = ?", setClause)

	// Execute query
	result := db.Exec(tsql, args...)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return updatePenjualan, gorm.ErrRecordNotFound
		}
		return updatePenjualan, result.Error
	}

	// Retrieve the updated penyimpanan
	row := db.Raw("SELECT id, id_user, id_barang, kuantitas, harga_satuan, total_harga, tanggal_penjualan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM penjualan WHERE id = ?", updatePenjualan.ID).Row()
	if err := row.Scan(&updatePenjualan.ID, &updatePenjualan.IDUser, &updatePenjualan.IDBarang, &updatePenjualan.Kuantitas, &updatePenjualan.HargaSatuan, &updatePenjualan.TotalHarga, &updatePenjualan.TanggalPenjualan, &updatePenjualan.CreatedAt, &updatePenjualan.UpdatedAt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return updatePenjualan, gorm.ErrRecordNotFound
		}
		return updatePenjualan, err
	}

	return updatePenjualan, nil
}

func DeletePenjualan(id int) error {
	db := database.GetDB()

	// Query to delete penjualan by ID
	tsql := `DELETE FROM penjualan WHERE id = ?`

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

// service penyimpanan
func CreatePemesanan(ctx *gin.Context) (models.Pemesanan, error) {

	db := database.GetDB()

	var pemesanan models.Pemesanan
	if err := ctx.ShouldBindJSON(&pemesanan); err != nil {
		return pemesanan, err
	}

	hargaSatuanFloat := float64(pemesanan.HargaSatuan)
	totalHarga := float64(pemesanan.Kuantitas) * hargaSatuanFloat
	totalHarga = float64(int(totalHarga))

	pemesanan.TotalHarga = models.Float64OrString(totalHarga)
	pemesanan.TotalBiayaPemesanan = pemesanan.BiayaTelepon + pemesanan.BiayaAdm + pemesanan.BiayaTransportasi

	// Set the CreatedAt and UpdatedAt fields
	now := time.Now()
	pemesanan.CreatedAt = now
	pemesanan.UpdatedAt = now

	// Query to insert supplier and return the id
	tsql := fmt.Sprintf(`
		INSERT INTO pemesanan (id_user, id_barang, kuantitas, harga_satuan, total_harga, biaya_telepon, biaya_adm, biaya_transportasi, total_biaya_pemesanan, tanggal_pemesanan, created_at, updated_at) 
		VALUES ('%d', '%d', '%d', '%f', '%f','%f','%f','%f','%f','%s', '%s', '%s') RETURNING id`,
		pemesanan.IDUser, pemesanan.IDBarang, pemesanan.Kuantitas, pemesanan.HargaSatuan, pemesanan.TotalHarga, pemesanan.BiayaTelepon, pemesanan.BiayaAdm, pemesanan.BiayaTransportasi, pemesanan.TotalBiayaPemesanan, pemesanan.TanggalPemesanan.Format(time.RFC3339), pemesanan.CreatedAt.Format(time.RFC3339), pemesanan.UpdatedAt.Format(time.RFC3339))

	// Execute query and get the returned ID
	var pemesananID int
	err := db.Raw(tsql).Row().Scan(&pemesananID)
	if err != nil {
		return pemesanan, err
	}

	pemesanan.ID = pemesananID

	return pemesanan, nil

}

// GetPemesanans service to get all Penyimpanans
func GetPemesanans(ctx *gin.Context) ([]models.PemesananWithBarangWithSupplier, error) {
	db := database.GetDB()
	var pemesanans []models.PemesananWithBarangWithSupplier

	// Query to get all pemesanans with inner join and aliases
	tsql := `SELECT 
                 p.id AS id, 
                 p.id_user AS id_user, 
                 p.id_barang AS id_barang, 
                 b.nama_barang AS barang_nama, 
                 s.nama AS supplier_nama,
                 p.kuantitas AS kuantitas, 
                 p.harga_satuan AS harga_satuan, 
                 p.total_harga AS total_harga, 
                 p.biaya_telepon AS biaya_telepon, 
                 p.biaya_adm AS biaya_adm, 
                 p.biaya_transportasi AS biaya_transportasi, 
                 p.total_biaya_pemesanan AS total_biaya_pemesanan, 
                 p.tanggal_pemesanan AS tanggal_pemesanan, 
                 COALESCE(p.created_at, NOW()) AS created_at, 
                 COALESCE(p.updated_at, NOW()) AS updated_at
            FROM pemesanan p
            INNER JOIN barang b ON p.id_barang = b.id
            INNER JOIN supplier s ON b.id_supplier = s.id`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pemesanan models.PemesananWithBarangWithSupplier
		if err := rows.Scan(
			&pemesanan.ID,
			&pemesanan.IDUser,
			&pemesanan.IDBarang,
			&pemesanan.BarangNama,
			&pemesanan.SupplierNama,
			&pemesanan.Kuantitas,
			&pemesanan.HargaSatuan,
			&pemesanan.TotalHarga,
			&pemesanan.BiayaTelepon,
			&pemesanan.BiayaAdm,
			&pemesanan.BiayaTransportasi,
			&pemesanan.TotalBiayaPemesanan,
			&pemesanan.TanggalPemesanan,
			&pemesanan.CreatedAt,
			&pemesanan.UpdatedAt,
		); err != nil {
			return nil, err
		}

		pemesanans = append(pemesanans, pemesanan)
	}

	return pemesanans, nil
}

// GetPemesananbyId service to get a pemesanan by ID
func GetPemesananbyId(id int) (models.Pemesanan, error) {
	db := database.GetDB()

	var pemesanan models.Pemesanan

	// Query to get penyimpanan by ID
	tsql := `SELECT id, id_user, id_barang, kuantitas, harga_satuan, total_harga, biaya_telepon, biaya_adm, biaya_transportasi, total_biaya_pemesanan, tanggal_pemesanan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM pemesanan WHERE id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(&pemesanan.ID, &pemesanan.IDUser, &pemesanan.IDBarang, &pemesanan.Kuantitas, &pemesanan.HargaSatuan, &pemesanan.TotalHarga, &pemesanan.BiayaTelepon, &pemesanan.BiayaAdm, &pemesanan.BiayaTransportasi, &pemesanan.TotalBiayaPemesanan, &pemesanan.TanggalPemesanan, &pemesanan.CreatedAt, &pemesanan.UpdatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pemesanan, gorm.ErrRecordNotFound
		}
		return pemesanan, err
	}

	return pemesanan, nil
}

// UpdatePemesanan service to update pemesanan with optional fields
func UpdatePemesanan(updatePemesanan models.Pemesanan) (models.Pemesanan, error) {
	db := database.GetDB()

	// Remove ID from the map to prevent updating it
	updatedFields := map[string]interface{}{}
	if updatePemesanan.IDUser != 0 {
		updatedFields["id_user"] = updatePemesanan.IDUser
	}
	if updatePemesanan.IDBarang != 0 {
		updatedFields["id_barang"] = updatePemesanan.IDBarang
	}
	if updatePemesanan.Kuantitas != 0 {
		updatedFields["kuantitas"] = updatePemesanan.Kuantitas
	}
	if updatePemesanan.HargaSatuan != 0 {
		updatedFields["harga_satuan"] = updatePemesanan.HargaSatuan
	}
	if updatePemesanan.Kuantitas != 0 || updatePemesanan.HargaSatuan != 0 {
		updatedFields["total_harga"] = updatePemesanan.Kuantitas * models.IntOrString(updatePemesanan.HargaSatuan)
	}
	// if updatePemesanan.TotalHarga != 0 {
	// 	updatedFields["total_harga"] = updatePemesanan.TotalHarga
	// }
	if updatePemesanan.BiayaTelepon != 0 {
		updatedFields["biaya_telepon"] = updatePemesanan.BiayaTelepon
	}
	if updatePemesanan.BiayaAdm != 0 {
		updatedFields["biaya_adm"] = updatePemesanan.BiayaAdm
	}
	if updatePemesanan.BiayaTransportasi != 0 {
		updatedFields["biaya_transportasi"] = updatePemesanan.BiayaTransportasi
	}
	if updatePemesanan.BiayaTelepon != 0 || updatePemesanan.BiayaAdm != 0 || updatePemesanan.BiayaTransportasi != 0 {
		updatedFields["total_biaya_pemesanan"] = updatePemesanan.BiayaTelepon + updatePemesanan.BiayaAdm + updatePemesanan.BiayaTransportasi
	}
	// if updatePemesanan.TotalBiayaPemesanan != 0 {
	// 	updatedFields["total_biaya_pemesanan"] = updatePemesanan.TotalBiayaPemesanan
	// }
	if !updatePemesanan.TanggalPemesanan.IsZero() {
		updatedFields["tanggal_pemesanan"] = updatePemesanan.TanggalPemesanan
	}
	updatedFields["updated_at"] = time.Now()

	if len(updatedFields) == 0 {
		return updatePemesanan, errors.New("no fields to update")
	}

	// Query to update penyimpanan by ID
	setClause := ""
	args := []interface{}{}
	for field, value := range updatedFields {
		setClause += fmt.Sprintf("%s = ?, ", field)
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-2] // remove the last comma and space
	args = append(args, updatePemesanan.ID)

	tsql := fmt.Sprintf("UPDATE pemesanan SET %s WHERE id = ?", setClause)

	// Execute query
	result := db.Exec(tsql, args...)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return updatePemesanan, gorm.ErrRecordNotFound
		}
		return updatePemesanan, result.Error
	}

	// Retrieve the updated penyimpanan
	row := db.Raw("SELECT id, id_user, id_barang, kuantitas, harga_satuan, total_harga, biaya_telepon, biaya_adm, biaya_transportasi, total_biaya_pemesanan, tanggal_pemesanan, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM pemesanan WHERE id = ?", updatePemesanan.ID).Row()
	if err := row.Scan(&updatePemesanan.ID, &updatePemesanan.IDUser, &updatePemesanan.IDBarang, &updatePemesanan.Kuantitas, &updatePemesanan.HargaSatuan, &updatePemesanan.TotalHarga, &updatePemesanan.BiayaTelepon, &updatePemesanan.BiayaAdm, &updatePemesanan.BiayaTransportasi, &updatePemesanan.TotalBiayaPemesanan, &updatePemesanan.TanggalPemesanan, &updatePemesanan.CreatedAt, &updatePemesanan.UpdatedAt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return updatePemesanan, gorm.ErrRecordNotFound
		}
		return updatePemesanan, err
	}

	return updatePemesanan, nil
}

func DeletePemesanan(id int) error {
	db := database.GetDB()

	// Query to delete pemesanan by ID
	tsql := `DELETE FROM pemesanan WHERE id = ?`

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

func Login(ctx *gin.Context) (models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		return user, err
	}

	// Prepare query to find the user by username and password
	query := `
		SELECT id, nama, username, posisi, hp, alamat, created_at, updated_at 
		FROM "user" 
		WHERE username = ? AND password = ?`

	// Execute query and scan the result into the user struct
	err := db.Raw(query, user.Username, user.Password).Scan(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, gorm.ErrRecordNotFound
		}
		return user, err
	}

	// Return the user object
	return user, nil
}

func CreateUser(ctx *gin.Context) (models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		return user, err
	}

	// Set the CreatedAt and UpdatedAt fields
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Prepare query to insert user and return the id
	query := `
		INSERT INTO "user" (nama, username, password, posisi, hp, alamat, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	// Execute query and get the returned ID
	var userID int
	err := db.Raw(query, user.Nama, user.Username, user.Password, user.Posisi, user.HP, user.Alamat, user.CreatedAt, user.UpdatedAt).Scan(&userID).Error
	if err != nil {
		return user, err
	}

	user.ID = userID

	return user, nil
}

// GetUsers service to get all GetUsers
func GetUsers(ctx *gin.Context) ([]models.User, error) {
	db := database.GetDB()
	var users []models.User

	// Query to get all users
	tsql := `SELECT id, nama, username, password, posisi, hp, alamat, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM "user"`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nama, &user.Username, &user.Password, &user.Posisi, &user.HP, &user.Alamat, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserbyId service to get a user by ID
func GetUserbyId(id int) (models.User, error) {
	db := database.GetDB()

	var user models.User

	// Query to get user by ID
	tsql := `SELECT id, nama, username, password, posisi, hp, alamat, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM "user" WHERE id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(&user.ID, &user.Nama, &user.Username, &user.Password, &user.Posisi, &user.HP, &user.Alamat, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, gorm.ErrRecordNotFound
		}
		return user, err
	}

	return user, nil
}

// UpdateUser service to update User with optional fields
func UpdateUser(updatedUser models.User) (models.User, error) {
	db := database.GetDB()

	// Remove ID from the map to prevent updating it
	updatedFields := map[string]interface{}{}
	if updatedUser.Nama != "" {
		updatedFields["nama"] = updatedUser.Nama
	}
	if updatedUser.Username != "" {
		updatedFields["username"] = updatedUser.Username
	}
	if updatedUser.Password != "" {
		updatedFields["password"] = updatedUser.Password
	}
	if updatedUser.Posisi != "" {
		updatedFields["posisi"] = updatedUser.Posisi
	}
	if updatedUser.HP != "" {
		updatedFields["hp"] = updatedUser.HP
	}
	if updatedUser.Alamat != "" {
		updatedFields["alamat"] = updatedUser.Alamat
	}
	updatedFields["updated_at"] = time.Now()

	if len(updatedFields) == 0 {
		return updatedUser, errors.New("no fields to update")
	}

	// Query to update user by ID
	setClause := ""
	args := []interface{}{}
	for field, value := range updatedFields {
		setClause += fmt.Sprintf("%s = ?, ", field)
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-2] // remove the last comma and space
	args = append(args, updatedUser.ID)

	tsql := fmt.Sprintf(`UPDATE "user" SET %s WHERE id = ?`, setClause)

	// Execute query
	result := db.Exec(tsql, args...)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return updatedUser, gorm.ErrRecordNotFound
		}
		return updatedUser, result.Error
	}

	// Retrieve the updated user
	row := db.Raw(`SELECT id, nama, username, password, posisi, hp, alamat, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at FROM "user" WHERE id = ?`, updatedUser.ID).Row()
	if err := row.Scan(&updatedUser.ID, &updatedUser.Nama, &updatedUser.Username, &updatedUser.Password, &updatedUser.Posisi, &updatedUser.HP, &updatedUser.Alamat, &updatedUser.CreatedAt, &updatedUser.UpdatedAt); err != nil {
		if err == gorm.ErrRecordNotFound {
			return updatedUser, gorm.ErrRecordNotFound
		}
		return updatedUser, err
	}

	return updatedUser, nil
}

func DeleteUser(id int) error {
	db := database.GetDB()

	// Query to delete user by ID
	tsql := `DELETE FROM "user" WHERE id = ?`

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

	// Insert barang and return the id using GORM
	err := db.Raw(`
		INSERT INTO barang (id_supplier, nama_barang, created_at, updated_at) 
		VALUES (?, ?, ?, ?) RETURNING id`,
		barang.IDSupplier, barang.NamaBarang, barang.CreatedAt.Format(time.RFC3339), barang.UpdatedAt.Format(time.RFC3339)).
		Scan(&barang.ID).Error
	if err != nil {
		return barang, err
	}

	return barang, nil
}

// GetBarangs service to get all Barangs with supplier details
func GetBarangs(ctx *gin.Context) ([]BarangWithSupplier, error) {
	db := database.GetDB()
	var barangs []BarangWithSupplier

	// Query to get all Barangs with supplier details
	tsql := `SELECT 
                b.id, 
                b.id_supplier, 
                b.nama_barang, 
                COALESCE(b.created_at, NOW()) AS created_at, 
                COALESCE(b.updated_at, NOW()) AS updated_at,
                s.Nama AS supplier_nama,
                s.perusahaan AS supplier_perusahaan,
                s.kontak AS supplier_kontak,
                s.alamat AS supplier_alamat
            FROM 
                Barang b
            INNER JOIN 
                Supplier s ON b.id_supplier = s.id;`

	// Execute query
	rows, err := db.Raw(tsql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var barang BarangWithSupplier
		if err := rows.Scan(
			&barang.ID,
			&barang.IDSupplier,
			&barang.NamaBarang,
			&barang.CreatedAt,
			&barang.UpdatedAt,
			&barang.SupplierNama,
			&barang.SupplierPerusahaan,
			&barang.SupplierKontak,
			&barang.SupplierAlamat,
		); err != nil {
			return nil, err
		}
		barangs = append(barangs, barang)
	}

	return barangs, nil
}

// GetBarangByID service to get a barang by ID
func GetBarangByID(id int) (BarangWithSupplier, error) {
	db := database.GetDB()

	var barang BarangWithSupplier

	// Query to get barang by ID with supplier details
	tsql := `SELECT 
                b.id, 
                b.id_supplier, 
                b.nama_barang, 
                COALESCE(b.created_at, NOW()) AS created_at, 
                COALESCE(b.updated_at, NOW()) AS updated_at,
                s.Nama AS supplier_nama,
                s.perusahaan AS supplier_perusahaan,
                s.kontak AS supplier_kontak,
                s.alamat AS supplier_alamat
            FROM 
                Barang b
            INNER JOIN 
                Supplier s ON b.id_supplier = s.id
            WHERE 
                b.id = ?`

	// Execute query
	row := db.Raw(tsql, id).Row()
	if err := row.Scan(
		&barang.ID,
		&barang.IDSupplier,
		&barang.NamaBarang,
		&barang.CreatedAt,
		&barang.UpdatedAt,
		&barang.SupplierNama,
		&barang.SupplierPerusahaan,
		&barang.SupplierKontak,
		&barang.SupplierAlamat,
	); err != nil {
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

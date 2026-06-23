package main

import (
	"backend/config"
	"backend/models"
	"fmt"
	"time"
)

func SeedData() {
	// Koneksi ke Database MySQL
	config.ConnectDatabase()

	// Pastikan tabel sudah termigrasi
	config.DB.AutoMigrate(
		&models.User{},
		&models.Appointment{},
		&models.MedicalResult{},
		&models.Prescription{},
		&models.PrescriptionItem{},
		&models.Medicine{},
	)

	fmt.Println("🚀 Memulai Simulasi Dokter untuk Demo UTS...")

	// 1. Ambil User Pertama (Akun Demo)
	var user models.User
	if err := config.DB.First(&user).Error; err != nil {
		fmt.Println("❌ Gagal: Tidak ada user di database. Harap register dulu di aplikasi.")
		return
	}

	// 2. Simulasi Janji Temu Selesai
	appointment := models.Appointment{
		UserID:           user.ID,
		DoctorID:         1,
		DoctorScheduleID: 1,
		Status:           "completed",
		Notes:            "Keluhan pusing dan mual selesai diperiksa. Pasien disarankan istirahat.",
	}
	config.DB.Create(&appointment)

	// 3. Simulasi Hasil Lab Medis (Medical Result)
	medicalResult := models.MedicalResult{
		UserID:        user.ID,
		AppointmentID: appointment.ID,
		DoctorName:    "dr. Budi Santoso, Sp.A",
		TestType:      "Laboratorium",
		TestName:      "Cek Darah Lengkap",
		Result:        "Hemoglobin normal. Trombosit sedikit menurun.",
		Notes:         "Perbanyak minum air putih dan istirahat cukup.",
		FileURL:       "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
		ResultDate:    time.Now(),
	}
	config.DB.Create(&medicalResult)

	// 4. Simulasi Katalog Obat (Pastikan ada obat Amoxicillin)
	medicine := models.Medicine{
		Name:        "Paracetamol 500mg",
		Description: "Obat penurun panas dan pereda nyeri.",
		Price:       15000,
		Stock:       100,
		ImageURL:    "https://via.placeholder.com/150",
	}
	config.DB.FirstOrCreate(&medicine, models.Medicine{Name: "Paracetamol 500mg"})

	// 5. Simulasi E-Resep Obat (Prescription)
	prescription := models.Prescription{
		UserID:           user.ID,
		DoctorName:       "dr. Budi Santoso, Sp.A",
		DoctorSpecialty:  "Spesialis Anak",
		Diagnosis:        "Demam ringan disertai gejala flu.",
		PrescriptionDate: time.Now().Format("2006-01-02"),
		ExpiryDate:       time.Now().AddDate(0, 0, 7).Format("2006-01-02"), // Berlaku 7 hari
		Status:           "Active",
	}
	config.DB.Create(&prescription)

	// Tambahkan item obat ke dalam resep
	item := models.PrescriptionItem{
		PrescriptionID:    prescription.ID,
		MedicineID:        medicine.ID,
		DosageInstruction: "3 x 1 sehari sesudah makan",
		Duration:          "5 hari",
		Quantity:          15,
	}
	config.DB.Create(&item)

	fmt.Println("✅ SUKSES! Data Simulasi (Appointment Selesai, Hasil Medis, E-Resep) berhasil ditambahkan ke akun:", user.Email)
	fmt.Println("👉 Sekarang buka aplikasi Flutter Anda, periksa menu 'Riwayat' dan 'E-Resep'. Aplikasi sudah 100% Dinamis!")
}

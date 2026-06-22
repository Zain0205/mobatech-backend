package main

import (
	"log"
	"time"

	"backend/config"
	"backend/models"
)

func main() {
	log.Println("Connecting to Database...")
	config.ConnectDatabase()
	db := config.DB

	log.Println("Seeding Doctors...")
	doctors := []models.Doctor{
		{Name: "Dr. Andi Setiawan", Specialization: "Spesialis Penyakit Dalam", ContactInfo: "081234567890", Description: "Berpengalaman 10 tahun di RS Hermina.", IsActive: true},
		{Name: "Dr. Budi Santoso", Specialization: "Spesialis Anak", ContactInfo: "081234567891", Description: "Berpengalaman 8 tahun, sangat disukai anak-anak.", IsActive: true},
		{Name: "Dr. Citra Lestari", Specialization: "Spesialis Kandungan", ContactInfo: "081234567892", Description: "Ramah dan komunikatif.", IsActive: true},
		{Name: "Dr. Dewi Anggraini", Specialization: "Dokter Gigi", ContactInfo: "081234567893", Description: "Ahli estetika gigi dan bedah mulut.", IsActive: true},
	}
	for _, d := range doctors {
		db.Where(models.Doctor{Name: d.Name}).FirstOrCreate(&d)
	}

	log.Println("Fetching Users...")
	var users []models.User
	db.Find(&users)

	if len(users) == 0 {
		log.Println("No users found. Creating a dummy user...")
		dummyUser := models.User{FullName: "Dummy User", Email: "dummy@test.com"}
		db.Create(&dummyUser)
		users = append(users, dummyUser)
	}

	log.Println("Clearing old medical results and reminders...")
	db.Exec("DELETE FROM medical_results")
	db.Exec("DELETE FROM reminders")

	log.Println("Seeding Medical Results and Reminders for all users...")
	now := time.Now()
	for _, u := range users {
		uid := u.ID
		results := []models.MedicalResult{
			{
				UserID:     uid,
				DoctorName: "dr. Andi Setiawan, Sp.PD",
				TestType:   "Patologi Klinik",
				TestName:   "Pemeriksaan Hematologi Rutin & Profil Lipid",
				Result:     "Hemoglobin: 14.5 g/dL (Normal). Leukosit: 7.200 /uL (Normal). Kolesterol Total: 215 mg/dL (Borderline High). Trigliserida: 140 mg/dL (Normal).",
				Notes:      "Terdapat peningkatan ringan pada profil lipid (Hiperkolesterolemia ringan). Disarankan modifikasi gaya hidup (Diet rendah lemak jenuh) dan evaluasi ulang dalam 3 bulan. Tidak memerlukan intervensi farmakologis saat ini.",
				FileURL:    "https://hermina-hospitals.com/reports/lab_001.pdf",
				ResultDate: now.Add(-24 * 7 * time.Hour),
			},
			{
				UserID:     uid,
				DoctorName: "dr. Citra Lestari, Sp.Rad",
				TestType:   "Radiologi",
				TestName:   "Ultrasonografi (USG) Abdomen",
				Result:     "Hepar: Ukuran dan ekostruktur normal. Vesika Biliaris: Tidak tampak batu/sludge. Ginjal Kanan/Kiri: Batas kortikomeduler tegas, sistem pelviokalises tidak melebar.",
				Notes:      "Kesan: Organ intra-abdomen dalam batas normal. Tidak ditemukan kelainan patologis pada pemeriksaan sonografi saat ini.",
				FileURL:    "https://hermina-hospitals.com/reports/rad_002.pdf",
				ResultDate: now.Add(-24 * 30 * time.Hour),
			},
		}
		for _, r := range results {
			db.Create(&r)
		}

		reminders := []models.Reminder{
			{UserID: uid, Title: "Pengingat Medikasi", Message: "Waktu pemberian: Paracetamol 500mg (1 tablet) pasca makan. Mohon patuhi dosis anjuran.", ReminderDate: now.Add(2 * time.Hour), IsRead: false, Type: "medicine"},
			{UserID: uid, Title: "Jadwal Evaluasi Klinis", Message: "Mengingatkan jadwal kontrol Anda bersama dr. Andi Setiawan, Sp.PD esok hari. Harap hadir 15 menit sebelum waktu konsultasi.", ReminderDate: now.Add(24 * time.Hour), IsRead: false, Type: "appointment"},
			{UserID: uid, Title: "Hasil Laboratorium Tersedia", Message: "Hasil Pemeriksaan Hematologi Rutin Anda telah dirilis oleh instalasi laboratorium. Silakan tinjau pada menu Rekam Medis.", ReminderDate: now.Add(-2 * time.Hour), IsRead: true, Type: "system"},
		}
		for _, rem := range reminders {
			db.Create(&rem)
		}
	}

	log.Println("✅ Seeding completed successfully! Data injected to DB.")
}

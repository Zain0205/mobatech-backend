#!/bin/bash

echo "1. Registering Test User..."
curl -s -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{
  "full_name": "Test User",
  "email": "test@mcboy.com",
  "phone_number": "08123456789",
  "password": "password123"
}'

echo -e "\n\n2. Login to get Token..."
# Menggunakan python untuk mengekstrak token dari response JSON dengan aman
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{
  "email": "test@mcboy.com",
  "password": "password123"
}' | python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))")

echo "Token: $TOKEN"

echo -e "\n\n3. Creating Doctor McBoy..."
# Asumsikan akan mendapatkan ID tertentu, tapi kita pakai jq/python untuk ambil ID-nya
DOC_ID=$(curl -s -X POST http://localhost:8080/api/admin/doctors -H "Content-Type: application/json" -d '{
  "name": "dr. Mcboy, Sp.PD",
  "specialization": "Penyakit Dalam",
  "contact_info": "mcboy@example.com",
  "description": "Dokter Spesialis Penyakit Dalam.",
  "image_url": "https://api.dicebear.com/7.x/avataaars/svg?seed=Mcboy",
  "is_active": true
}' | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('ID', 1))")

echo "Doctor ID: $DOC_ID"

echo -e "\n\n4. Creating Schedule for Doctor McBoy..."
SCHED_ID=$(curl -s -X POST http://localhost:8080/api/admin/schedules -H "Content-Type: application/json" -d "{
  \"doctor_id\": $DOC_ID,
  \"date\": \"2026-06-25T00:00:00Z\",
  \"start_time\": \"09:00\",
  \"end_time\": \"12:00\",
  \"quota\": 10,
  \"is_available\": true
}" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('ID', 1))")

echo "Schedule ID: $SCHED_ID"

echo -e "\n\n5. Booking Appointment for Doctor McBoy..."
curl -s -X POST http://localhost:8080/api/appointments -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "{
  \"doctor_id\": $DOC_ID,
  \"doctor_schedule_id\": $SCHED_ID,
  \"notes\": \"Konsultasi pertama dengan dr. McBoy\"
}"

echo -e "\n\nSelesai! Janji temu dengan dr. Mcboy berhasil dibuat."

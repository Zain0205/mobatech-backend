#!/bin/bash

# Doctors
echo "Creating doctors..."
curl -s -X POST http://localhost:8080/api/admin/doctors -H "Content-Type: application/json" -d '{
  "name": "dr. Andi Saputra, Sp.PD",
  "specialization": "Penyakit Dalam",
  "contact_info": "andi@example.com",
  "description": "Dokter Spesialis Penyakit Dalam dengan pengalaman 10 tahun.",
  "image_url": "https://api.dicebear.com/7.x/avataaars/svg?seed=Andi",
  "is_active": true
}'

echo -e "\n"
curl -s -X POST http://localhost:8080/api/admin/doctors -H "Content-Type: application/json" -d '{
  "name": "dr. Siti Aminah, Sp.A",
  "specialization": "Anak",
  "contact_info": "siti@example.com",
  "description": "Spesialis Anak yang ramah dan peduli pada tumbuh kembang balita.",
  "image_url": "https://api.dicebear.com/7.x/avataaars/svg?seed=Siti",
  "is_active": true
}'

echo -e "\n"
curl -s -X POST http://localhost:8080/api/admin/doctors -H "Content-Type: application/json" -d '{
  "name": "dr. Budi Santoso, Sp.JP",
  "specialization": "Jantung",
  "contact_info": "budi@example.com",
  "description": "Spesialis Jantung yang ahli dalam kardiologi intervensi.",
  "image_url": "https://api.dicebear.com/7.x/avataaars/svg?seed=Budi",
  "is_active": true
}'

echo -e "\n"
# Schedules
echo "Creating schedules..."
curl -s -X POST http://localhost:8080/api/admin/schedules -H "Content-Type: application/json" -d '{
  "doctor_id": 1,
  "date": "2026-06-20T00:00:00Z",
  "start_time": "09:00",
  "end_time": "12:00",
  "quota": 10,
  "is_available": true
}'
echo -e "\n"
curl -s -X POST http://localhost:8080/api/admin/schedules -H "Content-Type: application/json" -d '{
  "doctor_id": 1,
  "date": "2026-06-21T00:00:00Z",
  "start_time": "13:00",
  "end_time": "16:00",
  "quota": 15,
  "is_available": true
}'
echo -e "\n"
curl -s -X POST http://localhost:8080/api/admin/schedules -H "Content-Type: application/json" -d '{
  "doctor_id": 2,
  "date": "2026-06-20T00:00:00Z",
  "start_time": "08:00",
  "end_time": "11:00",
  "quota": 12,
  "is_available": true
}'

echo -e "\nDone seeding doctors and schedules!"

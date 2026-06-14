#!/bin/bash

# Categories
echo "Creating categories..."
curl -s -X POST http://localhost:8080/api/admin/pharmacy/categories -H "Content-Type: application/json" -d '{"name": "Obat Bebas", "description": "Obat yang dapat dibeli tanpa resep dokter", "icon": "medication"}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/categories -H "Content-Type: application/json" -d '{"name": "Obat Keras / Resep", "description": "Obat yang memerlukan resep dokter", "icon": "prescription"}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/categories -H "Content-Type: application/json" -d '{"name": "Vitamin & Suplemen", "description": "Suplemen untuk menjaga kesehatan tubuh", "icon": "vitamin"}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/categories -H "Content-Type: application/json" -d '{"name": "Alat Kesehatan", "description": "Peralatan medis dan perlengkapan kesehatan", "icon": "medical_devices"}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/categories -H "Content-Type: application/json" -d '{"name": "Obat Herbal", "description": "Obat tradisional dari bahan alami", "icon": "herbal"}'
echo -e "\n"

# Medicines
echo "Creating medicines..."
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 1,
  "name": "Paracetamol 500mg",
  "generic_name": "Paracetamol",
  "description": "Obat pereda nyeri dan penurun demam",
  "dosage": "500mg",
  "unit": "Strip",
  "price": 2500,
  "stock": 500,
  "requires_prescription": false,
  "manufacturer": "Kimia Farma",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 2,
  "name": "Amoxicillin 500mg",
  "generic_name": "Amoxicillin",
  "description": "Antibiotik golongan penisilin",
  "dosage": "500mg",
  "unit": "Strip",
  "price": 5000,
  "stock": 200,
  "requires_prescription": true,
  "manufacturer": "Phapros",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 2,
  "name": "Omeprazole 20mg",
  "generic_name": "Omeprazole",
  "description": "Obat untuk menurunkan asam lambung",
  "dosage": "20mg",
  "unit": "Strip",
  "price": 8000,
  "stock": 150,
  "requires_prescription": true,
  "manufacturer": "Dexa Medica",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 3,
  "name": "Vitamin C 1000mg",
  "generic_name": "Ascorbic Acid",
  "description": "Vitamin C dosis tinggi untuk daya tahan tubuh",
  "dosage": "1000mg",
  "unit": "Botol",
  "price": 15000,
  "stock": 300,
  "requires_prescription": false,
  "manufacturer": "Enervon",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 3,
  "name": "Vitamin D3 1000 IU",
  "generic_name": "Cholecalciferol",
  "description": "Vitamin D untuk kesehatan tulang",
  "dosage": "1000 IU",
  "unit": "Botol",
  "price": 25000,
  "stock": 200,
  "requires_prescription": false,
  "manufacturer": "Sido Muncul",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 1,
  "name": "Ibuprofen 400mg",
  "generic_name": "Ibuprofen",
  "description": "Obat anti inflamasi non steroid",
  "dosage": "400mg",
  "unit": "Strip",
  "price": 3500,
  "stock": 400,
  "requires_prescription": false,
  "manufacturer": "Kimia Farma",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 1,
  "name": "Cetirizine 10mg",
  "generic_name": "Cetirizine Hydrochloride",
  "description": "Obat anti alergi antihistamin",
  "dosage": "10mg",
  "unit": "Strip",
  "price": 4000,
  "stock": 250,
  "requires_prescription": false,
  "manufacturer": "Kalbe",
  "image_url": ""
}'
curl -s -X POST http://localhost:8080/api/admin/pharmacy/medicines -H "Content-Type: application/json" -d '{
  "category_id": 2,
  "name": "Metformin 500mg",
  "generic_name": "Metformin",
  "description": "Obat untuk diabetes tipe 2",
  "dosage": "500mg",
  "unit": "Strip",
  "price": 6000,
  "stock": 180,
  "requires_prescription": true,
  "manufacturer": "Dexa Medica",
  "image_url": ""
}'
echo -e "\n"

# Prescription
echo "Creating prescription for user 1..."
curl -s -X POST http://localhost:8080/api/admin/pharmacy/prescriptions -H "Content-Type: application/json" -d '{
  "user_id": 1,
  "doctor_name": "dr. Anisa Putri, Sp.PD",
  "doctor_specialty": "Spesialis Penyakit Dalam",
  "diagnosis": "GERD (Gastroesophageal Reflux Disease)",
  "prescription_date": "2026-06-14",
  "expiry_date": "2026-07-14",
  "status": "Active",
  "notes": "Hindari makanan pedas dan asam. Kontrol 2 minggu lagi.",
  "items": [
    {
      "medicine_id": 3,
      "dosage_instruction": "1x1 sehari sebelum makan pagi",
      "duration": "14 hari",
      "quantity": 14,
      "notes": "Diminum 30 menit sebelum makan"
    },
    {
      "medicine_id": 1,
      "dosage_instruction": "3x1 sehari sesudah makan jika nyeri",
      "duration": "7 hari",
      "quantity": 21,
      "notes": "Hanya jika diperlukan"
    }
  ]
}'
echo -e "\nDone!"

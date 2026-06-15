#!/bin/bash

BASE_DOC="http://127.0.0.1:8080/api/admin/doctors"
BASE_SCH="http://127.0.0.1:8080/api/admin/schedules"

echo "Seeding additional doctors..."

# 1. dr. Tio
RES1=$(curl -s -X POST $BASE_DOC \
  -H "Content-Type: application/json" \
  -d '{"name": "dr. Tio, Sp.PD", "specialization": "Spesialis Penyakit Dalam", "contact_info": "08123456781", "description": "Dokter spesialis penyakit dalam dengan pengalaman 10 tahun.", "image_url": "https://cdn-icons-png.flaticon.com/512/3774/3774299.png", "is_active": true}')
echo "Response Tio: $RES1"
ID1=$(echo $RES1 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID1" ]; then
  curl -s -X POST $BASE_SCH -H "Content-Type: application/json" -d "{\"doctor_id\": $ID1, \"date\": \"2026-06-18T00:00:00Z\", \"start_time\": \"09:00\", \"end_time\": \"13:00\", \"quota\": 15}"
fi

# 2. dr. Tirta
RES2=$(curl -s -X POST $BASE_DOC \
  -H "Content-Type: application/json" \
  -d '{"name": "dr. Tirta, Sp.A", "specialization": "Spesialis Anak", "contact_info": "08123456782", "description": "Dokter anak yang ramah dan menyenangkan.", "image_url": "https://cdn-icons-png.flaticon.com/512/3774/3774299.png", "is_active": true}')
echo "Response Tirta: $RES2"
ID2=$(echo $RES2 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID2" ]; then
  curl -s -X POST $BASE_SCH -H "Content-Type: application/json" -d "{\"doctor_id\": $ID2, \"date\": \"2026-06-19T00:00:00Z\", \"start_time\": \"10:00\", \"end_time\": \"14:00\", \"quota\": 20}"
fi

# 3. dr. Boyke
RES3=$(curl -s -X POST $BASE_DOC \
  -H "Content-Type: application/json" \
  -d '{"name": "dr. Boyke, Sp.OG", "specialization": "Spesialis Kandungan", "contact_info": "08123456783", "description": "Ahli kebidanan dan kandungan terpercaya.", "image_url": "https://cdn-icons-png.flaticon.com/512/3774/3774299.png", "is_active": true}')
echo "Response Boyke: $RES3"
ID3=$(echo $RES3 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID3" ]; then
  curl -s -X POST $BASE_SCH -H "Content-Type: application/json" -d "{\"doctor_id\": $ID3, \"date\": \"2026-06-20T00:00:00Z\", \"start_time\": \"08:00\", \"end_time\": \"12:00\", \"quota\": 10}"
fi

# 4. dr. Richard Lee
RES4=$(curl -s -X POST $BASE_DOC \
  -H "Content-Type: application/json" \
  -d '{"name": "dr. Richard Lee, Sp.KK", "specialization": "Spesialis Kulit & Kelamin", "contact_info": "08123456784", "description": "Pakar kecantikan dan kesehatan kulit ternama.", "image_url": "https://cdn-icons-png.flaticon.com/512/3774/3774299.png", "is_active": true}')
echo "Response Richard: $RES4"
ID4=$(echo $RES4 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID4" ]; then
  curl -s -X POST $BASE_SCH -H "Content-Type: application/json" -d "{\"doctor_id\": $ID4, \"date\": \"2026-06-21T00:00:00Z\", \"start_time\": \"13:00\", \"end_time\": \"17:00\", \"quota\": 12}"
fi

echo "Done seeding extra doctors!"

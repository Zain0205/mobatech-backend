#!/bin/bash

BASE_URL="http://127.0.0.1:8080/api/admin/polyclinics"

echo "Seeding Polyclinics..."

# 1. Poli Anak
echo "Adding Poli Anak..."
RES1=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"name": "Poli Anak", "description": "Layanan spesialis kesehatan bayi dan anak.", "image_url": "https://cdn-icons-png.flaticon.com/512/3063/3063205.png", "is_active": true}')
echo "Response: $RES1"
ID1=$(echo $RES1 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID1" ]; then
  curl -s -X POST "$BASE_URL/$ID1/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Senin", "start_time": "08:00", "end_time": "12:00", "is_available": true}'
  curl -s -X POST "$BASE_URL/$ID1/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Rabu", "start_time": "08:00", "end_time": "14:00", "is_available": true}'
fi

# 2. Poli Gigi
echo "Adding Poli Gigi..."
RES2=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"name": "Poli Gigi", "description": "Layanan pemeriksaan dan perawatan gigi terpadu.", "image_url": "https://cdn-icons-png.flaticon.com/512/3014/3014282.png", "is_active": true}')
echo "Response: $RES2"
ID2=$(echo $RES2 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID2" ]; then
  curl -s -X POST "$BASE_URL/$ID2/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Selasa", "start_time": "09:00", "end_time": "15:00", "is_available": true}'
  curl -s -X POST "$BASE_URL/$ID2/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Kamis", "start_time": "09:00", "end_time": "15:00", "is_available": true}'
fi

# 3. Poli Penyakit Dalam
echo "Adding Poli Penyakit Dalam..."
RES3=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"name": "Poli Penyakit Dalam", "description": "Konsultasi dan penanganan organ dalam oleh spesialis.", "image_url": "https://cdn-icons-png.flaticon.com/512/3003/3003011.png", "is_active": true}')
echo "Response: $RES3"
ID3=$(echo $RES3 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID3" ]; then
  curl -s -X POST "$BASE_URL/$ID3/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Senin", "start_time": "10:00", "end_time": "16:00", "is_available": true}'
  curl -s -X POST "$BASE_URL/$ID3/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Jumat", "start_time": "13:00", "end_time": "17:00", "is_available": true}'
fi

# 4. Poli Kulit & Kelamin
echo "Adding Poli Kulit & Kelamin..."
RES4=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"name": "Poli Kulit & Kelamin", "description": "Pemeriksaan dan penanganan masalah kulit serta kelamin.", "image_url": "https://cdn-icons-png.flaticon.com/512/3655/3655570.png", "is_active": true}')
echo "Response: $RES4"
ID4=$(echo $RES4 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID4" ]; then
  curl -s -X POST "$BASE_URL/$ID4/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Rabu", "start_time": "10:00", "end_time": "14:00", "is_available": true}'
  curl -s -X POST "$BASE_URL/$ID4/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Sabtu", "start_time": "09:00", "end_time": "12:00", "is_available": true}'
fi

# 5. Poli Kandungan
echo "Adding Poli Kandungan..."
RES5=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"name": "Poli Kandungan", "description": "Pemeriksaan kehamilan dan organ reproduksi wanita.", "image_url": "https://cdn-icons-png.flaticon.com/512/5031/5031911.png", "is_active": true}')
echo "Response: $RES5"
ID5=$(echo $RES5 | grep -o '"ID":[0-9]*' | head -1 | cut -d':' -f2)

if [ ! -z "$ID5" ]; then
  curl -s -X POST "$BASE_URL/$ID5/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Selasa", "start_time": "08:00", "end_time": "13:00", "is_available": true}'
  curl -s -X POST "$BASE_URL/$ID5/schedules" -H "Content-Type: application/json" -d '{"day_of_week": "Jumat", "start_time": "14:00", "end_time": "18:00", "is_available": true}'
fi

echo "Done seeding Polyclinics!"

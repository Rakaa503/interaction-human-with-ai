# AVIGO — AI Interaction System

AVIGO adalah sistem AI interaction yang dirancang untuk memahami input pengguna, menganalisis intent, emotion, dan topic, membangun context percakapan, menentukan decision/action, kemudian meneruskan proses ke layer response.

Arsitektur AVIGO menggunakan kombinasi:

- Go Fiber sebagai API Gateway
- Python FastAPI sebagai Machine Learning Service
- PostgreSQL sebagai database
- GORM sebagai ORM
- Scikit-learn sebagai machine learning engine
- TF-IDF sebagai text feature extraction
- Rule-Based Analyzer sebagai baseline/fallback
- Context Engine untuk memahami riwayat percakapan
- Decision Engine untuk menentukan tindakan AI

---

## Architecture

```text
User
 │
 ▼
Go Fiber API
 │
 ▼
Conversation
 │
 ▼
Interaction
 │
 ▼
ML Analyzer
 │
 ├───────────────┐
 │               │
 ▼               ▼
Python ML      Rule-Based
Service        Analyzer
 │
 ▼
Intent / Emotion / Topic
 │
 ▼
Context Engine
 │
 ▼
Decision Engine
 │
 ▼
Action
 │
 ├── Greeting
 ├── Answer Question
 ├── Solve Problem
 ├── Execute Request
 ├── General Conversation
 └── Clarify
 │
 ▼
Response Layer
 │
 ▼
User

AviGo/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   │
│   ├── config/
│   │   └── ...
│   │
│   ├── database/
│   │   └── ...
│   │
│   ├── env/
│   │   └── ...
│   │
│   ├── users/
│   │   └── ...
│   │
│   ├── conversation/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   ├── interaction/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── ml_analyzer.go
│   │   └── ml_analyzer_test.go
│   │
│   ├── context/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── service_test.go
│   │
│   └── decision/
│       ├── model.go
│       ├── service.go
│       └── service_test.go
│
├── ml/
│   │
│   ├── api/
│   │   └── main.py
│   │
│   ├── dataset/
│   │   └── interactions.csv
│   │
│   ├── training/
│   │   ├── validate_dataset.py
│   │   ├── train.py
│   │   ├── predict.py
│   │   ├── cross_validate.py
│   │   ├── benchmark.py
│   │   ├── benchmark_heldout.py
│   │   └── error_analysis.py
│   │
│   ├── models/
│   │   └── intent_model.joblib
│   │
│   └── evaluation/
│       ├── confusion_matrix.png
│       ├── intent_metrics.json
│       └── error_analysis.csv
│
├── migrations/
│   └── ...
│
├── .env
├── go.mod
├── go.sum
└── README.md


Core Components
1. Go Fiber Gateway

Go berfungsi sebagai gateway utama AVIGO.

Tanggung jawab Go:

menerima request dari client
melakukan validasi request
mengelola conversation
menyimpan interaction
memanggil ML service
membangun context
menjalankan decision engine
mengatur routing API
menghubungkan seluruh komponen AI

Go tidak melakukan proses machine learning secara langsung.

Go bertindak sebagai orchestration layer.

2. Python ML Service

Python digunakan khusus untuk machine learning.

ML Service dibuat menggunakan FastAPI dan berjalan pada:

http://127.0.0.1:8000

Endpoint utama:

GET /health
POST /predict

Contoh request:

{
  "input": "Aplikasi saya mengalami bug saat login"
}

Contoh response:

{
  "intent": "problem_solving",
  "confidence": 0.5078
}
3. Intent Classification

AVIGO saat ini memiliki 5 intent utama:

greeting
question
problem_solving
request
general

Contoh:

"Halo AVIGO"
→ greeting


"Bagaimana cara membuat REST API?"
→ question


"Aplikasi saya mengalami bug"
→ problem_solving


"Tolong buatkan script Python"
→ request


"Saya sedang belajar teknologi"
→ general
4. Emotion Classification

Dataset AVIGO memiliki beberapa kategori emotion:

positive
neutral
frustrated
negative

Emotion digunakan sebagai informasi tambahan bagi sistem untuk memahami kondisi pengguna.

Contoh:

"Terima kasih, sangat membantu!"
→ positive


"Database saya error."
→ frustrated / neutral
5. Topic Classification

Topic digunakan untuk mengetahui domain pembicaraan.

Kategori saat ini:

general
technical
education
business

Contoh:

"Bagaimana cara membuat REST API?"
→ technical


"Saya sedang membuat proyek kampus."
→ education


"Saya ingin membuat bisnis digital."
→ business
Machine Learning Pipeline

Pipeline ML AVIGO:

Dataset
   │
   ▼
Dataset Validation
   │
   ▼
Train / Test Split
   │
   ▼
TF-IDF Vectorization
   │
   ▼
ML Classifier
   │
   ▼
Evaluation
   │
   ├── Accuracy
   ├── Precision
   ├── Recall
   └── F1 Score
   │
   ▼
Model
   │
   ▼
FastAPI
   │
   ▼
Go Gateway
Dataset

Dataset utama berada di:

ml/dataset/interactions.csv

Dataset saat ini memiliki:

Total samples: 100

Distribusi intent:

greeting:          20
question:          20
problem_solving:   20
request:           20
general:           20

Dataset divalidasi menggunakan:

python ml\training\validate_dataset.py

Expected result:

Validation : PASS
Training

Training model dilakukan menggunakan:

python ml\training\train.py

Model menghasilkan:

ml/models/intent_model.joblib

Metrics disimpan pada:

ml/evaluation/intent_metrics.json
Model Evaluation

Hasil training terakhir:

Accuracy: 0.9500

Classification report:

                 precision    recall  f1-score


general             1.00       0.75      0.86
greeting            0.80       1.00      0.89
problem_solving     1.00       1.00      1.00
question            1.00       1.00      1.00
request             1.00       1.00      1.00
Cross Validation

AVIGO juga menggunakan 5-fold cross validation.

Jalankan:

python ml\training\cross_validate.py

Hasil terakhir:

Accuracy : 0.9700 +/- 0.0400
Precision: 0.9740 +/- 0.0356
Recall   : 0.9700 +/- 0.0400
F1 Score : 0.9697 +/- 0.0404

Cross validation digunakan untuk mengetahui apakah performa model stabil pada beberapa pembagian dataset.

Fair Benchmark

AVIGO membandingkan:

Rule-Based
vs
Machine Learning

Benchmark dilakukan menggunakan held-out test set.

Jalankan:

python ml\training\benchmark_heldout.py

Hasil terakhir:

                 Rule-Based     ML


Accuracy           0.7000      0.9500
Precision          0.6000      0.9600
Recall             0.7000      0.9500
F1 Score           0.6400      0.9492

Hasil tersebut menunjukkan bahwa model ML saat ini memiliki performa lebih baik dibandingkan rule-based baseline pada held-out test set.

Error Analysis

Error analysis digunakan untuk mengetahui data yang salah diklasifikasikan oleh model.

Jalankan:

python ml\training\error_analysis.py

Output:

Total samples : 100
Errors        : 2
Accuracy      : 0.9800

Hasil error analysis disimpan pada:

ml/evaluation/error_analysis.csv
ML API

FastAPI digunakan sebagai service terpisah.

Install dependency:

pip install fastapi uvicorn scikit-learn joblib

Menjalankan ML service:

python -m uvicorn ml.api.main:app --host 127.0.0.1 --port 8000

Jika berhasil:

Uvicorn running on http://127.0.0.1:8000
ML Health Check

PowerShell:

Invoke-RestMethod `
    -Method GET `
    -Uri "http://127.0.0.1:8000/health"

Expected:

success service   status
------- -------   ------
True    avigo-ml  ok
ML Prediction Test

PowerShell:

$body = @{
    input = "Aplikasi saya mengalami bug saat login"
} | ConvertTo-Json


Invoke-RestMethod `
    -Method POST `
    -Uri "http://127.0.0.1:8000/predict" `
    -ContentType "application/json" `
    -Body $body

Expected:

intent           confidence
------           ----------
problem_solving  0.5078
Go Gateway

Go menjalankan API utama AVIGO.

Contoh konfigurasi ML analyzer:

interactionAnalyzer := interaction.NewMLAnalyzer(
    "http://127.0.0.1:8000",
)

Dengan architecture:

Client
  │
  ▼
Go Fiber :8080
  │
  ▼
Interaction Service
  │
  ▼
ML Analyzer
  │
  ▼
FastAPI :8000
  │
  ▼
ML Model
Running Go Server

Pastikan PostgreSQL aktif dan environment variable sudah tersedia.

Kemudian jalankan:

go run ./cmd/server

Go API berjalan sesuai konfigurasi:

http://localhost:8080

Health check:

Invoke-RestMethod `
    -Method GET `
    -Uri "http://localhost:8080/health"

Expected:

{
  "status": "ok"
}
Interaction API

Create interaction:

POST /api/v1/interactions

Request:

{
  "userId": 3,
  "conversationId": 5,
  "input": "Aplikasi saya mengalami bug saat login"
}

PowerShell:

$body = @{
    userId = 3
    conversationId = 5
    input = "Aplikasi saya mengalami bug saat login"
} | ConvertTo-Json


Invoke-RestMethod `
    -Method POST `
    -Uri "http://localhost:8080/api/v1/interactions" `
    -ContentType "application/json" `
    -Body $body

Response akan berisi hasil analysis:

{
  "success": true,
  "data": {
    "id": 3,
    "userId": 3,
    "conversationId": 5,
    "input": "Aplikasi saya mengalami bug saat login",
    "intent": "problem_solving",
    "emotion": "neutral",
    "topic": "technical",
    "confidence": 0.5078
  }
}
Conversation System

Conversation digunakan untuk menyimpan percakapan pengguna.

Model utama:

Conversation
 ├── ID
 ├── UserID
 ├── Title
 ├── CreatedAt
 ├── UpdatedAt
 └── Messages

Message:

Message
 ├── ID
 ├── ConversationID
 ├── Role
 ├── Content
 └── CreatedAt

Role yang diperbolehkan:

system
user
assistant
Conversation API

Create conversation:

POST /api/v1/conversations

Get conversation:

GET /api/v1/conversations/:id

Add message:

POST /api/v1/conversations/:id/messages

Get interaction history:

GET /api/v1/conversations/:id/interactions
Context Engine

Context Engine bertugas menyediakan informasi percakapan sebelumnya kepada sistem.

Context bukan pengganti ML.

Context digunakan untuk menjawab pertanyaan:

Apa yang sedang dibicarakan?
Apa pesan sebelumnya?
Apa interaction terakhir?
Apa hubungan input sekarang dengan percakapan sebelumnya?

Flow:

Current Input
     │
     ▼
Interaction Analysis
     │
     ▼
Context Engine
     │
     ├── Current Message
     ├── Previous Messages
     ├── Previous Intent
     ├── Previous Topic
     └── Conversation History
     │
     ▼
Decision Engine
Decision Engine

Decision Engine bertugas menentukan apa yang harus dilakukan AI setelah input dianalisis.

Decision Engine tidak menentukan isi jawaban.

Decision Engine menentukan ACTION.

Intent:

greeting
    ↓
ActionGreeting


question
    ↓
ActionAnswerQuestion


problem_solving
    ↓
ActionSolveProblem


request
    ↓
ActionExecuteRequest


general
    ↓
ActionGeneralConversation


unknown
    ↓
ActionClarify
Decision Flow
User Input
    │
    ▼
ML Analyzer
    │
    ▼
Intent
Emotion
Topic
Confidence
    │
    ▼
Context Engine
    │
    ▼
Decision Engine
    │
    ▼
Action
    │
    ├── Greeting
    ├── Answer Question
    ├── Solve Problem
    ├── Execute Request
    ├── General Conversation
    └── Clarify
Example AI Interaction

Input:

"Hai AVIGO"

ML:

Intent: greeting

Decision:

ActionGreeting

Response layer nantinya dapat menghasilkan:

"Halo! 👋 Ada yang bisa AVIGO bantu?"

Input:

"Bagaimana cara membuat REST API?"

ML:

Intent: question
Topic: technical

Decision:

ActionAnswerQuestion

Input:

"Aplikasi saya mengalami bug saat login"

ML:

Intent: problem_solving
Topic: technical

Decision:

ActionSolveProblem
Database

AVIGO menggunakan PostgreSQL.

Database menyimpan:

users
conversations
messages
interactions

Relationship:

User
 │
 ├── Conversations
 │       │
 │       └── Messages
 │
 └── Interactions
         │
         └── Conversation
Database Migration

Migration menggunakan Goose.

Contoh menjalankan migration:

goose -dir migrations postgres "$env:DATABASE_URL" up

Rollback:

goose -dir migrations postgres "$env:DATABASE_URL" down

Pastikan Goose sudah terinstall sebelum menjalankan command tersebut.

Environment

Buat file:

.env

Contoh:

APP_NAME=AVIGO
APP_PORT=8080


DATABASE_URL=postgres://username:password@localhost:5432/avigo


ML_SERVICE_URL=http://127.0.0.1:8000

Jangan commit credential asli ke repository.

Gunakan .env.example untuk konfigurasi yang aman dibagikan.

Development Setup
Requirements

Pastikan environment memiliki:

Go
Python
PostgreSQL
Git

Python dependency:

FastAPI
Uvicorn
scikit-learn
joblib
pandas
numpy
Clone Repository
git clone https://github.com/Rakaa503/interaction-human-with-ai.git

Masuk ke project:

cd interaction-human-with-ai

Python Environment

Buat virtual environment:

python -m venv .venv

Aktifkan:

.venv\Scripts\activate

Install dependency:

pip install fastapi uvicorn scikit-learn joblib pandas numpy
Train ML Model

Setelah environment aktif:

python ml\training\validate_dataset.py

Kemudian:

python ml\training\train.py

Optional evaluation:

python ml\training\cross_validate.py

Benchmark:

python ml\training\benchmark_heldout.py

Error analysis:

python ml\training\error_analysis.py

Prediction test:

python ml\training\predict.py
Run AVIGO

AVIGO membutuhkan dua service utama.

Terminal 1 — ML Service
.venv\Scripts\activate


python -m uvicorn ml.api.main:app --host 127.0.0.1 --port 8000

Expected:

Uvicorn running on http://127.0.0.1:8000
Terminal 2 — Go Gateway
go run ./cmd/server

Expected:

AVIGO API running on :8080

Architecture runtime:

┌──────────────┐
│    Client    │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│ Go Fiber :8080   │
│ API Gateway      │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ Interaction      │
│ Service          │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ Python FastAPI   │
│ :8000            │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ ML Model         │
│ TF-IDF + Model   │
└──────────────────┘
Testing

Run seluruh Go test:

go test ./...

Run interaction test:

go test ./internal/interaction -v

Run context test:

go test ./internal/context -v

Run decision test:

go test ./internal/decision -v

Format Go:

go fmt ./...
Current Development Status
Completed
 Go Fiber API Gateway
 PostgreSQL integration
 Conversation module
 Message system
 Interaction module
 Rule-Based Analyzer
 Machine Learning Analyzer
 Python FastAPI ML Service
 Intent classification
 Emotion classification dataset
 Topic classification dataset
 Dataset validation
 ML training
 Model evaluation
 Cross validation
 Fair held-out benchmark
 Error analysis
 Go → Python ML integration
 Context Engine
 Decision Engine
 Unit testing
In Development
 Response Generation Engine
 Real conversational response
 Greeting response
 Question answering
 Problem solving response
 Request execution
 Context-aware response
 Conversation memory improvement
 AI response orchestration
 Production deployment
Final AI Architecture

Target architecture AVIGO:

                         ┌───────────────┐
                         │     USER      │
                         └───────┬───────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │   GO FIBER API    │
                       │     GATEWAY       │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │   CONVERSATION    │
                       │      MODULE       │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │    INTERACTION    │
                       │      MODULE       │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │    ML ANALYZER    │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │   PYTHON FASTAPI  │
                       │    ML SERVICE     │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │    ML MODEL       │
                       │ TF-IDF + CLASSIFIER│
                       └─────────┬─────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ INTENT / EMOTION / TOPIC│
                    └────────────┬─────────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │   CONTEXT ENGINE  │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │   DECISION ENGINE │
                       └─────────┬─────────┘
                                 │
                                 ▼
                       ┌───────────────────┐
                       │  RESPONSE ENGINE  │
                       └─────────┬─────────┘
                                 │
                                 ▼
                         ┌───────────────┐
                         │     USER      │
                         └───────────────┘
Project Goal

AVIGO bukan hanya sistem classification.

Machine Learning digunakan untuk memahami input.

Context Engine digunakan untuk memahami percakapan.

Decision Engine digunakan untuk menentukan tindakan.

Response Engine nantinya digunakan untuk menghasilkan jawaban.

Dengan demikian:

ML
+
Context
+
Decision
+
Response
=
AI Interaction System

Target akhir AVIGO adalah sistem AI yang mampu melakukan interaksi conversational secara context-aware, bukan hanya mengklasifikasikan input pengguna.

Author

Rakha Avilla

GitHub:

https://github.com/Rakaa503

Repository:

https://github.com/Rakaa503/interaction-human-with-ai

License

This project is currently under development.

License and usage terms will be defined in a future release.

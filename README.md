# PlantPal API

Plant care and health monitoring application, powered by **Google Gemini AI** for plant identification, disease diagnosis, and conversational plant-care guidance.

## OpenAPI Spec

The full Swagger/OpenAPI 2.0 spec is available for auto-generating clients, types, and API docs:

| Resource | URL |
|----------|-----|
| **Swagger JSON** | https://plant-pal-api-ohhx.onrender.com/swagger.json |
| **Swagger UI** | https://plant-pal-api-ohhx.onrender.com/docs/ |

### Auto-generate a frontend client

```bash
# TypeScript / Axios client
npx @openapitools/openapi-generator-cli generate \
  -i https://plant-pal-api-ohhx.onrender.com/swagger.json \
  -g typescript-axios \
  -o ./generated/client

# Python client
npx @openapitools/openapi-generator-cli generate \
  -i https://plant-pal-api-ohhx.onrender.com/swagger.json \
  -g python \
  -o ./generated/client

# Dart / Flutter
npx @openapitools/openapi-generator-cli generate \
  -i https://plant-pal-api-ohhx.onrender.com/swagger.json \
  -g dart \
  -o ./generated/client

# Or use openapi-typescript for TS types only
npx openapi-typescript https://plant-pal-api-ohhx.onrender.com/swagger.json -o ./generated/api-types.ts
```

---

## Tech Stack

- **Language:** Go 1.26.3
- **Router:** Gorilla Mux
- **Database:** PostgreSQL (GORM)
- **AI:** Google Gemini 2.0 Flash
- **Image Storage:** Cloudinary
- **Auth:** JWT (HS256, 24h access / 15d refresh)
- **Docs:** Swagger at `/docs/`

## Setup

Copy `.env.example` or create a `.env` with:

```
PORT=8080
DB_URL=postgres://...
JWT_SECRET=...
CLOUDINARY_URL=...
GEMINI_API_KEY=...
```

```bash
go run cmd/main.go
```

Server starts on `http://localhost:8080`. All 13 database tables auto-migrate on startup.

---

## Authentication

All protected endpoints require a **Bearer token** in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

JWT access tokens expire after **24 hours**. Use `/refresh` to obtain a new pair. Refresh tokens are **single-use** and expire after **15 days**.

---

## Base URL

```
http://localhost:8080
```

Production: `https://plant-pal-api-ohhx.onrender.com`

---

## Error Responses

All errors follow:

```json
{ "error": "string" }
```

| Status | Meaning |
|--------|---------|
| 400 | Bad request / validation error |
| 401 | Unauthorized (missing or invalid token) |
| 403 | Forbidden |
| 404 | Resource not found |
| 409 | Conflict (e.g. duplicate email) |
| 429 | Rate limit exceeded |
| 500 | Internal server error |

---

## API Endpoints

### Table of Contents

1. [Health](#1-get-health)
2. [Register](#2-post-register)
3. [Login](#3-post-login)
4. [Refresh Token](#4-post-refresh)
5. [Logout](#5-post-logout)
6. [List Plants](#6-get-plants)
7. [Create Plant](#7-post-plants)
8. [Get Plant](#8-get-plantsid)
9. [Update Plant](#9-put-plantsid)
10. [Delete Plant](#10-delete-plantsid)
11. [Get Care Plan](#11-get-plantsidcare-plan)
12. [Update Care Plan](#12-put-plantsidcare-plan)
13. [Get Plant Reminders](#13-get-plantsidreminders)
14. [Get Plant Activities](#14-get-plantsidactivities)
15. [Log Activity](#15-post-plantsidactivities)
16. [Get Growth Metrics](#16-get-plantsidgrowth)
17. [Log Growth](#17-post-plantsidgrowth)
18. [Today's Reminders](#18-get-reminderstoday)
19. [Complete / Snooze Reminder](#19-put-remindersid)
20. [Get Notification Settings](#20-get-notifications)
21. [Update Notification Settings](#21-put-notifications)
22. [AI Scan (Identify Plant)](#22-post-scan)
23. [Get Scan Details](#23-get-scanid)
24. [Confirm Scan](#24-post-scanidconfirm)
25. [AI Diagnosis (Disease)](#25-post-diagnosis)
26. [Get Diagnosis Session](#26-get-diagnosissession_id)
27. [Diagnosis Chat](#27-post-diagnosissession_idchat)

---

### 1. `GET /health`

Health check endpoint.

**Auth:** None

**Response** `200 OK`
```json
{ "status": "ok" }
```

---

### 2. `POST /register`

Create a new user account.

**Auth:** None

**Request Body**
```json
{
  "full_name": "string (required)",
  "email": "string (required, must contain @ and .)",
  "password": "string (required, min 8 chars)"
}
```

**Response** `201 Created`
```json
{
  "access_token": "string (JWT, 24h)",
  "refresh_token": "string (JWT, 15d)"
}
```

---

### 3. `POST /login`

Authenticate with email and password.

**Auth:** None

**Request Body**
```json
{
  "email": "string (required)",
  "password": "string (required)"
}
```

**Response** `200 OK`
```json
{
  "access_token": "string (JWT, 24h)",
  "refresh_token": "string (JWT, 15d)"
}
```

---

### 4. `POST /refresh`

Exchange a refresh token for a new access + refresh pair. The old refresh token is revoked.

**Auth:** None

**Request Body**
```json
{
  "refresh_token": "string (required)"
}
```

**Response** `200 OK`
```json
{
  "access_token": "string (JWT, 24h)",
  "refresh_token": "string (JWT, 15d)"
}
```

---

### 5. `POST /logout`

Revoke a refresh token to end a session.

**Auth:** None

**Request Body**
```json
{
  "refresh_token": "string (required)"
}
```

**Response** `200 OK`
```json
{ "message": "logged out successfully" }
```

---

### 6. `GET /plants`

List all plants for the authenticated user.

**Auth:** Bearer JWT

**Response** `200 OK`
```json
[
  {
    "id": 1,
    "created_at": "2026-07-18T10:00:00Z",
    "updated_at": "2026-07-18T10:00:00Z",
    "user_id": 1,
    "species_id": 1,
    "nickname": "Fern",
    "location": "Living Room",
    "health_score": 85,
    "status": "good",
    "species": {
      "id": 1,
      "common_name": "Boston Fern",
      "scientific_name": "Nephrolepis exaltata",
      "family": "Nephrolepidaceae",
      "origin": "Tropical Americas",
      "difficulty_level": "easy"
    },
    "care_plans": [
      {
        "id": 1,
        "plant_id": 1,
        "watering_frequency_days": 7,
        "watering_amount": "200ml",
        "watering_method": "top watering",
        "watering_tips": "Keep soil moist",
        "light_requirement": "bright indirect light",
        "humidity_requirement": "60-80%"
      }
    ]
  }
]
```

---

### 7. `POST /plants`

Add a new plant to your collection.

**Auth:** Bearer JWT

**Request Body**
```json
{
  "species_id": 1,
  "nickname": "My Fern (required)",
  "location": "Living Room"
}
```

**Response** `201 Created`
```json
{
  "id": 2,
  "created_at": "2026-07-18T10:00:00Z",
  "updated_at": "2026-07-18T10:00:00Z",
  "user_id": 1,
  "species_id": 1,
  "nickname": "My Fern",
  "location": "Living Room",
  "health_score": 0,
  "status": "",
  "species": { ... }
}
```

---

### 8. `GET /plants/{id}`

Get detailed info about a specific plant, including all relations.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Response** `200 OK`
```json
{
  "id": 1,
  "user_id": 1,
  "species_id": 1,
  "nickname": "Fern",
  "location": "Living Room",
  "health_score": 85,
  "status": "good",
  "species": { ... },
  "care_plans": [ ... ],
  "reminders": [
    {
      "id": 1,
      "plant_id": 1,
      "task_type": "water",
      "scheduled_time": "2026-07-25T10:00:00Z",
      "is_completed": false,
      "completed_at": "0001-01-01T00:00:00Z",
      "snooze_count": 0
    }
  ],
  "growth_metrics": [
    {
      "id": 1,
      "plant_id": 1,
      "recorded_date": "2026-07-18T10:00:00Z",
      "height_cm": 25.5,
      "growth_rate_status": "moderate"
    }
  ],
  "activity_logs": [
    {
      "id": 1,
      "plant_id": 1,
      "activity_type": "watered",
      "logged_date": "2026-07-18T10:00:00Z",
      "notes": "Watered thoroughly",
      "photo_url": "https://..."
    }
  ],
  "scans": [ ... ]
}
```

---

### 9. `PUT /plants/{id}`

Update details of an existing plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Request Body** (all fields optional)
```json
{
  "nickname": "string",
  "location": "string",
  "health_score": 85,
  "status": "good | needs_attention"
}
```

**Response** `200 OK` - Returns updated plant object with preloaded `Species`.

---

### 10. `DELETE /plants/{id}`

Remove a plant from your collection.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Response** `200 OK`
```json
{ "message": "plant deleted successfully" }
```

---

### 11. `GET /plants/{id}/care-plan`

Retrieve the care plan for a specific plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Response** `200 OK`
```json
{
  "id": 1,
  "plant_id": 1,
  "watering_frequency_days": 7,
  "watering_amount": "200ml",
  "watering_method": "top watering",
  "watering_tips": "Keep soil consistently moist but not waterlogged",
  "light_requirement": "bright indirect light",
  "humidity_requirement": "60-80%"
}
```

---

### 12. `PUT /plants/{id}/care-plan`

Update the care plan for a specific plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Request Body** (all fields optional)
```json
{
  "watering_frequency_days": 7,
  "watering_amount": "200ml",
  "watering_method": "top watering",
  "watering_tips": "string",
  "light_requirement": "bright indirect light",
  "humidity_requirement": "60-80%"
}
```

**Response** `200 OK` - Returns updated `CarePlan` object.

---

### 13. `GET /plants/{id}/reminders`

Get all reminders for a specific plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Response** `200 OK`
```json
[
  {
    "id": 1,
    "plant_id": 1,
    "task_type": "water",
    "scheduled_time": "2026-07-25T10:00:00Z",
    "is_completed": false,
    "completed_at": "0001-01-01T00:00:00Z",
    "snooze_count": 0
  }
]
```

`task_type` values: `"water"`, `"fertilize"`, `"mist"`, `"rotate"`, `"repot"`

---

### 14. `GET /plants/{id}/activities`

Get all activity logs for a specific plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Response** `200 OK`
```json
[
  {
    "id": 1,
    "plant_id": 1,
    "activity_type": "watered",
    "logged_date": "2026-07-18T14:00:00Z",
    "notes": "Deep watered",
    "photo_url": "https://..."
  }
]
```

`activity_type` values: `"watered"`, `"fertilized"`, `"repotted"`, `"photo_node"`, `"milestone"`

---

### 15. `POST /plants/{id}/activities`

Record a new activity for a plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Request Body**
```json
{
  "activity_type": "watered (required)",
  "notes": "optional notes",
  "photo_url": "https://..."
}
```

Valid `activity_type` values: `"watered"`, `"fertilized"`, `"repotted"`, `"photo_node"`, `"milestone"`

**Response** `201 Created` - Returns `ActivityLog` object.

---

### 16. `GET /plants/{id}/growth`

Get all growth metrics for a specific plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Response** `200 OK`
```json
[
  {
    "id": 1,
    "plant_id": 1,
    "recorded_date": "2026-07-18T14:00:00Z",
    "height_cm": 25.5,
    "growth_rate_status": "moderate"
  }
]
```

`growth_rate_status` values: `"slow"`, `"moderate"`, `"fast"`

---

### 17. `POST /plants/{id}/growth`

Record a new growth measurement for a plant.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Plant ID

**Request Body**
```json
{
  "height_cm": 25.5,
  "growth_rate_status": "moderate"
}
```

**Response** `201 Created` - Returns `GrowthMetric` object.

---

### 18. `GET /reminders/today`

Get all incomplete reminders scheduled for today across all of the user's plants. Ordered by `scheduled_time` ascending.

**Auth:** Bearer JWT

**Response** `200 OK` - Array of `Reminder` objects (same shape as [endpoint 13](#13-get-plantsidreminders)).

---

### 19. `PUT /reminders/{id}`

Mark a reminder as completed or snooze it (+15 minutes).

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Reminder ID

**Request Body** (at least one field required)
```json
{
  "is_completed": true,
  "snooze": false
}
```

- `is_completed: true` - marks the reminder done and increments `total_task_done` on the user.
- `snooze: true` - postpones the reminder by 15 minutes and increments `snooze_count`.

**Response** `200 OK` - Returns updated `Reminder` object.

---

### 20. `GET /notifications`

Get notification preferences for the authenticated user. If none exist, creates defaults and returns them.

**Auth:** Bearer JWT

**Response** `200 OK`
```json
{
  "id": 1,
  "user_id": 1,
  "notification_enabled": true,
  "daily_summary_enabled": false,
  "sound_alert_enabled": true,
  "vibration_enabled": true,
  "preferred_notification_time": "0001-01-01T00:00:00Z",
  "default_snooze_duration_minute": 15
}
```

---

### 21. `PUT /notifications`

Update notification preferences for the authenticated user.

**Auth:** Bearer JWT

**Request Body** (all fields optional)
```json
{
  "notification_enabled": true,
  "daily_summary_enabled": false,
  "sound_alert_enabled": true,
  "vibration_enabled": true,
  "preferred_notification_time": "08:00:00",
  "default_snooze_duration_minute": 15
}
```

**Response** `200 OK` - Returns updated `Notification` object.

---

### 22. `POST /scan`

Upload a plant image for AI-powered identification. Returns a preview with confidence score for the user to confirm before saving.

**Auth:** Bearer JWT
**Content-Type:** `multipart/form-data`
**Rate Limit:** 5 scans per day per user

**Form Params**
| Field | Type | Required |
|-------|------|----------|
| `image` | file | Yes (max 10 MB) |

**Response** `200 OK`
```json
{
  "scan_id": 1,
  "retake": false,
  "confidence_score": 0.92,
  "identification": {
    "common_name": "Monstera Deliciosa",
    "scientific_name": "Monstera deliciosa",
    "family": "Araceae",
    "origin": "Central America",
    "confidence_score": 0.92,
    "health_assessment": "The plant appears healthy...",
    "detected_symptoms": ["none"],
    "primary_assessment": "Healthy Monstera Deliciosa",
    "treatment_steps": [],
    "care_recommendations": {
      "watering_frequency_days": 7,
      "watering_amount": "200ml",
      "watering_method": "top watering",
      "watering_tips": "Allow top inch of soil to dry",
      "light_requirement": "bright indirect light",
      "humidity_requirement": "60-80%"
    }
  },
  "captured_image_url": "https://res.cloudinary.com/..."
}
```

---

### 23. `GET /scan/{id}`

Retrieve details of a specific scan including the AI analysis result.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Scan ID

**Response** `200 OK`
```json
{
  "scan": {
    "id": 1,
    "user_id": 1,
    "plant_id": 1,
    "analysis_id": 1,
    "captured_image_url": "https://...",
    "retake": false,
    "confidence_score": 0.92,
    "json_identification_payload": "...",
    "plant": { ... }
  },
  "analysis": {
    "id": 1,
    "scan_id": 1,
    "ai_model_version": "gemini-2.0-flash",
    "confidence_score": 0.92,
    "analysis_type": "identification",
    "detected_symptoms": ["none"],
    "primary_assessment": "Healthy Monstera Deliciosa",
    "treatment_plan_steps": "No treatment needed",
    "metadata_payload": "..."
  }
}
```

---

### 24. `POST /scan/{id}/confirm`

Confirm an identified scan. Creates a Species, Plant, CarePlan, AiAnalysisResult, and initial Reminders (water, fertilize, rotate) in a single transaction.

**Auth:** Bearer JWT

**Path Params:** `id` (uint) - Scan ID

**Request Body**
```json
{
  "nickname": "My Monstera",
  "location": "Bedroom"
}
```

Both fields optional. `nickname` defaults to the identified `common_name`.

**Response** `200 OK`
```json
{
  "scan_id": 1,
  "plant": { ... },
  "species": {
    "common_name": "Monstera Deliciosa",
    "scientific_name": "Monstera deliciosa",
    "family": "Araceae",
    "origin": "Central America",
    "difficulty_level": "medium"
  },
  "analysis": { ... },
  "care_plan": { ... },
  "captured_image_url": "https://..."
}
```

---

### 25. `POST /diagnosis`

Upload a plant image for AI-powered disease diagnosis. Creates a chat session with an initial AI diagnosis message.

**Auth:** Bearer JWT
**Content-Type:** `multipart/form-data`
**Rate Limit:** 5 diagnoses per day per user

**Form Params**
| Field | Type | Required |
|-------|------|----------|
| `image` | file | Yes (max 10 MB) |

**Response** `200 OK`
```json
{
  "session_id": 1,
  "image_url": "https://res.cloudinary.com/...",
  "diagnosis": {
    "plant_type": "Monstera Deliciosa",
    "issue_description": "Leaf yellowing with brown spots",
    "severity": "medium",
    "causes": ["Overwatering", "Poor drainage"],
    "solutions": ["Reduce watering frequency", "Check drainage holes"],
    "prevention_tips": ["Water only when top inch is dry", "Ensure pot has drainage"]
  },
  "chat_history": [
    {
      "id": 1,
      "session_id": 1,
      "sender_type": "ai",
      "message_body": "Plant Type: Monstera Deliciosa\n\nIssue: Leaf yellowing..."
    }
  ],
  "scan_id": 1
}
```

---

### 26. `GET /diagnosis/{session_id}`

Retrieve the full chat history for a diagnosis session.

**Auth:** Bearer JWT

**Path Params:** `session_id` (uint) - Chat session ID

**Response** `200 OK`
```json
{
  "session_id": 1,
  "status": "active",
  "chat_history": [
    {
      "id": 1,
      "session_id": 1,
      "sender_type": "ai",
      "message_body": "Plant Type: Monstera...\n\nIssue: ..."
    },
    {
      "id": 2,
      "session_id": 1,
      "sender_type": "user",
      "message_body": "How often should I water it?"
    },
    {
      "id": 3,
      "session_id": 1,
      "sender_type": "ai",
      "message_body": "For a Monstera Deliciosa..."
    }
  ]
}
```

`status` values: `"active"`, `"archived"`
`sender_type` values: `"user"`, `"ai"`

---

### 27. `POST /diagnosis/{session_id}/chat`

Send a follow-up message in a diagnosis chat session. The message and full history are sent to Gemini which returns an AI reply. Both messages are persisted.

**Auth:** Bearer JWT

**Path Params:** `session_id` (uint) - Chat session ID (must be `"active"`)

**Request Body**
```json
{
  "message": "What fertilizer should I use? (required)"
}
```

**Response** `200 OK`
```json
{
  "session_id": 1,
  "user_message": {
    "id": 3,
    "session_id": 1,
    "sender_type": "user",
    "message_body": "What fertilizer should I use?"
  },
  "ai_message": {
    "id": 4,
    "session_id": 1,
    "sender_type": "ai",
    "message_body": "For Monstera deliciosa, I recommend..."
  },
  "chat_history": [
    { "id": 1, "sender_type": "ai", "message_body": "..." },
    { "id": 2, "sender_type": "user", "message_body": "..." },
    { "id": 3, "sender_type": "user", "message_body": "What fertilizer should I use?" },
    { "id": 4, "sender_type": "ai", "message_body": "For Monstera deliciosa, I recommend..." }
  ]
}
```

---

## Quick Reference Table

| # | Method | Endpoint | Auth | Rate Limit | Description |
|---|--------|----------|------|------------|-------------|
| 1 | `GET` | `/health` | No | - | Health check |
| 2 | `POST` | `/register` | No | - | Create account |
| 3 | `POST` | `/login` | No | - | Get tokens |
| 4 | `POST` | `/refresh` | No | - | Refresh tokens |
| 5 | `POST` | `/logout` | No | - | Revoke refresh token |
| 6 | `GET` | `/plants` | JWT | - | List all plants |
| 7 | `POST` | `/plants` | JWT | - | Add a plant |
| 8 | `GET` | `/plants/{id}` | JWT | - | Get plant details |
| 9 | `PUT` | `/plants/{id}` | JWT | - | Update a plant |
| 10 | `DELETE` | `/plants/{id}` | JWT | - | Delete a plant |
| 11 | `GET` | `/plants/{id}/care-plan` | JWT | - | Get care plan |
| 12 | `PUT` | `/plants/{id}/care-plan` | JWT | - | Update care plan |
| 13 | `GET` | `/plants/{id}/reminders` | JWT | - | Get plant reminders |
| 14 | `GET` | `/plants/{id}/activities` | JWT | - | Get activity logs |
| 15 | `POST` | `/plants/{id}/activities` | JWT | - | Log an activity |
| 16 | `GET` | `/plants/{id}/growth` | JWT | - | Get growth metrics |
| 17 | `POST` | `/plants/{id}/growth` | JWT | - | Log growth measurement |
| 18 | `GET` | `/reminders/today` | JWT | - | Today's reminders |
| 19 | `PUT` | `/reminders/{id}` | JWT | - | Complete/snooze reminder |
| 20 | `GET` | `/notifications` | JWT | - | Get notification prefs |
| 21 | `PUT` | `/notifications` | JWT | - | Update notification prefs |
| 22 | `POST` | `/scan` | JWT | 5/day | AI plant identification |
| 23 | `GET` | `/scan/{id}` | JWT | - | Get scan details |
| 24 | `POST` | `/scan/{id}/confirm` | JWT | - | Confirm scan -> create plant |
| 25 | `POST` | `/diagnosis` | JWT | 5/day | AI disease diagnosis |
| 26 | `GET` | `/diagnosis/{session_id}` | JWT | - | Get diagnosis chat history |
| 27 | `POST` | `/diagnosis/{session_id}/chat` | JWT | - | Send follow-up message |

---

## Enum Values Reference

| Field | Allowed Values |
|-------|----------------|
| `status` (Plant) | `"good"`, `"needs_attention"` |
| `task_type` (Reminder) | `"water"`, `"fertilize"`, `"mist"`, `"rotate"`, `"repot"` |
| `activity_type` (Activity) | `"watered"`, `"fertilized"`, `"repotted"`, `"photo_node"`, `"milestone"` |
| `growth_rate_status` (Growth) | `"slow"`, `"moderate"`, `"fast"` |
| `severity` (Diagnosis) | `"low"`, `"medium"`, `"high"`, `"critical"` |
| `sender_type` (Chat) | `"user"`, `"ai"` |
| `status` (Session) | `"active"`, `"archived"` |

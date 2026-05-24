# LE-STUDIO Integration Guide

## Overview
This document describes the complete integration of the le-studio ecosystem: **Backend API**, **Admin Frontend**, and **Client Frontend**, all connected to the live Render deployment.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        LE-STUDIO ECOSYSTEM                              │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────────────────┐         ┌──────────────────────┐
│  le-studio-admin     │         │  le-studio-front     │
│  (Next.js)           │         │  (React/Vite)        │
├──────────────────────┤         ├──────────────────────┤
│ Admin Panel UI       │         │ Client Portal UI     │
│ Users management     │         │ Login / Register     │
│ Coaches CRUD         │         │ Schedule browsing    │
│ Schedule builder     │         │ Bookings & payments  │
│ Analytics dashboard  │         │ Profile management   │
└──────────────────────┘         └──────────────────────┘
           │                              │
           │   HTTP (JWT Bearer)          │
           └──────────┬───────────────────┘
                      │
                      ▼
    ┌──────────────────────────────────────┐
    │   LE-STUDIO API (Render)              │
    │   https://le-studio-api.onrender.com  │
    │   (Go/Gin + PostgreSQL)               │
    ├──────────────────────────────────────┤
    │ Routes:                               │
    │ • /api/v1/auth/*       (Public)      │
    │ • /api/v1/users/*      (Protected)   │
    │ • /api/v1/bookings/*   (Protected)   │
    │ • /api/v1/admin/*      (Admin only)  │
    │ • /api/v1/coaches      (Public)      │
    │ • /api/v1/schedules    (Public)      │
    │ • /api/v1/pack-templates (Public)    │
    └──────────────────────────────────────┘
                      │
                      ▼
         ┌──────────────────────────┐
         │  Render PostgreSQL DB    │
         │  (Production Database)   │
         └──────────────────────────┘
```

---

## Environment Configuration

### Admin Frontend (`le-studio-admin`)
**File:** `.env.local`
```env
NEXT_PUBLIC_ADMIN_API_BASE=https://le-studio-api.onrender.com/api/v1
NEXT_PUBLIC_API_TIMEOUT=30000
```

### Client Frontend (`le-studio-front`)
**File:** `.env`
```env
VITE_API_BASE=https://le-studio-api.onrender.com/api/v1
VITE_API_TIMEOUT=30000
```

### Local Development
To test against local backend, update env files:
```env
# Admin (le-studio-admin/.env.local)
NEXT_PUBLIC_ADMIN_API_BASE=http://localhost:8080/api/v1

# Client (le-studio-front/.env)
VITE_API_BASE=http://localhost:8080/api/v1
```

---

## Authentication Flow

### JWT Token System
- **Access Token:** 15 minutes  
- **Refresh Token:** 7 days  
- **Format:** Bearer token in Authorization header  
- **Stored:** localStorage (`gym.admin.auth`, `gym.auth.token`)

### Admin Authentication
```
1. Admin enters credentials (admin@example.com / admin123)
   ↓
2. POST /api/v1/admin/auth/login
   ↓
3. Backend validates and returns access_token + refresh_token
   ↓
4. Admin frontend stores tokens in localStorage
   ↓
5. All subsequent requests include: Authorization: Bearer {access_token}
```

**Implementation:**  
File: `le-studio-admin/src/admin/auth/adminAuthService.js`
- `loginAdmin()` - Real API call (removed mock)
- `fetchAdminProfile()` - Calls GET /admin/auth/me
- `refreshAdminToken()` - Calls POST /admin/auth/refresh
- `logoutAdmin()` - Calls POST /admin/auth/logout

### User Authentication  
```
1. User enters credentials (email/password)
   ↓
2. POST /api/v1/auth/login
   ↓
3. Backend validates and returns access_token + refresh_token
   ↓
4. Client frontend stores tokens in localStorage
   ↓
5. All subsequent requests include: Authorization: Bearer {access_token}
```

**Implementation:**  
File: `le-studio-front/src/services/authService.js`
- `registerUser()` - Real API call (removed mock)
- `loginUser()` - Real API call (removed mock)
- `refreshUserToken()` - Calls POST /api/v1/auth/refresh
- `logoutUser()` - Calls POST /api/v1/auth/logout

---

## API Endpoints Summary

### Public Routes (No Authentication)
```
POST   /api/v1/auth/register              Register new user
POST   /api/v1/auth/login                 Login user
POST   /api/v1/auth/refresh               Refresh access token
POST   /api/v1/auth/logout                Logout user
POST   /api/v1/auth/guest                 Guest checkout flow

GET    /api/v1/coaches                    List all coaches
GET    /api/v1/coaches/:id                Get coach details
GET    /api/v1/pack-templates             List packages
GET    /api/v1/pack-templates/:id         Get package details
GET    /api/v1/training-types             List training types
GET    /api/v1/schedules                  List published schedules
GET    /api/v1/schedules/:id              Get schedule details
GET    /api/v1/schedules/:id/slots        Get weekly slots
```

### Protected Routes (User)
```
GET    /api/v1/users/me                   Current user profile
PATCH  /api/v1/users/me                   Update user profile
GET    /api/v1/user-packs                 List user's purchases
POST   /api/v1/user-packs                 Purchase package
GET    /api/v1/bookings                   List user's bookings
POST   /api/v1/bookings                   Create booking
GET    /api/v1/bookings/:id               Get booking details
PATCH  /api/v1/bookings/:id/cancel        Cancel booking
```

### Admin Routes (Admin + Authentication)
```
GET    /api/v1/admin/auth/me              Current admin profile
PATCH  /api/v1/admin/auth/me              Update admin profile
POST   /api/v1/admin/auth/login           Admin login
POST   /api/v1/admin/auth/logout          Admin logout
POST   /api/v1/admin/auth/refresh         Refresh admin token

GET    /api/v1/admin/users                List all users
GET    /api/v1/admin/users/:id            Get user details
POST   /api/v1/admin/users                Create user
PATCH  /api/v1/admin/users/:id            Update user

GET    /api/v1/admin/coaches              List coaches
POST   /api/v1/admin/coaches              Create coach
PATCH  /api/v1/admin/coaches/:id          Update coach

GET    /api/v1/admin/schedules            List schedules
POST   /api/v1/admin/schedules            Create schedule
POST   /api/v1/admin/schedules/:id/publish Publish schedule

GET    /api/v1/admin/stats/overview       Dashboard stats
GET    /api/v1/admin/bookings             List all bookings
GET    /api/v1/admin/attendance           Get attendance records
POST   /api/v1/admin/attendance           Mark attendance
```

---

## Frontend Implementation Status

### Admin Frontend (le-studio-admin)
| Component | Status | Notes |
|-----------|--------|-------|
| **Authentication** | ✅ Connected | Now calls real `/admin/auth/login` endpoint |
| **API Client** | ✅ Ready | `adminApiClient.js` properly configured |
| **Environment** | ✅ Configured | `.env.local` set with API base URL |
| **Login Page** | ✅ Working | Connects to real backend, tries `admin@example.com` / `admin123` |
| **Dashboard** | ✅ Connected | Uses `/admin/stats/overview` and `/admin/bookings` |
| **Users Page** | ✅ Connected | Uses `/admin/users` |
| **Coaches Page** | ✅ Connected | Uses `/admin/coaches` |
| **Schedule Page** | ✅ Connected | Uses `/admin/schedules` and `/admin/training-types` |
| **Bookings Page** | ✅ Connected | Uses `/admin/bookings` |

**What's Done:**
- Environment configuration
- Real JWT auth flow (no more fake tokens)
- API client setup with Bearer token
- Mock data fixtures removed from `src/admin/data/*`

**What's Next:**
- Add richer optimistic UI and inline validation for CRUD forms
- Add integration/e2e tests for admin CRUD flows

### Client Frontend (le-studio-front)
| Component | Status | Notes |
|-----------|--------|-------|
| **Authentication** | ✅ Connected | Now calls real `/auth/login` and `/auth/register` endpoints |
| **API Client** | ✅ Ready | `apiClient.js` properly configured |
| **Environment** | ✅ Configured | `.env` set with API base URL |
| **Login Page** | ✅ Working | Connects to real backend |
| **Register Page** | ✅ Working | Connects to real backend, creates real user accounts |
| **Profile Page** | ✅ Connected | Uses API-backed dashboard service (`/users/me`, `/bookings`) |
| **Classes Page** | ✅ Connected | Uses real API-backed content |
| **Coaches Page** | ✅ Connected | Uses real API-backed content |
| **Booking Page** | ✅ Connected | Uses live booking resources |

**What's Done:**
- Environment configuration
- Real JWT auth flow
- API client setup with Bearer token
- Login/Register connected to live API
- Mock dashboard service replaced with API dashboard service
- Automatic token refresh and request retry on `401`

**What's Next:**
- Add stronger UX fallback states for flaky network conditions
- Add e2e coverage for profile + booking lifecycle

---

## Testing Endpoints

### 1. Test Admin Login
```bash
curl -X POST https://le-studio-api.onrender.com/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'

# Response:
# {
#   "data": {
#     "access_token": "eyJhbGc...",
#     "refresh_token": "1a979...",
#     "admin": {
#       "ID": 1,
#       "Name": "Default Admin",
#       "Email": "admin@example.com"
#     }
#   },
#   "success": true
# }
```

### 2. Test User Registration
```bash
curl -X POST https://le-studio-api.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name":"John",
    "last_name":"Doe",
    "email":"john@example.com",
    "password":"SecurePass123",
    "phone":"+1234567890"
  }'
```

### 3. Test Protected Endpoint (Get Current User)
```bash
curl -X GET https://le-studio-api.onrender.com/api/v1/users/me \
  -H "Authorization: Bearer {access_token}"
```

### 4. Test Public Endpoint (List Coaches)
```bash
curl -X GET https://le-studio-api.onrender.com/api/v1/coaches
```

---

## File Modifications Summary

### Admin Frontend
- **Modified:** `le-studio-admin/.env.local` (created)
  - Added `NEXT_PUBLIC_ADMIN_API_BASE` environment variable

- **Modified:** `le-studio-admin/src/admin/auth/adminAuthService.js`
  - Replaced mock token generation with real API calls
  - `loginAdmin()` now calls `POST /admin/auth/login`
  - `fetchAdminProfile()` now calls `GET /admin/auth/me`
  - `refreshAdminToken()` now calls `POST /admin/auth/refresh`
  - `logoutAdmin()` now calls `POST /admin/auth/logout`

### Client Frontend
- **Created:** `le-studio-front/.env`
  - Added `VITE_API_BASE` environment variable

- **Modified:** `le-studio-front/src/services/authService.js`
  - Replaced all mock authentication with real API calls
  - `registerUser()` now calls `POST /auth/register`
  - `loginUser()` now calls `POST /auth/login`
  - `refreshUserToken()` now calls `POST /auth/refresh`
  - `logoutUser()` now calls `POST /auth/logout`

- **Modified:** `le-studio-front/src/services/dashboardService.js`
   - Replaced mock dashboard data with API-based dashboard loader
   - Reads `/users/me` and `/bookings` and maps into profile dashboard cards

- **Modified:** `le-studio-front/src/services/apiClient.js`
   - Added automatic token refresh flow on `401`
   - Retries the original request once after refresh

- **Modified:** `le-studio-front/src/context/AuthContext.jsx`
   - Persists refresh token alongside access token

- **Modified:** `le-studio-front/src/services/authStorage.js`
   - Added refresh token storage helpers

### CI/CD
- **Created:** `le-studio-api/.github/workflows/ci.yml`
   - Runs `go test ./...` and `go build ./cmd/server` on push/PR

- **Created:** `le-studio-api/.github/workflows/production-smoke.yml`
   - Runs scheduled/manual production smoke checks against live API

- **Created:** `le-studio-api/scripts/production_smoke_test.sh`
   - Public endpoint checks + optional protected admin check

- **Created:** `le-studio-front/.github/workflows/ci.yml`
   - Runs `npm ci`, `npm run lint`, `npm run build`

- **Created:** `le-studio-admin/.github/workflows/ci.yml`
   - Runs `npm ci`, `npm run lint`, `npm run build`

---

## CI/CD Setup

### Repository Secrets
Configure these GitHub repository secrets:

For `le-studio-api`:
- `API_BASE_URL` (example: `https://le-studio-api.onrender.com`)
- `ADMIN_EMAIL` (optional, for protected smoke checks)
- `ADMIN_PASSWORD` (optional, for protected smoke checks)

For `le-studio-front`:
- `VITE_API_BASE` (example: `https://le-studio-api.onrender.com/api/v1`)

For `le-studio-admin`:
- `NEXT_PUBLIC_ADMIN_API_BASE` (example: `https://le-studio-api.onrender.com/api/v1`)

### Workflows Added
- API CI: `.github/workflows/ci.yml`
- API production smoke: `.github/workflows/production-smoke.yml`
- Frontend CI: `.github/workflows/ci.yml`
- Admin CI: `.github/workflows/ci.yml`

### Local Verification Commands
```bash
cd le-studio-api && go test ./...
cd le-studio-api && BASE_URL="https://le-studio-api.onrender.com" bash scripts/production_smoke_test.sh
cd le-studio-front && npm ci && npm run build
cd le-studio-admin && npm ci && npm run build
```

## Next Steps
1. Add e2e tests for user booking lifecycle and admin management workflows.
2. Add branch protection requiring all CI workflows to pass before merge.
3. Add deploy checks for frontend hosts (Vercel previews + health checks).

---

## Troubleshooting

### Issue: "Cannot find module" or imports failing
**Solution:** Ensure `VITE_API_BASE` or `NEXT_PUBLIC_ADMIN_API_BASE` are set in `.env` or `.env.local`

### Issue: 401 Unauthorized on protected endpoints
**Solution:** 
- Check token is being stored correctly: `localStorage.getItem('gym.auth.token')`
- Verify token format: Should be a valid JWT starting with `eyJ...`
- Check Authorization header: Should be exactly `Bearer {token}`

### Issue: CORS errors
**Solution:**
- Backend CORS is configured to allow requests from any origin
- If still failing, verify request includes `Content-Type: application/json` header
- Check browser console for exact error message

### Issue: API returns 404
**Solution:**
- Verify endpoint path matches backend routes exactly
- Check URL has `/api/v1` prefix
- Ensure ID parameters in URLs are correct (integers, not null)

---

## CodeArchitecture

### Clean Architecture Principles Maintained

1. **Separation of Concerns**
   - Auth logic: `authService.js` / `adminAuthService.js`
   - API requests: `apiClient.js` / `adminApiClient.js`
   - Storage: `authStorage.js` / `adminAuthStorage.js`
   - Components: Only handle UI, delegate business logic to services

2. **Environment Configuration**
   - Backend URL managed via environment variables
   - No hardcoded URLs in source code
   - Easy to switch between dev/staging/production environments

3. **Error Handling**
   - All API calls wrapped in try-catch
   - User-friendly error messages
   - Proper HTTP status code handling

4. **State Management**
   - Auth state stored in localStorage
   - Auth context for React components (le-studio-front)
   - Clear separation between authenticated and unauthenticated states

---

## Performance Notes

1. **Token Caching:** Tokens cached in localStorage to avoid API calls on every request
2. **Request Size:** Minimal payload sizes - only essential fields included
3. **Connection Pool:** Backend uses connection pooling for database
4. **Rate Limiting:** Backend enforces 200 requests/minute per IP
5. **CORS:** Configured to minimize overhead while maintaining security

---

## Security Considerations

1. **JWT Tokens:**
   - Stored securely in localStorage (browser handles security)
   - Short expiry time (15 min access, 7 day refresh)
   - Refresh tokens can be revoked by backend

2. **HTTPS Only:**
   - All API calls use HTTPS in production
   - Tokens never transmitted over HTTP

3. **Password Security:**
   - Passwords hashed with bcrypt on backend
   - Never logged or exposed in API responses
   - Password fields stripped from responses

4. **CORS:**
   - Specific origins allowed for requests
   - Credentials included only when necessary

---

## Support & Questions

For issues or questions:
1. Check this document first
2. Review backend API logs at Render dashboard
3. Check browser console for client-side errors
4. Contact backend team for API-specific issues
5. Contact frontend team for UI/UX issues

---

**Last Updated:** May 24, 2026  
**Status:** Integration Complete - Ready for data endpoint replacement

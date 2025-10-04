# Bloggo REST API Documentation

## 📋 Table of Contents
- [API Overview](#api-overview)
- [Authentication](#authentication)
- [Base URL and Endpoints](#base-url-and-endpoints)
- [Common Response Formats](#common-response-formats)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [API Endpoints](#api-endpoints)
- [Pagination](#pagination)
- [Filtering and Sorting](#filtering-and-sorting)
- [API Versioning](#api-versioning)

## 🌐 API Overview

The Bloggo REST API provides programmatic access to all blog management functionality. It follows RESTful principles and uses JSON for data exchange.

### Key Features

- **RESTful Design**: Follows REST architectural principles
- **JSON Communication**: All requests and responses use JSON
- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Granular permission system
- **Rate Limiting**: Built-in protection against abuse
- **Comprehensive Error Handling**: Structured error responses
- **Pagination Support**: Efficient data retrieval
- **Filtering and Sorting**: Flexible data querying

### API Architecture

```mermaid
graph TB
    subgraph "API Layers"
        CLIENT[Client Applications]
        AUTH[Authentication Layer]
        RATE[Rate Limiting]
        API[API Endpoints]
        BIZ[Business Logic]
        DB[Database Layer]
    end

    CLIENT --> AUTH
    AUTH --> RATE
    RATE --> API
    API --> BIZ
    BIZ --> DB
```

## 🔐 Authentication

### JWT Token-Based Authentication

Bloggo uses JSON Web Tokens (JWT) for authentication. The API uses two types of tokens:

#### 1. Access Token
- **Purpose**: Access protected API endpoints
- **Lifetime**: 15 minutes (configurable)
- **Format**: Bearer token in Authorization header

#### 2. Refresh Token
- **Purpose**: Obtain new access tokens
- **Lifetime**: 7 days (configurable)
- **Storage**: Server-side storage with revocation support

### Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Auth
    participant DB

    Client->>API: POST /api/auth/login
    API->>Auth: Validate credentials
    Auth->>DB: Check user credentials
    DB-->>Auth: User data
    Auth-->>API: Generate tokens
    API-->>Client: Return access + refresh tokens

    Client->>API: GET /api/posts (with Authorization header)
    API->>Auth: Validate access token
    Auth-->>API: Token valid + user claims
    API-->>Client: Return requested data

    Client->>API: POST /api/auth/refresh
    API->>Auth: Validate refresh token
    Auth->>DB: Check refresh token
    DB-->>Auth: Token valid
    Auth-->>API: Generate new access token
    API-->>Client: Return new access token
```

### Making Authenticated Requests

```http
GET /api/posts HTTP/1.1
Host: localhost:8723
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

### Token Structure

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "rid": 123456789,
    "uid": 987654321,
    "iat": 1516239022,
    "exp": 1516242622
  }
}
```

**Payload Claims**:
- `rid`: Role ID (user's permission level)
- `uid`: User ID (unique identifier)
- `iat`: Issued at timestamp
- `exp`: Expiration timestamp

## 🌐 Base URL and Endpoints

### Development Environment
```
http://localhost:8723/api
```

### Production Environment
```
https://your-domain.com/api
```

### API Endpoint Categories

| Category | Base Path | Description |
|----------|-----------|-------------|
| Authentication | `/api/auth` | Login, logout, token management |
| Posts | `/api/posts` | Post CRUD operations |
| Users | `/api/users` | User management |
| Categories | `/api/categories` | Category management |
| Tags | `/api/tags` | Tag management |
| Statistics | `/api/statistics` | Analytics and metrics |
| Storage | `/api/storage` | File upload and management |
| Search | `/api/search` | Content search |
| Webhooks | `/api/webhooks` | Webhook management |
| Audit | `/api/audit` | Audit log access |
| Dashboard | `/api/dashboard` | Admin dashboard data |
| Health | `/health` | System health checks |

## 📄 Common Response Formats

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data varies by endpoint
  },
  "meta": {
    "timestamp": "2025-10-04T12:00:00Z",
    "version": "1.0.0"
  }
}
```

### Paginated Response

```json
{
  "success": true,
  "data": [
    // Array of items
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8,
    "hasNext": true,
    "hasPrev": false
  },
  "meta": {
    "timestamp": "2025-10-04T12:00:00Z",
    "version": "1.0.0"
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": {
      "field": "title",
      "reason": "Title is required and must be between 1 and 500 characters"
    }
  },
  "meta": {
    "timestamp": "2025-10-04T12:00:00Z",
    "version": "1.0.0",
    "requestId": "req_123456789"
  }
}
```

## ⚠️ Error Handling

### HTTP Status Codes

| Status Code | Meaning | Usage |
|-------------|---------|-------|
| `200 OK` | Success | Request completed successfully |
| `201 Created` | Created | Resource created successfully |
| `204 No Content` | No Content | Request successful, no content returned |
| `400 Bad Request` | Bad Request | Invalid request parameters |
| `401 Unauthorized` | Unauthorized | Authentication required or failed |
| `403 Forbidden` | Forbidden | Insufficient permissions |
| `404 Not Found` | Not Found | Resource not found |
| `409 Conflict` | Conflict | Resource conflict (duplicate) |
| `422 Unprocessable Entity` | Validation Error | Request validation failed |
| `429 Too Many Requests` | Rate Limited | Too many requests |
| `500 Internal Server Error` | Server Error | Internal server error |

### Error Codes

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `VALIDATION_ERROR` | Request validation failed | 400 |
| `UNAUTHORIZED` | Authentication required | 401 |
| `FORBIDDEN` | Insufficient permissions | 403 |
| `NOT_FOUND` | Resource not found | 404 |
| `CONFLICT` | Resource already exists | 409 |
| `RATE_LIMITED` | Too many requests | 429 |
| `INTERNAL_ERROR` | Internal server error | 500 |
| `DATABASE_ERROR` | Database operation failed | 500 |
| `FILE_ERROR` | File operation failed | 500 |

### Common Error Scenarios

#### Authentication Errors
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

#### Validation Errors
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      {
        "field": "title",
        "message": "Title is required"
      },
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

#### Permission Errors
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Insufficient permissions to perform this action"
  }
}
```

## 🚦 Rate Limiting

### Global Rate Limiting

- **Default Limit**: 100 requests per minute per IP
- **Window**: 60 minutes sliding window
- **Response Headers**:
  - `X-RateLimit-Limit`: Request limit per window
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time when limit resets

### Rate Limit Response

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1633316400

{
  "success": false,
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests in a short period"
  }
}
```

### Endpoint-Specific Limits

Some endpoints have custom rate limits:

| Endpoint | Limit | Purpose |
|----------|-------|---------|
| `POST /api/auth/login` | 5/min | Prevent brute force |
| `POST /api/auth/register` | 3/min | Prevent spam registration |
| `POST /api/storage/upload` | 10/min | Prevent upload abuse |
| `POST /api/posts` | 20/min | Reasonable content creation |

## 🔗 API Endpoints

### Authentication Endpoints

#### POST /api/auth/login
Authenticate user and return tokens.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "userpassword"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "def50200e3b7b8f4...",
    "expiresIn": 900,
    "user": {
      "id": 123,
      "name": "John Doe",
      "email": "user@example.com",
      "role": {
        "id": 2,
        "name": "author"
      }
    }
  }
}
```

#### POST /api/auth/refresh
Refresh access token using refresh token.

**Request Body**:
```json
{
  "refreshToken": "def50200e3b7b8f4..."
}
```

#### POST /api/auth/logout
Logout user and invalidate refresh token.

**Request Headers**:
```
Authorization: Bearer {access_token}
```

#### POST /api/auth/logout-all
Logout user from all devices.

**Request Headers**:
```
Authorization: Bearer {access_token}
```

### Post Endpoints

#### GET /api/posts
Retrieve paginated list of posts.

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 20, max: 100)
- `status` (string, optional): Filter by status (draft, published, archived)
- `authorId` (integer, optional): Filter by author
- `categoryId` (integer, optional): Filter by category
- `tagId` (integer, optional): Filter by tag
- `search` (string, optional): Search in title and content
- `sort` (string, optional): Sort field (publishedAt, createdAt, viewCount)
- `order` (string, optional): Sort order (asc, desc)

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Sample Post",
      "slug": "sample-post",
      "excerpt": "This is a sample post excerpt...",
      "coverImage": "/uploads/cover/sample.jpg",
      "author": {
        "id": 1,
        "name": "John Doe",
        "avatar": "/uploads/avatar/john.jpg"
      },
      "status": "published",
      "publishedAt": "2025-10-01T10:00:00Z",
      "viewCount": 150,
      "readTime": 5,
      "categories": [
        {
          "id": 1,
          "name": "Technology",
          "slug": "technology"
        }
      ],
      "tags": [
        {
          "id": 1,
          "name": "go",
          "slug": "go"
        }
      ],
      "createdAt": "2025-10-01T09:00:00Z",
      "updatedAt": "2025-10-01T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8,
    "hasNext": true,
    "hasPrev": false
  }
}
```

#### GET /api/posts/{id}
Retrieve specific post by ID.

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Sample Post",
    "slug": "sample-post",
    "content": "# Sample Post\n\nThis is the full content...",
    "excerpt": "This is a sample post excerpt...",
    "coverImage": "/uploads/cover/sample.jpg",
    "author": {
      "id": 1,
      "name": "John Doe",
      "avatar": "/uploads/avatar/john.jpg"
    },
    "status": "published",
    "publishedAt": "2025-10-01T10:00:00Z",
    "viewCount": 150,
    "readTime": 5,
    "categories": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      }
    ],
    "tags": [
      {
        "id": 1,
        "name": "go",
        "slug": "go"
      }
    ],
    "createdAt": "2025-10-01T09:00:00Z",
    "updatedAt": "2025-10-01T10:00:00Z"
  }
}
```

#### POST /api/posts
Create new post.

**Request Headers**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "title": "New Post Title",
  "content": "# New Post\n\nThis is the content...",
  "excerpt": "Post excerpt...",
  "coverImage": "/uploads/cover/new-post.jpg",
  "status": "draft",
  "categoryIds": [1, 2],
  "tagNames": ["go", "web-development"]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 151,
    "title": "New Post Title",
    "slug": "new-post-title",
    "status": "draft",
    "createdAt": "2025-10-04T12:00:00Z"
  }
}
```

#### PUT /api/posts/{id}
Update existing post.

**Request Headers**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "title": "Updated Post Title",
  "content": "# Updated Post\n\nUpdated content...",
  "excerpt": "Updated excerpt...",
  "coverImage": "/uploads/cover/updated.jpg",
  "status": "published",
  "categoryIds": [1],
  "tagNames": ["go"]
}
```

#### DELETE /api/posts/{id}
Delete post (soft delete).

**Request Headers**:
```
Authorization: Bearer {access_token}
```

#### POST /api/posts/{id}/versions
Create new version of post.

**Request Headers**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "title": "New Version Title",
  "content": "# New Version\n\nContent for new version...",
  "excerpt": "New version excerpt...",
  "coverImage": "/uploads/cover/new-version.jpg"
}
```

#### PUT /api/posts/{id}/versions/{version}/approve
Approve post version.

**Request Headers**:
```
Authorization: Bearer {access_token}
```

### User Endpoints

#### POST /api/users/register
Register new user account.

**Request Body**:
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "securepassword"
}
```

#### GET /api/users/profile
Get current user profile.

**Request Headers**:
```
Authorization: Bearer {access_token}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "avatar": "/uploads/avatar/john.jpg",
    "role": {
      "id": 2,
      "name": "author",
      "permissions": [
        "posts.create",
        "posts.edit",
        "posts.delete"
      ]
    },
    "createdAt": "2025-09-01T10:00:00Z",
    "updatedAt": "2025-10-01T12:00:00Z"
  }
}
```

#### PUT /api/users/profile
Update current user profile.

**Request Headers**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "Updated Name",
  "avatar": "/uploads/avatar/updated.jpg"
}
```

#### GET /api/users
List all users (admin only).

**Request Headers**:
```
Authorization: Bearer {access_token}
```

**Query Parameters**:
- `page` (integer, optional): Page number
- `limit` (integer, optional): Items per page
- `roleId` (integer, optional): Filter by role
- `search` (string, optional): Search by name or email

### Category Endpoints

#### GET /api/categories
List all categories.

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Technology",
      "slug": "technology",
      "description": "Technology related posts",
      "parentId": null,
      "postCount": 25,
      "createdAt": "2025-09-01T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Go Programming",
      "slug": "go-programming",
      "description": "Posts about Go language",
      "parentId": 1,
      "postCount": 15,
      "createdAt": "2025-09-01T11:00:00Z"
    }
  ]
}
```

#### POST /api/categories
Create new category (admin/editor only).

**Request Headers**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "New Category",
  "description": "Category description",
  "parentId": 1
}
```

#### PUT /api/categories/{id}
Update category (admin/editor only).

#### DELETE /api/categories/{id}
Delete category (admin only).

### Tag Endpoints

#### GET /api/tags
List all tags.

**Query Parameters**:
- `search` (string, optional): Search tag names
- `limit` (integer, optional): Maximum number of tags

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "go",
      "slug": "go",
      "postCount": 25,
      "createdAt": "2025-09-01T10:00:00Z"
    }
  ]
}
```

#### POST /api/tags
Create new tag (admin/editor only).

### Statistics Endpoints

#### GET /api/statistics/posts/{id}
Get detailed statistics for a specific post.

**Response**:
```json
{
  "success": true,
  "data": {
    "postId": 1,
    "viewCount": 150,
    "readTime": 5,
    "desktopViews": 120,
    "mobileViews": 30,
    "tabletViews": 0,
    "viewsByDate": [
      {
        "date": "2025-10-01",
        "views": 25
      },
      {
        "date": "2025-10-02",
        "views": 30
      }
    ],
    "deviceBreakdown": {
      "desktop": 80.0,
      "mobile": 20.0,
      "tablet": 0.0
    }
  }
}
```

#### GET /api/statistics/overview
Get overall statistics (admin only).

**Response**:
```json
{
  "success": true,
  "data": {
    "totalPosts": 150,
    "publishedPosts": 120,
    "totalUsers": 25,
    "totalViews": 5000,
    "viewsThisMonth": 1200,
    "topPosts": [
      {
        "id": 1,
        "title": "Most Popular Post",
        "viewCount": 500
      }
    ],
    "growthMetrics": {
      "postGrowth": 15.5,
      "viewGrowth": 25.2,
      "userGrowth": 8.3
    }
  }
}
```

#### POST /api/statistics/posts/{id}/view
Track post view (used internally).

### Storage Endpoints

#### POST /api/storage/upload
Upload file to storage.

**Request Headers**:
```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data
```

**Request Body**:
```
file: [binary file data]
```

**Response**:
```json
{
  "success": true,
  "data": {
    "filename": "image_123456789.jpg",
    "originalName": "my-image.jpg",
    "filePath": "/uploads/2025/10/image_123456789.jpg",
    "fileSize": 1024000,
    "mimeType": "image/jpeg",
    "url": "http://localhost:8723/uploads/2025/10/image_123456789.jpg"
  }
}
```

#### GET /api/storage/{filename}
Serve uploaded file.

#### DELETE /api/storage/{filename}
Delete uploaded file.

### Search Endpoints

#### GET /api/search/posts
Search posts with advanced filtering.

**Query Parameters**:
- `q` (string, required): Search query
- `categoryIds` (string, optional): Comma-separated category IDs
- `tagIds` (string, optional): Comma-separated tag IDs
- `authorId` (integer, optional): Author ID
- `dateFrom` (string, optional): Date from (YYYY-MM-DD)
- `dateTo` (string, optional): Date to (YYYY-MM-DD)
- `page` (integer, optional): Page number
- `limit` (integer, optional): Results per page

**Response**:
```json
{
  "success": true,
  "data": {
    "query": "go programming",
    "results": [
      {
        "id": 1,
        "title": "Go Programming Tutorial",
        "excerpt": "Learn Go programming language...",
        "contentHighlight": "...<mark>Go</mark> <mark>programming</mark>...",
        "author": "John Doe",
        "publishedAt": "2025-10-01T10:00:00Z",
        "relevanceScore": 0.95
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 15,
      "totalPages": 1
    },
    "facets": {
      "categories": [
        {
          "id": 1,
          "name": "Technology",
          "count": 10
        }
      ],
      "tags": [
        {
          "id": 1,
          "name": "go",
          "count": 15
        }
      ]
    }
  }
}
```

#### GET /api/search/suggestions
Get search suggestions.

**Query Parameters**:
- `q` (string, required): Partial query

**Response**:
```json
{
  "success": true,
  "data": {
    "suggestions": [
      "go programming",
      "go tutorials",
      "go web development"
    ]
  }
}
```

### Health Endpoints

#### GET /health
Basic health check.

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2025-10-04T12:00:00Z",
    "uptime": 86400,
    "version": "1.0.0"
  }
}
```

#### GET /health/detailed
Detailed health check (admin only).

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2025-10-04T12:00:00Z",
    "uptime": 86400,
    "version": "1.0.0",
    "checks": {
      "database": {
        "status": "healthy",
        "responseTime": 5,
        "details": "Connection successful"
      },
      "storage": {
        "status": "healthy",
        "availableSpace": "10GB",
        "details": "File system accessible"
      },
      "memory": {
        "status": "healthy",
        "used": "256MB",
        "available": "768MB",
        "percentage": 25.0
      }
    }
  }
}
```

## 📄 Pagination

### Pagination Parameters

All list endpoints support pagination with these parameters:

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `page` | integer | 1 | - | Page number (1-based) |
| `limit` | integer | 20 | 100 | Items per page |

### Pagination Response Format

```json
{
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8,
    "hasNext": true,
    "hasPrev": false
  }
}
```

### Pagination Headers

```http
X-Total-Count: 150
X-Page-Count: 8
X-Current-Page: 1
X-Per-Page: 20
```

## 🔍 Filtering and Sorting

### Common Filter Parameters

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `search` | string | `q=go programming` | Search query |
| `status` | string | `status=published` | Filter by status |
| `authorId` | integer | `authorId=1` | Filter by author |
| `categoryId` | integer | `categoryId=1` | Filter by category |
| `tagId` | integer | `tagId=1` | Filter by tag |
| `dateFrom` | string | `dateFrom=2025-10-01` | Filter from date |
| `dateTo` | string | `dateTo=2025-10-31` | Filter to date |

### Sorting Parameters

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `sort` | string | `sort=publishedAt` | Sort field |
| `order` | string | `order=desc` | Sort order (asc/desc) |

### Available Sort Fields

- `publishedAt`: Publication date
- `createdAt`: Creation date
- `updatedAt`: Last update date
- `title`: Title (alphabetical)
- `viewCount`: View count
- `name`: Name (for users, categories, tags)

## 📊 API Versioning

### Current Version: v1.0.0

The API uses semantic versioning. Version is included in response metadata:

```json
{
  "meta": {
    "version": "1.0.0"
  }
}
```

### Versioning Strategy

- **Major Version**: Breaking changes
- **Minor Version**: New features, backward compatible
- **Patch Version**: Bug fixes, backward compatible

### Backward Compatibility

- Fields are only added, not removed
- Enum values are only added, not removed
- Required fields remain required
- Optional fields remain optional

### Deprecation Policy

- Deprecated features are announced 3 months in advance
- Deprecated endpoints continue to work for 6 months
- Deprecation warnings are included in response headers

```http
X-API-Deprecated: true
X-API-Deprecation-Date: 2026-01-01
X-API-Alternative-Endpoint: /api/v2/posts
```

---

**Document Version**: 1.0.0
**Last Updated**: October 4, 2025
**Author**: Bloggo Development Team
**Reviewers**: API Architecture Committee
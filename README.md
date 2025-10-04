<div align="center">
  <img src="./bloggo.webp" width="196" heigth="196" />

# Bloggo: Blog CMS

![Go](https://img.shields.io/badge/Go-blue.svg?logo=go&logoColor=white&style=for-the-badge)
![SQLite](https://img.shields.io/badge/SQLite-gold.svg?logo=sqlite&logoColor=black&style=for-the-badge)
![GitHub License](https://img.shields.io/github/license/Elagoht/bloggo?style=for-the-badge)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Elagoht/bloggo/release.yaml?style=for-the-badge)

</div>

A modern, feature-rich blog management system built with Go. Bloggo provides a complete content management solution with post versioning, role-based permissions, statistics tracking, and comprehensive analytics.

## ✨ Features

### 📝 Content Management

- **Post Creation & Editing** - Rich content management with cover images and metadata
- **Version Control** - Track and manage multiple versions of posts with approval workflow
- **Draft System** - Save drafts and publish when ready
- **Categories & Tags** - Organize content with flexible categorization
- **Cover Images** - Upload and serve automatically optimized cover images (resized and saved in WebP format)

### 👥 User Management

- **Role-Based Access Control** - Admin, Editor, and Author roles with granular permissions
- **User Authentication** - JWT-based secure authentication
- **Profile Management** - User profiles with avatar support
- **Permission System** - Fine-grained permission control for different actions

### 📊 Analytics & Statistics

- **View Tracking** - Track timestamp based post views with user agent analysis
- **Device Analytics** - Desktop vs mobile traffic analysis
- **Operating System Detection** - Detailed OS statistics
- **Browser Analytics** - Browser usage statistics
- **Read Time Calculation** - Automatic reading time estimation
- **Performance Optimized** - Denormalized read counts for fast retrieval

### 🔧 Technical Features

- **Modular Architecture** - Clean, maintainable codebase with module-based design
- **SQLite Database** - Lightweight, embedded database with full SQL support
- **File Storage** - Organized file storage system for uploads
- **Rate Limiting** - Built-in rate limiting for API protection
- **Caching Headers** - Optimized caching for static assets

## 🚀 Quick Start

### Prerequisites

Nothing! Bloggo is pre-compiled via github workflows and does not requires any additional software. Just download and run the exacurable from [Releases Page](/releases)

## 📁 Project Structure

```
bloggo/
├── cli/                    # Application entry point
├── internal/
│   ├── app/               # Application core
│   ├── config/            # Configuration management
│   ├── db/                # Database initialization and queries
│   ├── infrastructure/    # Infrastructure components
│   │   ├── bucket/        # File storage abstraction
│   │   ├── permissions/   # Permission system
│   │   └── tokens/        # Token management
│   ├── middleware/        # HTTP middlewares
│   ├── module/            # Feature modules
│   │   ├── category/      # Category management
│   │   ├── post/          # Post management
│   │   ├── session/       # Authentication
│   │   ├── statistics/    # Analytics
│   │   ├── storage/       # File serving
│   │   ├── tag/           # Tag management
│   │   └── user/          # User management
│   └── utils/             # Utility packages
├── uploads/               # File storage directory
├── .env                  # Environment variables (create from .env.example)
├── bloggo.sqlite         # SQLite database
├── go.mod                # Go module definition
└── LICENSE               # GPLv3 License
```

## 🔧 Configuration

Configuration is managed through environment variables. Create a `.env` file from `.env.example`:

```bash
cp .env.example .env
```

Key environment variables include:

- **PORT** - Server port (default: 8723)
- **JWT_SECRET** - JWT signing secret (required, min 32 characters)
- **ACCESS_TOKEN_DURATION** - Access token lifetime in seconds (default: 900)
- **REFRESH_TOKEN_DURATION** - Refresh token lifetime in seconds (default: 604800)
- **GEMINI_API_KEY** - Google Gemini API key (optional, for AI features)
- **TRUSTED_FRONTEND_KEY** - Key for trusted frontend requests (required, min 32 characters)

## 🗄️ Database Schema

The application uses SQLite with the following main entities:

- **Users** - User accounts and profiles
- **Roles & Permissions** - Role-based access control
- **Posts & Post Versions** - Content with version control
- **Categories & Tags** - Content organization
- **Post Views** - Analytics and tracking
- **Sessions** - Authentication management

## 🛡️ Security Features

- **JWT Authentication** - Secure token-based authentication
- **Password Hashing** - bcrypt password hashing
- **Rate Limiting** - Protection against abuse
- **Input Validation** - Comprehensive input validation
- **SQL Injection Protection** - Parameterized queries

## 📈 Performance

- **Optimized Queries** - Efficient database queries with proper indexing
- **Denormalized Counters** - Fast read counts without expensive aggregations
- **Caching Headers** - Proper HTTP caching for static assets
- **Connection Pooling** - Efficient database connection management
- **Modular Loading** - Lazy loading of modules and dependencies

## 🧪 Development

### Code Structure

- **Modular Design** - Each feature is a self-contained module
- **Repository Pattern** - Clean separation of data access
- **Service Layer** - Business logic abstraction
- **Handler Layer** - HTTP request handling

### Dev Dependencies

- **Chi Router** - Fast HTTP router
- **SQLite3** - Embedded database
- **JWT** - JSON Web Tokens for authentication
- **Validator** - Request validation
- **WebP** - Image processing and optimization

## 📜 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📞 Support

If you encounter any issues or have questions, please open an issue in the repository.

---

**Bloggo** - A modern blog management system built with ❤️ and Go

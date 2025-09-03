# Tubely - Video Sharing Platform

A fully-fledged video sharing API server built from scratch in Go. A YouTube-like platform where users can upload, manage, and share videos with thumbnail support and cloud storage integration. Hit the project with a star if you find it useful.

## Fork Information

This project is forked from the [Boot.dev Learn File Storage S3 Golang Starter](https://github.com/bootdotdev/learn-file-storage-s3-golang-starter) repository. The original repo serves as starter code for the "Learn File Servers and CDNs with S3 and CloudFront" course on Boot.dev.

## Motivation

This project demonstrates building a complete video sharing platform from the ground up using Go's standard library with cloud storage integration. It showcases modern Go web development practices including file upload handling, video processing, AWS S3 integration, authentication, database management, and clean architecture patterns.

## Goal

The goal with Tubely is to provide a complete example of a production-ready Go web API for video sharing that includes:

- User authentication and authorization with JWT tokens
- Video upload and processing with FFMPEG
- Thumbnail upload and management
- AWS S3 cloud storage integration
- Video aspect ratio detection and organization
- RESTful API design with proper HTTP methods and status codes
- Database integration with SQLite
- Secure password hashing with bcrypt
- Refresh token management
- Clean, maintainable code structure
- Comprehensive error handling
- Video optimization for fast streaming

## ⚙️ Installation

### Prerequisites

- Go 1.23 or higher
- SQLite 3 database
- FFMPEG (both `ffmpeg` and `ffprobe` required in PATH)
- AWS CLI configured with credentials
- Git

### Setup

1. Clone the repository:
```bash
git clone https://github.com/dmitriy-zverev/tubely.git
cd tubely
```

2. Install dependencies:
```bash
go mod download
```

3. Install FFMPEG:
```bash
# Linux
sudo apt update
sudo apt install ffmpeg

# macOS
brew update
brew install ffmpeg
```

4. Install SQLite 3:
```bash
# Linux
sudo apt update
sudo apt install sqlite3

# macOS
brew update
brew install sqlite3
```

5. Set up your environment variables by creating a `.env` file:
```bash
cp .env.example .env
```

Update the `.env` file with your configuration:
```env
DB_PATH="./tubely.db"
JWT_SECRET="your-secret-key-here"
PLATFORM="dev"
FILEPATH_ROOT="./app"
ASSETS_ROOT="./assets"
S3_BUCKET="your-s3-bucket-name"
S3_REGION="us-east-2"
S3_CF_DISTRO="your-cloudfront-distribution"
PORT="8091"
```

6. Configure AWS credentials:
```bash
aws configure
```

7. Download sample videos and images:
```bash
./samplesdownload.sh
```

8. Start the server:
```bash
go run .
```

The server will be available at `http://localhost:8091/app/`

## 🚀 Quick Start

### Creating a User
```bash
curl -X POST http://localhost:8091/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### User Login
```bash
curl -X POST http://localhost:8091/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### Creating a Video Metadata
```bash
curl -X POST http://localhost:8091/api/videos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title": "My First Video", "description": "This is my first video upload!"}'
```

### Uploading a Video File
```bash
curl -X POST http://localhost:8091/api/video_upload/{videoID} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "video=@path/to/your/video.mp4"
```

### Uploading a Thumbnail
```bash
curl -X POST http://localhost:8091/api/thumbnail_upload/{videoID} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "thumbnail=@path/to/your/thumbnail.jpg"
```

### Getting All Videos
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8091/api/videos
```

## 📁 Project Structure

```
tubely/
├── main.go                    # Application entry point with server configuration
├── handler_*.go              # HTTP handlers for different endpoints
├── assets.go                 # Asset management utilities
├── cache.go                  # Caching middleware
├── json.go                   # JSON response utilities
├── reset.go                  # Database reset functionality
├── internal/
│   ├── auth/                 # Authentication and JWT handling
│   ├── database/             # Database models and queries
│   │   ├── database.go       # Database client
│   │   ├── users.go          # User management
│   │   ├── videos.go         # Video management
│   │   └── refresh_tokens.go # Token management
├── app/                      # Frontend web application
│   ├── index.html           # Main web interface
│   ├── app.js               # Frontend JavaScript
│   └── styles.css           # Styling
├── assets/                   # Local file storage directory
├── samples/                  # Sample videos and images
└── README.md
```

## 🔧 API Endpoints

### Authentication
- `POST /api/users` - Create a new user
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh JWT token
- `POST /api/revoke` - Revoke refresh token

### Videos
- `GET /api/videos` - Get all user videos
- `GET /api/videos/{videoID}` - Get a specific video
- `POST /api/videos` - Create video metadata (requires authentication)
- `DELETE /api/videos/{videoID}` - Delete a video (requires authentication)
- `POST /api/video_upload/{videoID}` - Upload video file (requires authentication)
- `POST /api/thumbnail_upload/{videoID}` - Upload thumbnail (requires authentication)

### Admin
- `POST /admin/reset` - Reset database (dev environment only)

### Static Files
- `/app/*` - Frontend web application
- `/assets/*` - Uploaded video and thumbnail files

## 🔒 Authentication

Tubely uses JWT (JSON Web Tokens) for authentication. After logging in, include the token in the Authorization header:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

Refresh tokens are also supported for maintaining long-term sessions.

## 🗄️ Database

The project uses SQLite with custom SQL queries for data management. The database schema includes:

- Users table with hashed passwords and email authentication
- Videos table with metadata, file URLs, and user relationships
- Refresh tokens for secure authentication management

## 🎥 Video Processing

Tubely includes advanced video processing features:

- **Aspect Ratio Detection**: Automatically detects and categorizes videos as 16:9 (landscape), 9:16 (portrait), or other ratios
- **Video Optimization**: Uses FFMPEG to optimize videos for fast streaming with "fast start" processing
- **Cloud Storage**: Automatically uploads processed videos to AWS S3 with organized folder structure
- **Thumbnail Support**: Separate thumbnail upload and management system
- **File Validation**: Ensures only MP4 video files are accepted

## ☁️ Cloud Integration

### AWS S3 Storage
Videos are automatically organized in S3 buckets by aspect ratio:
- `landscape/` - 16:9 aspect ratio videos
- `portrait/` - 9:16 aspect ratio videos  
- `other/` - All other aspect ratios

### CloudFront CDN
Integration ready for CloudFront distribution for global content delivery.

## 🧪 Testing

Run the tests:
```bash
go test ./...
```

## 📝 Configuration

The application can be configured using environment variables:

- `DB_PATH` - SQLite database file path
- `JWT_SECRET` - Secret key for JWT signing
- `PLATFORM` - Environment (dev/prod)
- `FILEPATH_ROOT` - Frontend application root directory
- `ASSETS_ROOT` - Local assets storage directory
- `S3_BUCKET` - AWS S3 bucket name
- `S3_REGION` - AWS S3 region
- `S3_CF_DISTRO` - CloudFront distribution ID
- `PORT` - Server port (default: 8091)

## 🚀 Deployment

The application is designed to be easily deployable to various cloud platforms. Ensure all environment variables are properly set and AWS credentials are configured in your production environment.

## 🛠️ Technologies Used

- **Go** - Programming language
- **SQLite** - Database
- **AWS S3** - Cloud storage
- **FFMPEG** - Video processing
- **JWT** - Authentication tokens
- **bcrypt** - Password hashing
- **UUID** - Unique identifiers
- **Standard Library** - Minimal external dependencies

## 💻 What We Built

This implementation extends the original Boot.dev starter with several key enhancements:

### Core Features Added
1. **Complete Video Upload Pipeline**: Full video upload, processing, and storage workflow
2. **Advanced Video Processing**: FFMPEG integration for video optimization and aspect ratio detection
3. **Cloud Storage Integration**: AWS S3 upload with organized folder structure
4. **Thumbnail Management**: Separate thumbnail upload and management system
5. **User Authentication**: Complete JWT-based authentication with refresh tokens
6. **Database Management**: SQLite integration with proper data models
7. **Frontend Interface**: Web application for video management

### Technical Improvements
1. **Structured Architecture**: Clean separation of concerns with internal packages
2. **Error Handling**: Comprehensive error handling throughout the application
3. **Security**: Secure file upload validation and authentication
4. **Performance**: Video optimization for streaming and caching middleware
5. **Scalability**: Cloud-ready architecture with S3 integration

### Why These Choices
- **SQLite**: Chosen for simplicity and ease of deployment while maintaining ACID compliance
- **AWS S3**: Industry-standard cloud storage with excellent scalability and reliability
- **FFMPEG**: Industry-standard video processing for optimization and metadata extraction
- **JWT**: Stateless authentication perfect for API-based applications
- **Go Standard Library**: Minimal dependencies for better security and maintainability

## 💬 Contact

- GitHub: [@dmitriy-zverev](https://github.com/dmitriy-zverev)
- Submit an issue here on GitHub

## 👏 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is fully open source. Feel free to use it as you wish.

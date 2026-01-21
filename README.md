# Personal Blog Generator

A modern, static site generator for personal blogs built with Go. Features a clean admin interface for content management, automatic static site generation, and easy deployment.

## Features

- 📝 **Content Management**: Create and manage blog posts, portfolio items, and static pages
- 🎨 **Modern Admin UI**: Clean, responsive interface built with Tailwind CSS
- ⚡ **Static Site Generation**: Fast, SEO-friendly static HTML output
- 🚀 **Easy Deployment**: Simple deployment with nginx and systemd

## Quick Start

### Local Development

1. **Prerequisites**
   - Go 1.22 or later
   - Node.js 16+ (for Tailwind CSS)
   - SQLite3

2. **Installation**
   ```bash
   # Clone the repository
   git clone <repository-url>
   cd personal-blog-generator

   # Install Go dependencies
   go mod tidy

   # Install Node.js dependencies
   npm install
   ```

3. **Environment Setup**
   ```bash
   # Create configuration directory
   mkdir -p ~/.personal-blog-generator

   # Create .env file with your settings
   cat > ~/.personal-blog-generator/.env << EOF
   # Database configuration
   DB_PATH=~/.personal-blog/-generatorblog.db

   # Server configuration
   APP_PORT=8080
   EOF

   # Edit .env if needed
   nano ~/.personal-blog-generator/.env
   ```

4. **Run Locally**
   ```bash
   # Build Tailwind CSS (in background)
   npx tailwindcss -i ./admin-files/input.css -o ./admin-files/css/style.css --watch &

   # Run the application
   go run main.go
   ```

5. **Access the Application**
   - Blog: http://localhost:8080
   - Admin: http://localhost:8080/admin/dashboard

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_PATH` | Path to SQLite database | `blog.db` |
| `APP_PORT` | Port for the application | `8080` |
| `TEMPLATE_PATH` | Path to HTML templates directory | `./templates` |
| `OUTPUT_PATH` | Path to generated static site directory | `./html-outputs` |
| `DEPLOY_HOST` | SSH host for deployment | - |
| `DEPLOY_USER` | SSH user for deployment | - |
| `DEPLOY_PATH` | Remote path for deployment | - |
| `BACKUP_PATH` | Path for database backups | - |

> Deployment is still an idea for the future. This configuration is unused at the moment.

## Admin Interface

### Login
Access the admin at `/admin/dashboard` with basic authentication.

### Features
- **Posts**: Create, edit, and manage blog posts with markdown support
- **Portfolio**: Showcase your projects and work
- **Pages**: Create static pages for your site
- **Publishing**: Generate and deploy your static site
- **Backup**: Automatic database backups

### Content Management
- Rich text editing with markdown support
- Tag management for posts
- Image upload and management
- SEO-friendly URLs and metadata

## Deployment

Deployment is still an idea for the future.


### Automated Deployment

1. **Configure Environment**
   ```bash
   # Update ~/.personal-blog-generator/.env with deployment settings
   cat >> ~/.personal-blog-generator/.env << EOF
   DEPLOY_HOST=your-server.com
   DEPLOY_USER=your-user
   DEPLOY_PATH=/var/www/blog
   BACKUP_PATH=./backups
   EOF
   ```

2. **Deploy**
   ```bash
   ./deploy.sh
   ```

### Manual Server Setup

#### 1. Server Preparation
```bash
# Create user and directories
sudo useradd -m -s /bin/bash blog
sudo mkdir -p /var/www/blog
sudo chown -R blog:blog /var/www/blog

# Install dependencies
sudo apt update
sudo apt install nginx sqlite3 certbot python3-certbot-nginx
```

#### 2. Application Setup
```bash
# Copy files to server
scp main blog.db .env blog@your-server.com:/var/www/blog/

# Set permissions
sudo chown -R blog:blog /var/www/blog
sudo chmod 600 /var/www/blog/.env
```

#### 3. Nginx Configuration
```bash
# Copy nginx config
sudo cp nginx.conf /etc/nginx/sites-available/blog
sudo ln -s /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/

# Update server_name and root path in nginx.conf
sudo nano /etc/nginx/sites-available/blog

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

#### 4. SSL Certificate (Let's Encrypt)
```bash
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

#### 5. Systemd Service
```bash
# Copy service file
sudo cp personal-blog-generator.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable and start service
sudo systemctl enable personal-blog-generator
sudo systemctl start personal-blog-generator

# Check status
sudo systemctl status personal-blog-generator
```

#### 6. Service Management
```bash
# View logs
sudo journalctl -u personal-blog-generator -f

# Restart service
sudo systemctl restart personal-blog-generator

# Stop service
sudo systemctl stop personal-blog-generator
```

## Project Structure

```
personal-blog-generator/
├── admin-files/          # Admin interface static files
├── html-outputs/         # Generated static site
├── internal/            # Go application code
│   ├── db/             # Database utilities
│   ├── generator/      # Static site generator
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   └── repository/     # Data access layer
├── static/             # Public static assets
├── templates/          # HTML templates
├── main.go             # Application entry point
├── nginx.conf          # Production nginx config
├── nginx-dev.conf      # Development nginx config
├── personal-blog-generator.service  # Systemd service file
└── README.md           # This file
```

## Development

### Building
```bash
# Build the application
go build -o main .

# Build for production
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .
```

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   # Find process using port 8080
   sudo lsof -i :8080
   # Kill the process
   sudo kill -9 <PID>
   ```

2. **Database connection failed**
   - Check `DB_PATH` in `.env`
   - Ensure database file exists and has correct permissions
   - Run database migrations if needed

3. **Admin login not working**
   - Check `.htpasswd` file exists and has correct format
   - Verify nginx configuration includes auth directives

4. **Static site not generating**
   - Check file permissions on `html-outputs/` directory
   - Ensure all required templates exist
   - Check application logs for errors

5. **Service not starting**
   ```bash
   # Check service status
   sudo systemctl status personal-blog-generator

   # View logs
   sudo journalctl -u personal-blog-generator -f
   ```

### Logs
- Application logs: Check systemd journal
- Nginx logs: `/var/log/nginx/`
- Database logs: Check SQLite file permissions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue on GitHub
- Check the troubleshooting section
- Review the documentation

# Local Development Setup with Nginx

This guide explains how to set up and use the Nginx configuration for local development testing of the blog generator system.

## Prerequisites

- Nginx installed on your system
- Go application running on port 8080
- Generated static files in the `html-outputs/` directory

## Configuration Files

The `nginx-dev.conf` file is configured to:
- Serve static blog files from `html-outputs/`
- Proxy admin interface requests to the Go application
- Proxy API requests to the Go application

The Go application itself implements HTTP Basic Authentication for all `/admin` and `/api` routes. You can configure the credentials using the `ADMIN_USERNAME` and `ADMIN_PASSWORD` environment variables in your `.env` file. If not set, they default to `admin`/`admin`.

## Setup Instructions

1. **Configure Admin Credentials (Optional):**
   Add the following to your `~/.personal-blog-generator/.env` file if you want to change from the defaults:
   ```bash
   ADMIN_USERNAME=your_username
   ADMIN_PASSWORD=your_password
   ```

2. **Update nginx-dev.conf paths:**
   - Replace `/path/to/your/project` with the absolute path to your `personal-blog-generator` directory
   - Ensure the paths in the `root` directive and `proxy_pass` URLs are correct

3. **Include the configuration in your main nginx.conf:**
   ```bash
   # Add this line to your main nginx.conf inside the http block:
   include /mnt/c/MyDriver/WorkDocuments/Personal/ariefbayu-dev-blog-generator/personal-blog-generator/nginx-dev.conf;
   ```

   Or create a symlink:
   ```bash
   sudo ln -s /mnt/c/MyDriver/WorkDocuments/Personal/ariefbayu-dev-blog-generator/personal-blog-generator/nginx-dev.conf /etc/nginx/sites-enabled/blog-dev.conf
   ```

4. **Start the Go application:**
   ```bash
   go run main.go
   ```

## Access Points

- **Blog Site:** http://localhost/
- **Admin Interface:** http://localhost/admin/ (requires basic auth)
- **API Endpoints:** http://localhost/api/ (proxied to Go app)

## Testing

1. Visit http://localhost/ to view the generated blog
2. Visit http://localhost/admin/ to access the admin interface (login with admin credentials)
3. Test API endpoints through the admin interface or directly

## Stopping Services

```bash
# Stop Nginx
nginx -s stop

# Stop Go application (Ctrl+C in the terminal running go run)
```

## Troubleshooting

- Ensure the Go application is running on port 8080
- Check that `html-outputs/` contains generated static files
- Verify `.htpasswd` file exists and has correct permissions
- Check Nginx error logs: `sudo nginx -t` and `sudo systemctl status nginx`
- Make sure the nginx-dev.conf is properly included in your main nginx.conf
- If using symlink method, ensure the link exists: `ls -la /etc/nginx/sites-enabled/blog-dev.conf`
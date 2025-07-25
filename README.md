# Image Resizer API in Go

A fast, cloud-ready HTTP API for on-the-fly image resizing, cropping, and format optimization.
Built with Go, Gin, and bimg (libvips), this service fetches images from URLs, processes them according to query parameters, and returns optimized images in modern formats (JPEG, PNG, WebP, AVIF).
Perfect for use in web apps, CMS, and image-heavy sites.

---

## Features
- Resize images by width (and optionally height)
- Crop images
- Optimize and convert to modern formats (JPEG, PNG, WebP, AVIF)
- Fetch images from remote URLs
- Fast and efficient (libvips under the hood)

## Requirements
- Go 1.16+
- [libvips](https://libvips.github.io/libvips/) and [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) installed on your system

## Getting Started

1. **Install dependencies:**
   ```sh
   brew install pkg-config vips  # macOS
   # or use your system's package manager for Linux
   ```
2. **Clone the repo and run:**
   ```sh
   git clone <your-repo-url>
   cd imagery
   go run api/main.go
   ```

## Usage

Send a GET request to `/api/res` with the following query parameters:

- `url` (required): Image URL to process
- `w` (required): Width in pixels
- `h` (optional): Height in pixels (if omitted, height is scaled automatically)
- `c` (optional): Crop (`true` or `false`, default: `false`)
- `f` (optional): Output format (`jpeg`, `png`, `webp`, `avif`, or `auto`, default: `auto`)
- `b` (optional): Blur amount (float, e.g., 1.5)
- `g` (optional): Grayscale (`true` or `false`)

### Example Request

```
curl "http://localhost:8080/api/res?url=https://example.com/image.jpg&w=300&h=200&c=true&f=webp" --output resized.webp
```

- This fetches the image, resizes to 300x200, crops, and returns as WebP.

### Example: Only Width (auto height)
```
curl "http://localhost:8080/api/res?url=https://example.com/image.jpg&w=400" --output resized.jpg
```

### Example: Blur and Grayscale

```
curl "http://localhost:8080/api/res?url=https://example.com/image.jpg&w=400&b=2.5&g=true" --output blurred-gray.jpg
```

This fetches the image, resizes to 400px width, applies a blur with sigma 2.5, and converts it to grayscale.

## Deploying to the Cloud

You can deploy this API to platforms like Render, Fly.io, DigitalOcean, or Google Cloud Run. Make sure to install `libvips` and `pkg-config` in your deployment environment (see Dockerfile example in the repo).

## License

MIT 
# Image Upload System Documentation

## Overview

This module provides a complete image upload system for the e-commerce platform, handling:
- **Product Images** - Multiple images per product with primary image support
- **Review Images** - Customer review photo attachments
- **User Avatars** - Profile picture uploads

---

## Table of Contents

1. [Business Flow](#business-flow)
2. [Supported Image Types](#supported-image-types)
3. [Database Schema](#database-schema)
4. [API Endpoints](#api-endpoints)
5. [File Validation](#file-validation)
6. [Security](#security)
7. [Usage Examples](#usage-examples)

---

## Business Flow

### Image Upload Process

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   1. User   │────▶│  2. Frontend │────▶│  3. Backend │
│  Selects    │     │  Sends File  │     │  Receives   │
│  Image      │     │  (multipart) │     │  Request    │
└─────────────┘     └──────────────┘     └─────────────┘
                                              │
                                              ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  8. Returns │◀────│  7. Saves    │◀────│  4. Validates│
│  Image URL  │     │  Metadata    │     │  File Type  │
│  to Client  │     │  in Database │     │  & Size     │
└─────────────┘     └──────────────┘     └─────────────┘
                                              │
                                              ▼
                                     ┌─────────────┐
                                     │  5. Generates│
                                     │  Unique     │
                                     │  Filename   │
                                     └─────────────┘
                                              │
                                              ▼
                                     ┌─────────────┐
                                     │  6. Saves   │
                                     │  File to    │
                                     │  Storage    │
                                     └─────────────┘
```

### Step-by-Step Explanation

| Step | Description | Details |
|------|-------------|---------|
| 1 | User selects image | User chooses image from device |
| 2 | Frontend sends file | Multipart/form-data POST request |
| 3 | Backend receives | Gin parses multipart form |
| 4 | Validation | Check file type, size, MIME type |
| 5 | Generate filename | UUID-based unique filename |
| 6 | Save file | Write to `/uploads/{type}/` directory |
| 7 | Save metadata | Store URL in database table |
| 8 | Return URL | JSON response with image URL |

---

## Supported Image Types

### Image Categories

| Type | Endpoint | Max Size | Allowed Extensions | Storage Path |
|------|----------|----------|-------------------|--------------|
| Product Images | `/api/upload/product` | 5MB | jpg, jpeg, png, webp | `/uploads/products/` |
| Review Images | `/api/upload/review` | 5MB | jpg, jpeg, png, webp | `/uploads/reviews/` |
| User Avatars | `/api/upload/avatar` | 2MB | jpg, jpeg, png, webp | `/uploads/avatars/` |

### Storage Folder Structure

```
uploads/
├── products/
│   ├── 550e8400-e29b-41d4-a716-446655440000.jpg
│   ├── 550e8400-e29b-41d4-a716-446655440001.png
│   └── ...
├── reviews/
│   ├── 20260310123456_a1b2c3d4e5f6.jpg
│   └── ...
└── avatars/
    ├── user123_avatar.jpg
    └── ...
```

---

## Database Schema

### ProductImages Table

```sql
CREATE TABLE product_images (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    product_id BIGINT NOT NULL,
    url NVARCHAR(500) NOT NULL,
    alt_text NVARCHAR(255),
    is_primary BIT DEFAULT 0,
    sort_order INT DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT FK_product_images_product 
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX IX_product_images_product_id ON product_images(product_id);
```

### ReviewImages Table

```sql
CREATE TABLE review_images (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    review_id BIGINT NOT NULL,
    url NVARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT FK_review_images_review 
        FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE
);

CREATE INDEX IX_review_images_review_id ON review_images(review_id);
```

### UserAvatars Table

```sql
CREATE TABLE user_avatars (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    url NVARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT FK_user_avatars_user 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IX_user_avatars_user_id ON user_avatars(user_id);
```

### ImageUploadLogs Table (Audit Trail)

```sql
CREATE TABLE image_upload_logs (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT,
    image_type NVARCHAR(50),
    original_name NVARCHAR(255),
    stored_name NVARCHAR(255),
    file_size BIGINT,
    ip_address NVARCHAR(45),
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    INDEX IX_image_upload_logs_user_id (user_id),
    INDEX IX_image_upload_logs_created_at (created_at)
);
```

---

## API Endpoints

### 1. Upload Product Image

```http
POST /api/upload/product
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| file | File | Yes | Image file (jpg, jpeg, png, webp) |
| product_id | int64 | Yes | Product ID |
| is_primary | bool | No | Set as primary image |

**Response:**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "url": "/uploads/products/550e8400-e29b-41d4-a716-446655440000.jpg",
    "filename": "550e8400-e29b-41d4-a716-446655440000.jpg",
    "size": 245678,
    "image_type": "product"
  }
}
```

### 2. Upload Multiple Product Images

```http
POST /api/upload/product/multiple
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| files | File[] | Yes | Multiple image files |
| product_id | int64 | Yes | Product ID |

**Response:**
```json
{
  "success": true,
  "message": "Uploaded 3 images successfully",
  "data": {
    "images": [
      {
        "url": "/uploads/products/img1.jpg",
        "filename": "img1.jpg",
        "size": 123456,
        "image_type": "product"
      }
    ],
    "count": 3
  }
}
```

### 3. Upload Review Image

```http
POST /api/upload/review
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| file | File | Yes | Image file |
| review_id | int64 | Yes | Review ID |

### 4. Upload User Avatar

```http
POST /api/upload/avatar
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| file | File | Yes | Avatar image file |

**Note:** User ID is automatically extracted from JWT token.

### 5. Get Product Images

```http
GET /api/upload/product/images?product_id=123
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "images": [
      {
        "id": 1,
        "product_id": 123,
        "url": "/uploads/products/img1.jpg",
        "alt_text": "Product Image 1",
        "is_primary": true,
        "sort_order": 0,
        "created_at": "2026-03-10T10:00:00Z"
      }
    ],
    "count": 1
  }
}
```

### 6. Delete Product Image

```http
DELETE /api/upload/product/:id
Authorization: Bearer {token}
```

### 7. Serve Static Images

```http
GET /uploads/products/{filename}
GET /uploads/reviews/{filename}
GET /uploads/avatars/{filename}
```

**Example:**
```
http://localhost:8080/uploads/products/550e8400-e29b-41d4-a716-446655440000.jpg
```

---

## File Validation

### Validation Rules

| Check | Description | Implementation |
|-------|-------------|----------------|
| File Size | Maximum size enforced | 5MB (products/reviews), 2MB (avatars) |
| File Extension | Whitelist validation | .jpg, .jpeg, .png, .webp |
| MIME Type | Content-type verification | image/jpeg, image/png, image/webp |
| Filename | Sanitization | Remove path traversal characters |

### Validation Code Example

```go
func ValidateFile(file *multipart.FileHeader, maxSize int64, allowedExts []string) error {
    // Check file size
    if file.Size > maxSize {
        return ErrFileTooLarge
    }

    // Check file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    isAllowed := false
    for _, allowedExt := range allowedExts {
        if ext == allowedExt {
            isAllowed = true
            break
        }
    }
    if !isAllowed {
        return ErrInvalidFileType
    }

    // Validate MIME type
    src, _ := file.Open()
    defer src.Close()
    buffer := make([]byte, 512)
    src.Read(buffer)
    contentType := http.DetectContentType(buffer)
    // ... validate content type

    return nil
}
```

---

## Security

### Security Measures

| Measure | Description | Implementation |
|---------|-------------|----------------|
| **File Type Validation** | Prevent non-image uploads | Extension + MIME type check |
| **File Size Limits** | Prevent DoS attacks | Max 5MB per file |
| **Filename Sanitization** | Prevent directory traversal | Remove `../`, `\`, null bytes |
| **UUID Filenames** | Prevent filename conflicts | `uuid.New().String()` |
| **Authentication** | Require login for uploads | JWT middleware |
| **Audit Logging** | Track all uploads | `ImageUploadLog` table |
| **Directory Permissions** | Restrict file access | `os.MkdirAll(dir, 0755)` |

### Security Code Example

```go
func SanitizeFilename(filename string) string {
    // Remove path separators to prevent directory traversal
    filename = strings.ReplaceAll(filename, "/", "")
    filename = strings.ReplaceAll(filename, "\\", "")
    filename = strings.ReplaceAll(filename, "..", "")
    
    // Remove null bytes
    filename = strings.ReplaceAll(filename, "\x00", "")
    
    // Limit filename length
    if len(filename) > 255 {
        ext := filepath.Ext(filename)
        filename = filename[:255-len(ext)] + ext
    }
    
    return filename
}
```

---

## Usage Examples

### Frontend: Upload Product Image (JavaScript)

```javascript
async function uploadProductImage(productId, file, isPrimary = false) {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('product_id', productId);
  formData.append('is_primary', isPrimary);

  const response = await fetch('/api/upload/product', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: formData
  });

  const result = await response.json();
  
  if (result.success) {
    console.log('Image URL:', result.data.url);
    return result.data.url;
  } else {
    throw new Error(result.message);
  }
}

// Usage
const fileInput = document.querySelector('input[type="file"]');
const file = fileInput.files[0];

uploadProductImage(123, file, true)
  .then(url => {
    document.getElementById('preview').src = url;
  })
  .catch(err => console.error(err));
```

### Frontend: Upload Multiple Images

```javascript
async function uploadMultipleProductImages(productId, files) {
  const formData = new FormData();
  formData.append('product_id', productId);
  
  for (const file of files) {
    formData.append('files', file);
  }

  const response = await fetch('/api/upload/product/multiple', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: formData
  });

  return await response.json();
}

// Usage
const files = document.querySelector('input[type="file"]').files;
uploadMultipleProductImages(123, files)
  .then(result => console.log('Uploaded:', result.data.count, 'images'));
```

### Frontend: Upload Avatar

```javascript
async function uploadAvatar(file) {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch('/api/upload/avatar', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: formData
  });

  const result = await response.json();
  
  if (result.success) {
    // Update avatar in UI
    document.getElementById('user-avatar').src = result.data.url;
  }
  
  return result;
}
```

### React Component Example

```jsx
function ProductImageUpload({ productId }) {
  const [uploading, setUploading] = useState(false);
  const [preview, setPreview] = useState(null);

  const handleFileChange = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setUploading(true);
    try {
      const url = await uploadProductImage(productId, file, true);
      setPreview(url);
    } catch (error) {
      alert('Upload failed: ' + error.message);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <input 
        type="file" 
        accept="image/jpeg,image/png,image/webp"
        onChange={handleFileChange}
        disabled={uploading}
      />
      {uploading && <p>Uploading...</p>}
      {preview && <img src={preview} alt="Preview" />}
    </div>
  );
}
```

---

## Error Handling

### Error Response Format

```json
{
  "success": false,
  "message": "File size exceeds maximum allowed size",
  "error": "FILE_TOO_LARGE"
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `FILE_TOO_LARGE` | 400 | File exceeds size limit |
| `INVALID_FILE_TYPE` | 400 | File type not allowed |
| `NO_FILE` | 400 | No file in request |
| `INVALID_FILE` | 400 | Corrupted or invalid file |

---

## Testing

### cURL Examples

**Upload Product Image:**
```bash
curl -X POST http://localhost:8080/api/upload/product \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/image.jpg" \
  -F "product_id=123" \
  -F "is_primary=true"
```

**Upload Avatar:**
```bash
curl -X POST http://localhost:8080/api/upload/avatar \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/avatar.jpg"
```

**Get Product Images:**
```bash
curl -X GET "http://localhost:8080/api/upload/product/images?product_id=123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Configuration

### Environment Variables

```env
# Upload Configuration
UPLOAD_MAX_SIZE=5242880        # 5MB in bytes
UPLOAD_ALLOWED_TYPES=jpg,jpeg,png,webp
UPLOAD_DIR=./uploads
```

---

## Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| 413 Payload Too Large | Check file size limits in nginx/proxy |
| 403 Forbidden | Verify JWT token is valid |
| File not found | Check `/uploads` directory permissions |
| CORS error | Add frontend origin to `CORS_ALLOWED_ORIGINS` |

---

## Best Practices

1. **Always validate on both client and server side**
2. **Compress images before upload** (recommend: Sharp.js on frontend)
3. **Use CDN for production** image serving
4. **Implement lazy loading** for product galleries
5. **Generate thumbnails** for faster previews
6. **Clean up orphaned files** when deleting products

---

## Future Enhancements

- [ ] Image compression on server
- [ ] Thumbnail generation
- [ ] WebP auto-conversion
- [ ] Cloud storage (S3, Azure Blob)
- [ ] Image CDN integration
- [ ] Batch delete functionality
- [ ] Image editing (crop, rotate)

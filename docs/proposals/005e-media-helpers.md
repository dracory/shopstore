# Proposal: Media Business Logic Helpers

## Status: Draft

## Overview

Add business logic helpers to the Media entity for attachment checks and URL validation.

## Proposed Helpers

```go
// Check if media is attached to an entity
func (media *Media) IsAttached() bool {
    return media.EntityID() != ""
}
```

### Package Helpers

```go
// IsValidURL checks if URL has http/https prefix
func IsValidURL(url string) bool {
    return url != "" && (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"))
}

// GetFileExtension extracts file extension from URL string
func GetFileExtension(url string) string {
    parts := strings.Split(url, ".")
    if len(parts) > 1 {
        return parts[len(parts)-1]
    }
    return ""
}
```

Usage:
```go
if shopstore.IsValidURL(media.URL()) {
    ext := shopstore.GetFileExtension(media.URL())
}
```
```

## Interface Updates

Add to `MediaInterface`:

```go
type MediaInterface interface {
    // ... existing methods ...
    
    IsAttached() bool
}
```

Note: `IsValidURL()` and `GetFileExtension()` are package-level helpers, not methods on Media.

## Implementation Notes

### Edge Cases
- `IsValidURL()` only checks for http/https prefixes - does not validate URL format
- `GetFileExtension()` returns raw extension without validation

### URL Edge Cases to Test
- Protocol-relative URLs (`//example.com/image.jpg`)
- URLs with query strings (`image.jpg?size=large`)
- URLs with fragments (`image.jpg#section`)
- Data URIs (`data:image/png;base64,...`)

### Future Extension
Consider MIME type helper:
```go
func (media *Media) MIMEType() string {
    ext := media.SuggestedExtension()
    mimeMap := map[string]string{
        "jpg": "image/jpeg",
        "png": "image/png",
        // ...
    }
    return mimeMap[ext]
}
```

## Testing Requirements

Cover:
- Attachment check with empty and non-empty EntityID
- Valid URL detection (http, https)
- Invalid URL detection (empty, ftp, file, etc.)
- Extension extraction from various URL formats

## Acceptance Criteria

- [ ] `IsAttached()` implemented and tested
- [ ] `IsValidURL()` helper implemented and tested
- [ ] `GetFileExtension()` helper implemented and tested
- [ ] Interface updated

Note: `IsValidURL()` and `GetFileExtension()` are package-level helpers, not methods on Media.

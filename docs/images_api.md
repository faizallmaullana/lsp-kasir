# Images API Documentation

## Overview
Manage image uploads/downloads. Files are stored under `storages/images` and metadata is stored in DB. Use this when you need to upload via multipart or base64 and reference filenames from items.

Base prefix: `/api/images`

Authentication: Uploads and delete require JWT. Downloads are public.

---

## 1) Upload Image (Multipart)
- Method: POST
- Path: `/api/images/upload`
- Auth: Bearer JWT
- Request: multipart/form-data with `file`

Responses
- 201 Created
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "created",
  "DATA": { "id_image": "uuid", "file_name": "generated-name.png" }
}
```
- 400 Bad Request: missing file
- 401 Unauthorized
- 500 Internal Server Error: failed to upload

---

## 2) Upload Image (Base64)
- Method: POST
- Path: `/api/images/upload/base64`
- Auth: Bearer JWT
- Request (JSON):
```json
{ "file_name": "suggested-name.png", "content_type": "image/png", "data_base64": "..." }
```
Responses same as multipart, returns `id_image` and `file_name`.

---

## 3) Download Image (Blob)
- Method: GET
- Path: `/api/images/:id`
- Returns: binary data with appropriate content-type.

---

## 4) Download Image (Base64)
- Method: GET
- Path: `/api/images/:id/base64`
- Returns (200 OK):
```json
{ "MESSAGE": "SUCCESS", "STATUS": "OK", "DATA": { "content_type": "image/png", "data_base64": "..." } }
```

---

## 5) Delete Image
- Method: DELETE
- Path: `/api/images/:id`
- Auth: Bearer JWT

Responses
- 200 OK
```json
{ "MESSAGE": "SUCCESS", "STATUS": "deleted", "DATA": { "id": "uuid" } }
```

## Notes
- Stored filename is returned; use it in `image_url` of items if you serve the `storages/images` folder statically.
- Accepted common types: image/png, image/jpeg, image/gif. Unknown types fall back to extension derived from provided name.

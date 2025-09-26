# Report API Documentation

## Overview
The Report API provides read-only endpoints to retrieve transaction reports by month/year, for today, and for an exact date. These endpoints aggregate transactions and return totals and line items.

Base prefix: `/api/reports`

Authentication: Not required (read-only). If you need to protect them, put behind JWT later.

---

## 1) Report by Month and Year

- Method: GET
- Path: `/api/reports/:bulan/:tahun`
- Description: Returns all transactions in the specified month (1-12) and year (>= 1970), with total count, revenue sum, total products sold, AOV, min/max order, average items/transaction, and top items.

Request
- Path Params:
  - `bulan` — integer month (1-12)
  - `tahun` — integer year (e.g., 2025)

Example
```
GET /api/reports/9/2025
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "month": 9,
    "year": 2025,
    "total_transactions": 3,
    "sum_total_price": 450000,
    "total_products_sold": 12,
    "average_order_value": 150000,
    "min_order_value": 50000,
    "max_order_value": 250000,
    "average_items_per_transaction": 4,
    "top_items": [
      { "id_item": "item-uuid-1", "item_name": "Product A", "image_url": "https://.../a.jpg", "quantity_sold": 7, "revenue": 210000 }
    ],
    "transactions": [
      {
        "id_transaction": "tx-uuid-1",
        "total_price": 150000,
        "buyer_contact": "0812-xxxx",
        "timestamp": "2025-09-26T10:30:00Z"
      }
    ]
  }
}
```
- 400 Bad Request
```json
{ "STATUS": "BAD_REQUEST", "MESSAGE": "invalid month/year" }
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

---

## 2) Report for Today

- Method: GET
- Path: `/api/reports/today`
- Description: Returns today's transactions with total count, revenue sum, total products sold, AOV, min/max order, average items/transaction, and top items.

Example
```
GET /api/reports/today
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "date": "2025-09-26",
    "total_transactions": 2,
    "sum_total_price": 300000,
    "total_products_sold": 6,
    "average_order_value": 150000,
    "min_order_value": 100000,
    "max_order_value": 200000,
    "average_items_per_transaction": 3,
    "top_items": [
      { "id_item": "item-uuid-1", "item_name": "Product A", "image_url": "https://.../a.jpg", "quantity_sold": 4, "revenue": 160000 }
    ],
    "transactions": [
      {
        "id_transaction": "tx-uuid-2",
        "total_price": 200000,
        "buyer_contact": "0813-xxxx",
        "timestamp": "2025-09-26T12:00:00Z"
      }
    ]
  }
}
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

---

## 3) Report by Exact Date (dd/mm/yyyy)

- Method: GET
- Path: `/api/reports/date/:dd/:mm/:yyyy`
- Description: Returns transactions for the exact calendar date with total count, revenue sum, total products sold, AOV, min/max order, average items/transaction, and top items.

Request
- Path Params:
  - `dd` — day (1-31)
  - `mm` — month (1-12)
  - `yyyy` — year (>= 1970)

Example
```
GET /api/reports/date/26/09/2025
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "date": "2025-09-26",
    "total_transactions": 1,
    "sum_total_price": 150000,
    "total_products_sold": 2,
    "average_order_value": 150000,
    "min_order_value": 150000,
    "max_order_value": 150000,
    "average_items_per_transaction": 2,
    "top_items": [
      { "id_item": "item-uuid-3", "item_name": "Product C", "image_url": "https://.../c.jpg", "quantity_sold": 2, "revenue": 150000 }
    ],
    "transactions": [
      {
        "id_transaction": "tx-uuid-3",
        "total_price": 150000,
        "buyer_contact": "0812-xxxx",
        "timestamp": "2025-09-26T08:00:00Z"
      }
    ]
  }
}
```
- 400 Bad Request
```json
{ "STATUS": "BAD_REQUEST", "MESSAGE": "invalid dd/mm/yyyy" }
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

---

## 4) Today Summary (Revenue, Counts, Insights)

- Method: GET
- Path: `/api/reports/today/summary`
- Description: Returns today's revenue (sum_total_price), total transactions, total products sold, average order value (AOV), min/max order value, average items per transaction, and top items.

Example
```
GET /api/reports/today/summary
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "date": "2025-09-26",
    "total_transactions": 5,
    "total_products_sold": 18,
    "sum_total_price": 750000,
    "average_order_value": 150000,
    "min_order_value": 50000,
    "max_order_value": 300000,
    "average_items_per_transaction": 3.6,
    "top_items": [
      { "id_item": "item-uuid-1", "item_name": "Product A", "image_url": "https://.../a.jpg", "quantity_sold": 10, "revenue": 250000 },
      { "id_item": "item-uuid-2", "item_name": "Product B", "image_url": "https://.../b.jpg", "quantity_sold": 8, "revenue": 180000 }
    ]
  }
}
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

## Data Shapes

Monthly Response
```json
{
  "month": 9,
  "year": 2025,
  "total_transactions": 0,
  "sum_total_price": 0,
  "total_products_sold": 0,
  "average_order_value": 0,
  "min_order_value": 0,
  "max_order_value": 0,
  "average_items_per_transaction": 0,
  "top_items": [],
  "transactions": [
    { "id_transaction": "", "total_price": 0, "buyer_contact": "", "timestamp": "" }
  ]
}
```

Today/Exact Date Response
```json
{
  "date": "YYYY-MM-DD",
  "total_transactions": 0,
  "sum_total_price": 0,
  "total_products_sold": 0,
  "average_order_value": 0,
  "min_order_value": 0,
  "max_order_value": 0,
  "average_items_per_transaction": 0,
  "top_items": [],
  "transactions": [
    { "id_transaction": "", "total_price": 0, "buyer_contact": "", "timestamp": "" }
  ]
}
```

## Notes
- These endpoints currently perform in-memory filtering on paginated results; for large datasets, consider DB-side date range queries.
- Timestamps are ISO 8601 strings.
- Top items are capped at 5 and sorted by quantity sold (desc) then revenue (desc).
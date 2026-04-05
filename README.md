# Shop Store <a href="https://gitpod.io/#https://github.com/dracory/shopstore" style="float:right"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/shopstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/shopstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/shopstore)](https://goreportcard.com/report/github.com/dracory/shopstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/shopstore)](https://pkg.go.dev/github.com/dracory/shopstore)

Shop Store is a Go package that provides a database-backed store for common commerce entities such as products, orders, discounts, media, and categories. It builds on the [`dataobject`](https://github.com/dracory/dataobject) pattern to give you rich domain objects with change tracking, metadata handling, soft deletion, and query builders out of the box.

## Table of contents

1. [Features](#features)
2. [Installation](#installation)
3. [Quick start](#quick-start)
4. [Product variants](#product-variants)
5. [Domain entities](#domain-entities)
6. [Query builders](#query-builders)
7. [Metadata & soft deletion](#metadata--soft-deletion)
8. [Debugging & observability](#debugging--observability)
9. [Testing](#testing)
10. [Development](#development)
11. [License](#license)

## Features

- **Composable store** – instantiate a `Store` with your own table names, database connection, and migration settings.
- **Rich domain objects** – `Category`, `Discount`, `Media`, `Order`, `OrderLineItem`, and `Product` types expose defaults, helpers, predicates, and getter/setter chains.
- **Product variants** – support for both simple products and parent/child product variants (e.g., size/color combinations).
- **Change tracking** – entities track dirty fields and ensure updates persist only modified values.
- **Metadata support** – uniform `metas` JSON helpers (`SetMetas`, `UpsertMetas`, `Meta`) across all entities.
- **Soft deletion** – soft-delete helpers hide records unless explicitly requested.
- **Query builders** – fluent `Query` helpers for filtering, pagination, sorting, and counting.
- **Auto-migration** – optional schema bootstrap via [`sb`](https://github.com/dracory/sb) builders.
- **Observability** – opt-in SQL logging and configurable timeouts.

## Installation

The module path is `github.com/dracory/shopstore`.

```bash
go get github.com/dracory/shopstore
```

## Quick start

```go
package main

import (
    "context"
    "database/sql"
    _ "modernc.org/sqlite" // or your preferred driver

    "github.com/dracory/shopstore"
)

func main() {
    db, err := sql.Open("sqlite", "file:shop.db?_pragma=busy_timeout(5000)&cache=shared")
    if err != nil {
        panic(err)
    }

    store, err := shopstore.NewStore(shopstore.NewStoreOptions{
        DB:                     db,
        CategoryTableName:      "shop_category",
        DiscountTableName:      "shop_discount",
        MediaTableName:         "shop_media",
        OrderTableName:         "shop_order",
        OrderLineItemTableName: "shop_order_line_item",
        ProductTableName:       "shop_product",
        AutomigrateEnabled:     true,
        DebugEnabled:           true,
    })
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    product := shopstore.NewProduct().
        SetTitle("Cascade T-Shirt").
        SetDescription("Premium cotton tee").
        SetQuantityInt(25).
        SetPriceFloat(19.99)

    if err := store.ProductCreate(ctx, product); err != nil {
        panic(err)
    }

    list, err := store.ProductList(ctx, shopstore.NewProductQuery().
        SetStatus(shopstore.PRODUCT_STATUS_ACTIVE).
        SetLimit(10))
    if err != nil {
        panic(err)
    }

    for _, p := range list {
        println(p.ID(), p.Title(), p.PriceFloat())
    }
}
```

### Creating related entities

All entity constructors set sensible defaults (IDs, timestamps, status, and empty metadata). For example:

```go
order := shopstore.NewOrder().
    SetCustomerID("customer_123").
    SetStatus(shopstore.ORDER_STATUS_PENDING)

if err := store.OrderCreate(ctx, order); err != nil {
    // handle error
}
```

### Product variants

The store supports both **simple products** (single SKU) and **product variants** (parent/child matrix for size, color, etc.).

**Simple product** (no variants):
```go
product := shopstore.NewProduct().
    SetTitle("Cascade T-Shirt").
    SetDescription("Premium cotton tee").
    SetSKU("TS-CAS-001").
    SetQuantityInt(25).
    SetPriceFloat(19.99)

if err := store.ProductCreate(ctx, product); err != nil {
    panic(err)
}
```

**Parent product with variants** (matrix for size/color):
```go
// 1. Create parent (display-only, defines variant dimensions)
parent := shopstore.NewProduct().
    SetTitle("Nike Air Max").
    SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

// Define variant dimensions (schema)
_ = parent.SetVariantMatrixSchema(shopstore.VariantSchema{
    {Name: "color", Required: true, Options: []string{"red", "blue", "black"}},
    {Name: "size", Required: true, Options: []string{"8", "9", "10", "11"}},
})

if err := store.ProductCreate(ctx, parent); err != nil {
    panic(err)
}

// 2. Create variant (specific color/size combination)
variant := shopstore.NewProduct().
    SetTitle("Nike Air Max").
    SetParentID(parent.GetID()).  // Links to parent
    SetSKU("NAM-RED-9").
    SetPriceFloat(129.99).
    SetQuantityInt(15).
    SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

// Store variant's dimension values
_ = variant.SetVariantMatrixValues(map[string]string{"color": "red", "size": "9"})

if err := store.ProductCreate(ctx, variant); err != nil {
    panic(err)
}

// 3. List all variants for a parent
variants, err := store.ProductVariantList(ctx, parent.GetID())
if err != nil {
    panic(err)
}

for _, v := range variants {
    println(v.GetID(), v.GetSKU(), v.GetPriceFloat())
}
```

**Querying products:**
```go
// Get all top-level products (no parent_id)
parents, err := store.ProductList(ctx, shopstore.NewProductQuery().
    SetParentID("0"))

// Get variants of a specific parent
variants, err := store.ProductList(ctx, shopstore.NewProductQuery().
    SetParentID(parentID))

// Check if product is a parent (has variant dimensions)
isParent, err := store.ProductIsParent(ctx, productID)

// Get parent of a variant
parent, err := store.ProductGetParent(ctx, variantID)
```

## Domain entities

Each entity embeds `dataobject.DataObject`, enabling fluent setters and change tracking. Key helpers include:

| Entity | Highlights |
| --- | --- |
| `Product` | `IsActive`, `IsDraft`, slug generation, price/quantity helpers, **parent/child variants support**. |
| `Order` | Rich status predicates (awaiting shipment, refunded, etc.). |
| `OrderLineItem` | Links products to orders, maintains quantity and price helpers. |
| `Discount` | Code generator, amount/percent handling, start/end scheduling. |
| `Category` | Parent/child relationships, active/draft state, meta helpers. |
| `Media` | Sequence positioning, media type/URL helpers for assets. |

Getter and setter methods mirror column names (`Title()`, `SetTitle(string)`, `PriceFloat()`, `SetPriceFloat(float64)`, etc.) and accept both string/typed values where appropriate.

## Query builders

The `New<Category|Discount|Media|Order|OrderLineItem|Product>Query` helpers expose fluent filters:

- `SetID`, `SetIDIn`, `SetStatus`, `SetStatusIn` for equality filters.
- `SetCreatedAtGte`/`SetCreatedAtLte` for time windows.
- `SetTitleLike` for partial matches.
- `SetLimit`, `SetOffset`, `SetOrderBy`, `SetSortDirection` for pagination and sorting.
- `SetSoftDeletedIncluded(true)` to include soft-deleted rows.
- `SetCountOnly(true)` to build count queries.

Queries validate their input (`Validate()`) so you get fast feedback on missing or invalid parameters before hitting the database.

## Metadata & soft deletion

Every entity stores arbitrary key/value pairs in a JSON `metas` column. Use:

```go
_ = product.SetMetas(map[string]string{"color": "navy"})
_ = product.UpsertMetas(map[string]string{"size": "L"})
size := product.Meta("size")
```

Soft deletion is handled via the `soft_deleted_at` column. Standard list operations exclude soft-deleted rows unless `SetSoftDeletedIncluded(true)` is used. Helpers such as `ProductSoftDelete` and `CategorySoftDelete` set the column to the current timestamp.

## Debugging & observability

- Enable SQL logging with `store.EnableDebug(true, slogLogger)`.
- Customize store behaviour through `NewStoreOptions` (`AutomigrateEnabled`, `DebugEnabled`, `DbDriverName`, custom timeouts).
- `AutoMigrate()` runs table creation statements generated through [`sb`](https://github.com/dracory/sb) column builders when `AutomigrateEnabled` is set.

## Testing

Comprehensive unit tests cover entity defaults, predicate helpers, metadata operations, and store behaviour. Run them with:

```bash
go test ./...
```

## Development

- Tasks are defined in `Taskfile.yml` for common workflows.
- The project targets Go `1.25` (see `go.mod`) and uses SQLite for integration tests via `modernc.org/sqlite`.
- Contributions are welcome—please open an issue or pull request with proposed changes.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt).

For commercial use, please use the [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

# Security Review Report
Date: 2026-04-04
Reviewer: Senior Principal Golang Engineer
Codebase: github.com/dracory/shopstore

## Executive Summary

The `shopstore` module is a Go-based e-commerce data store that manages products, orders, categories, discounts, and media. The codebase is currently in a **partially migrated state** (transitioning from an older `goqu`/`dataobject` stack to the `neat` ORM), which has left compilation issues and inconsistent security patterns.

**Overall Risk Level: HIGH**

The most critical issues are:
1. **SQL Injection vulnerabilities** in migration functions and potentially in dynamically constructed `ORDER BY` clauses
2. **Complete absence of authentication/authorization** at the store layer
3. **Sensitive data exposure** through debug logging of raw SQL parameters
4. **Broken build state** with missing `go.sum` entries and unresolved imports

| Severity | Count |
|----------|-------|
| Critical | 3 |
| High | 3 |
| Medium | 4 |
| Low | 6 |
| **Total** | **16** |

---

## Critical Findings (Severity: Critical)

### Finding #1: SQL Injection in Migration Functions
- **Location**: `sqls.go:11`, `sqls.go:23`, `sqls.go:30`
- **Description**: Migration functions concatenate table names and column names directly into raw SQL strings using string concatenation (`+`). While `store.productTableName` is typically controlled by the developer, if it is ever influenced by user input (e.g., multi-tenant table naming), this becomes a direct SQL injection vector. More importantly, this pattern violates secure coding practices by bypassing parameterization entirely.
- **Impact**: An attacker who can influence the `productTableName` value could execute arbitrary SQL commands, including data exfiltration, modification, or schema destruction.
- **Recommendation**: Use parameterized queries for all dynamic SQL. If table/column names must be dynamic, validate them against a strict allow-list of known identifiers before interpolation.
- **Code Example**:
```go
// Vulnerable code (sqls.go:11)
_, err = db.ExecContext(context.Background(), "ALTER TABLE "+store.productTableName+" ADD COLUMN "+COLUMN_PARENT_ID+" TEXT DEFAULT '0'")
```
- **Suggested Fix**:
```go
// Secure implementation
q := store.db.Query().Table(store.productTableName)
// Use schema builder for migrations instead of raw SQL
// or validate identifiers against allow-list
```
- **References**: CWE-89 (SQL Injection), OWASP A03:2021 – Injection

---

### Finding #2: SQL Injection via ORDER BY Clause
- **Location**: `store_category.go:202`, `store_product.go:218`, `store_discount.go:247`, `store_order.go:205`, `store_media.go` (indirectly via `OrderBy`)
- **Description**: The query builders pass user-controlled `OrderBy` and `SortDirection` strings directly to `q.OrderBy(options.OrderBy(), sortOrder)`. SQL `ORDER BY` clauses cannot be parameterized in most database drivers, and if the `neat` ORM does not strictly validate the column name against a whitelist, an attacker could inject arbitrary SQL (e.g., `id; DROP TABLE products; --`).
- **Impact**: Potential SQL injection allowing data exfiltration, unauthorized modification, or denial of service.
- **Recommendation**: Implement a strict allow-list of sortable columns. Reject any `OrderBy` value that is not in the predefined list. Additionally, validate `SortDirection` to only accept `ASC` or `DESC` (case-insensitive).
- **Code Example**:
```go
// Vulnerable code (store_product.go:217-218)
sortOrder := lo.Ternary(options.HasSortDirection(), options.SortDirection(), "desc")
if options.HasOrderBy() {
    q = q.OrderBy(options.OrderBy(), sortOrder)
}
```
- **Suggested Fix**:
```go
// Secure implementation
allowedColumns := map[string]bool{"id": true, "title": true, "created_at": true, "price": true}
if options.HasOrderBy() {
    orderBy := options.OrderBy()
    if !allowedColumns[orderBy] {
        return nil, errors.New("invalid order_by column")
    }
    sortOrder := strings.ToUpper(lo.Ternary(options.HasSortDirection(), options.SortDirection(), "desc"))
    if sortOrder != "ASC" && sortOrder != "DESC" {
        return nil, errors.New("invalid sort_direction")
    }
    q = q.OrderBy(orderBy, sortOrder)
}
```
- **References**: CWE-89 (SQL Injection), OWASP A03:2021 – Injection

---

### Finding #3: Missing Authentication and Authorization
- **Location**: Entire codebase (`store_category.go`, `store_order.go`, `store_product.go`, `store_discount.go`, `store_media.go`)
- **Description**: The `Store` struct exposes all CRUD operations (Create, Read, Update, Delete, SoftDelete) as public methods without any authentication, authorization, or tenant isolation checks. Any code that holds a reference to a `Store` instance can read any customer's orders, modify product prices, delete categories, or apply arbitrary discounts.
- **Impact**: Complete bypass of access control. A bug or malicious code in any part of the application that has access to the store can compromise all e-commerce data.
- **Recommendation**: Introduce an authorization middleware or context-based permission system. At minimum:
  1. Pass an authenticated `Actor` or `UserID` through `context.Context`
  2. Validate that the caller has permission for the requested operation
  3. Implement row-level security for multi-tenant scenarios (e.g., `customer_id` filtering for orders)
- **Code Example**:
```go
// Vulnerable code - no auth checks anywhere
func (store *Store) OrderDeleteByID(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("order id is empty")
    }
    _, err := store.db.Query().Table(store.orderTableName).Where(COLUMN_ID+" = ?", id).Delete()
    return err
}
```
- **Suggested Fix**:
```go
// Secure implementation
func (store *Store) OrderDeleteByID(ctx context.Context, id string) error {
    actor, ok := auth.ActorFromContext(ctx)
    if !ok || !actor.HasPermission("order:delete") {
        return auth.ErrForbidden
    }
    // ... existing validation ...
    return store.db.Query().Table(store.orderTableName).Where(COLUMN_ID+" = ?", id).Delete()
}
```
- **References**: CWE-306 (Missing Authentication for Critical Function), OWASP A01:2021 – Broken Access Control

---

## High Severity Findings

### Finding #4: Inadequate LIKE Clause Escaping
- **Location**: `store_category.go:199-202`, `store_product.go:175-178`, `store_media.go:199-202`
- **Description**: The `TitleLike` search functionality uses manual string replacement (`strings.ReplaceAll`) to escape single quotes, percent signs, and underscores. This approach is:
  1. **Dialect-specific**: Doubling single quotes (`''`) works for SQLite but not consistently across PostgreSQL, MySQL, or SQL Server
  2. **Incomplete**: It does not escape backslashes, which are escape characters in many SQL dialects
  3. **Bypassable**: Depending on the underlying driver's handling of the `LIKE` parameter, complex Unicode sequences or alternative encodings may bypass the escaping
- **Impact**: Potential SQL injection or unintended wildcard behavior in search queries, leading to information disclosure or DoS via expensive wildcard queries
- **Recommendation**: Rely entirely on the ORM's parameter binding for the pattern value. Do not manually concatenate `%` wildcards. Instead, pass the pattern as a bound parameter and let the ORM handle escaping:
```go
q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+options.TitleLike()+"%")
```
If the ORM does not safely handle this, use a prepared statement with a validated search term.
- **Code Example**:
```go
// Vulnerable code (store_product.go:175-178)
searchTerm := strings.ReplaceAll(options.TitleLike(), "'", "''")
searchTerm = strings.ReplaceAll(searchTerm, "%", "\\%")
searchTerm = strings.ReplaceAll(searchTerm, "_", "\\_")
q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+searchTerm+"%")
```
- **Suggested Fix**:
```go
// Secure implementation - let the ORM handle escaping
if options.HasTitleLike() {
    pattern := "%" + options.TitleLike() + "%"
    q = q.Where(COLUMN_TITLE+" LIKE ?", pattern)
}
```
- **References**: CWE-89 (SQL Injection), CWE-943 (Improper Neutralization of Special Elements in Data Query Logic)

---

### Finding #5: Sensitive Data Exposure in Debug Logs
- **Location**: `store.go:29-37`
- **Description**: The `logSql` method logs all SQL statements and their parameters using `slog.Any("params", params)`. When debug mode is enabled, this logs **all query parameters** including `customer_id`, order details, discount codes, and potentially PII. If debug logs are shipped to a centralized logging system (ELK, Splunk, CloudWatch), sensitive customer data is permanently stored in plaintext.
- **Impact**: GDPR/CCPA violation potential; exposure of customer PII, order details, and discount codes to anyone with log access.
- **Recommendation**: 
  1. Never log sensitive parameters (customer IDs, email addresses, order contents) in debug mode
  2. Implement a `SensitiveFields` redaction list that masks known sensitive columns in log output
  3. Use structured logging with explicit field allow-lists for SQL params
- **Code Example**:
```go
// Vulnerable code (store.go:35)
store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
```
- **Suggested Fix**:
```go
// Secure implementation
func redactParams(params []interface{}) []interface{} {
    // Return a redacted copy for logging
    // Never log actual customer/order data
}
```
- **References**: CWE-532 (Insertion of Sensitive Information into Log File), GDPR Article 32 (Security of Processing)

---

### Finding #6: Raw Database Connection Exposure
- **Location**: `store.go:84-87`
- **Description**: The `DB()` method exposes the underlying `*sql.DB` connection without any access control. Callers can bypass all store abstractions and execute arbitrary SQL, circumventing soft-delete filters, audit logging, and business logic.
- **Impact**: Complete circumvention of the store layer's security controls. Any caller can read soft-deleted records, modify internal columns, or drop tables.
- **Recommendation**: Remove the `DB()` method from the public interface, or return a read-only wrapper. If raw access is absolutely required for migrations, restrict it to an internal `migrate` package with build tags.
- **Code Example**:
```go
// Vulnerable code (store.go:84-87)
func (store *Store) DB() *sql.DB {
    db, _ := store.db.DB()
    return db
}
```
- **Suggested Fix**:
```go
// Secure implementation - remove from public API or restrict access
func (store *Store) dbInternal() *sql.DB {
    db, _ := store.db.DB()
    return db
}
```
- **References**: CWE-749 (Exposed Dangerous Method or Function)

---

## Medium Severity Findings

### Finding #7: Unused Query Timeout (Denial of Service)
- **Location**: `store_new.go:71`, `store.go:22`
- **Description**: The `Store` struct declares a `timeoutSeconds` field initialized to 7200 seconds (2 hours), but this value is **never used** in any database operation. Queries can run indefinitely, and a malicious or accidental unbounded query (e.g., a search with no filters on a large product catalog) could consume all database connections and CPU.
- **Impact**: Denial of service via resource exhaustion.
- **Recommendation**: Apply `context.WithTimeout` to all store operations, or configure the `neat` ORM's connection pool with query timeouts.
- **Code Example**:
```go
// Vulnerable code - timeoutSeconds is set but never used
store.timeoutSeconds = 2 * 60 * 60 // 2 hours
```
- **Suggested Fix**:
```go
// Secure implementation
ctx, cancel := context.WithTimeout(ctx, time.Duration(store.timeoutSeconds)*time.Second)
defer cancel()
return store.db.Query().WithContext(ctx).Table(store.productTableName)...
```
- **References**: CWE-400 (Uncontrolled Resource Consumption), CWE-770 (Allocation of Resources Without Limits or Throttling)

---

### Finding #8: Unvalidated Update Data via DataChanged()
- **Location**: `store_product.go:136-138`, `store_order.go:130-132`, `store_category.go:146-148`, `store_discount.go:150-152`, `store_media.go:146-148`
- **Description**: Update methods use `DataChanged()` from `dataobject.DataObject`, then only remove `id`, `hash`, and `data` keys before passing the map directly to the database. If a malicious actor can inject additional keys into the `DataObject` (e.g., via `Hydrate` from untrusted input), those keys will be included in the `UPDATE` statement. The `delete` calls only protect three specific keys.
- **Impact**: Potential privilege escalation or data corruption if internal fields (e.g., `status`, `price`, `customer_id`) are modified by an attacker-controlled entity.
- **Recommendation**: Use an explicit allow-list of updateable fields instead of relying on `DataChanged()`. Validate each key in the changed map against the allow-list before constructing the update.
- **Code Example**:
```go
// Vulnerable code (store_product.go:134-138)
dataChanged := product.DataChanged()
delete(dataChanged, COLUMN_ID) // ID is not updateable
delete(dataChanged, "hash")      // Hash is not updateable
delete(dataChanged, "data")      // Data is not updateable
_, err := store.db.Query().Table(store.productTableName).Where(COLUMN_ID+" = ?", product.GetID()).Update(dataChanged)
```
- **Suggested Fix**:
```go
// Secure implementation
allowed := map[string]bool{"title": true, "description": true, "price": true, "quantity": true}
for k := range dataChanged {
    if !allowed[k] {
        delete(dataChanged, k)
    }
}
```
- **References**: CWE-20 (Improper Input Validation), CWE-915 (Improperly Controlled Modification of Dynamically-Determined Object Attributes)

---

### Finding #9: Missing MarkAsNotDirty in CategoryCreate
- **Location**: `store_category.go:27-43`
- **Description**: Unlike all other `Create` methods (ProductCreate, OrderCreate, DiscountCreate, MediaCreate), `CategoryCreate` does not call `category.MarkAsNotDirty()` after successful creation. This inconsistency can cause subsequent `CategoryUpdate` calls to include the initial creation fields in the update, potentially overwriting data.
- **Impact**: Unexpected data overwrites and potential race conditions in concurrent update scenarios.
- **Recommendation**: Add `category.MarkAsNotDirty()` after successful creation, consistent with all other entity types.
- **Code Example**:
```go
// Vulnerable code (store_category.go:38-43)
err := store.db.Query().Table(store.categoryTableName).Create(data)
if err != nil {
    return err
}
return nil // Missing MarkAsNotDirty()
```
- **Suggested Fix**:
```go
category.MarkAsNotDirty()
return nil
```
- **References**: CWE-665 (Improper Initialization)

---

### Finding #10: Broken Build State with Missing Dependencies
- **Location**: `go.mod`, `go.sum`
- **Description**: The codebase imports `github.com/dracory/neat` and `github.com/dracory/neat/contracts/database/orm` but the `go.sum` file is missing entries for these packages, causing `go build` and vulnerability scanning tools (`govulncheck`) to fail. Additionally, some type files still import the old `github.com/dracory/sb` and `github.com/dracory/dataobject` packages, suggesting an incomplete migration.
- **Impact**: Inability to build, test, or scan for vulnerabilities. The codebase cannot be reliably deployed or audited.
- **Recommendation**: Run `go mod tidy` and `go mod download` to resolve missing entries. Complete the migration by removing all legacy imports (`dracory/sb`, `dracory/dataobject`) and ensuring the build passes cleanly.
- **Code Example**:
```
// Build errors
govulncheck: missing go.sum entry for module providing package github.com/dracory/neat
store.go:8:2: could not import github.com/dracory/neat (invalid package name: "")
```
- **Suggested Fix**:
```bash
go mod tidy
go mod download
go build ./...
```
- **References**: CWE-1104 (Use of Unmaintained Third-Party Components)

---

## Low Severity Findings

### Finding #11: Unstructured Logging in Migrations
- **Location**: `sqls.go:15`, `sqls.go:27`, `sqls.go:34`
- **Description**: Migration functions use `fmt.Println` for status output instead of structured logging. In production environments, `fmt.Println` writes to stdout without log levels, timestamps, or correlation IDs, making it difficult to monitor and alert on migration failures.
- **Impact**: Operational difficulty; log ingestion pipelines may miss or misclassify migration output.
- **Recommendation**: Replace `fmt.Println` with the store's `sqlLogger` or a dedicated migration logger.
- **References**: CWE-117 (Improper Output Neutralization for Logs)

---

### Finding #12: SQLite-Specific Error Message Parsing
- **Location**: `sqls.go:12`, `sqls.go:24`, `sqls.go:31`
- **Description**: Migration functions check for the string `"duplicate column name"` in error messages to determine if a column already exists. This string is SQLite-specific and will not match on PostgreSQL, MySQL, or other databases, causing migrations to fail on non-SQLite backends.
- **Impact**: Non-portable migrations; potential false positives/negatives on different database backends.
- **Recommendation**: Use schema introspection (`HasTable`, `HasColumn`) to check column existence before attempting `ALTER TABLE`, or catch specific driver error codes rather than string matching.
- **Code Example**:
```go
// Vulnerable code (sqls.go:12)
if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
    return err
}
```
- **References**: CWE-393 (Return of Wrong Status Code), CWE-755 (Improper Handling of Exceptional Conditions)

---

### Finding #13: No URL Validation on Media URLs
- **Location**: `type_media.go:294-296`
- **Description**: The `SetURL` method accepts any string without validating that it is a valid, safe URL. Malicious URLs (e.g., `javascript:alert(1)`, `file:///etc/passwd`, or internal service endpoints) could be stored and later served to users, leading to XSS or SSRF vulnerabilities in the presentation layer.
- **Impact**: XSS or SSRF in downstream consumers that render media URLs without additional validation.
- **Recommendation**: Validate URLs against an allow-list of protocols (`http`, `https`) and optionally check that the domain is trusted.
- **Code Example**:
```go
// Vulnerable code (type_media.go:294-296)
func (m *Media) SetURL(url string) MediaInterface {
    m.Set(COLUMN_MEDIA_URL, url)
    return m
}
```
- **Suggested Fix**:
```go
func (m *Media) SetURL(urlStr string) MediaInterface {
    u, err := url.Parse(urlStr)
    if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
        return m // or return error
    }
    m.Set(COLUMN_MEDIA_URL, urlStr)
    return m
}
```
- **References**: CWE-20 (Improper Input Validation), CWE-601 (URL Redirection to Untrusted Site)

---

### Finding #14: No Rate Limiting on Discount Code Lookup
- **Location**: `store_discount.go:81-100`
- **Description**: `DiscountFindByCode` performs an exact-match lookup on discount codes without any rate limiting, brute-force protection, or caching. An attacker could enumerate discount codes by making rapid queries with guessed values.
- **Impact**: Information disclosure of valid discount codes; potential financial impact if codes are used fraudulently.
- **Recommendation**: Implement rate limiting at the API layer (e.g., per-IP or per-user limits on discount code validation). Consider adding a short-lived cache to reduce database load.
- **References**: CWE-307 (Improper Restriction of Excessive Authentication Attempts), OWASP A07:2021 – Identification and Authentication Failures

---

### Finding #15: Go Version Mismatch in CI/CD
- **Location**: `.github/workflows/tests.yml:22`
- **Description**: The GitHub Actions workflow specifies Go 1.22, but `go.mod` requires Go 1.26. This mismatch could lead to build failures or the use of language features not supported by the CI runner.
- **Impact**: CI builds may fail unexpectedly or not catch issues that appear in local development.
- **Recommendation**: Update the GitHub Actions workflow to use Go 1.26 or the latest stable version.
- **Code Example**:
```yaml
# Vulnerable code (.github/workflows/tests.yml)
- uses: actions/setup-go@v5
  with:
    go-version: '1.22'
```
- **Suggested Fix**:
```yaml
- uses: actions/setup-go@v5
  with:
    go-version-file: 'go.mod'
```
- **References**: CWE-1104 (Use of Unmaintained Third-Party Components)

---

### Finding #16: Missing Security Scanning in CI/CD
- **Location**: `.github/workflows/tests.yml`
- **Description**: The CI/CD pipeline only runs `go build` and `go test`. It does not include static analysis (`go vet`), vulnerability scanning (`govulncheck`), code coverage thresholds, or linting (`golangci-lint`).
- **Impact**: Security vulnerabilities and code quality issues may go undetected until they reach production.
- **Recommendation**: Add the following steps to the CI pipeline:
  1. `go vet ./...`
  2. `govulncheck ./...`
  3. `golangci-lint run`
  4. Enforce a minimum code coverage threshold
- **References**: CWE-1071 (Empty or Incomplete Code Block), OWASP CI/CD Security

---

## Best Practice Recommendations

1. **Implement Row-Level Security**: For multi-tenant deployments, ensure `OrderList` and `OrderFindByID` automatically filter by the authenticated customer's ID unless the caller has admin privileges.

2. **Use Immutable IDs**: The `GenerateShortID()` function produces IDs based on timestamps, which are somewhat predictable. Consider using cryptographically secure random IDs for sensitive entities like orders and discounts.

3. **Validate All Foreign Keys**: `OrderLineItemCreate` accepts `product_id` and `order_id` without verifying that the referenced records exist. Add foreign key constraints or explicit existence checks.

4. **Sanitize JSON Metadata**: The `metas` fields store arbitrary JSON. Ensure that downstream consumers properly escape this data when rendering to prevent stored XSS.

5. **Enable SQL Prepared Statements**: Verify that the `neat` ORM uses prepared statements for all parameterized queries to maximize injection resistance.

6. **Implement Audit Logging**: Add an `audit_log` table to track who created, modified, or deleted orders, products, and discounts, including before/after snapshots.

7. **Soft Delete Consistency**: Ensure all `SoftDelete` operations are irreversible without admin privileges. Some `FindByID` methods may still return soft-deleted records if `SoftDeletedIncluded` is not explicitly set to false.

8. **Input Length Limits**: Enforce maximum length limits on string fields (e.g., `title`, `description`, `memo`) to prevent storage abuse and DoS via oversized payloads.

9. **Use HTTPS for Media URLs**: Enforce that all media URLs use `https://` to prevent mixed-content issues and man-in-the-middle attacks.

10. **Graceful Shutdown**: The `Store` does not implement a `Close()` method. Ensure the underlying `*sql.DB` connection pool is properly closed when the application shuts down.

---

## Dependencies Analysis

- **Total dependencies**: 14 direct + transitive
- **Dependencies with known vulnerabilities**: Unable to determine (build is broken; `govulncheck` failed)
- **Outdated dependencies**: 
  - `golang.org/x/crypto v0.52.0` – check for newer versions
  - `modernc.org/sqlite v1.50.1` – check for newer versions
  - `github.com/dracory/neat v0.19.0` – verify against latest
- **Critical observation**: The codebase imports `github.com/dracory/neat` but the `go.sum` is incomplete, and some files still reference the old `github.com/dracory/sb` and `github.com/dracory/dataobject` packages. This hybrid state is a significant supply chain risk.

### Dependency Risk Matrix

| Package | Version | Risk Level | Notes |
|---------|---------|------------|-------|
| `github.com/dracory/neat` | v0.19.0 | High | Build broken; missing go.sum entries |
| `github.com/dracory/dataobject` | v1.7.0 | Medium | Legacy package; migration in progress |
| `github.com/dracory/sb` | unknown | High | Still imported but may be removed |
| `modernc.org/sqlite` | v1.50.1 | Low | Well-maintained; no known CVEs at time of review |
| `golang.org/x/crypto` | v0.52.0 | Low | Part of Go extended crypto; monitor for updates |

---

## Compliance Considerations

### GDPR (General Data Protection Regulation)
- **Article 32 (Security of Processing)**: The debug logging of SQL parameters containing customer IDs and order data violates the principle of data minimization and secure processing.
- **Article 5(1)(f) (Integrity and Confidentiality)**: The lack of access controls means personal data is not adequately protected against unauthorized access.

### PCI-DSS (if handling card data)
- This codebase does not appear to handle raw card data, but if order metadata stores payment references, those fields must be encrypted at rest.
- **Recommendation**: Ensure `memo` and `metas` fields do not store card numbers, CVV codes, or magnetic stripe data.

### SOC 2 (Security and Availability)
- The missing query timeouts and lack of rate limiting are availability concerns.
- The broken CI/CD pipeline (Go version mismatch) affects the change management control.

---

## Summary Statistics

| Severity | Count |
|----------|-------|
| Critical | 3 |
| High | 3 |
| Medium | 4 |
| Low | 6 |
| **Total** | **16** |

---

## Next Steps

1. **Immediate (Critical)**
   - [ ] Fix the build state: run `go mod tidy`, remove all legacy imports (`dracory/sb`, `dracory/dataobject`), and ensure `go build ./...` passes
   - [ ] Refactor migration functions in `sqls.go` to use the `neat` schema builder instead of raw SQL concatenation
   - [ ] Implement an `OrderBy` allow-list in all query builders to prevent SQL injection via sort columns

2. **Short-term (High)**
   - [ ] Remove or restrict the `DB()` public method to prevent raw SQL bypass
   - [ ] Redact sensitive parameters from debug SQL logging
   - [ ] Implement a validation layer for `DataChanged()` updates using field allow-lists
   - [ ] Replace manual `LIKE` escaping with ORM-parameterized patterns

3. **Medium-term (Medium)**
   - [ ] Add `context.WithTimeout` to all store operations using the configured `timeoutSeconds`
   - [ ] Add `MarkAsNotDirty()` to `CategoryCreate` for consistency
   - [ ] Replace `fmt.Println` in migrations with structured logging
   - [ ] Validate media URLs to allow only `http`/`https` protocols

4. **Long-term (Low / Best Practices)**
   - [ ] Introduce an authentication/authorization context system at the store layer
   - [ ] Add `govulncheck`, `go vet`, and `golangci-lint` to CI/CD
   - [ ] Fix the GitHub Actions Go version to match `go.mod`
   - [ ] Implement audit logging for all CRUD operations
   - [ ] Add rate limiting for discount code lookups and other enumeration-prone endpoints

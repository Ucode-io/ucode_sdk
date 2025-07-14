# Ucode SDK for Go

A comprehensive Go SDK for interacting with the Ucode API, providing easy-to-use methods for CRUD operations, querying, and data management.

## Table of Contents
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Basic Usage](#basic-usage)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Installation

```bash
go get github.com/ucode-io/ucode_sdk
```

## Quick Start

```go
import (
    sdk "github.com/ucode-io/ucode_sdk"
)

// Initialize the SDK
newsdk := sdk.New(&sdk.Config{
    BaseURL:   "https://api.your-domain.com",
    AppId:     "your-app-id",
    ProjectId: "your-project-id",
})
```

## Configuration

The SDK requires the following configuration parameters:

| Parameter | Type | Description |
|-----------|------|-------------|
| `BaseURL` | string | The base URL of your Ucode API |
| `AppId` | string | Your application ID |
| `ProjectId` | string | Your project ID |

## Basic Usage

### Creating Records

```go
// Create a new record
body := map[string]any{
    "title": "New Order",
    "status": "pending",
}

createResp, _, err := newsdk.Items("order").Create(body).DisableFaas(true).Exec()
if err != nil {
    // Handle error
}
```

### Updating Records

```go
// Update an existing record
updateBody := map[string]any{
    "title": "Updated Order",
    "status": "processed",
    "guid":  recordGuid, // Required for updates
}

updateResp, _, err := newsdk.Items("order").Update(updateBody).DisableFaas(true).ExecSingle()
if err != nil {
    // Handle error
}
```

### Deleting Records

```go
// Delete a record by GUID
_, err = newsdk.Items("order").Delete().Single(recordGuid).DisableFaas(true).Exec()
if err != nil {
    // Handle error
}
```

### Querying Records

#### Basic List Query
```go
// Get a list of records with pagination
getListResp, _, err := newsdk.Items("order").
    GetList().
    Page(1).
    Limit(20).
    Sort(map[string]any{"created_at": -1}).
    Filter(map[string]any{"status": []string{"new"}}).
    Exec()
```

#### Advanced Filtering
```go
// Query with complex filters
getListResp, _, err := newsdk.Items("order_product").
    GetList().
    Page(1).
    Limit(20).
    Filter(map[string]any{
        "quantity": map[string]any{
            "$gte": 4, // Greater than or equal to 4
        },
    }).
    Exec()
```

## API Reference

### SDK Methods

#### `Items(tableName string)`
Returns an Items instance for the specified table.

#### Items Methods

##### Create Operations
- `Create(body map[string]any)` - Creates a new record
- `DisableFaas(disable bool)` - Disables FaaS execution
- `Exec()` - Executes the operation

##### Update Operations
- `Update(body map[string]any)` - Updates a record
- `ExecSingle()` - Executes update for a single record

##### Delete Operations
- `Delete()` - Prepares delete operation
- `Single(guid string)` - Specifies single record to delete

##### Query Operations
- `GetList()` - Prepares list query
- `Page(page int)` - Sets pagination page
- `Limit(limit int)` - Sets result limit
- `Sort(sort map[string]any)` - Sets sorting criteria
- `Filter(filter map[string]any)` - Sets filtering criteria

### Filter Operators

The SDK supports MongoDB-style query operators:

| Operator | Description | Example |
|----------|-------------|---------|
| `$gte` | Greater than or equal | `{"quantity": {"$gte": 4}}` |
| `$gt` | Greater than | `{"price": {"$gt": 100}}` |
| `$lte` | Less than or equal | `{"quantity": {"$lte": 10}}` |
| `$lt` | Less than | `{"price": {"$lt": 1000}}` |
| `$in` | In array | `{"status": {"$in": ["new", "pending"]}}` |
| `$nin` | Not in array | `{"status": {"$nin": ["cancelled"]}}` |

### Sorting

Sorting uses MongoDB-style syntax:
- `1` for ascending order
- `-1` for descending order

```go
Sort(map[string]any{
    "created_at": -1,  // Descending
    "title": 1,        // Ascending
})
```

## Examples

### Complete CRUD Example

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/spf13/cast"
    sdk "github.com/ucode-io/ucode_sdk"
)

func main() {
    // Initialize SDK
    newsdk := sdk.New(&sdk.Config{
        BaseURL:   "https://api.your-domain.com",
        AppId:     "your-app-id",
        ProjectId: "your-project-id",
    })
    
    // Create
    body := map[string]any{
        "title": fmt.Sprintf("Order_%d", time.Now().Unix()),
        "status": "new",
    }
    
    createResp, _, err := newsdk.Items("order").Create(body).DisableFaas(true).Exec()
    if err != nil {
        log.Fatal(err)
    }
    
    guid := cast.ToString(createResp.Data.Data["guid"])
    fmt.Printf("Created record with GUID: %s\n", guid)
    
    // Update
    updateBody := map[string]any{
        "title": fmt.Sprintf("Updated_Order_%d", time.Now().Unix()),
        "status": "processed",
        "guid":  guid,
    }
    
    updateResp, _, err := newsdk.Items("order").Update(updateBody).DisableFaas(true).ExecSingle()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Updated record: %+v\n", updateResp)
    
    // Query
    listResp, _, err := newsdk.Items("order").
        GetList().
        Page(1).
        Limit(10).
        Sort(map[string]any{"created_at": -1}).
        Filter(map[string]any{"status": "processed"}).
        Exec()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d records\n", len(listResp.Data.Data))
    
    // Delete
    _, err = newsdk.Items("order").Delete().Single(guid).DisableFaas(true).Exec()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Record deleted successfully")
}
```

### HTTP Handler Example

```go
func Handle() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        newsdk := sdk.New(&sdk.Config{
            BaseURL:   os.Getenv("UCODE_BASE_URL"),
            AppId:     os.Getenv("UCODE_APP_ID"),
            ProjectId: os.Getenv("UCODE_PROJECT_ID"),
        })
        
        // Your business logic here
        body := map[string]any{
            "title": fmt.Sprintf("Order_%d", time.Now().Unix()),
        }
        
        createResp, _, err := newsdk.Items("order").Create(body).DisableFaas(true).Exec()
        if err != nil {
            handleError(w, "Failed to create order", err)
            return
        }
        
        response := sdk.Response{
            Status: "success",
            Data:   createResp.Data,
        }
        
        handleResponse(w, response, http.StatusOK)
    }
}
```

## Error Handling

Always check for errors when making API calls:

```go
createResp, _, err := newsdk.Items("order").Create(body).DisableFaas(true).Exec()
if err != nil {
    // Log the error
    log.Printf("API Error: %v", err)
    
    // Return appropriate error response
    return sdk.Response{
        Status: "error",
        Data:   map[string]interface{}{"message": "Failed to create record", "error": err.Error()},
    }
}
```

## Best Practices

1. **Environment Variables**: Store sensitive configuration in environment variables
   ```go
   newsdk := sdk.New(&sdk.Config{
       BaseURL:   os.Getenv("UCODE_BASE_URL"),
       AppId:     os.Getenv("UCODE_APP_ID"),
       ProjectId: os.Getenv("UCODE_PROJECT_ID"),
   })
   ```

2. **Error Handling**: Always handle errors appropriately
   ```go
   if err != nil {
       log.Printf("Operation failed: %v", err)
       // Handle error appropriately
   }
   ```

3. **Use DisableFaas**: When you don't need FaaS execution, disable it for better performance
   ```go
   .DisableFaas(true)
   ```

4. **Pagination**: Always use pagination for large datasets
   ```go
   .Page(1).Limit(100)
   ```

5. **Specific Filters**: Use specific filters to reduce data transfer
   ```go
   .Filter(map[string]any{"status": "active"})
   ```

6. **Proper Sorting**: Always specify sorting for consistent results
   ```go
   .Sort(map[string]any{"created_at": -1})
   ```

## Response Structure

All API responses follow this structure:

```go
type Response struct {
    Status string                 `json:"status"`
    Data   map[string]interface{} `json:"data"`
}
```

## Support

For issues, feature requests, or questions, please refer to the [Ucode documentation](https://ucode.gitbook.io/ucode-docs) or contact support.

## License

This SDK is licensed under the MIT License.
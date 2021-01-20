# JUMPER
[![Github go.mod](https://img.shields.io/github/go-mod/go-version/verzth/jumper?style=for-the-badge)](https://golang.org)
[![Github Release](https://img.shields.io/github/v/release/verzth/jumper?style=for-the-badge)](https://github.com/verzth/jumper)
[![Github Pre-Release](https://img.shields.io/github/v/tag/verzth/jumper?include_prereleases&sort=semver&style=for-the-badge)](https://github.com/verzth/jumper)

Jumper is Go Module to help Developer handling HTTP Request & Response.

```bash
go get git.teknoku.digital/teknoku/jumper
```

##### Usage
###### Request Parser
```go
func SomeHandler(w http.ResponseWriter, r *http.Request) {
    var req = jumper.PlugRequest(r, w) // Request Parser

    if req.HasHeader("X-Custom") {
        // Check whether 'X-Custom' header exist without check the value
    }
    if req.HeaderFilled("X-Custom") {
        // Check whether 'X-Custom' header exist and filled
    }
    customHeader := req.Header("X-Custom") // Get header value

    // http://localhost/service/{id:[0-9]+}/{segment...}
    id := req.GetSegment("id") // Named segment from mux router
    id := req.GetSegmentUint64("id") // Named segment from mux router
    id := req.GetSegmentUint32("id") // Named segment from mux router
    id := req.GetSegmentUint("id") // Named segment from mux router
    id := req.GetSegmentInt64("id") // Named segment from mux router
    id := req.GetSegmentInt32("id") // Named segment from mux router
    id := req.GetSegmentInt("id") // Named segment from mux router

    if req.Has("name") {
        // Check whether 'name' exist without check the value
    }
    if req.Filled("name") {
        // Check whether 'name' exist and filled
    }

    name := req.Get("name") // Get name value as string
    id := req.GetUint64("id") // Get id value as uint64
    id := req.GetUint32("id") // Get id value as uint32
    id := req.GetUint("id") // Get id value as uint
    id := req.GetInt64("id") // Get id value as int64
    id := req.GetInt32("id") // Get id value as int32
    id := req.GetInt("id") // Get id value as int
    price := req.GetFloat64("price") // Get price value as float64
    price := req.GetFloat("price") // Get price value as float32
    status := req.GetBool("active") // Get active value as bool
    birthdate, err := req.GetTime("birthdate") // Get birthdate value as *time.Time with Error handler
    birthdate := req.GetTimeNE("birthdate") // Get birthdate value as *time.Time with No Error
    ids := req.GetArray("ids") // Get ids value as Array of interface{}
    ids := req.GetArrayUniquify("ids") // Get ids value as Array of interface{} and uniquify if possible
    obj := req.GetMap("object") // Get object value as Map of map[string]interface{}
    obj := req.GetStruct("object") // Get object value as struct of interface{}
    json := req.GetJSON("jsonstring") // Get jsonstring value as jumper.JSON
    file := req.GetFile("file") // Get file value as jumper.File
    files := req.GetFiles("files") // Get files value as Array of jumper.File
}
```

###### Response Writer
Response Failed sample:
```json
{
  "status": 0,
  "status_number": "1000001",
  "status_code": "ABCDEF",
  "status_message": "Error occurred",
  "data": null
}
```
Response Success sample:
```json
{
  "status": 1,
  "status_number": "F000002",
  "status_code": "SSSSSS",
  "status_message": "Success",
  "data": {
    "id": 1,
    "name": "json"
  }
}
```

Plug Response Writer
```go
package mypackage

import (
    // SOME PACKAGES
	"git.teknoku.digital/teknoku/jumper"
    // SOME PACKAGES
)

func SomeHandler(w http.ResponseWriter, r *http.Request) {
    var res = jumper.PlugResponse(w) // Response Writer
    var data interface{}

    res.SetHttpCode(200) // Set HTTP Response Code. HTTP/1.1 standard (RFC 7231)

    res.Reply(0, "1000001", "ABCDEF", "Error Occurred")
    res.Reply(1, "F000002", "SSSSSS", "Success", data)
    res.ReplyFailed("1000001", "ABCDEF", "Error Occurred")
    res.ReplySuccess("F000002", "SSSSSS", "Success", data)
}
```

Demo Link
```
http://localhost:9999/?list={"obj":{"id":[1,2,3]}}
```
# JUMPER
[![Github go.mod](https://img.shields.io/github/go-mod/go-version/verzth/jumper?style=for-the-badge)](https://golang.org)
[![Github Release](https://img.shields.io/github/v/release/verzth/jumper)](https://github.com/verzth/jumper)

Jumper is Go Module to help Developer handling HTTP Request, its style like [Laravel](https://laravel.com/docs/5.8/requests) Request.

```bash
go get github.com/verzth/jumper
```

##### Usage
```go
func SomeHandler(w http.ResponseWriter, r *http.Request) {
    var req = jumper.Parse(r, w)

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
    birthdate, err := req.GetTime("birthdate") // Get birthdate value as *time.Time with Error handler
    birthdate := req.GetTimeNE("birthdate") // Get birthdate value as *time.Time with No Error
    ids := req.GetArray("ids") // Get ids value as Array of interface{}
    obj := req.GetMap("object") // Get object value as Map of map[string]interface{}
    obj := req.GetStruct("object") // Get object value as struct of interface{}
    json := req.GetJSON("jsonstring") // Get jsonstring value as jumper.JSON
    file := req.GetFile("file") // Get file value as jumper.File
    files := req.GetFiles("files") // Get files value as Array of jumper.File
    
    res.ReplySuccess("0000000", "SSSSSS", "Success", nil)
}
```
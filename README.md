# 🧠 Oxmies — The Universal Object Mapper for Go

Oxmies is a modern, extensible Go library that unifies **ORM**, **ODM**, and **OHM** patterns into a single, cohesive interface.  
It allows you to define your data models once and interact seamlessly with **SQL**, **NoSQL**, or **in-memory** databases — all using the same API.

### ⚙️ Supported Backends
- 🐘 PostgreSQL / MySQL (ORM)
- 🍃 MongoDB (ODM)
- 🔥 Redis (OHM)
- More coming soon…

---

## 🚀 Usage

Define your model (fr, it's that easy):

```go
import "github.com/oxmies/oxmies"

type User struct {
    oxmies.Model
    ID    int    `orm:"primary_key,column:id"`
    Name  string `orm:"column:name"`
    Email string `orm:"column:email"`
}
```

Initialize Oxmies with a SQL adapter (just plug and play, no stress):

```go
cfg := map[string]any{
    "db": oxmies.SQLConfig{
        Driver:   "postgres",
        User:     "user",
        Password: "pass",
        Host:     "localhost",
        Port:     5432,
        DBName:   "testdb",
        SSLMode:  "disable",
        Params:   map[string]string{"search_path": "public"},
        Debug:    true,
        OxmiesDbConfig: oxmies.OxmiesDbConfig{
            Models: []any{&User{}},
        },
    },
}
oxmies.Initialize(cfg)
```

> That's it! Your app is ready. 
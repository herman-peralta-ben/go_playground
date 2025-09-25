# shared

A local Go shared package that can be used from other files.

## Usage

1. Importing `env_loader.go` on e.g. `main.go`:

```
/
├── shared
│   └── env_loader.go
└── stripe
    └── one_click_checkout
        └── main.go
```
2. Run:

```bash
cd ../stripe/one_click_checkout
go mod edit -replace local/shared=../../shared
go get local/shared
```

3. Import and use on `main.go`:

```go
import(
    "local/shared"
)
func main() {
    shared.LoadDotEnv("../.env")
}
```

4. Run `main.go`:

```bash
go run .
```

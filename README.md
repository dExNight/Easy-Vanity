# TON Vanity Address Generator

Multi-threaded TON vanity address generator (Golang)

## Usage

```bash
# Generate address ending with "CAFE"
go run cmd/main.go -suffix CAFE

# Generate address in workchain -1
go run cmd/main.go -suffix ABC -workchain -1

# Case-insensitive search
go run cmd/main.go -suffix test -case=false
```

## Building
```bash
go build -o tonvanity cmd/main.go
```
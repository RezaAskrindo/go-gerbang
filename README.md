
# go-gerbang

Go Gerbang adalah Api Gateway menggunakan bahasa Go

## Persyaratan

1. Menggunakan framework [Gofiber](https://gofiber.io)
2. Database menggunakan PostgreSQL dengan nama database apigateway

## Alamat Migrasi

```bash
{PATH}:{PORT}/migration
```

## Swagger

untuk mengupdate swagger gunakan syntax berikut:
```bash
swag init
```

## Untuk push tanpa mengganti Git setup
```bash
git push https://github.com/RezaAskrindo/go-gerbang.git main
```
dimana main adalah branch-name

## Setting Env
Untuk penggunaan ENV silahkan cek email atau WA

## Deploy For LINUX
env GOOS=linux GOARCH=amd64 go build -o apigateway-9000

## Running in background in LINUX
./golangapp/apigateway-9000 & disown


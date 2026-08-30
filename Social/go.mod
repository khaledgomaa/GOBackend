// This file equivalent to csproj in .NET
module github.com/gobackend/social

go 1.26.5

require github.com/go-chi/chi v1.5.5

require (
	github.com/go-chi/chi/v5 v5.3.2
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.3
)

tool github.com/joho/godotenv/cmd/godotenv

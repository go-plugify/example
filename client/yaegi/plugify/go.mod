module yaegi_plugify

go 1.23.10

replace example.com/server => ./../../../server

require (
	example.com/server v0.0.0-00010101000000-000000000000
	github.com/go-plugify/go-plugify v0.0.0-20251031032447-9fafe1a59671
)

require github.com/traefik/yaegi v0.16.1 // indirect

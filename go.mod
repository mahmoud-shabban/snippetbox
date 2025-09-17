module github.com/mahmoud-shabban/snippetbox

go 1.22.3

replace github.com/mahmoud-shabban/snippetbox => .

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/justinas/alice v1.2.0
)

require filippo.io/edwards25519 v1.1.0 // indirect

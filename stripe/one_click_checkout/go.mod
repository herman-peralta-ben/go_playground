module stripe_backend

go 1.24.3

require (
	github.com/stripe/stripe-go v70.15.0+incompatible // indirect
	github.com/stripe/stripe-go/v78 v78.12.0 // indirect
	local/shared v0.0.0-00010101000000-000000000000 // indirect
)

// $ go mod edit -replace local/shared=../../shared
// $ go get local/shared
replace local/shared => ../../shared

module github.com/syntio/aquarium-persistor-gcp/pull
 
go 1.13
 
replace github.com/syntio/aquarium-persistor-gcp/lib => ../lib
 
require (
    cloud.google.com/go v0.70.0 // indirect
    github.com/syntio/aquarium-persistor-gcp/lib v1.2.3
)
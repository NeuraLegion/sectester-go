module github.com/NeuraLegion/sectester-go

go 1.19

require (
	github.com/gofrs/uuid/v5 v5.0.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Due to the workflow misconfiguration, the initial development release has been tagged by v1.0.0 instead of the expected v0.1.0
// For details please refer to the GitHub issue at https://github.com/golang/go/issues/58852
retract v1.0.0

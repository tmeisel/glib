module github.com/tmeisel/glib/clients/redis

go 1.22

toolchain go1.23.1

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	github.com/tmeisel/glib/error v0.0.2
	github.com/tmeisel/glib/exec v0.0.1
	github.com/tmeisel/glib/testing v0.0.2
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.4.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.34.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runc v1.1.15 // indirect
	github.com/ory/dockertest v3.3.5+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tmeisel/glib/utils v0.0.1 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/tmeisel/glib/error => ../../error
	github.com/tmeisel/glib/exec => ../../exec
	github.com/tmeisel/glib/testing => ../../testing
	github.com/tmeisel/glib/utils => ../../utils
)

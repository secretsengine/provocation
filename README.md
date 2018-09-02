# Provocation

Provocation is a library for the provisioning and revocation of
secrets/credentials.

Supported secrets engines:
- Consul
- Password
- PostgreSQL
- RabbitMQ

Examples:
```
// Postgres secrets engine
import (
	"fmt"

	"github.com/secretsengine/provocation/engines/postgresql"
)

func main() {
	engine := postgresql.Engine{
		URI: "postgres://localhost:5432/database",
		Username: "postgres",
		Password: "secret",
		Creation: []string{
			`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}'`,
			`GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}"`,
		},
	}

	revocation, credentials, err := engine.Provision(context.TODO(), "foo", "bar")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated credentials: %v:%v\n", credentials["username"], credentials["password"])
	// Use credentials
	// ...

	// Revoke credentials
	if err = engine.Revoke(context.TODO(), revocation); err != nil {
		panic(err)
	}
}
```

```
// Password secrets engine
import (
    "fmt"

    "github.com/secretsengine/provocation/engines/password"
)

func main() {
    engine := password.Engine{
	    Length: 16,
    }

    _, credentials, err := engine.Provision(context.TODO(), "", "")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated password: %v\n", credentials["password"])
}
```
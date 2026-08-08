# Helix SDK (Go)

Go client for a running [Helix](https://github.com/letslego/helix) console — the stable `/helix/v1/*` HTTP API.

```bash
go get github.com/letslego/helix-sdk-go@latest
```

## Quick start

```go
package main

import (
	"fmt"
	"log"

	"github.com/letslego/helix-sdk-go/helix"
)

func main() {
	client := helix.New("http://127.0.0.1:8787")

	turn, err := client.Chat("Plan a weekend trip to Paris", "", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(turn.Reply)

	again, err := client.Chat("Make it cheaper", turn.SessionID, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(again.Reply)
}
```

## API

| Method | HTTP |
| --- | --- |
| `Agent()` | `GET /helix/v1/agent` |
| `Stack()` | `GET /helix/v1/stack` |
| `Sessions()` | `GET /helix/v1/sessions` |
| `Workflows(sessionID)` | `GET /helix/v1/workflows` |
| `Events(sessionID)` | `GET /helix/v1/events` |
| `Run(opts)` / `Chat(...)` | `POST /helix/v1/sessions` |
| `ResolveApproval(...)` | `POST /helix/v1/approvals` |
| `RunSchedule(name)` | `POST /helix/v1/schedules/run` |

## Ecosystem

- Framework: [helix](https://github.com/letslego/helix)
- TypeScript SDK: [helix-sdk](https://github.com/letslego/helix-sdk)
- Python SDK: [helix-sdk-python](https://github.com/letslego/helix-sdk-python)
- Overview: [helix-ecosystem](https://letslego.github.io/helix-ecosystem/)

## License

Apache-2.0 © LetsLego

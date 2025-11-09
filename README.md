# ☁️ northwind-api
A modern Go + PostgreSQL REST API for the Northwind dataset, providing structured access to customers, orders, products, suppliers, and analytical summaries.  
This API powers Power Automate REST API connectors for Copilot Studio agents, and MCP-integrated copilots (like GitHub Copilot Chat via `northwind-mcp-layer`).

## Architecture Overview
northwind-api/
├── cmd/
│   └── api/               # Application entrypoint (starts the Gin HTTP server)
├── internal/
│   ├── db/                # Database connection management
│   ├── handlers/          # All REST API route logic grouped by domain
│   └── models/            # Model response structures
├── schema/                # OpenAPI schema
├── go.mod / go.sum        # Go module dependencies
└── README.md              # Project documentation

## Example Endpoints
```text
| Method | Endpoint                    | Description              | Example Parameters               |
| ------ | --------------------------- | ------------------------ | -------------------------------- |
| `GET`  | `/customers`                | Retrieve customers       | `country=Germany`, `city=London` |
| `GET`  | `/orders`                   | Retrieve orders          | `year=1998`, `customer_id=ALFKI` |
| `GET`  | `/summary/sales-by-country` | Sales by country         | `year=1998`                      |
| `GET`  | `/analytics/top-customers`  | Top customers by revenue | `limit=10`, `country=USA`        |

## Tech Stack
Language: Go 1.23+
Framework: Gin
Database: PostgreSQL (via Supabase)

## Related Project
https://github.com/nicholasraynes/northwind-mcp-layer
- MCP bridge that connects this API to GitHub Copilot Chat or other AI systems.

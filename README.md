# ☁️ northwind-api
A Go + PostgreSQL REST API for the Northwind dataset, providing structured access to customers, orders, products, suppliers, and analytical summaries.  
This API powers Power Automate REST API connectors for Copilot Studio agents, and MCP-integrated copilots (like GitHub Copilot Chat via `northwind-mcp-layer`).

## Architecture Overview
```text
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
```

## Example Endpoints
| Method | Endpoint                    | Description              | Example Parameters               |
| ------ | --------------------------- | ------------------------ | -------------------------------- |
| `GET`  | `/customers`                | Retrieve customers       | `country=Germany`, `city=Berlin` |
| `GET`  | `/orders`                   | Retrieve orders          | `year=1998`, `customer_id=ALFKI` |
| `GET`  | `/summary/sales-by-country` | Sales by country         | `year=1998`                      |
| `GET`  | `/analytics/top-customers`  | Top customers by revenue | `country=USA`                    |

## Example Response Structure `/customers?country=Germany`
```text
{
  "filters": {
    "country": "Germany"
  },
  "count": 1,
  "data": [
    {
      "customer_id": "ALFKI",
      "company_name": "Alfreds Futterkiste",
      "contact_name": "Maria Anders",
      "contact_title": "Sales Representative",
      "address": "Obere Str. 57",
      "city": "Berlin",
      "region": null,
      "postal_code": "12209",
      "country": "Germany",
      "phone": "030-0074321",
      "fax": "030-0076545"
    }
  ]
}
```

## Tech Stack
- Language: Go 1.23+
- Framework: Gin
- Database: PostgreSQL (via Supabase)
- Deployment: Render Web Services

## Companion Project
https://github.com/nicholasraynes/northwind-mcp-layer
- MCP bridge that connects this API to GitHub Copilot Chat or other AI systems.

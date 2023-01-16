# db

`db` is a Go database client library for SRC. We use [sqlc](https://github.com/kyleconroy/sqlc) to auto-generate a SQL database client based off predefined queries.

We also manually maintain a http database client that leverages Supabase's auto-generated PostgREST API. In the future, we should be able to auto-generate the http client as well.

## Contributing

1. Install [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)

2. Write and name your query in `query.sql`

3. Re-generate the `/db` package

```bash
sqlc generate
```

4. Manually update the http client in `http.go`. For information on PostgREST's APIs, checkout their [documentation](https://postgrest.org/en/stable/).

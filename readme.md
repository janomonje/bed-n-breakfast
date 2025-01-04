# Bed and Breakfast project

This is the repository for m practice project bed and breakfast

- Build in Go and javaScript
- Uses bootstrap version
- Uses [chi router](https://github.com/go-chi/chi)
- Uses [alex edwards SCS](https://github.com/alexedwards/scs) router session management
- Uses [nosurf](https://github.com/justinas/nosurf)
- Uses [pop](https://gobuffalo.io/documentation/database/pop/) for database migrations

# Installing soda

- go install github.com/gobuffalo/pop/soda@latest
- export PATH=$PATH:$(go env GOPATH)/bin (if not done already)
- soda -v

# Fizz

To generate migration file `soda generate fizz table_name`

- Create table up migration
  create_table("users") {
  t.Column("id", "integer", {primary: true})
  t.Column("first_name", "string", {default: ""})
  t.Column("last_name", "string", {default: ""})
  t.Column("email", "string", {})
  t.Column("password", "string", {"size": 60})
  t.Column("access_level", "integer", {"default": 1})
  }

cmd: `soda migrate` or `soda migrate down`

- Generating foreign key
 `soda generate fizz FKForReservationTable`

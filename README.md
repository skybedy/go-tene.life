# go-tene.life

## Database

This project uses versioned SQL migrations for MariaDB.

- Schema snapshot: `db/schema.sql`
- Migration files: `migrations/*.up.sql` and `migrations/*.down.sql`
- Migrator: `golang-migrate`

### Prerequisites

- MariaDB running and reachable with credentials from `.env`
- Go installed (used to run the migrate CLI)
- `mysqldump` installed (for schema export)

### Run migrations

```bash
make migrate-up
```

Rollback last migration:

```bash
make migrate-down
```

Show current migration version:

```bash
make migrate-status
```

### Generate schema snapshot

```bash
make dump-schema
```

This exports structure-only SQL for app tables (`weather*`) into `db/schema.sql`.

### Recommended workflow for DB changes

1. Add a new migration pair in `migrations/` with the next version number.
2. Apply it locally using `make migrate-up`.
3. Refresh snapshot using `make dump-schema`.
4. Commit migration files and `db/schema.sql` together.

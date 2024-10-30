```bash
# run local with hot reload
air

# for tailwind css update
tailwindcss -o static/main.css --watch

# Create migration file
migrate create -ext sql -dir db/migrations -seq create_users_table

# Run migrations
migrate -database "postgresql://postgres:secret@localhost:5433/crawlero?sslmode=disable" -path db/migrations up
```

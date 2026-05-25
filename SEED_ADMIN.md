# Seeding the First Admin

If your system has no admin users yet, create the first admin using the protected setup endpoint. This repository includes a helper script at `scripts/seed_first_admin.sh`.

Prerequisites:
- The API must be running and reachable via `API_BASE_URL` (defaults to `http://localhost:8080`).
- You must have the `ADMIN_SETUP_KEY` value configured (the bootstrap key present in your `config` or environment).
- `jq` must be installed (used to build the JSON payload).

Example usage:

```bash
ADMIN_SETUP_KEY=your_setup_key \
API_BASE_URL=http://localhost:8080 \
ADMIN_EMAIL=admin@example.com \
ADMIN_PASSWORD=supersecret \
ADMIN_NAME="Primary Admin" \
./scripts/seed_first_admin.sh
```

Notes:
- The script calls `POST /api/v1/admin/auth/register` with header `x-admin-setup-key`.
- After running, you should be able to log in via the admin UI and create additional admins via the `/admin/admins` page.
- For production, handle `ADMIN_SETUP_KEY` securely and remove or rotate it once initial bootstrap is complete.

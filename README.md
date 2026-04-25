# UniShop

A small e-commerce web app for selling university merchandise — hoodies, t-shirts, caps, jackets, and other branded apparel. Built with Go and PostgreSQL, no JavaScript framework involved. Server-rendered HTML the old-fashioned way.

The project is structured as a full mini-shop: registration and login, a product catalog, a shopping cart, a checkout flow with a (simulated) card payment, an order history, in-app notifications, and a separate admin panel for managing products, orders, and broadcast messages.

> **Note:** This is a learning/demo project. The "payment" is fake — card details are stored as plain text to keep the flow simple for educational purposes. **Do not deploy this as a real store.**

---

## What it does

- Allows users to browse the university merchandise catalog and select sizes (S / M / L / XL) and quantities.
- Saves each user’s shopping cart contents to the database, ensuring the cart is preserved even after logging out.
- Guides the user through the checkout page, where they enter shipping information and card details, and can save card information for future use.
- Creates the order, clears the cart, and sends a confirmation notification—all within a single SQL transaction.
- Supports two roles: `USER` (default) and `ADMIN`. Administrators are directed to a separate control panel after logging in and can manage the entire store.

---

## Main features

**For shoppers**
- Registration and login using an email address and password (passwords are hashed using bcrypt).
- Persistent shopping cart with AJAX-based quantity updates.
- Order checkout with delivery address, phone number, and credit card information.
- “Save this card” option for one-click reuse on the next order.
- Order history showing status, total amount, and delivery details.
- A notifications page combining personal alerts (e.g., “Your order has been placed”) and general announcements from administrators.
- Additional icons to highlight products, such as *Limited Edition* or *Trending*.

**For admins**
- A dashboard displaying data on total revenue, orders, products, users, and the 10 most recent orders.
- Product management (CRUD): add, edit, delete, with the ability to highlight text and toggle views.
- Full order view, including each item, customer name, and shipping information.
- Sending general notifications (promotions, announcements) to all users.

**Auth and sessions**
- A session is a JSON object encoded in Base64 format, stored in a cookie with the `HttpOnly` and `SameSite=Lax` attributes, and valid for 7 days.
- Two levels of middleware—`requireAuth` and `requireAdmin`—control access to the corresponding routes.

---

## Tech stack

- **Language:** Go 1.24
- **Web layer:** Go standard library (`net/http`, `html/template`) — no third-party router or framework
- **Database:** PostgreSQL, accessed via [`github.com/lib/pq`](https://github.com/lib/pq) and raw SQL (no ORM)
- **Password hashing:** `golang.org/x/crypto/bcrypt`
- **Frontend:** Plain HTML templates, vanilla CSS (~1300 lines, all in `static/css/style.css`), and a tiny vanilla JS file for cart updates

---

## Project structure

```
unishop/
├── main.go                  # Entry point, DB connection, server start
├── router.go                # Route registration + auth/admin middleware
├── session.go               # Cookie-based session helpers
├── models.go                # Struct definitions (User, Product, Order, etc.)
├── auth_handlers.go         # Register, login, logout, home redirect
├── product_handlers.go      # Catalog and product detail pages
├── cart_handlers.go         # Cart view, add/remove/update items
├── checkout_handlers.go     # Checkout page, order creation, saved cards
├── order_handlers.go        # User order history
├── user_handlers.go         # Profile and notifications pages
├── admin_handlers.go        # Admin dashboard, product mgmt, orders, broadcasts
│
├── go.mod / go.sum          # Module definition and locked dependencies
├── .env                     # Local DB credentials (used as defaults in main.go)
│
├── migrations/
│   └── unishop.sql          # Full PostgreSQL schema + seed data (pg_dump)
│
├── static/
│   ├── css/style.css        # All styling
│   ├── js/cart.js           # AJAX quantity updater for cart
│   └── images/              # Product images (referenced by the seed data)
│
└── templates/               # html/template files
    ├── login.html, register.html
    ├── catalog.html, product.html
    ├── cart.html, checkout.html
    ├── orders.html, profile.html, notifications.html
    └── admin_*.html         # Dashboard, products list, product form,
                             # orders, notifications page
```

---

## Installation

### Prerequisites
- Go 1.24 or newer
- PostgreSQL 14+ (the dump was produced with version 16.10, so anything reasonably modern should work)

### 1. Clone and grab dependencies
```bash
git clone <your-repo-url> unishop
cd unishop
go mod download
```

### 2. Create the database
Create an empty database called `unishop` (or whatever you set in your env):
```bash
createdb unishop
```

### 3. Load the schema and seed data
The repo ships with a full `pg_dump` file that creates every table and inserts 8 demo products plus a few demo accounts:
```bash
psql -U postgres -d unishop -f migrations/unishop.sql
```

> The dump includes a `\restrict ...` line at the top from a recent `pg_dump`. If your `psql` doesn't recognize it, just delete that line — the rest is standard SQL.

### 4. Make sure product images exist
The seed data references files like `/static/images/University_Hoodie.png`, `/static/images/Classic_TShirt.png`, etc. The `static/images/` directory ships empty in this snapshot, so drop matching PNGs in there or update the `image_url` values in the database. Missing images won't break the app — they just won't render.

---

## Running the project

From the project root:
```bash
go run .
```

You should see:
```
Successfully connected to PostgreSQL database
Server starting on http://localhost:8080
```

Then open **http://localhost:8080** in your browser.

To build a binary instead:
```bash
go build -o unishop
./unishop
```

---

## Configuration

All configuration is set via environment variables, and the default values are defined in the `main.go` file. The included `.env` file contains the corresponding values, but **please note that `main.go` does not load the `.env` file automatically**—it only reads `os.Getenv`. You can either export the variables manually, run the program using a tool such as `direnv`, or simply use the default values.
| Variable      | Default     | Description                  |
|---------------|-------------|------------------------------|
| `DB_HOST`     | `localhost` | PostgreSQL host              |
| `DB_PORT`     | `5432`      | PostgreSQL port              |
| `DB_USER`     | `postgres`  | DB user                      |
| `DB_PASSWORD` | `753159`    | DB password (change this)    |
| `DB_NAME`     | `unishop`   | Database name                |
| `PORT`        | `8080`      | HTTP port the app listens on |

Example:
```bash
export DB_PASSWORD=mysecret
export PORT=3000
go run .
```

---

## Demo accounts

The seed data in `migrations/unishop.sql` includes a few accounts. The login page itself prints these test credentials:

| Email                  | Password      | Role   |
|------------------------|---------------|--------|
| `test@university.edu`  | `password123` | USER   |
| `admin@university.edu` | *(see below)* | ADMIN  |

The administrator account is present in the dump, but its password is not hardcoded in plain text. If you are unable to recover it, the easiest solution is to create a new user and then elevate their privileges in SQL:```sql
UPDATE users SET role = 'ADMIN' WHERE email = 'you@example.com';
```

---

## Routes

**Public**

| Method     | Path        | Purpose                              |
|------------|-------------|--------------------------------------|
| GET        | `/`         | Redirects to `/login` or `/catalog`  |
| GET / POST | `/register` | Show form / create account           |
| GET / POST | `/login`    | Show form / authenticate             |
| GET        | `/logout`   | Clear session and bounce to login    |

**Authenticated user**

| Method | Path                | Purpose                                          |
|--------|---------------------|--------------------------------------------------|
| GET    | `/catalog`          | Product grid                                     |
| GET    | `/product/{id}`     | Product detail with add-to-cart form             |
| GET    | `/cart`             | Cart contents and totals                         |
| POST   | `/cart/add`         | Add a product (size + qty) to the cart           |
| POST   | `/cart/remove`      | Remove a cart line item                          |
| POST   | `/cart/update`      | AJAX — update an item's quantity                 |
| GET    | `/checkout`         | Checkout form, prefilled with saved card if any  |
| POST   | `/checkout/process` | Create the order, charge, clear cart             |
| POST   | `/payment/card`     | Save / update a card outside checkout            |
| GET    | `/orders`           | Order history                                    |
| GET    | `/profile`          | Account info + order count                       |
| GET    | `/notifications`    | Personal + global notifications                  |

**Admin only**

| Method     | Path                          | Purpose                                  |
|------------|-------------------------------|------------------------------------------|
| GET        | `/admin`                      | Dashboard with stats                     |
| GET        | `/admin/products`             | Product list                             |
| GET / POST | `/admin/products/add`         | Create product                           |
| GET / POST | `/admin/products/edit?id=...` | Edit product                             |
| POST       | `/admin/products/delete`      | Delete product                           |
| GET        | `/admin/orders`               | All orders with line items               |
| GET        | `/admin/notifications`        | Manage broadcasts                        |
| POST       | `/admin/notifications/send`   | Push a global notification to every user |

---

## How a typical flow looks

1. Sign up on the `/register` page, then log in on the `/login` page.
2. Browse the catalog on the `/catalog` page, select an item, specify the size and quantity, and add it to the cart.
3. Open the `/cart` page—change the quantity using the “plus” and “minus” buttons (these actions trigger a call to `/cart/update` after the data request).
4. Click “Checkout,” enter your address and card details, check “Save this card” if desired, and submit the form.
5. The handler opens a transaction, updates or inserts the saved card details if necessary, inserts a row into the `orders` table, copies the items from the cart to `order_items`, clears the cart, creates a notification, and commits the changes.
6. You are redirected to the `/orders?success=true` page, where the new order is displayed at the top.

---

## Database schema (short version)

The `migrations/unishop.sql` file is the source of truth. The main tables:

- **users** — `id`, `name`, `email`, `password_hash`, `role` (`USER` / `ADMIN`)
- **products** — name, description, price, image URL, optional `highlight_text` + `highlight_enabled` flag
- **carts** — one per user
- **cart_items** — `cart_id`, `product_id`, `size`, `quantity`
- **orders** — total, status, delivery address, phone, payment method
- **order_items** — snapshot of each line at purchase time (product name + price are copied so historical orders stay correct even if the product changes)
- **saved_cards** — cardholder, number, expiry, cvv (plain text, demo only)
- **notifications** — personal (`user_id` set) or global (`is_global = true`)

---

## Troubleshooting

**`Failed to connect to the database` on startup**
Make sure PostgreSQL is running and that the environment variables / default values in the `main.go` file actually match your local configuration. The default password (`753159`) is almost certainly not yours—set the `DB_PASSWORD` value before running.

**`pq: relation ‘users’ does not exist`**
You skipped the migration step. Run `psql -U postgres -d unishop -f migrations/unishop.sql` and try again.

**Unable to log in with the test account**
The dump contains real bcrypt hashes. If `test@university.edu` / `password123` don’t work, the safest way is to register a new account via `/register`.

**Product images are not displayed**
The `static/images/` folder in this snapshot is empty, but the source data references PNG files by name. Either place the corresponding files there, or edit the `image_url` column in the `products` table. When an administrator creates a product, the file `/static/images/placeholder.jpg` is also used if the URL is left blank.

**`Access denied. Administrator privileges required.`**
Your account's `role` is `USER`. Upgrade it in SQL (see *Demo Accounts*).

---
CREATE TABLE IF NOT EXISTS users (
  user_id INTEGER PRIMARY KEY,
  user_name TEXT,
  user_email TEXT,
  user_password TEXT
);

CREATE TABLE IF NOT EXISTS drinks (
  drink_id INTEGER PRIMARY KEY,
  drink_name TEXT,
  price INTEGER,
  stock INTEGER
);

CREATE TABLE IF NOT EXISTS sales (
  sale_id INTEGER PRIMARY KEY,
  user_id INTEGER,
  drink_id INTEGER,
  quantity INTEGER,
  total_price INTEGER,
  sale_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_data (
  sync_id INTEGER PRIMARY KEY,
  server_id INTEGER,
  data_type INTEGER,
  data TEXT,
  sync_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS financial_transactions (
  transaction_id INTEGER PRIMARY KEY,
  transaction_type INTEGER CHECK(transaction_type IN (1, 2)),
  transaction_date TIMESTAMP,
  payment_amount INTEGER CHECK(payment_amount >= 0),
  coin_fill_amount INTEGER CHECK(coin_fill_amount >= 0),
  check_constraints TEXT
);

CREATE TABLE IF NOT EXISTS alert_low_stock (
  alert_id INTEGER PRIMARY KEY,
  drink_id INTEGER,
  alert_date TIMESTAMP,
  restock_request_date TIMESTAMP
);

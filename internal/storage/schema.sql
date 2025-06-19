CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- 자동 증가하는 primary key
    user_id TEXT NOT NULL UNIQUE,          -- 사용자 고유 ID (예: 이메일, 사용자명)
    password TEXT NOT NULL                 -- 비밀번호
);

CREATE TABLE IF NOT EXISTS drinks (
  drink_id INTEGER PRIMARY KEY,          -- 음료 ID
  drink_name TEXT,                        -- 음료 이름
  price INTEGER,                          -- 가격
  stock INTEGER                           -- 재고
);

CREATE TABLE IF NOT EXISTS sales (
  sale_id INTEGER PRIMARY KEY,            -- 매출 기록 ID
  user_id INTEGER,                        -- 사용자 ID (users 테이블과 연결)
  drink_id INTEGER,                       -- 음료 ID (drinks 테이블과 연결)
  quantity INTEGER,                       -- 판매된 음료의 수량
  total_price INTEGER,                    -- 음료의 총 매출액
  sale_date TIMESTAMP,                    -- 구매 일시
  FOREIGN KEY (user_id) REFERENCES users(id), -- 사용자 ID 참조
  FOREIGN KEY (drink_id) REFERENCES drinks(drink_id) -- 음료 ID 참조
);

CREATE TABLE IF NOT EXISTS sync_data (
  sync_id INTEGER PRIMARY KEY,            -- 동기화 기록 ID
  server_id INTEGER,                      -- 서버 ID
  data_type INTEGER,                      -- 데이터 유형
  data TEXT,                              -- 동기화된 데이터
  sync_time TIMESTAMP                     -- 동기화 일시
);

CREATE TABLE IF NOT EXISTS financial_transactions (
  transaction_id INTEGER PRIMARY KEY,     -- 거래 기록 ID
  transaction_type INTEGER CHECK(transaction_type IN (1, 2)),  -- 거래 유형 (1: 수금, 2: 동전 채우기)
  transaction_date TIMESTAMP,             -- 거래 일시
  payment_amount INTEGER CHECK(payment_amount >= 0), -- 수금 금액 (0 이상)
  coin_fill_amount INTEGER CHECK(coin_fill_amount >= 0), -- 동전 채우기 금액 (0 이상)
  FOREIGN KEY (transaction_id) REFERENCES sync_data(sync_id) -- sync_data 테이블과 연결
);

CREATE TABLE IF NOT EXISTS inventory_transactions (
  transaction_id INTEGER PRIMARY KEY,     -- 재고 기록 ID
  transaction_type INTEGER CHECK(transaction_type IN (1, 2)), -- 재고 유형 (1: 보충, 2: 삭제)
  transaction_date TIMESTAMP,             -- 거래 일시
  stock_quantity INTEGER,                 -- 재고량
  change_amount INTEGER,                  -- 변화량
  FOREIGN KEY (transaction_id) REFERENCES sync_data(sync_id) -- sync_data 테이블과 연결
);

CREATE TABLE IF NOT EXISTS alert_low_stock (
  alert_id INTEGER PRIMARY KEY,           -- 솔드아웃 알림 ID
  drink_id INTEGER,                       -- 음료 ID (drinks 테이블과 연결)
  alert_date TIMESTAMP,                   -- 알림 일자
  restock_request_date TIMESTAMP,         -- 재고 보충 요청 일자
  FOREIGN KEY (drink_id) REFERENCES drinks(drink_id) -- 음료 ID 참조
);

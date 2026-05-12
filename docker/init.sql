-- =============================================================
-- KFC Forecast Service — Database Initialization
-- Runs automatically on first postgres container start.
-- =============================================================

-- -------------------------
-- Schema
-- -------------------------

CREATE TABLE IF NOT EXISTS stores (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    location   VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100)   NOT NULL,
    category   VARCHAR(50)    NOT NULL,
    price      NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS historical_sales (
    id         SERIAL  PRIMARY KEY,
    store_id   INTEGER NOT NULL REFERENCES stores(id)   ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sale_date  DATE    NOT NULL,
    hour       SMALLINT NOT NULL CHECK (hour BETWEEN 0 AND 23),
    quantity   INTEGER  NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS forecasts (
    id                 SERIAL         PRIMARY KEY,
    store_id           INTEGER        NOT NULL REFERENCES stores(id)   ON DELETE CASCADE,
    product_id         INTEGER        NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    forecast_date      DATE           NOT NULL,
    hour               SMALLINT       NOT NULL CHECK (hour BETWEEN 0 AND 23),
    predicted_quantity NUMERIC(10, 2) NOT NULL,
    generated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- -------------------------
-- Indexes
-- -------------------------

CREATE INDEX IF NOT EXISTS idx_historical_sales_store_product_date
    ON historical_sales (store_id, product_id, sale_date, hour);

CREATE INDEX IF NOT EXISTS idx_forecasts_store_date
    ON forecasts (store_id, forecast_date);

-- Prevent duplicate forecasts for same store/product/date/hour
CREATE UNIQUE INDEX IF NOT EXISTS uq_forecasts_slot
    ON forecasts (store_id, product_id, forecast_date, hour);

-- -------------------------
-- Seed: Stores
-- -------------------------

INSERT INTO stores (name, location) VALUES
    ('KFC Tel Aviv Center',    'Dizengoff Street 50, Tel Aviv'),
    ('KFC Jerusalem Downtown', 'Jaffa Road 12, Jerusalem'),
    ('KFC Haifa Bay',          'HaAtzmaut Road 3, Haifa')
ON CONFLICT DO NOTHING;

-- -------------------------
-- Seed: Products
-- -------------------------

INSERT INTO products (name, category, price) VALUES
    ('Zinger Burger',       'Burgers',   39.90),
    ('Crunchy Burger',      'Burgers',   34.90),
    ('Bucket 6 Pieces',     'Buckets',   79.90),
    ('Bucket 10 Pieces',    'Buckets',  119.90),
    ('French Fries Regular','Sides',     18.90),
    ('French Fries Large',  'Sides',     24.90),
    ('Coleslaw',            'Sides',     14.90),
    ('Pepsi 500ml',         'Drinks',    12.90)
ON CONFLICT DO NOTHING;

-- -------------------------
-- Seed: Historical Sales (30 days, hours 10–23)
-- Uses generate_series + realistic hourly demand pattern.
-- -------------------------

DO $$
DECLARE
    v_store     RECORD;
    v_product   RECORD;
    v_day       DATE;
    v_hour      SMALLINT;
    v_base      INTEGER;
    v_qty       INTEGER;
    -- product-level base multiplier (some items sell more than others)
    v_multiplier NUMERIC;
BEGIN
    FOR v_store IN SELECT id FROM stores LOOP
        FOR v_product IN SELECT id FROM products LOOP

            -- Give each product a stable demand multiplier
            v_multiplier := CASE v_product.id
                WHEN 1 THEN 1.6   -- Zinger Burger  (most popular)
                WHEN 2 THEN 1.2   -- Crunchy Burger
                WHEN 3 THEN 1.0   -- Bucket 6pc
                WHEN 4 THEN 0.7   -- Bucket 10pc    (higher price, fewer sales)
                WHEN 5 THEN 1.8   -- Fries Regular  (highest volume)
                WHEN 6 THEN 1.1   -- Fries Large
                WHEN 7 THEN 0.9   -- Coleslaw
                WHEN 8 THEN 2.0   -- Pepsi          (drinks are high volume)
                ELSE 1.0
            END;

            FOR v_day IN
                SELECT d::DATE
                FROM generate_series(
                    CURRENT_DATE - INTERVAL '30 days',
                    CURRENT_DATE - INTERVAL '1 day',
                    INTERVAL '1 day'
                ) AS d
            LOOP
                FOR v_hour IN 10..23 LOOP

                    -- Base demand by hour
                    v_base := CASE
                        WHEN v_hour BETWEEN 12 AND 14 THEN 14  -- lunch peak
                        WHEN v_hour BETWEEN 17 AND 20 THEN 18  -- dinner peak
                        WHEN v_hour IN (11, 15, 16)   THEN 9   -- shoulder
                        WHEN v_hour IN (21, 22)        THEN 7   -- late evening
                        WHEN v_hour = 23               THEN 4   -- closing
                        ELSE 5                                   -- morning / quiet
                    END;

                    -- Scale by product multiplier and add ±30% random noise
                    v_qty := GREATEST(
                        1,
                        ROUND(
                            v_base * v_multiplier
                            * (0.70 + random() * 0.60)  -- 0.70 – 1.30 range
                        )::INTEGER
                    );

                    INSERT INTO historical_sales
                        (store_id, product_id, sale_date, hour, quantity)
                    VALUES
                        (v_store.id, v_product.id, v_day, v_hour, v_qty);

                END LOOP; -- hour
            END LOOP; -- day
        END LOOP; -- product
    END LOOP; -- store
END $$;

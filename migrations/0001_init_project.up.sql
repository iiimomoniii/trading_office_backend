-- 0001_init_project.up.sql
-- Initial schema for Trading Office AI Dashboard

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS trading_office;

-- market_prices: เก็บราคาตลาดจาก Binance
CREATE TABLE IF NOT EXISTS trading_office.market_prices (
    id         BIGSERIAL PRIMARY KEY,
    symbol     VARCHAR(20)    NOT NULL,
    price      NUMERIC(20, 8) NOT NULL,
    timestamp  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    -- audit columns
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100)             NOT NULL DEFAULT 'tdo-system',
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by VARCHAR(100),
    is_deleted BOOLEAN                  NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS idx_market_prices_symbol     ON trading_office.market_prices (symbol);
CREATE INDEX IF NOT EXISTS idx_market_prices_timestamp  ON trading_office.market_prices (timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_market_prices_is_deleted ON trading_office.market_prices (is_deleted);
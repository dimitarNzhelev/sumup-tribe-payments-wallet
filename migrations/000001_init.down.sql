--- Your reverse (down) migrations go here
--- Drop Triggers
DROP TRIGGER IF EXISTS trigger_wallets_updated_at ON wallets;

--- Drop Trigger Function
DROP FUNCTION IF EXISTS update_timestamp;

--- Drop Tables (reverse order of creation)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;

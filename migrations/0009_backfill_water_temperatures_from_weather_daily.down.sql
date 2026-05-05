-- Reverts backfilled legacy rows only (keeps new manual measurements).
DELETE FROM water_temperatures WHERE source = 'legacy_backfill';

-- Adds an index for fast lookup of the latest station weather measurement.
CREATE INDEX weather_measured_at_index ON weather (measured_at);

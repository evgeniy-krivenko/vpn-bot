ALTER TABLE connections ADD COLUMN is_active boolean DEFAULT false;
ALTER TABLE connections ADD COLUMN last_activate timestamp;
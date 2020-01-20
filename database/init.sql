-------------------- Create tables

DROP SCHEMA public cascade;
CREATE SCHEMA IF NOT EXISTS public;

\i accounts_tbl.sql
\i transactions_tbl.sql

-------------------- Populate data
INSERT INTO accounts (id, balance) VALUES
(1, 118.91) ON CONFLICT DO NOTHING;

INSERT INTO transactions (tx_id, prev_balance, amount, balance, state, source, status) VALUES
('3e83a51f-9407-4b0e-bbc2-68f22c5b7fdc', 0, 10.1, 10.1, 'win', 'client', 'complete'),
('4e2d0dbe-4792-44b6-a8a1-299679d63738', 10.1, 5.1, 15.2, 'win', 'client', 'complete'),
('3d432c9d-5818-4044-bafd-4675483c02fa', 15.2, 10.1, 25.3, 'win', 'client', 'complete'),
('f5ffb6ed-90fb-43fd-af08-60b15888d8fb', 25.3, -12.1, 13.2, 'lost', 'client', 'complete'),
('e5139345-302f-4b75-b79c-26dee76c8c66', 13.2, 20.1, 33.3,'win', 'client', 'complete'),
('0be3a564-00cf-411c-b5ab-9e17768c3529', 33.3, 10.1, 43.4, 'win', 'client', 'complete'),
('81e540cb-b490-4b2f-83a6-294c53831cf3', 43.4, 40.1, 83.5, 'win', 'client', 'complete'),
('3225b6e9-b89d-4a60-ad6c-f9b64e5a80ee', 83.5, 10.11, 93.61, 'win', 'client', 'complete'),
('193b4980-819d-4ec5-ba82-8d9712939720', 93.61, 15.1, 108.71, 'win', 'client', 'complete'),
('c475b953-eda3-4aab-9372-90281192c186', 108.71, 10.2, 118.91, 'win', 'client', 'complete')

ON CONFLICT DO NOTHING;




echo "start create schema, tables, indexes..."

cat init.sql | PGPASSWORD=fh439vbk psql -h localhost -p 5432 -d payments_db -U app_user

echo "end create"

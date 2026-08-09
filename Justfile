migrations-up TARGET_DATABASE_PATH REGISTRY_DATABASE_PATH:
  sqitch deploy --plan-file sql/migrations/sqitch.plan --registry {{ REGISTRY_DATABASE_PATH }}  --db-name {{ TARGET_DATABASE_PATH }}

migrations-down TARGET_DATABASE_PATH REGISTRY_DATABASE_PATH:
  sqitch revert --plan-file sql/migrations/sqitch.plan --registry {{ REGISTRY_DATABASE_PATH }} {{ TARGET_DATABASE_PATH }}

migrations-add CHANGE:
  sqitch add --chdir sql/migrations $(date +%Y%m%d%H%M%S)_{{ CHANGE }}

warehouse-migrations-up:
  sqitch --target tacohiro-warehouse deploy

warehouse-migrations-down:
  sqitch --target tacohiro-warehouse revert

warehouse-migrations-add CHANGE:
  sqitch add --chdir sql/warehouse_migrations {{ CHANGE }}

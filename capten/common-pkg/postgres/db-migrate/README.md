# Usage of the Postgresql DB migration library

- This library used to apply the DB schema upgrade, so that main logic can direct use the new tables.

- This library supports file and BinData based sources to apply the DB schema changes.

# Code Layout

- migration: DB migration specific logic to apply the migration.

- source: DB schema SQL queries in order of applying with ascending sequential numbers. Examples can be found postgres/test/migrations/postres

- postgres: Main interface to use by applications to apply the DB schema changes. Refer the test cases for the usage.

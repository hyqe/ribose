# Migrations

All migrations MUST only ever add to existing functionality. New columns/indexes may be added to a table, but existing columns must never be modified.

For performance reasons it is better to duplicate data, then it is to create relationships. Foreign Keys are strongly discouraged. Joins are strong discouraged. when designing new functionality do so with these principles in mind.

SQL file names should take the form `<version>_<title>.<up|down>.sql`. Each `version` should have both `up` and `down` scripts.
- `version` is a sequential numeric value that indicates the order the script should be run in.
- `title` a friendly description of that the script will do.
- `up` a scripts that creates/adds functionality
- `down` a scripts that can undo its corresponding `up` script. 

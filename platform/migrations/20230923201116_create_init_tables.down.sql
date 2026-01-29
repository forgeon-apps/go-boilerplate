-- Drop child tables first (FK dependencies)
DROP TABLE IF EXISTS public.tasks CASCADE;
DROP TABLE IF EXISTS public.projects CASCADE;
DROP TABLE IF EXISTS public.book CASCADE;

-- Then drop parent table
DROP TABLE IF EXISTS public.users CASCADE;

-- Drop indexes (safe even if tables are gone; CASCADE above usually removes them anyway)
DROP INDEX IF EXISTS public.idx_projects_owner_user_id;
DROP INDEX IF EXISTS public.idx_tasks_project_id;
DROP INDEX IF EXISTS public.idx_tasks_status;

DROP INDEX IF EXISTS public.active_users;
DROP INDEX IF EXISTS public.active_books;

-- Optional: drop uuid extension (only if you are sure nothing else uses it)
-- DROP EXTENSION IF EXISTS "uuid-ossp";

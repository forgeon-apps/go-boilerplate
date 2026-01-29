-- Make this repeatable (optional)
-- TRUNCATE TABLE public.tasks, public.projects, public.book, public.users RESTART IDENTITY CASCADE;

-- ----------------------------
-- 1) USERS (5 records)
-- ----------------------------
WITH inserted_users AS (
  INSERT INTO public.users (
    username, email, name, password, first_name, last_name,
    is_active, is_deleted, is_admin, created_at, updated_at
  )
  VALUES
    ('yaahtze', 'yaahtze@example.com', 'Mr. Yaahtze', 'bcrypt$hash_demo_01', 'Mr.', 'Yaahtze', TRUE, FALSE, TRUE,  NOW(), NOW()),
    ('nami',    'nami@example.com',    'Nami',       'bcrypt$hash_demo_02', 'Nami', 'Navigator', TRUE, FALSE, FALSE, NOW(), NOW()),
    ('luna',    'luna@example.com',    'Luna',       'bcrypt$hash_demo_03', 'Luna', 'Moon', TRUE, FALSE, FALSE, NOW(), NOW()),
    ('zoro',    'zoro@example.com',    'Roronoa Zoro','bcrypt$hash_demo_04','Roronoa','Zoro', TRUE, FALSE, FALSE, NOW(), NOW()),
    ('sanji',   'sanji@example.com',   'Vinsmoke Sanji','bcrypt$hash_demo_05','Vinsmoke','Sanji', TRUE, FALSE, FALSE, NOW(), NOW())
  ON CONFLICT (username) DO NOTHING
  RETURNING id, username
)
SELECT * FROM inserted_users;

-- ----------------------------
-- 2) BOOKS (5 records) - each must reference users(id)
-- ----------------------------
WITH u AS (
  SELECT id, username FROM public.users
  WHERE username IN ('yaahtze','nami','luna','zoro','sanji')
),
inserted_books AS (
  INSERT INTO public.book (
    user_id, title, author, status, meta, created_at, updated_at, is_deleted
  )
  VALUES
    ((SELECT id FROM u WHERE username='yaahtze'), 'Forgeon: Deploy Like a God', 'Yaahtze', 1, '{"tags":["devops","paas"],"pages":120,"lang":"en"}'::jsonb, NOW(), NOW(), FALSE),
    ((SELECT id FROM u WHERE username='nami'),    'The Navigator’s Map',        'Nami',   1, '{"tags":["product","strategy"],"pages":88,"lang":"en"}'::jsonb, NOW(), NOW(), FALSE),
    ((SELECT id FROM u WHERE username='luna'),    'Moonlit Databases',          'Luna',   2, '{"tags":["db","postgres"],"pages":200,"lang":"en"}'::jsonb, NOW(), NOW(), FALSE),
    ((SELECT id FROM u WHERE username='zoro'),    'Three-Sword Refactoring',    'Zoro',   1, '{"tags":["go","clean-code"],"pages":150,"lang":"en"}'::jsonb, NOW(), NOW(), FALSE),
    ((SELECT id FROM u WHERE username='sanji'),   'Cookbook for APIs',          'Sanji',  3, '{"tags":["api","fiber"],"pages":95,"lang":"en"}'::jsonb, NOW(), NOW(), FALSE)
  ON CONFLICT DO NOTHING
  RETURNING id, title
)
SELECT * FROM inserted_books;

-- ----------------------------
-- 3) PROJECTS (5 records)
-- ----------------------------
WITH u AS (
  SELECT id, username FROM public.users
  WHERE username IN ('yaahtze','nami','luna','zoro','sanji')
),
inserted_projects AS (
  INSERT INTO public.projects (
    owner_user_id, name, description, created_at, updated_at
  )
  VALUES
    ((SELECT id FROM u WHERE username='yaahtze'), 'Forgeon Core',      'Main platform services + gateway', NOW(), NOW()),
    ((SELECT id FROM u WHERE username='nami'),    'DX Playground',     'UI/UX and developer experience experiments', NOW(), NOW()),
    ((SELECT id FROM u WHERE username='luna'),    'Billing System',    'Plans, usage records, invoices', NOW(), NOW()),
    ((SELECT id FROM u WHERE username='zoro'),    'Refactor Sprint',   'Tech debt cleanup and consistency', NOW(), NOW()),
    ((SELECT id FROM u WHERE username='sanji'),   'API Kitchen',       'Fiber endpoints + docs + swagger', NOW(), NOW())
  ON CONFLICT DO NOTHING
  RETURNING id, name
)
SELECT * FROM inserted_projects;

-- ----------------------------
-- 4) TASKS (>=5 records) - each must reference projects(id)
-- ----------------------------
WITH p AS (
  SELECT id, name FROM public.projects
  WHERE name IN ('Forgeon Core','DX Playground','Billing System','Refactor Sprint','API Kitchen')
),
inserted_tasks AS (
  INSERT INTO public.tasks (
    project_id, title, status, due_at, created_at, updated_at
  )
  VALUES
    ((SELECT id FROM p WHERE name='Forgeon Core'),    'Wire pooler DATABASE_URL in prod', 'done',  NOW() - interval '1 day', NOW(), NOW()),
    ((SELECT id FROM p WHERE name='Forgeon Core'),    'Add healthcheck endpoint',         'doing', NOW() + interval '1 day', NOW(), NOW()),
    ((SELECT id FROM p WHERE name='DX Playground'),   'Implement users/projects/tasks pages', 'todo', NOW() + interval '3 days', NOW(), NOW()),
    ((SELECT id FROM p WHERE name='Billing System'),  'Create usage rollup job',          'todo',  NOW() + interval '7 days', NOW(), NOW()),
    ((SELECT id FROM p WHERE name='Refactor Sprint'), 'Normalize repo patterns',          'doing', NOW() + interval '2 days', NOW(), NOW()),
    ((SELECT id FROM p WHERE name='API Kitchen'),     'Add swagger comments + examples',  'todo',  NOW() + interval '4 days', NOW(), NOW())
  ON CONFLICT DO NOTHING
  RETURNING id, title
)
SELECT * FROM inserted_tasks;

create schema if not exists forger;

create table if not exists forger.routes (
  id serial primary key,
  uuid uuid not null default gen_random_uuid() unique,
  name varchar(255) not null,
  path varchar(255) not null,
  prefix varchar(255) not null,
  method varchar(10) not null,
  response_type varchar(10) not null,
  response_status_code int not null,
  response_body text not null,
  response_headers jsonb not null,
  response_delay bigint not null default 0,
  is_active boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_routes_uuid on forger.routes (uuid);
create index if not exists idx_routes_prefix on forger.routes (prefix);
create unique index if not exists idx_routes_path_method on forger.routes(path, method);


insert into forger.routes 
(name, path, prefix, method, response_type, response_status_code, response_body, response_headers)
values
('Get all users', '/users', '/users', 'GET', 'STATIC', 200, '[]', '{}'),
('Get user by id', '/users/:id', '/users', 'GET', 'STATIC', 200, '{}', '{}');

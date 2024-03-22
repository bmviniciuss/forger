CREATE TABLE IF NOT EXISTS routes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  path VARCHAR(255) NOT NULL,
  prefix VARCHAR(255) NOT NULL,
  method VARCHAR(10) NOT NULL,
  response_type VARCHAR(10) NOT NULL,
  response_status_code INTEGER NOT NULL,
  response_body TEXT NOT NULL,
  response_headers TEXT NOT NULL,
  response_delay INTEGER NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_routes_uuid ON routes (uuid);

CREATE INDEX IF NOT EXISTS idx_routes_prefix ON routes (prefix);

CREATE UNIQUE INDEX IF NOT EXISTS idx_routes_path_method ON routes(path, method);

INSERT INTO
  routes (
    uuid,
    name,
    path,
    prefix,
    method,
    response_type,
    response_status_code,
    response_body,
    response_headers
  )
VALUES
  (
    'b0ce05e4-99a4-420e-81b8-d580b0bd3ab2',
    'Get all items',
    '/items',
    '/items',
    'GET',
    'DYNAMIC',
    200,
    '{"id":1,"name":"Item 1"}',
    '{"Content-Type": "application/json"}'
  ),
  (
    '5b90c04a-1ef6-4bd5-9d70-4af46e6616d9',
    'Get item by id',
    '/items/{id}',
    '/items',
    'GET',
    'DYNAMIC',
    200,
    '{"id": "{{ requestVar "id"}}","page": "{{ requestQuery "page" }}", "client_id": "{{ requestHeader "client-id" }}", "random_uuid": "{{ uuid "ulid" }}", "time": "{{ time "iso8601" }}"}',
    '{"Content-Type":"application/json","Item-ID":"{{ requestVar \"id\"}}","page":"{{ requestQuery \"page\" }}","res-client-id":"{{ requestHeader \"client-id\" }}"}'
  );

create table sites (
  id serial primary key not null,
  user_id serial not null references users (id),
  name text not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

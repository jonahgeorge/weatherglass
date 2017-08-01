create table events (
  id serial primary key not null,
  site_id serial references sites (id) on delete cascade,
  resource text,
  referrer text,
  title text,
  user_agent text,
  created_at timestamp without time zone not null default now()
);

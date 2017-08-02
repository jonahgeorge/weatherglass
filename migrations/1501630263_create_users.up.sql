create table users (
  id serial,
  name text not null,
  email text not null,
  password_digest text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  email_confirmation_token text,
  is_email_confirmed bool not null default false,
  primary key(id),
  unique(email)
);

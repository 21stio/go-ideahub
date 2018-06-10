create table ideas
(
	title text not null,
	description_markdown text not null,
	created_at timestamp not null,
	updated_at timestamp,
	created_by text not null,
	description_html text,
	id serial not null
		constraint ideas_id_pk
			primary key,
	slug text not null,
	badges text
)
;

create unique index ideas_slug_uindex
	on ideas (slug)
;

create table users
(
	id text not null
		constraint users_id_pk
			primary key,
	username text not null,
	about text,
	email text,
	avatar_url text,
	linkedin_url text,
	website_url text,
	twitter_url text,
	github_url text,
	hackernews_url text not null,
	medium_url text,
	created_at timestamp not null
)
;

create table upvotes
(
	id serial not null
		constraint upvotes_pkey
			primary key,
	user_id text not null,
	idea_id integer
)
;

create unique index upvotes_user_id_idea_id_uindex
	on upvotes (user_id, idea_id)
;

create table comments
(
	parent_id integer,
	idea_id integer not null,
	comment text not null,
	created_at timestamp not null,
	id serial not null
		constraint comments_id_pk
			primary key,
	user_id text not null,
	username text not null
)
;

create table visits
(
	city text not null,
	country_code text not null,
	continent_code text,
	path text not null,
	visitor_id text not null,
	id serial not null
		constraint visits_pkey
			primary key,
	created_at timestamp not null
)
;

create unique index visits_path_visitor_id_uindex
	on visits (path, visitor_id)
;


create table users (
    user_id serial primary key,
    name varchar(256) not null,
    email varchar(256) not null,
    password varchar(256) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

create table access_token (
    user_id integer not null,
    access_token varchar not null,
    refresh_token varchar not null,
    expiry timestamp not null,
    constraint google_calendar_token_user_id foreign key (user_id) references users(user_id)
);

create table google_calendar_token (
    user_id integer not null,
    access_token varchar not null,
    token_type varchar not null,
    refresh_token varchar not null,
    expiry varchar not null,
    constraint google_calendar_token_user_id foreign key (user_id) references users(user_id)
);

create table meeting (
    id serial primary key,
    user_id integer not null,
    summary varchar(256) not null,
    start_date timestamp not null,
    end_date timestamp not null,
    duration varchar(128) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint meeting_user_id foreign key (user_id) references users(user_id)
);

create table meeting_range (
    id serial primary key,
    user_id integer not null,
    summary varchar(256) not null,
    start_time time not null,
    end_time time not null,
    duration varchar(128) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint meeting_user_id foreign key (user_id) references users(user_id)
);

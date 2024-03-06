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
        on delete cascade
        on update cascade
);

create table google_calendar_token (
    user_id integer not null,
    access_token varchar not null,
    token_type varchar not null,
    refresh_token varchar not null,
    expiry varchar not null,
    constraint google_calendar_token_user_id foreign key (user_id) references users(user_id)
        on delete cascade
        on update cascade
);

create table meetings (
    id serial primary key,
    user_id integer not null,
    summary varchar(256) not null,
    start_date timestamp not null,
    end_date timestamp not null,
    duration varchar(128) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint meetings_user_id foreign key (user_id) references users(user_id)
        on delete cascade
        on update cascade
);

create table meetings_ranges (
    id serial primary key,
    user_id integer not null,
    summary varchar(256) not null,
    start_time time not null,
    end_time time not null,
    duration varchar(128) not null,
    code varchar(128) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint meetings_user_id foreign key (user_id) references users(user_id)
        on delete cascade
        on update cascade
);

create table meetings_ranges_emails (
    user_id integer not null,
    meetings_id integer not null,
    email varchar not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint emails_meetings_user_id foreign key (user_id) references users(user_id)
        on delete cascade
        on update cascade,

    constraint emails_meetings_user_id foreign key (user_id) references users(user_id)
        on delete cascade
        on update cascade
);
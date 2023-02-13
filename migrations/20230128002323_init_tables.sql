-- +goose Up
-- +goose StatementBegin
create table users (
    id        serial primary key,
    login     text not null unique,
    password  text not null,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
create index "users_login_idx" ON users ("login");

create table data_creds (
                             "id"  serial primary key,
                             "user_id"   int not null references users on delete cascade,
                             "title" character varying not null,
                             "username"  character varying not null,
                             "password"  character varying not null,
                             "meta"  character varying not null,
                             "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                UNIQUE(user_id,title,username,password,meta)
);

create table data_files
(
    "id"         serial primary key,
    "user_id"     integer REFERENCES users (id) ON DELETE CASCADE,
    "title" character varying not null,
    "filename" character varying not null,
    "path"     character varying not null,
    "meta" character varying not null,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id,title,filename,path,meta)
);

create table data_cards (
                       "id"   serial primary key,
                       "user_id"   int not null references users on delete cascade,
                       "title" character varying not null,
                       "number"    character varying not null,
                       "date"      character varying not null,
                       "cvv"       character varying not null,
                       "meta"  character varying not null,
                       "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       UNIQUE(user_id,title,number,date,cvv,meta)
);

create table data_text (
                            "id"   serial primary key,
                            "user_id"   int not null references users on delete cascade,
                            "title" character varying not null,
                            "text"      character varying not null,
                            "meta"  character varying not null,
                            "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE(user_id,title,text,meta)
);


create table refresh_tokens (
                           "id"   serial primary key,
                           "user_id"   int not null references users on delete cascade,
                           "token" character varying not null,
                           "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "expired_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  "data_cards";
DROP TABLE  "data_files";
DROP TABLE  "data_creds";
DROP TABLE  "data_text";
DROP TABLE  "refresh_tokens";
DROP TABLE  "users";
-- +goose StatementEnd

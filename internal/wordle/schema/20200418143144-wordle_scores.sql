-- +migrate Up
CREATE TABLE wordle_scores
(
    id          BIGSERIAL PRIMARY KEY,
    discord_id  text      NOT NULL references accounts (discord_id),
    game_id     date       not null,
    game_type   int        not null default 5,
    guesses        int       not null,
    created_at  timestamp not null default now()
);

create index on wordle_scores (discord_id);
create index on wordle_scores (created_at);
create index on wordle_scores (game_id);
create index on wordle_scores (game_type);

create unique index
    on wordle_scores
    (
    discord_id,
    game_id,
    game_type
    );

-- +migrate Down
drop table wordle_scores;
-- name: GetScoreHistoryByAccount :many
SELECT *
FROM wordle_scores
         inner join nicknames nick on wordle_scores.discord_id = nick.discord_id
WHERE nick.discord_id = $1
  and nick.server_id = $2
order by game_id;

-- name: ListScores :many
SELECT *
FROM wordle_scores
ORDER BY created_at;

-- name: CountScoresByDiscordId :one
SELECT count(*)
FROM wordle_scores
where discord_id = $1;

-- name: CreateScore :one
INSERT INTO wordle_scores (discord_id, game_id, guesses)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateScore :one
update wordle_scores
set guesses = $2
where discord_id = $1
  and game_id = $3
returning *;

-- name: DeleteScoresForUser :exec
DELETE
FROM wordle_scores
WHERE discord_id = $1;

-- name: GetScoresByServerId :many
with max_game_week as (select max((game_id - cast('1970-01-05' as date)) / 7) game_week
                       from wordle_scores
                                inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                       where n2.server_id = $1
)
select dense_rank() over (partition by n.server_id order by sum((7 - s.guesses) ^ 2) desc) as position,
       n.nickname,
       n.discord_id,
       json_agg(guesses order by s.game_id)             guesses_per_game,
       json_agg(json_build_object(game_id, guesses) order by s.game_id) game_guesses,
       json_agg((7 - s.guesses) ^ 2 order by s.game_id) points_per_game,
       count(distinct game_id)                          games_count,
       sum((7 - s.guesses) ^ 2)                         total
from wordle_scores s
         inner join nicknames n on s.discord_id = n.discord_id
         inner join max_game_week g on g.game_week = ((s.game_id - cast('1970-01-05' as date)) / 7)
where n.server_id = $1
group by n.server_id, n.nickname, n.discord_id
order by sum((7 - s.guesses) ^ 2) desc, count(distinct game_id) desc, nickname;

-- name: GetScoresByServerIdPreviousWeek :many
with max_game_week as (select (max((game_id - cast('1970-01-05' as date)) / 7)) - 1 game_week
                       from wordle_scores
                                inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                       where n2.server_id = $1
)
select dense_rank() over (partition by n.server_id order by sum((7 - s.guesses) ^ 2) desc) as position,
       n.nickname,
       n.discord_id,
       json_agg(guesses order by s.game_id)             guesses_per_game,
       json_agg(json_build_object(game_id, guesses) order by s.game_id) game_guesses,
       json_agg((7 - s.guesses) ^ 2 order by s.game_id) points_per_game,
       count(distinct game_id)                          games_count,
       sum((7 - s.guesses) ^ 2)                         total
from wordle_scores s
         inner join nicknames n on s.discord_id = n.discord_id
         inner join max_game_week g on g.game_week = ((game_id - cast('1970-01-05' as date)) / 7)
where n.server_id = $1
group by n.server_id, n.nickname, n.discord_id
order by sum((7 - s.guesses) ^ 2) desc, count(distinct game_id) desc, nickname;

-- name: GetExpectedWeekGames :many
with max_game_week as (select max((game_id - cast('1970-01-05' as date)) / 7) game_week
                       from wordle_scores
                                inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                       where n2.server_id = $1
),
     current_week_games as (select distinct game_id
                            from wordle_scores
                                     inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                                     inner join max_game_week on game_week = (game_id - cast('1970-01-05' as date)) / 7
                            where n2.server_id = $1)
select * from current_week_games;

-- name: GetExpectedPreviousWeekGames :many
with max_game_week as (select max((game_id - cast('1970-01-05' as date)) / 7) -1 game_week
                       from wordle_scores
                                inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                       where n2.server_id = $1
),
     current_week_games as (select distinct game_id
                            from wordle_scores
                                     inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                                     inner join max_game_week on game_week = (game_id - cast('1970-01-05' as date)) / 7
                            where n2.server_id = $1)
select * from current_week_games;
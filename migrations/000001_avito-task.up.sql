create table teams
(
    team_name text primary key
);
create table users
(
    user_id text primary key,
    username text not null,
    team_name text not null references teams(team_name),
    is_active bool not null default true
);
create table pull_requests
(
    pull_request_id text primary key,
    pull_request_name text not null,
    author_id text not null references users(user_id),
    status text not null CHECK (status IN ('OPEN', 'MERGED')),
    created_at timestamp not null default now(),
    merged_at timestamp
);

create table pr_reviewers
(
    pull_request_id text not null references pull_requests(pull_request_id) on delete cascade,
    reviewer_id text not null references users(user_id),
    primary key (pull_request_id,reviewer_id)
);

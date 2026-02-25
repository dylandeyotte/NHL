# NHL Tracker
NHL Tracker is a RESTful backend API server for tracking and following your favourite players and team.

## Motivation
I wanted an easier way to montior the stats of players without regularly searching for player or sifting through pages of players I don't care about. I created NHL Tracker to have a condensed list of players I care about, whether that's surging rookies or my favourite players, with easy to read stats that give me a broad outlook on how they are performing.

## Quick Start

### Prerequisites

Tested on:
- [Go 1.21](https://go.dev/doc/install)
- [PostgreSQL 15](https://www.postgresql.org)

### Step 1
Create database
```
createdb nhl // Or whatever name you like
```

### Step 2
- Set environment variables
- Create .env file with the following:

```
DB_URL = "postgres://[USER]:@localhost:5432/[DATABASE NAME]?sslmode=disable"
SECRET = [JWT Secret]
PLATFORM = "dev"
```
- Secret can be generated with:
```
openssl rand -base64 64
```

### Step 3
Make migrations:
```
export DATABASE_URL=[YOUR DB_URL]
make migrate
```

Alternatively, you can run migrations with goose if you have it

### Step 4
Load teams into the database with:
```
go run ./cmd/load_teams
```

### Step 5
Run the server
```
go run .
```

## Usage

### API Endpoints

- `POST /api/login` To login
- `POST /api/users` To create a new user
- `POST /api/teams/{tricode}/follow` To follow a team
- `DELETE /api/teams/{tricode}/follow` To unfollow a team
- `POST /api/players/{playerID}/follow` To follow a player
- `DELETE /api/players/{playerID}/follow` To unfollow a player
- `GET /api/players/search` To search for players
- `GET /api/following` To get a list of following players
- `GET /api/home` To see the stats of followed players and team, the homepage
- `POST /api/refresh` To get new access token
- `POST /api/revoke` To revoke refresh token
- `POST /admin/reset` To reset database

## Contributing

If you'd like to suggest improvements or report any bugs, you can open an issue or make a pull request.
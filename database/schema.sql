CREATE TABLE IF NOT EXISTS users (
  id       INTEGER      PRIMARY KEY,
  username VARCHAR(100) NOT NULL UNIQUE,
  passhash VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS api_keys (
  id      INTEGER     PRIMARY KEY,
  api_key VARCHAR(40) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS api_key_scopes (
  id         INTEGER PRIMARY KEY,
  api_key_id INTEGER NOT NULL,

  FOREIGN KEY (api_key_id) REFERENCES api_keys(id)
);

CREATE TABLE IF NOT EXISTS coverage (
  id      INTEGER      PRIMARY KEY,
  project VARCHAR(250) NOT NULL,
  cover   VARCHAR(10)  NOT NULL
);

CREATE TABLE IF NOT EXISTS user_sessions (
  id        INTEGER     PRIMARY KEY,
  sessionid VARCHAR(40) NOT NULL UNIQUE,
  user_id   INTEGER     NOT NULL,
  expires   INTEGER     NOT NULL,

  FOREIGN KEY (user_id) REFERENCES users(id)
);
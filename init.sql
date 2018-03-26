CREATE TABLE IF NOT EXISTS users(
  id  INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(32) UNIQUE NOT NULL DEFAULT "",
  password VARCHAR(64) NOT NULL DEFAULT "",
  email VARCHAR(64) NOT NULL DEFAULT "",
  api_key VARCHAR(64) NOT NULL DEFAULT "",
  state INT NOT NULL DEFAULT 0,
  emailed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);

INSERT OR IGNORE into users (username, password, email, api_key) values ("kitty", "kitty", "329365307@qq.com", "apikey");

CREATE TABLE IF NOT EXISTS userchains(
  userid  INTEGER NOT NULL DEFAULT 0,
  chain VARCHAR(64) NOT NULL DEFAULT "",
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (userid, chain)
);

CREATE TABLE IF NOT EXISTS records (
  id  INTEGER PRIMARY KEY AUTOINCREMENT,
  chain VARCHAR(64) NOT NULL DEFAULT "",
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP
);

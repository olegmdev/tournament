
CREATE TABLE IF NOT EXISTS players (
  id character varying(255) NOT NULL,
  balance double precision NOT NULL DEFAULT 0,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  CONSTRAINT players_pkey PRIMARY KEY (id)
) WITH (
  OIDS=FALSE
);

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'enum_status') THEN
    CREATE TYPE enum_status AS ENUM ('active', 'closed');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS tournaments (
  id character varying(255) NOT NULL,
  deposit double precision NOT NULL DEFAULT 0,
  status enum_status,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  CONSTRAINT tournaments_pkey PRIMARY KEY (id)
) WITH (
  OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS bets (
  id character varying(255) NOT NULL,
  player_id character varying(255),
  tournament_id character varying(255),
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  CONSTRAINT bets_pkey PRIMARY KEY (id),
  CONSTRAINT bets_player_id_foreign
    FOREIGN KEY (player_id)
    REFERENCES players(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT bets_tournament_id_foreign
    FOREIGN KEY (tournament_id)
    REFERENCES tournaments(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE IF NOT EXISTS backers (
  bet_id character varying(255),
  player_id character varying(255),
  CONSTRAINT backers_pkey PRIMARY KEY (bet_id, player_id),
  CONSTRAINT backers_player_id_foreign
    FOREIGN KEY (player_id)
    REFERENCES players(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT backers_bet_id_foreign
    FOREIGN KEY (bet_id)
    REFERENCES bets(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED
);

CREATE UNIQUE INDEX bets_player_id_tournament_id ON bets USING btree (player_id, tournament_id);

CREATE UNIQUE INDEX backers_bet_id_player_id ON backers USING btree (bet_id, player_id);

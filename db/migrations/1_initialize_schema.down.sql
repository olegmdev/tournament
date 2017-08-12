DROP INDEX IF EXISTS bets_player_id_tournament_id;
DROP INDEX IF EXISTS backers_bet_id_player_id;

DROP TABLE IF EXISTS backers;
DROP TABLE IF EXISTS bets;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS tournaments;

DROP TYPE IF EXISTS enum_status;

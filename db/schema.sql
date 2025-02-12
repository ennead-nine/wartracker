-- Warzones
CREATE TABLE IF NOT EXISTS "warzone"(
  "id" TEXT NOT NULL UNIQUE,
  "server" INTEGER NOT NULL,
  PRIMARY KEY("id")
)
-- Alliances
CREATE TABLE IF NOT EXISTS "alliance"(
  "id"	TEXT NOT NULL UNIQUE,
  "warzone_id"	TEXT,
  "tag"	TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("warzone_id") REFERENCES "warzone"("id"),
);
CREATE TABLE IF NOT EXISTS "alliance_data"(
  "name"	TEXT,
  "date"	TEXT NOT NULL,
  "power"	INTEGER,
  "gift_level"	INTEGER,
  "member_count"	INTEGER,
  "r5_id"	TEXT,
  "alliance_id"	TEXT NOT NULL,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("r5_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "alliance_alias"(
  "alias"	TEXT,
  "tag"	TEXT,
  "preferred"	BOOLEAN NOT NULL 
    CHECK (
      "preferred" IN (0, 1)
    ),
  "alliance_id"	TEXT
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id")
);
-- Commanders
CREATE TABLE IF NOT EXISTS "commander"(
  "id"	TEXT NOT NULL,
  "name"	TEXT,
  "warzone_id"	TEXT,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "commander_data"(
  "date"	TEXT NOT NULL,
  "pfp"	BLOB,
  "hq_level"	INTEGER,
  "likes"	INTEGER,
  "hq_power"	INTEGER,
  "kills"	INTEGER,
  "profession_level"	INTEGER,
  "total_hero_power"	INTEGER,
  "commander_id"	TEXT,
  "alliance_id"	TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("commander_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "commander_alias"(
  "alias"	TEXT,
  "tag"	TEXT,
  "preferred"	BOOLEAN NOT NULL 
    CHECK (
      "preferred" IN (0, 1)
    ),
  "commander_id"	TEXT,
  FOREIGN KEY("commander_id") REFERENCES "commander"("id")
);
-- VsDuels
CREATE TABLE IF NOT EXISTS "vsduel"(
  "id"	TEXT NOT NULL,
  "date"	TEXT NOT NULL,
  "league_level"	TEXT,
  "league_id"	TEXT,
  "tournament_id" TEXT,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_week"(
  "vsduel_id" TEXT NOT NULL,
  "week"	INTEGER,
  "alliance1_id" TEXT NOT NULL,
  "alliance2_id" TEXT NOT NULL,
  FOREIGN KEY("vsduel_id") REFERENCES "vsduel"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_day"(
  "id"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "short_name"	TEXT NOT NULL,
  "day_of_week"	TEXT NOT NULL,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_data"(
  "id"	TEXT NOT NULL,
  "vsduel_day_id"	TEXT NOT NULL,
  "vsduel_id"	TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("vsduel_day_id") REFERENCES "vsduel_day"("id"),
  FOREIGN KEY("vsduel_id") REFERENCES "vsduel"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_alliance"(
  "points"	INTEGER,
  "alliance_id"	TEXT,
  "vsduel_data_id"	TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("vsduel_data_id") REFERENCES "vsduel_data"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_commander"(
  "points"	INTEGER NOT NULL,
  "rank"	INTEGER NOT NULL,
  "commander_id"	TEXT,
  "vsduel_data_id"	TEXT NOT NULL,
  FOREIGN KEY("commander_id") REFERENCES "commander"("id"),
  FOREIGN KEY("vsduel_data_id") REFERENCES "vsduel_data"("id")
);






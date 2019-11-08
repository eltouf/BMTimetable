CREATE TABLE IF NOT EXISTS  `sv_arret_p` (
    `gid` INTEGER PRIMARY KEY,
    `geo_point_2d` TEXT NOT NULL,
    `ident` TEXT NOT NULL,
    `group` TEXT NOT NULL,
    `label` TEXT NOT NULL,
    `city` TEXT NOT NULL,
    `zipcode` TEXT NOT NULL,
    `create_time` INTEGER NOT NULL,
    `update_time` INTEGER NOT NULL
);
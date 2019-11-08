.open --new ":memory:"

.mode csv
.separator ;
.import "tmp/2019/2019-10-18/sv_arret_p/sv_arret_p.csv" sv_arret_p


ATTACH DATABASE "database/bmtimetable.db" AS bmtimetable;

INSERT OR IGNORE INTO bmtimetable.sv_arret_p 
SELECT `gid`, `geo_point_2d`, `ident`, `groupe`, `libelle`, `commune`, `code_commune`, `cdate`, `mdate` FROM main.sv_arret_p;
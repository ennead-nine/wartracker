-- All Data
SELECT vcd.points, ca.alias, a.tag, vd.day_of_week 
FROM vsduel_commander_data AS vcd
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
INNER JOIN commander AS c
	ON vcd.commander_id=c.id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
WHERE vw.vsweek_number=3 AND
	ca.preferred=1
ORDER BY vdy.day_number, vcd.points DESC

--P4K Weekly Total
SELECT 
	SUM(vcd.points), 
	ca.alias
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=3
AND a.tag="P4K"
AND ca.preferred=1
GROUP BY ca.alias

--P4K Individual Daily
SELECT 
	vcd.points,
	ca.alias
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=3
AND a.tag="P4K"
AND ca.preferred=1
AND vd.day_of_week="Saturday"
ORDER BY vcd.points

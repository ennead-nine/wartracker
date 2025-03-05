--Individual Daily Top 10
SELECT 
	ca.alias,
	a.tag,
	vcd.points
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
AND vd.day_of_week="Saturday"
ORDER BY vcd.points DESC
LIMIT 10

--Individual Daily Bottom 10
SELECT 
	ca.alias,
	a.tag,
	vcd.points
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
AND vd.day_of_week="Saturday"
ORDER BY vcd.points
LIMIT 10

--Individual Daily Under 7.2
SELECT 
	ca.alias,
	a.tag,
	vcd.points
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
AND vd.day_of_week="Saturday"
AND vcd.points < 7200000
ORDER BY vcd.points

--Individual Weekly Top 10
SELECT 
	ca.alias,
	ca.tag,
	SUM(vcd.points) AS t
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
GROUP BY ca.alias
ORDER BY t DESC
LIMIT 10

--Individual Weekly Bottom 10
SELECT 
	ca.alias,
	ca.tag,
	SUM(vcd.points) AS t
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
GROUP BY ca.alias
ORDER BY t
LIMIT 10

--Individual Weekly Under 45000000
SELECT 
	ca.alias,
	ca.tag,
	SUM(vcd.points) AS t
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
GROUP BY ca.alias
ORDER BY t
-- League Info
SELECT * FROM vsduel

-- All Data
SELECT vcd.points, ca.alias, a.tag, vd.day_of_week 
FROM vsduel_commander_data AS vcd
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
INNER JOIN commander AS c
	ON vcd.commander_id=c.id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
WHERE vw.vsweek_number=3
ORDER BY vd.day_of_week, vcd.points DESC

--Alliance Data
SELECT vad.points, vad.vsduel_points, a.tag, vd.day_of_week 
FROM vsduel_alliance_data AS vad
INNER JOIN alliance AS a
	ON vad.alliance_id=a.id
INNER JOIN vsduel_data AS vd
	ON vad.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
WHERE vw.vsweek_number=3

--Individual Totals
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

--Individual Daily
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
AND vd.day_of_week="Thursday"
ORDER BY vcd.points

--Individual Daily Top 10
SELECT 
	ca.alias,
	a.tag,
	vcd.points
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
AND ca.preferred=1
AND vd.day_of_week="Wednesday"
ORDER BY vcd.points
--LIMIT 10

--Individual Daily Bottom 10
SELECT 
	ca.alias,
	a.tag,
	COUNT(*)
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=1
AND ca.preferred=1
AND a.tag="P4K"
GROUP BY ca.alias
ORDER BY ca.alias


--Individual Daily Under 7.2
SELECT 
	ca.alias,
	a.tag,
	vcd.points,
	vd.day_of_week
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
LEFT JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=3
AND ca.preferred=1
AND vcd.points < 7200000
ORDER BY ca.commander_id, vd.day_of_week

--Individual Total Average Under 7.2
SELECT 
	ca.alias,
	a.tag,
	AVG(vcd.points) AS average
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
AND ca.preferred=1
GROUP BY ca.alias
ORDER BY average DESC
LIMIT 10


SELECT DISTINCT LOWER(ca.alias)
FROM commander_alias AS ca
INNER JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN alliance AS a
	ON vcd.alliance_id=a.id
WHERE vw.vsweek_number=1
AND a.tag="P4K"
AND ca.preferred=1
ORDER BY ca.alias


--Individual Total Top 10
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total, vsduel_week.vsweek_number
FROM commander, vsduel_commander_data, alliance, vsduel_week,vsduel_data
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_week.vsweek_number=2
GROUP BY commander.name, vsduel_week.vsweek_number
ORDER BY total DESC
LIMIT 10

--Individual Bottom 10
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total
FROM commander, vsduel_commander_data, alliance
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
GROUP BY commander.name
ORDER BY total
LIMIT 10

--Individual {DOW} Top 10
SELECT commander.name, alliance.tag AS alliance, vsduel_commander_data.points
FROM commander, vsduel_commander_data, alliance, vsduel_data, vsduel_week
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND vsduel_data.day_of_week='Monday'
AND vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_week.vsweek_number=2
GROUP BY commander.name
ORDER BY points

SELECT commander.name, alliance.tag AS alliance, vsduel_commander_data.points
FROM commander, vsduel_commander_data, alliance, vsduel_data, vsduel_week
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_week.vsweek_number=2
GROUP BY commander.name
ORDER BY points
LIMIT 10

--Individual {DOW} Bottom 10
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total
FROM commander, vsduel_commander_data, alliance, vsduel_data
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND vsduel_data.day_of_week='Saturday'
GROUP BY commander.name
ORDER BY total
LIMIT 10

--Individual Under 7.2 Million
SELECT commander.name, alliance.tag AS alliance,  AVG(vsduel_commander_data.points) AS average
FROM commander, vsduel_commander_data, alliance, vsduel_data, vsduel_week
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND (vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_week.vsweek_number=2)
GROUP BY commander.name
ORDER BY average


--Individual By {Alliance}/{Day}
SELECT commander.name, vsduel_commander_data.points
FROM vsduel_commander_data, alliance, commander, vsduel_data 
WHERE vsduel_commander_data.alliance_id=alliance.id 
AND vsduel_commander_data.commander_id=commander.id 
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND alliance.tag='szb'
AND vsduel_data.day_of_week='Monday'

SELECT commander.name, alliance.tag, vsduel_commander_data.points, vsduel_week.vsweek_number, vsduel_data.day_of_week
FROM commander, alliance, vsduel_commander_data, vsduel_week, vsduel_data
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data_id
AND vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_commander_data.alliance_id=alliance.id
AND alliance.tag="P4K"

SELECT commander.name
FROM commander, commander_data, alliance
WHERE commander_data.commander_id=commander.id
AND commander_data.alliance_id=alliance.id
AND alliance.tag="P4K"
ORDER BY commander.name


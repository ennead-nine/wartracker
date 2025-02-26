-- League Info
SELECT * FROM vsduel

-- All Data
SELECT vsduel_commander_data.points, commander.name, alliance.tag, vsduel_data.day_of_week 
FROM vsduel_commander_data, alliance, commander, vsduel_data, vsduel_week
WHERE vsduel_commander_data.alliance_id=alliance.id 
AND vsduel_commander_data.commander_id=commander.id 
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_week.vsweek_number=2

--Alliance Data
SELECT vsduel_alliance_data.points, vsduel_alliance_data.vsduel_points, alliance.tag, vsduel_data.day_of_week 
FROM vsduel_alliance_data, alliance, vsduel_data, vsduel_week
WHERE vsduel_alliance_data.alliance_id=alliance.id
AND vsduel_alliance_data.vsduel_data_id=vsduel_data.id
AND (vsduel_data.vsduel_week_id=vsduel_week.id
AND vsduel_week.vsweek_number=2)

--Individual Totals
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total
FROM commander, vsduel_commander_data, alliance
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
GROUP BY commander.name
ORDER BY total DESC

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


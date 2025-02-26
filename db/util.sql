-- League Info
SELECT * FROM vsduel

-- All Data
SELECT vsduel_commander_data.points, commander.name, alliance.tag, vsduel_data.day_of_week 
FROM vsduel_commander_data, alliance, commander, vsduel_data 
WHERE vsduel_commander_data.alliance_id=alliance.id 
AND vsduel_commander_data.commander_id=commander.id 
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id

--Alliance Data
SELECT vsduel_alliance_data.points, vsduel_alliance_data.vsduel_points, alliance.tag, vsduel_data.day_of_week 
FROM vsduel_alliance_data, alliance, vsduel_data
WHERE vsduel_alliance_data.alliance_id=alliance.id
AND vsduel_alliance_data.vsduel_data_id=vsduel_data.id

--Individual Totals
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total
FROM commander, vsduel_commander_data, alliance
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
GROUP BY commander.name
ORDER BY total DESC

--Individual Total Top 10
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total
FROM commander, vsduel_commander_data, alliance
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
GROUP BY commander.name
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
SELECT commander.name, alliance.tag AS alliance, SUM(vsduel_commander_data.points) AS total
FROM commander, vsduel_commander_data, alliance, vsduel_data
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
AND vsduel_data.day_of_week='Saturday'
GROUP BY commander.name
ORDER BY total DESC
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
SELECT commander.name, alliance.tag AS alliance, AVG(vsduel_commander_data.points) AS average
FROM commander, vsduel_commander_data, alliance, vsduel_data
WHERE vsduel_commander_data.commander_id=commander.id
AND vsduel_commander_data.alliance_id=alliance.id
AND vsduel_commander_data.vsduel_data_id=vsduel_data.id
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


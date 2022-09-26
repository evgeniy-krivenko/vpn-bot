package repository

const (
	GetConnectsByUserId = `SELECT connections.id, server_id, port, u.user_id, encrypted_secret, s.ip_address as "ip_address", s.location as "location" FROM %s LEFT JOIN servers s on s.id = connections.server_id LEFT JOIN users u on connections.user_id = u.id WHERE u.user_id=$1;`
	getAllServersSQL    = `SELECT * FROM %s;`
	getServerSQL        = `SELECT * FROM %s WHERE id=$1`
	getPortWithCountSQL = `SELECT port, count(port) FROM %s GROUP BY port ORDER BY port DESC LIMIT 1;`
	createConnectionSQL = `INSERT INTO %s (port, encrypted_secret, user_id, server_id) VALUES ($1, $2, $3, $4) RETURNING id;`
)

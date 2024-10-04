package sqlstorage

const (
	eventsCommandCreateTable string = `
		CREATE TABLE IF NOT EXISTS 'events' (
			id TEXT NOT NULL UNIQUE COLLATE NOCASE,
			title TEXT NOT NULL,
			date_time INTEGER UNIQUE DEFAULT (strftime('%s', 'now')),
			duration INTEGER NOT NULL,
			description TEXT,
			user_id TEXT NOT NULL COLLATE NOCASE,
			alarm_until INTEGER DEFAULT 0
		);
	`
	eventsCommandAdd string = `
		INSERT OR IGNORE INTO 'events'
			(id, title, date_time, duration, description, user_id, alarm_until)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	eventsCommandList string = `
		SELECT id, title, date_time, duration, description, user_id, alarm_until
		FROM 'events'
		ORDER BY date_time ASC;
	`

	eventsCommandRemove string = `
		DELETE FROM 'events'
		WHERE id = ?;
	`

	eventsCommandUpdate string = `
	UPDATE 'events'
	SET
		title = ?,
		date_time = ?,
		duration = ?,
		description = ?,
		user_id = ?,
		alarm_until = ?,
	WHERE id = ?;
	`
)

package fixtures

import (
	"database/sql"
)

func PopulateHistory(db *sql.DB) ([]byte, error) {

	_, err := db.Exec(`
INSERT INTO history (event_time, expression, result)
VALUES ('2020-11-22T14:00:00+03:00', '2 + 2', '4'),
       ('2020-11-22T14:10:00+03:00', '3 + 3', '6'),
       ('2020-11-22T14:20:00+03:00', '4 + 4', '8'),
       ('2020-11-22T14:30:00+03:00', '5 + 5', '10'),
       ('2020-11-22T14:40:00+03:00', '6 + 6', '12');`)
	return allHistory, err
}

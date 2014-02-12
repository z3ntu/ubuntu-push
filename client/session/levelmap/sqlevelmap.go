/*
 Copyright 2013-2014 Canonical Ltd.

 This program is free software: you can redistribute it and/or modify it
 under the terms of the GNU General Public License version 3, as published
 by the Free Software Foundation.

 This program is distributed in the hope that it will be useful, but
 WITHOUT ANY WARRANTY; without even the implied warranties of
 MERCHANTABILITY, SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR
 PURPOSE.  See the GNU General Public License for more details.

 You should have received a copy of the GNU General Public License along
 with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package levelmap

import (
	_ "code.google.com/p/gosqlite/sqlite3"
	"database/sql"
	"fmt"
)

type sqliteLevelMap struct {
	db *sql.DB
}

// NewSqliteLevelMap returns an implementation of LevelMap that
// persists the map in an sqlite database.
func NewSqliteLevelMap(filename string) (LevelMap, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite level map %#v: %v", filename, err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS level_map (level text primary key, top integer)")
	if err != nil {
		return nil, fmt.Errorf("cannot (re)create sqlite level map table: %v", err)
	}
	return &sqliteLevelMap{db}, nil
}

func (pm *sqliteLevelMap) Set(level string, top int64) error {
	_, err := pm.db.Exec("REPLACE INTO level_map (level, top) VALUES (?, ?)", level, top)
	if err != nil {
		return fmt.Errorf("cannot set %#v to %#v in level map: %v", level, top, err)
	}
	return nil
}
func (pm *sqliteLevelMap) GetAll() (map[string]int64, error) {
	rows, err := pm.db.Query("SELECT * FROM level_map")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve levels from sqlite level map: %v", err)
	}
	m := map[string]int64{}
	for rows.Next() {
		var level string
		var top int64
		err = rows.Scan(&level, &top)
		if err != nil {
			return nil, fmt.Errorf("cannot read level from sqlite level map: %v", err)
		}
		m[level] = top
	}
	return m, nil
}
package store

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
)

// TODO adjust locationId when changed on local db
const (
	getDeviceStmt = `SELECT uid as DeviceUID, name as DeviceName, tag as DeviceTag, 
       		groupUID, companyUID, parentUID
		FROM device
		WHERE type='power' AND device.uid = ?`
)

// DeviceSql is accessed both directly and by the device Cache, direct access is only for uplinks which show Compliance events
type DeviceSql struct {
	sqlDB          *sql.DB
	getDeviceByUID *sql.Stmt
}

func NewDeviceSql(db *sql.DB) (*DeviceSql, error) {
	sqlDB := &DeviceSql{
		sqlDB: db,
	}
	var err error

	if sqlDB.getDeviceByUID, err = db.Prepare(getDeviceStmt); err != nil {
		return nil, err
	}

	return sqlDB, nil
}

func (db DeviceSql) GetDevice(uid string) (*lib.Device, error) {
	log.Debug().Str("uid", uid).Msg("getting device from sql")
	row := db.getDeviceByUID.QueryRow(uid)

	serverDevice, err := scanDevice(row)
	if err != nil {
		return nil, err
	}

	return serverDevice, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanDevice(s rowScanner) (*lib.Device, error) {
	var (
		name       sql.NullString
		uid        sql.NullString
		tag        sql.NullString
		groupUID   sql.NullString
		companyUID sql.NullString
		parentUID  sql.NullString
	)

	err := s.Scan(&name, &uid, &tag, &groupUID, &companyUID, &parentUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	deviceInfo := lib.Device{
		DeviceUID: uid.String,
		DeviceMeta: &lib.DeviceMeta{
			DeviceName: name.String,
			DeviceTag:  tag.String,
		},
		Group: &lib.Group{
			GroupUID:   groupUID.String,
			CompanyUID: companyUID.String,
		},
	}

	return &deviceInfo, nil
}

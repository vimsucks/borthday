package birthday

import (
	"github.com/elgris/sqrl"
	"time"
)

type Birthday struct {
	ID            int64     `db:"id"`
	UID           int64     `db:"uid"`
	Name          string    `db:"name"`           // 名字
	LunarBirthday time.Time `db:"lunar_birthday"` // 农历生日
	SolarBirthday time.Time `db:"solar_birthday"` // 阳历生日
}

func CreateBirthday(birthday *Birthday) error {
	sql, args, err := sqrl.
		Insert("birthday").
		Columns("uid", "name", "lunar_birthday", "solar_birthday").
		Values(birthday.UID, birthday.Name, birthday.LunarBirthday, birthday.SolarBirthday).ToSql()
	if err != nil {
		return err
	}
	ret, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}
	birthday.ID, err = ret.LastInsertId()
	return err
}

func UpdateBirthday(birthday *Birthday) error {
	sql, args, err := sqrl.
		Update("birthday").
		Set("uid", birthday.UID).
		Set("name", birthday.Name).
		Set("lunar_birthday", birthday.LunarBirthday).
		Set("solar_birthday", birthday.SolarBirthday).
		Where(sqrl.Eq{"id": birthday.ID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = db.Exec(sql, args...)
	return err
}

func GetBirthday(query sqrl.Eq) ([]Birthday, error) {
	sqlBuilder := sqrl.Select("id", "uid", "name", "lunar_birthday", "solar_birthday").From("birthday")
	if len(query) != 0 {
		sqlBuilder = sqlBuilder.Where(query)
	}
	sql, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var birthday []Birthday
	err = db.Select(&birthday, sql, args...)
	if err != nil {
		return nil, err
	}
	return birthday, nil
}

func GetBirthdayByUser(uid int64, name string) (*Birthday, error) {
	query := sqrl.Eq{}
	query["uid"] = uid
	query["name"] = name
	sqlBuilder := sqrl.Select("id", "uid", "name", "lunar_birthday", "solar_birthday").From("birthday")
	if len(query) != 0 {
		sqlBuilder = sqlBuilder.Where(query)
	}
	sql, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var birthday []*Birthday
	err = db.Select(&birthday, sql, args...)
	if err != nil {
		return nil, err
	}
	if len(birthday) == 0 {
		return nil, nil
	}
	return birthday[0], nil
}

// from, to 格式： 10-01（%m-%d）
func GetBirthdaySolarBetween(from, to string) ([]Birthday, error) {
	sql, args, err := sqrl.Select("id", "uid", "name", "lunar_birthday", "solar_birthday").
		From("birthday").
		Where("strftime('%m-%d', solar_birthday) between ? and ?", from, to).
		ToSql()
	if err != nil {
		return nil, err
	}
	var birthday []Birthday
	err = db.Select(&birthday, sql, args...)
	if err != nil {
		return nil, err
	}
	return birthday, nil
}

// from, to 格式： 10-01（%m-%d）
func GetBirthdayLunarBetween(from, to string) ([]Birthday, error) {
	sql, args, err := sqrl.Select("id", "uid", "name", "lunar_birthday", "solar_birthday").
		From("birthday").
		Where("strftime('%m-%d', lunar_birthday) between ? and ?", from, to).
		ToSql()
	if err != nil {
		return nil, err
	}
	var birthday []Birthday
	err = db.Select(&birthday, sql, args...)
	if err != nil {
		return nil, err
	}
	return birthday, nil
}

func DeleteBirthdayById(id int64) error {
	sql, args, err := sqrl.Delete("birthday").Where(sqrl.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = db.Exec(sql, args...)
	return err
}

func DeleteBirthday(query sqrl.Eq) error {
	sql, args, err := sqrl.Delete("birthday").Where(query).ToSql()
	if err != nil {
		return err
	}
	_, err = db.Exec(sql, args...)
	return err
}

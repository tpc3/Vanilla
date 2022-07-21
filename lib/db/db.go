package db

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tpc3/Vanilla/lib/config"
)

var DB *sql.DB
var getDBverStmt *sql.Stmt
var setDBverStmt *sql.Stmt
var loadGuildStmt *sql.Stmt
var insertGuildStmt *sql.Stmt
var updateGuildStmt *sql.Stmt

var addLogStmt map[string]*sql.Stmt
var updateLogEmojiStmt map[string]*sql.Stmt
var cleanOldLogStmt map[string]*sql.Stmt
var rankingStmt map[string]*sql.Stmt
var rankingInvertStmt map[string]*sql.Stmt
var addEmojiStmt map[string]*sql.Stmt
var getEmojisStmt map[string]*sql.Stmt
var getEmojiStmt map[string]*sql.Stmt
var getEmojiByNameStmt map[string]*sql.Stmt
var syncValueStmt map[string]*sql.Stmt
var deleteEmojiStmt map[string]*sql.Stmt
var deleteLogStmt map[string]*sql.Stmt

const db_version = 1

func init() {
	var err error
	DB, err = sql.Open(config.CurrentConfig.Db.Kind, config.CurrentConfig.Db.Parameter)
	if err != nil {
		log.Fatal("DB load error:", err)
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS " + config.CurrentConfig.Db.Tableprefix + "guilds (" +
		"id BIGINT NOT NULL PRIMARY KEY," +
		"db_version INT DEFAULT " + strconv.Itoa(db_version) + "," +
		"prefix VARCHAR," +
		"lang VARCHAR," +
		"msgweight INT," +
		"newreactweight INT," +
		"addreactweight INT)")
	if err != nil {
		log.Fatal("Create guild table error:", err)
	}
	_, err = DB.Exec("ALTER TABLE " + config.CurrentConfig.Db.Tableprefix + "guilds DROP COLUMN bots")
	if err == nil {
		log.Print("WARN: Database update succesfully")
	} else {
		log.Print("FINE: Guilds DB up-to-date!")
	}
	getDBverStmt, err = DB.Prepare("SELECT db_version FROM " + config.CurrentConfig.Db.Tableprefix + "guilds WHERE " + "id = ?")
	if err != nil {
		log.Fatal("Prepere getDBverStmt error:", err)
	}
	setDBverStmt, err = DB.Prepare("UPDATE " + config.CurrentConfig.Db.Tableprefix + "guilds SET db_version = ? WHERE id = ?")
	if err != nil {
		log.Fatal("Prepere setDBverStmt error:", err)
	}
	loadGuildStmt, err = DB.Prepare("SELECT * FROM " + config.CurrentConfig.Db.Tableprefix + "guilds WHERE " + "id = ?")
	if err != nil {
		log.Fatal("Prepere loadGuildStmt error:", err)
	}
	insertGuildStmt, err = DB.Prepare("INSERT INTO " + config.CurrentConfig.Db.Tableprefix + "guilds(id,prefix,lang,msgweight,newreactweight,addreactweight) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepere insertGuildStmt error:", err)
	}
	updateGuildStmt, err = DB.Prepare("UPDATE " + config.CurrentConfig.Db.Tableprefix + "guilds " +
		"SET db_version = ?, prefix = ?, lang = ?, msgweight = ?, newreactweight = ?, addreactweight = ? " +
		"WHERE id = ?")
	if err != nil {
		log.Fatal("Prepere updateGuildStmt error:", err)
	}
	addLogStmt = map[string]*sql.Stmt{}
	updateLogEmojiStmt = map[string]*sql.Stmt{}
	cleanOldLogStmt = map[string]*sql.Stmt{}
	rankingStmt = map[string]*sql.Stmt{}
	rankingInvertStmt = map[string]*sql.Stmt{}
	addEmojiStmt = map[string]*sql.Stmt{}
	getEmojisStmt = map[string]*sql.Stmt{}
	getEmojiStmt = map[string]*sql.Stmt{}
	getEmojiByNameStmt = map[string]*sql.Stmt{}
	syncValueStmt = map[string]*sql.Stmt{}
	deleteEmojiStmt = map[string]*sql.Stmt{}
	deleteLogStmt = map[string]*sql.Stmt{}
}

func Close() {
	err := DB.Close()
	if err != nil {
		log.Fatal("DB close error:", err)
	}
}

func LoadGuild(id *string) *config.Guild {
	val, exists := guildCache[*id]
	if exists {
		return val
	}
	rows, err := getDBverStmt.Query(id)
	if err != nil {
		log.Fatal("GetDBver query error:" + err.Error())
	}
	emojisTable := config.CurrentConfig.Db.Tableprefix + *id + "_emojis"
	logsTable := config.CurrentConfig.Db.Tableprefix + *id + "_logs"
	var dbVersion int
	if rows.Next() {
		err := rows.Scan(&dbVersion)
		rows.Close()
		if err != nil {
			log.Fatal("GetDBver Scan error:" + err.Error())
		}
		if dbVersion <= 1 {
			log.Print("WARN: Updating guild from version 1")
			_, err = DB.Exec("ALTER TABLE " + logsTable + " ADD COLUMN bot BIT DEFAULT 0 NOT NULL")
			if err != nil {
				log.Fatal("Update log_DB from version 1 error 01: ", err)
			}
			_, err = DB.Exec("ALTER TABLE " + logsTable + " ADD COLUMN userid BIGINT DEFAULT 0 NOT NULL")
			if err != nil {
				log.Fatal("Update log_DB from version 1 error 02: ", err)
			}
			_, err = DB.Exec("ALTER TABLE " + logsTable + " ADD COLUMN channelid BIGINT DEFAULT 0 NOT NULL")
			if err != nil {
				log.Fatal("Update log_DB from version 1 error 03: ", err)
			}
			_, err = DB.Exec("ALTER TABLE " + logsTable + " ADD COLUMN messageid BIGINT DEFAULT 0 NOT NULL")
			if err != nil {
				log.Fatal("Update log_DB from version 1 error 04: ", err)
			}
			_, err := setDBverStmt.Exec(2, id)
			if err != nil {
				log.Fatal("Update guilds_DB from version 1 error: ", err)
			}
		}
	} else {
		rows.Close()
		log.Print("WARN: Guild not found, making row.")
		_, err = DB.Exec("CREATE TABLE IF NOT EXISTS " + emojisTable + " (" +
			"id BIGINT NOT NULL PRIMARY KEY," +
			"name VARCHAR," +
			"description VARCHAR)")
		if err != nil {
			log.Fatal("Create Emoji Table error:", err)
		}
		_, err = DB.Exec("CREATE TABLE IF NOT EXISTS " + logsTable + " (" +
			"emoji BIGINT NOT NULL," +
			"type TINYINT NOT NULL," +
			"value INT NOT NULL," +
			"timeat BIGINT NOT NULL," +
			"bot BIT DEFAULT 0 NOT NULL," +
			"userid BIGINT DEFAULT 0 NOT NULL," +
			"channelid BIGINT DEFAULT 0 NOT NULL," +
			"messageid BIGINT DEFAULT 0 NOT NULL)")
		if err != nil {
			log.Fatal("Create log table error:", err)
		}
		guild := config.CurrentConfig.Guild
		_, err = insertGuildStmt.Exec(id, guild.Prefix, guild.Lang, guild.Weight.Message, guild.Weight.Reactnew, guild.Weight.Reactadd)
		if err != nil {
			log.Fatal("LoadGuild insert error:" + err.Error())
		}
	}
	var guild config.Guild
	var guildID int64
	rows, err = loadGuildStmt.Query(id)
	if err != nil {
		log.Fatal("LoadGuild query error:" + err.Error())
	}
	defer rows.Close()
	if !rows.Next() {
		log.Fatal("LoadGuild next returned false")
	}
	err = rows.Scan(&guildID, &dbVersion, &guild.Prefix, &guild.Lang, &guild.Weight.Message, &guild.Weight.Reactnew, &guild.Weight.Reactadd)
	if err != nil {
		log.Fatal("LoadGuild scan error:" + err.Error())
	}
	defer rows.Close()

	addLogStmt[*id], err = DB.Prepare("INSERT INTO " + logsTable + "(emoji,type,value,timeat,bot,userid,channelid,messageid) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepere addLogStmt error: ", err)
	}
	updateLogEmojiStmt[*id], err = DB.Prepare("UPDATE " + logsTable + " " +
		"SET emoji = ? " +
		"WHERE emoji = ?")
	if err != nil {
		log.Fatal("Prepere updateLogEmojiStmt error: ", err)
	}
	cleanOldLogStmt[*id], err = DB.Prepare("DELETE FROM " + logsTable + " " +
		"WHERE timeat < ?")
	if err != nil {
		log.Fatal("Prepere cleanOldLogStmt error: ", err)
	}
	rankingStmt[*id], err = DB.Prepare("SELECT emojis.id,emojis.name,emojis.description,sum(logs.value) " +
		"FROM " + emojisTable + " emojis " +
		"LEFT OUTER JOIN " + logsTable + " logs " +
		"ON emojis.id = logs.emoji AND logs.timeat > ? " +
		"AND logs.bot <> ? " +
		"GROUP BY emojis.id " +
		"ORDER BY sum(logs.value) DESC " +
		"LIMIT ?,?")
	if err != nil {
		log.Fatal("Prepere rankingStmt error: ", err)
	}
	rankingInvertStmt[*id], err = DB.Prepare("SELECT emojis.id,emojis.name,emojis.description,sum(logs.value) " +
		"FROM " + emojisTable + " emojis " +
		"LEFT OUTER JOIN " + logsTable + " logs " +
		"ON emojis.id = logs.emoji AND logs.timeat > ? " +
		"AND logs.bot <> ? " +
		"GROUP BY emojis.id " +
		"ORDER BY sum(logs.value) ASC " +
		"LIMIT ?,?")
	if err != nil {
		log.Fatal("Prepere rankingInvertStmt error: ", err)
	}
	addEmojiStmt[*id], err = DB.Prepare("INSERT INTO " + emojisTable + "(id,name,description) VALUES(?,?,?)")
	if err != nil {
		log.Fatal("Prepere addEmojiStmt error: ", err)
	}
	deleteEmojiStmt[*id], err = DB.Prepare("DELETE " +
		"FROM " + emojisTable + " " +
		"WHERE id = ?")
	if err != nil {
		log.Fatal("Prepere deleteEmojiStmt error: ", err)
	}
	getEmojisStmt[*id], err = DB.Prepare("SELECT * " +
		"FROM " + emojisTable)
	if err != nil {
		log.Fatal("Prepere getEmojisStmt error: ", err)
	}
	getEmojiStmt[*id], err = DB.Prepare("SELECT * " +
		"FROM " + emojisTable + " " +
		"WHERE id = ?")
	if err != nil {
		log.Fatal("Prepere getEmojiStmt error: ", err)
	}
	getEmojiByNameStmt[*id], err = DB.Prepare("SELECT * " +
		"FROM " + emojisTable + " " +
		"WHERE name = ?")
	if err != nil {
		log.Fatal("Prepere getEmojiByNameStmt error: ", err)
	}
	syncValueStmt[*id], err = DB.Prepare("UPDATE " + logsTable + " " +
		"SET value = ? " +
		"WHERE type = ? AND value <> ?")
	if err != nil {
		log.Fatal("Prepere syncValueStmt error: ", err)
	}
	deleteLogStmt[*id], err = DB.Prepare("DELETE " +
		"FROM " + logsTable + " " +
		"WHERE emoji = ?")
	if err != nil {
		log.Fatal("Prepere deleteLogStmt error: ", err)
	}
	guildCache[*id] = &guild
	return &guild
}

func SaveGuild(id *string, guild *config.Guild) error {
	res, err := updateGuildStmt.Exec(1, guild.Prefix, guild.Lang, guild.Weight.Message, guild.Weight.Reactnew, guild.Weight.Reactadd, *id)
	res.RowsAffected()
	log.Print("WARN: SaveGuild error:", res)
	if err != nil {
		log.Print("WARN: SaveGuild error:", err)
	} else {
		delete(guildCache, *id)
	}
	return err
}

const (
	MSG      = 1
	REACTNEW = 2
	REACTADD = 3
)

func AddLog(guildId *string, actionType int, emojiId *string, bot bool, userid, channelid, messageid *string) {
	guild := LoadGuild(guildId)
	var value int
	switch actionType {
	case MSG:
		value = guild.Weight.Message
	case REACTNEW:
		value = guild.Weight.Reactnew
	case REACTADD:
		value = guild.Weight.Reactadd
	}
	_, err := addLogStmt[*guildId].Exec(emojiId, actionType, value, time.Now().Unix(), bot, userid, channelid, messageid)
	if err != nil {
		log.Fatal("Insert log error: ", err)
	}
}

func ChangeLogID(guildId *string, oldEmojiId *string, newEmojiId *string) (sql.Result, error) {
	return updateLogEmojiStmt[*guildId].Exec(newEmojiId, oldEmojiId)
}

func CleanLogEmoji(guildId *string, validEmojis []string) (*int64, error) {
	valid := "("
	for _, v := range validEmojis {
		valid += v + ","
	}
	valid = strings.TrimSuffix(valid, ",") + ")"
	res, err := DB.Exec("DELETE " +
		"FROM " + config.CurrentConfig.Db.Tableprefix + *guildId + "_logs " +
		"WHERE emoji NOT IN " + valid)
	if err != nil {
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &affect, nil
}

func CleanOldLog(guildId *string) (*int64, error) {
	res, err := cleanOldLogStmt[*guildId].Exec(time.Now().Unix() - config.CurrentConfig.LogPeriod)
	if err != nil {
		return nil, err
	}
	num, err := res.RowsAffected()
	return &num, err
}

// return rows contains emoji,sum(value)
func GetRanking(guildId *string, limit int, limitOffset int, period int64, invert bool, no_bot bool, bot bool) (*sql.Rows, error) {
	botNum := 2
	if !no_bot {
		botNum = 0
	}
	if !bot {
		botNum = 1
	}
	if invert {
		return rankingInvertStmt[*guildId].Query(time.Now().Unix()-period, botNum, limitOffset, limit)
	} else {
		return rankingStmt[*guildId].Query(time.Now().Unix()-period, botNum, limitOffset, limit)
	}
}

func UpdateValue(guildId *string, value map[int]int) (*int64, error) {
	updated := int64(0)
	for i, v := range value {
		res, err := syncValueStmt[*guildId].Exec(v, i, v)
		if err != nil {
			return nil, err
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		updated += affect
	}
	return &updated, nil
}

func AddEmoji(guildId *string, emojiId string, emojiName string, description string) (sql.Result, error) {
	return addEmojiStmt[*guildId].Exec(emojiId, emojiName, description)
}

func DeleteEmoji(guildId *string, emojiId string) (sql.Result, error) {
	return deleteEmojiStmt[*guildId].Exec(emojiId)
}

func GetEmojis(guildId *string) (*sql.Rows, error) {
	return getEmojisStmt[*guildId].Query()
}

func GetEmoji(guildId *string, emojiId string) *sql.Row {
	return getEmojiStmt[*guildId].QueryRow(emojiId)
}

func GetEmojiByName(guildId *string, emojiName string) *sql.Row {
	return getEmojiByNameStmt[*guildId].QueryRow(emojiName)
}

func DeleteLogEmoji(guildId *string, emojiId string) (sql.Result, error) {
	return deleteLogStmt[*guildId].Exec(emojiId)
}

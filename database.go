package main

import (
	"database/sql"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Connection String

func toConnectionString(database Database) string {
	return database.User + ":" + database.Password + "@tcp(" + database.Address + ":" + strconv.Itoa(database.Port) + ")/" + database.Database
}

// Tables

// Channels
func createTableChannels() (res sql.Result, err error) {

	// Declare channels
	channels := [...]string{
		tableWelcome,
		tablePresentation,
		tableFame,
	}

	// Create them
	for _, channel := range channels {
		res, err = db.Exec("create table if not exists `" + channel + "` (`server` varchar(32) primary key, `channel` varchar(32) not null) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
		if err != nil {
			return
		}
	}

	return
}

// Roles
func createTableRoles() (res sql.Result, err error) {

	// Declare roles
	roles := [...]string{
		tableAdmin,
		tableMod,
		tableLight,
		tableAbsynthe,
		tableObsidian,
		tableShadow,
		tableEel,
		tableNPC,
	}

	// Create them
	for _, role := range roles {
		res, err = db.Exec("create table if not exists `" + role + "` (`server` varchar(32) primary key, `role` varchar(32) not null) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
		if err != nil {
			return
		}
	}

	return
}

// Self-Assignable Roles
func createTableSAR() (res sql.Result, err error) {
	return db.Exec("create table if not exists `" + tableSAR + "` (`server` varchar(32) not null, `role` varchar(32) not null, constraint `pk_roles_sar` primary key (`server`, `role`)) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
}

// Pins
func createTablePin() (res sql.Result, err error) {
	return db.Exec("create table if not exists `" + tablePins + "` (`server` varchar(32) not null, `channel` varchar(32) not null, `message` varchar(32) primary key, `member` varchar(32) not null) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
}

// Minimum Reactions
func createTableMinimumReactions() (res sql.Result, err error) {
	return db.Exec("create table if not exists `" + tableMinimumReactions + "` (`channel` varchar(32) primary key, `minimum` int not null) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
}

// Status
func createTableStatus() (res sql.Result, err error) {
	return db.Exec("create table if not exists `" + tableStatus + "` (`id` int auto_increment primary key, `status` varchar(32) not null) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
}

// Pins Count
func createViewPinsCount() (res sql.Result, err error) {
	return db.Exec("create or replace view `pins-count` as select `server`, `member`, count(`message`) as `count` from `pins` group by `server`, `member` order by `server` asc, `count` desc, `member` asc;")
}

// Messages in hall of fame
func createTableMessageFame() (res sql.Result, err error) {
	return db.Exec("create table if not exists `" + tableMessagesFamed + "` (`message` varchar(32) primary key) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
}

// Welcome Channel

// Select Welcome Channel
func selectWelcomeChannel(g *discordgo.Guild) (id string, err error) {
	err = db.QueryRow("select `channel` from `"+tableWelcome+"` where server = ?;", g.ID).Scan(&id)
	return
}

// Insert Welcome Channel
func insertWelcomeChannel(g *discordgo.Guild, c *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("insert into `"+tableWelcome+"`(`server`, `channel`) values(?, ?);", g.ID, c.ID)
}

// Update Welcome Channel
func updateWelcomeChannel(g *discordgo.Guild, c *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("update `"+tableWelcome+"` set `channel` = ? where `server` = ?;", c.ID, g.ID)
}

// Pins

// Select Pin
func selectPin(m *discordgo.Message) (channelID string, messageID string, err error) {
	err = db.QueryRow("select `channel`, `message` from `pins` where `message` = ?;", m.ID).Scan(&channelID, &messageID)
	return
}

// Select all pins from a guild
func selectPins(g *discordgo.Guild) (rows *sql.Rows, err error) {
	rows, err = db.Query("select `channel`, `message` from `pins` where `server` = ? order by `message` asc, `server` asc, `channel` asc;", g.ID)
	return
}

func selectLovers(g *discordgo.Guild) (members []string, err error) {

	// Query
	rows, err := db.Query("select `member` from `pins` where `server` = ? group by `member` order by count(`message`) desc, `member` asc;", g.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	// Select
	for rows.Next() {

		// Get a member
		var member string
		err = rows.Scan(&member)
		if err != nil {
			return
		}

		// Place it in a list
		members = append(members, member)
	}

	// Check for errors
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Insert Pin
func insertPin(g *discordgo.Guild, c *discordgo.Channel, m *discordgo.Message) (res sql.Result, err error) {
	return db.Exec("insert into `pins`(`server`, `channel`, `message`, `member`) values(?, ?, ?, ?)", g.ID, c.ID, m.ID, m.Author.ID)
}

// Delete Pin
func deletePin(m *discordgo.Message) (res sql.Result, err error) {
	return db.Exec("delete from `pins` where message = ?;", m.ID)
}

// Roles

// Select Role

func selectRole(g *discordgo.Guild, table string) (id string, err error) {
	err = db.QueryRow("select `role` from `"+table+"` where server = ?;", g.ID).Scan(&id)
	return
}

func selectRoleAdmin(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableAdmin)
}

func selectRoleMod(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableMod)
}

func selectRoleLight(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableLight)
}

func selectRoleAbsynthe(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableAbsynthe)
}

func selectRoleObsidian(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableObsidian)
}

func selectRoleShadow(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableShadow)
}

func selectRoleEel(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableEel)
}

func selectRoleNPC(g *discordgo.Guild) (id string, err error) {
	return selectRole(g, tableNPC)
}

// Insert Role

func insertRole(g *discordgo.Guild, r *discordgo.Role, table string) (res sql.Result, err error) {
	return db.Exec("insert into `"+table+"`(server, role) values(?, ?);", g.ID, r.ID)
}

func insertRoleAdmin(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableAdmin)
}

func insertRoleMod(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableMod)
}

func insertRoleLight(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableMod)
}

func insertRoleAbsynthe(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableMod)
}

func insertRoleObsidian(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableMod)
}

func insertRoleShadow(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {

	return insertRole(g, r, tableMod)
}
func insertRoleEel(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableMod)
}

func insertRoleNPC(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return insertRole(g, r, tableMod)
}

// Update Role

func updateRole(g *discordgo.Guild, r *discordgo.Role, table string) (res sql.Result, err error) {
	return db.Exec("update `"+table+"` set role = ? where server = ?;", r.ID, g.ID)
}

func updateRoleAdmin(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableAdmin)
}

func updateRoleMod(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableMod)
}

func updateRoleLight(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableLight)
}

func updateRoleAbsynthe(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableAbsynthe)
}

func updateRoleObsidian(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableObsidian)
}

func updateRoleShadow(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableShadow)
}

func updateRoleEel(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableEel)
}

func updateRoleNPC(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return updateRole(g, r, tableNPC)
}

// Presentation Channel

// Select Presentation Channel
func selectPresentationChannel(g *discordgo.Guild) (id string, err error) {
	err = db.QueryRow("select `channel` from `"+tablePresentation+"` where server = ?;", g.ID).Scan(&id)
	return
}

// Insert Presentation Channel
func insertPresentationChannel(g *discordgo.Guild, c *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("insert into `"+tablePresentation+"`(`server`, `channel`) values(?, ?);", g.ID, c.ID)
}

// Update Presentation Channel
func updatePresentationChannel(g *discordgo.Guild, c *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("update `"+tablePresentation+"` set `channel` = ? where `server` = ?;", c.ID, g.ID)
}

// Minimum Pins

// Select Minimum Reactions
func selectMinimumReactions(c *discordgo.Channel) (minimum int, err error) {
	err = db.QueryRow("select `minimum` from `"+tableMinimumReactions+"` where channel = ?;", c.ID).Scan(&minimum)
	return
}

// Insert Minimum Reactions
func insertMinimumReactions(c *discordgo.Channel, minimum int) (res sql.Result, err error) {
	return db.Exec("insert into `"+tableMinimumReactions+"`(`channel`, `minimum`) values(?, ?);", c.ID, minimum)
}

// Update Minimum Reactions
func updateMinimumReactions(c *discordgo.Channel, minimum int) (res sql.Result, err error) {
	return db.Exec("update `"+tableMinimumReactions+"` set `minimum` = ? where `minimum` = ?;", c.ID, minimum)
}

// Self-Assignable Role

// Select Self-Assignable Role
func selectSAR(g *discordgo.Guild, r *discordgo.Role) (id string, err error) {
	err = db.QueryRow("select `role` from `"+tableSAR+"` where `server` = ? and `role` = ?;", g.ID, r.ID).Scan(&id)
	return
}

// Select Self-Assignable Roles
func selectSARs(g *discordgo.Guild) (roles []string, err error) {

	// Query
	rows, err := db.Query("select `role` from `"+tableSAR+"` where `server` = ?;", g.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	// Select
	for rows.Next() {

		// Get the role
		var role string
		err = rows.Scan(&role)
		if err != nil {
			return
		}

		// Place it in a list
		roles = append(roles, role)
	}

	// Check for errors
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Insert SAF
func insertSAR(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return db.Exec("insert into `"+tableSAR+"`(`server`, `role`) values(?, ?);", g.ID, r.ID)
}

// Delete SAF
func deleteSAR(g *discordgo.Guild, r *discordgo.Role) (res sql.Result, err error) {
	return db.Exec("delete from `"+tableSAR+"` where `server` = ? and `role` = ?;", g.ID, r.ID)
}

// Status

// Insert Status
func insertStatus(status string) (res sql.Result, err error) {
	return db.Exec("insert into `"+tableStatus+"`(`status`) values(?);", status)
}

// Delete Status
func deleteStatus(id int) (res sql.Result, err error) {
	return db.Exec("delete from `"+tableStatus+"` where `id` = ?;", id)
}

// Delete All Status
func emptyStatus() (res sql.Result, err error) {
	return db.Exec("delete from `" + tableStatus + "`;")
}

// Select Latest Status
func selectStatus() (status string, err error) {
	err = db.QueryRow("select `status` from `" + tableStatus + "` order by `id` desc limit 1;").Scan(&status)
	return
}

// Hall of Fame

// Select hall of hame channel
func selectFameChannel(g *discordgo.Guild) (id string, err error) {
	err = db.QueryRow("select `channel` from `"+tableFame+"` where `server` = ?;", g.ID).Scan(&id)
	return
}

// Insert hall of hame channel
func insertFameChannel(g *discordgo.Guild, c *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("insert into `"+tableFame+"`(`server`, `channel`) values(?, ?);", g.ID, c.ID)
}

// Update hall of hame channel
func updateFameChannel(g *discordgo.Guild, c *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("update `"+tableFame+"` set `channel` = ? where `server` = ?;", c.ID, g.ID)
}

// Messages in hall of fame

// Insert in hall of Fame
func insertMessagesFamed(g *discordgo.Guild, m *discordgo.Message) (res sql.Result, err error) {
	return db.Exec("insert into `"+tableMessagesFamed+"`(`message`) values(?);", m.ID)
}

// Select from hall of fame
func selectMessagesFamed(m *discordgo.Message) (id string, err error) {
	err = db.QueryRow("select `message` from `"+tableMessagesFamed+"` where `message` = ?;", m.ID).Scan(&id)
	return
}

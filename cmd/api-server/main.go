package main

import (
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"sync"
	"time"

	"bitbought.net/project-root-directory/cmd/database/install"
	"bitbought.net/project-root-directory/internal/hashcode"
	"github.com/couchbase/gocb/v2"
	"github.com/gorilla/sessions"
)

/***************
  server program
****************/

//  1 connect to database
//  2 start store webservice
//    start mobile webservice
//    start admin webservice
//    start system webservice
//     .monitor
//     .garbage
//     .memory
//     .latency
//     .security
//     .restart
//     .shutdown
//   3 wait for shutdown signal

/* todo:
implement an application for the offline payment processing
*/

// database connection context
const connectionString string = "localhost"
const bucketName string = "market"

// session store
// []byte(os.Getenv("SESSION_KEY"))
var (
	key = []byte("super-secret-key")
	cookieStore = sessions.NewCookieStore(key)
)

func main() {

	// connect to database
	cluster, err := databaseConnect()
	if err != nil {
		panic("database is not ready. server program exits.")
    } else {fmt.Printf("database %v connected\n", cluster)}
	
	// shopper webservice mutex
	serveMuxAt8080 := http.NewServeMux()
	serveMuxAt8080.HandleFunc("/", indexPage)

	serverAt8080 := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      serveMuxAt8080,
	}
    // store webservices mutex
	serveMuxAt8081 := http.NewServeMux()
	serveMuxAt8081.HandleFunc("/login",storeLoginPage)
	serveMuxAt8081.HandleFunc("/auth", storeAuthPage)
	serveMuxAt8081.HandleFunc("/setting", storeSettingPage)
	serveMuxAt8081.HandleFunc("/offer",offerPage)
	serveMuxAt8081.HandleFunc("/report",reportPage)
	serveMuxAt8081.HandleFunc("/support", supportPage)
	serveMuxAt8081.HandleFunc("/logout", logOutPage)
	

	serverAt8081 := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      serveMuxAt8081,
	}
    // admin webservices mutex
	serveMuxAt8082 := http.NewServeMux()
	serveMuxAt8082.HandleFunc("/", indexPage)
	serveMuxAt8082.HandleFunc("/login", adminLoginPage)
	serveMuxAt8082.HandleFunc("/auth", adminAuthPage)
	//serveMuxAt8082.HandleFunc("/logout", adminLogoutPage)
	


	serverAt8082 := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      serveMuxAt8082,
	}

	serveMuxAt8083 := http.NewServeMux()
	serveMuxAt8083.HandleFunc("/", indexPage)

	serverAt8083 := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      serveMuxAt8083,
	}

	
	// use waitgroup to wait goroutines finish
	wg := sync.WaitGroup{}
	wg.Add(4)

	// start shopper webservice (port 8080, authentificate,search,map service, payment,
	// transaction, accounts service, activity monitor)

	listenerAt8080, _ := net.Listen("tcp", ":8080")
	go func(server *http.Server, ln net.Listener) {
		defer wg.Done()
		defer fmt.Println("mobile service closed...at :8080")
		fmt.Println("mobile service start...at :8080")
		server.ServeTLS(ln, "go-server.crt", "go-server.key")
	}(serverAt8080, listenerAt8080)

	//start store webservice (port 8081,authentificate,listing,support)

	listenerAt8081, _ := net.Listen("tcp", ":8081")
	go func(server *http.Server, ln net.Listener) {
		defer wg.Done()
		defer fmt.Println("store service closed...at :8081")
		fmt.Println("store service start...at :8081")
		server.ServeTLS(ln, "go-server.crt", "go-server.key")
	}(serverAt8081, listenerAt8081)

	//start admin webservice (port 8082, account,support)

	listenerAt8082, _ := net.Listen("tcp", ":8082")
	go func(server *http.Server, ln net.Listener) {
		defer wg.Done()
		defer fmt.Println("admin service closed...at :8082")
		fmt.Println("admin service start...at :8082")
		server.ServeTLS(ln, "go-server.crt", "go-server.key")
	}(serverAt8082, listenerAt8082)

	// start system tasks (at port 8083, security, 
	// performance index, garbage, latency, restart, shutdown)

	listenerAt8083, _ := net.Listen("tcp", ":8083")
	go func(server *http.Server, ln net.Listener) {
		defer wg.Done()
		defer fmt.Println("system service closed...at :8083")
		fmt.Println("system service start...at :8083")
		server.ServeTLS(ln, "go-server.crt", "go-server.key")
	}(serverAt8083, listenerAt8083)

	/* close webservice calls
	listenerAt8080.Close() serverAt8080.Close()
	listenerAt8081.Close() serverAt8081.Close()
	listenerAt8082.Close() serverAt8082.Close()
	listenerAt8083.Close() serverAt8083.Close()
	*/
	wg.Wait()
	fmt.Println("All webservices are closed..")
}

func databaseConnect() (*gocb.Cluster, error){

	// For a secure cluster connection, use `couchbases://<your-cluster-ip>` instead.
	cluster, err := gocb.Connect("couchbase://"+connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "class_a_rw",
			Password: "123456",
		},
	})
	if err != nil {
		return nil, err
	}

	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(1*time.Second, nil)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "--index page--.")
}

func storeLoginPage(w http.ResponseWriter, r *http.Request) {
	    tmpl, err := template.ParseFiles("../../tmpl/store/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
}

func offerPage(w http.ResponseWriter, r *http.Request) {
	// session analysis
	session, err := cookieStore.Get(r, "store_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !session.IsNew {
		// session value leads to offer page template
		tmpl, err := template.ParseFiles("../../tmpl/store/offer.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		store, err := retrieveStoreSetting(session.Values["id"].(string))

		if err != nil {
			panic("store setting is not accessible")
		}
		err = tmpl.Execute(w, store)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// redirect to login page
		storeLoginPage(w,r)
	}
}

func storeSettingPage(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, "store_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if !session.IsNew {
		store, err := retrieveStoreSetting(session.Values["id"].(string))
		if err != nil {
			panic("store_setting is not accessible")
		}
		
		tmpl, err := template.ParseFiles("../../tmpl/store/setting.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = tmpl.Execute(w, store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// redirect to login page
		storeLoginPage(w,r)
	}
}

func reportPage(w http.ResponseWriter, r *http.Request) {
	// session analysis
	session, err := cookieStore.Get(r, "store_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !session.IsNew {
		// apply session value to report page template
		tmpl, err := template.ParseFiles("../../tmpl/store/report.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.Execute(w, session.Values["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// redirect to login page
		storeLoginPage(w,r)
	}
}

func supportPage(w http.ResponseWriter, r *http.Request) {
	//todo
}

func logOutPage(w http.ResponseWriter, r *http.Request) {
	// session analysis
	session, err := cookieStore.Get(r, "store_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options = &sessions.Options{
		MaxAge:   -1,
	}

	err = sessions.Save(r, w)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}

	tmpl, err := template.ParseFiles("../../tmpl/store/store_landing_page.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
}


func storeAuthPage(w http.ResponseWriter, r *http.Request) {
	
	r.ParseForm()
	
	store_id, err := authenticate_store(r.PostFormValue("user"), r.PostFormValue("password"))
	
	if err != nil{
		// authentication fails
		storeLoginPage(w,r)
    } else {
		// authentication succeeds
		
		session, err := cookieStore.Get(r, "store_session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["id"] = store_id
		err = sessions.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// after authentication, it shows store landing page
		tmpl, err := template.ParseFiles("../../tmpl/store/store_landing_page.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.Execute(w, session.Values["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
	}

}

func retrieveStoreSetting(store_id string) (database.Store, error) {
	cluster, err := databaseConnect()

	if err != nil {
		return database.Store{}, err
	}	
	
	col := cluster.Bucket("market").Scope("persona").Collection("store")
	result, err := col.Get(store_id, nil)
	// check query is successful or not
	if err != nil {
		return database.Store{}, err
	}

	var s database.Store
	err = result.Content(&s)

	if err != nil {
		return database.Store{}, err
	}

	return s, nil
}

// compare (user, password) with database authenticate table, success returns store_id
// failure returns error
func authenticate_store(user string, password string) (string, error) {
	cluster, err := databaseConnect()
	if err != nil {
		return "", err
	}	
	query := "SELECT x.store_id, x.pass_word FROM market.activity.authenticate_store x WHERE x.`user`=$user;"
	params := make(map[string]interface{},1)
	params["user"] = user
	rows, err := cluster.Query(query, &gocb.QueryOptions{NamedParameters: params, Adhoc: true})
    
	// check query is successful or not
	if err != nil {
		return "", err
	}

	type store_info struct {
		Store_id string `json:"store_id"`
		Password string `json:"pass_word"`
	} 
	var s store_info

	  // expect one or no queryresult
	err = rows.One(&s)

      if err != nil {
		return "", err
	  }

	  if passwordMatch := hashcode.VerifyPasswordHash(password, s.Password); passwordMatch {
	    fmt.Printf("store - %v logged in.\n", s)
	    return s.Store_id, nil
	  } else {
		return "", errors.New("password error")
	  }
}

// admin webservices
func adminLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../tmpl/admin/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminAuthPage(w http.ResponseWriter, r *http.Request) {
	
	r.ParseForm()
	
	admin_id, err := authenticate_admin(r.PostFormValue("user"), r.PostFormValue("password"))
	
	if err != nil{
		// authentication fails
		adminLoginPage(w,r)
    } else {
		// authentication succeeds
		session, err := cookieStore.Get(r, "admin_session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["id"] = admin_id
		err = sessions.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// after authentication, it shows admin or system admin landing page
		tmpl, err := template.ParseFiles("../../tmpl/admin/admin_landing_page.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.Execute(w, session.Values["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
	}

}

// compare (user, password) with database authenticate table, success returns admin_id
// failure returns error
func authenticate_admin(user string, password string) (string, error) {
	cluster, err := databaseConnect()
	if err != nil {
		return "", err
	}	
	query := "SELECT x.admin_id, x.pass_word FROM market.activity.authenticate_admin x WHERE x.`user`=$user;"
	params := make(map[string]interface{},1)
	params["user"] = user
	rows, err := cluster.Query(query, &gocb.QueryOptions{NamedParameters: params, Adhoc: true})
    
	// check query is successful or not
	if err != nil {
		return "", err
	}

	type admin_info struct {
		Admin_id string `json:"admin_id"`
		Password string `json:"pass_word"`
	} 
	var s admin_info

	  // expect one or no queryresult
	err = rows.One(&s)

      if err != nil {
		return "", err
	  }

	  if passwordMatch := hashcode.VerifyPasswordHash(password, s.Password); passwordMatch {
	    fmt.Printf("admin - %v logged in.\n", s.Admin_id)
		// todo: retrieve admin type and return id and type
	    return s.Admin_id, nil
	  } else {
		return "", errors.New("password error")
	  }
}


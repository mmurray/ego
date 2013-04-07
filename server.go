package ego

import (
	"github.com/murz/ego/http"
	// "github.com/murz/ego/tmpl"
	// "github.com/murz/ego/cfg"
	// "github.com/murz/ego/db"
	// "github.com/murz/ego/plugins"
	// "github.com/murz/ego/cache"
    // "go/build"
    "flag"
    // "regexp"
    "log"
    "fmt"
    "os"
    "strconv"
    "time"
    netHTTP "net/http"
)

type Server struct {
	// Package *build.Package
	PackageName string
	HTTPRouter *http.Router
	// Config *cfg.ConfigMap
}

func NewServer(pkgName string) *Server {
	// pkg, err := build.Default.Import(pkgName, "", build.FindOnly)
    // if err != nil {
    //     pkg = &build.Package{
    //     	Dir: "/app",
    //     }
    // }
	return &Server{
		// Package: pkg,
		PackageName: pkgName,
		HTTPRouter: http.GetDefaultRouter(),
	}
}

func (s *Server) Run() {
	startTime := time.Now().UnixNano() / 1000000
	// determine default port
	envport := os.Getenv("PORT")
	p, err := strconv.ParseInt(envport, 0, 0)
	if envport == "" || err != nil {
		p = 5000
	}

	// parse flags
	var isDev = flag.Bool("dev", false, "Start server in development mode.")
	var port = flag.Int("port", int(p), "HTTP server port.")
	flag.Parse()

	// read the routes file

	

	// allow plugins to do some intialization
	// for _, plugin := range plugins.All() {
	// 	if plugin.OnStart != nil {
	// 		plugin.OnStart()
	// 	}
	// }

	// parse the config files
	// cfg.ParseDir(s.Package.Dir + "/conf")

	// db.Connect(cfg.Get("db"))
	
	// cache.Init()

	// serve static assets from /public/
	netHTTP.Handle("/public/", netHTTP.StripPrefix("/public/", netHTTP.FileServer(netHTTP.Dir(s.PackageName+"/public/"))))

	// redirect favicon requests to /public/
	netHTTP.Handle("/favicon.ico", netHTTP.RedirectHandler("/public/favicon.ico", 301))

	// if actions.Count() == 0 {
	// 	// show the default page if there are no registered actions
	// 	netHTTP.HandleFunc("/", defaultHandler)
	// } else {

	// pipe all requests through the action dispatcher
	netHTTP.HandleFunc("/", http.ActionDispatchHandler(s.HTTPRouter))
	// }

	// parse mustache templates
	// tmpl.SetPackageName(s.Package.Dir)
 //    tmpl.ParseDir("/app/views")

    // call into plugins again now that everything is ready
	// for _, plugin := range plugins.All() {
	// 	if (plugin.OnReady != nil) {
	// 		plugin.OnReady()
	// 	}
	// }

	// listen and serve
	log.Print("____________ ______ ")
	log.Print("_  _ \\_  __ `/  __ \\")
	log.Print("/  __/  /_/ // /_/ /")
	log.Print("\\___/_\\__, / \\____/ ")
	log.Print("     /____/         ")
	log.Printf("ego server running on %v", *port)
	if (*isDev) {
		log.Printf(":: development mode ::")
	}
	log.Printf("## startup time: %dms ##", time.Now().UnixNano() / 1000000 - startTime)
	err = netHTTP.ListenAndServe(fmt.Sprintf(":%v", *port), nil)

	// give plugins a chance to cleanup
	// for _, plugin := range plugins.All() {
	// 	if plugin.OnStop != nil {
	// 		plugin.OnStop()
	// 	}
	// }

	log.Fatal(err)
}

func defaultHandler(w netHTTP.ResponseWriter, httpReq *netHTTP.Request) {
	fmt.Fprint(w, "ego rulz")
}


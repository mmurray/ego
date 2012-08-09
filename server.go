package ego

import (
	"github.com/murz/ego/http"
	"github.com/murz/ego/ws"
	"github.com/murz/ego/tmpl"
	"github.com/murz/ego/cfg"
	"github.com/murz/ego/db"
	"github.com/murz/ego/actions"
    "go/build"
    "github.com/murz/go-socket.io"
    // "regexp"
    "log"
    "fmt"
    netHTTP "net/http"
)

type Server struct {
	HTTPRouter *http.Router
	WSRouter *ws.Router
	Package *build.Package
	Config *cfg.ConfigMap
}

func NewServer(pkgName string) *Server {
	pkg, err := build.Default.Import(pkgName, "", build.FindOnly)
    if err != nil {
        pkg = &build.Package{
        	Dir: "/app",
        }
    }
	return &Server {
		Package: pkg,
		HTTPRouter: http.NewRouter(),
		WSRouter: ws.NewRouter(),
	}
}

func (s *Server) RegisterHTTPAction(action *http.Action) {
	s.HTTPRouter.Register(action)
}

func (s *Server) RegisterHTTPActions(actions []*http.Action) {
	for _, action := range actions {
		s.RegisterHTTPAction(action)
	}
}

func (s *Server) RegisterWSAction(action *ws.Action) {
	s.WSRouter.Register(action)
}

func (s *Server) RegisterWSActions(actions []*ws.Action) {
	for _, action := range actions {
		s.RegisterWSAction(action)
	}
}

func (s *Server) Run(p string) {
	// parse the config files
	cfg.ParseDir(s.Package.Dir + "/conf")

	db.Connect(cfg.Get("db"))

	// register all actions that were buffered in the action manager
	s.RegisterHTTPActions(actions.HTTPActions())
	s.RegisterWSActions(actions.WSActions()) 

	// serve static assets from /public/
	netHTTP.Handle("/public/", netHTTP.StripPrefix("/public/", netHTTP.FileServer(netHTTP.Dir(s.Package.Dir+"/public/"))))

	// redirect favicon requests to /public/
	netHTTP.Handle("/favicon.ico", netHTTP.RedirectHandler("/public/favicon.ico", 301))

	if actions.Count() == 0 {
		// show the default page if there are no registered actions
		netHTTP.HandleFunc("/", defaultHandler)
	} else {
		if len(actions.WSActions()) > 0 {
			// register the socket.io handler if there are any ws actions.
			sio := socketio.NewServer(nil)
			netHTTP.Handle("/socket.io/", netHTTP.StripPrefix("/socket.io/", sio.Handler(s.WSRouter.ActionDispatchHandler())))
		}
		// pipe all requests through the action dispatcher
		netHTTP.HandleFunc("/", s.HTTPRouter.ActionDispatchHandler())
	}

	// parse mustache templates
	tmpl.SetPackageName(s.Package.Dir)
    tmpl.ParseDir("/app/views")

	// listen and serve
	log.Print("____________ ______ ")
	log.Print("_  _ \\_  __ `/  __ \\")
	log.Print("/  __/  /_/ // /_/ /")
	log.Print("\\___/_\\__, / \\____/ ")
	log.Print("     /____/         ")
	log.Printf("ego server running on %v", p)
	netHTTP.ListenAndServe(p, nil)
}

func defaultHandler(w netHTTP.ResponseWriter, httpReq *netHTTP.Request) {
	fmt.Fprint(w, "ego rulz")
}


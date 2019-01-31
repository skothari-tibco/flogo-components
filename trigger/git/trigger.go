package git

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"github.com/skothari-tibco/flogo-components/trigger/git/cors"
)

const (
	CorsPrefix = "GIT_TRIGGER"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	err = checkGitPath(s.GitPath)
	if err != nil {
		return nil, err
	}
	err = initHooks(s.GitPath, s.Port)

	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, settings: s}, nil
}

func checkGitPath(path string) error {

	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		return err
	}
	return nil
}

func initHooks(gitpath string, port int) error {

	var shimFile = `package main
import (
	"fmt"
	"net/http"
	"time"
)


func main(){
	
	go func(){
		resp, err := http.Get("http://localhost:{{.PORT}}/test")
		defer resp.Body.Close()
		if err != nil{
			fmt.Println(err)
		}
	}()
	time.Sleep(2 * time.Second)
	fmt.Println("Successfully Committed")
}
`
	data := struct {
		PORT int
	}{
		port,
	}
	f, err := os.Create("post-commit.go")
	if err != nil {
		return err
	}
	RenderTemplate(f, shimFile, &data)
	f.Close()

	err = exec.Command("go", "build", "post-commit.go").Run()

	err = exec.Command("mv", "post-commit", filepath.Join(gitpath, ".git", "hooks")).Run()

	err = exec.Command("rm", "post-commit.go").Run()

	if err != nil {
		return err
	}
	return nil
}
func RenderTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

// Trigger REST trigger struct
type Trigger struct {
	server   *Server
	settings *Settings
	id       string
	logger   log.Logger
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.logger = ctx.Logger()

	router := httprouter.New()

	addr := ":" + strconv.Itoa(t.settings.Port)

	pathMap := make(map[string]string)

	preflightHandler := &PreflightHandler{logger: t.logger, c: cors.New(CorsPrefix, t.logger)}

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		method := s.Method
		path := s.Path

		t.logger.Debugf("Registering handler [%s: %s]", method, path)

		if _, ok := pathMap[path]; !ok {
			pathMap[path] = path
			router.OPTIONS(path, preflightHandler.handleCorsPreflight) // for CORS
		}

		//router.OPTIONS(path, handleCorsPreflight) // for CORS
		router.Handle(method, path, newActionHandler(t, handler))
	}

	t.logger.Debugf("Configured on port %d", t.settings.Port)
	t.server = NewServer(addr, router)

	return nil
}

func (t *Trigger) Start() error {
	return t.server.Start()
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	return t.server.Stop()
}

type PreflightHandler struct {
	logger log.Logger
	c      cors.Cors
}

// Handles the cors preflight request
func (h *PreflightHandler) handleCorsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	h.logger.Infof("Received [OPTIONS] request to CorsPreFlight: %+v", r)
	h.c.HandlePreflight(w, r)
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func newActionHandler(rt *Trigger, handler trigger.Handler) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		rt.logger.Info("Request receieved....")
		rt.logger.Debugf("Received request for id '%s'", rt.id)

		c := cors.New(CorsPrefix, rt.logger)
		c.WriteCorsActualRequestHeaders(w)

		out := &Output{}

		out.PathParams = make(map[string]string)
		for _, param := range ps {
			out.PathParams[param.Key] = param.Value
		}

		queryValues := r.URL.Query()
		out.QueryParams = make(map[string]string, len(queryValues))
		out.Headers = make(map[string]string, len(r.Header))

		for key, value := range r.Header {
			out.Headers[key] = strings.Join(value, ",")
		}

		for key, value := range queryValues {
			out.QueryParams[key] = strings.Join(value, ",")
		}

		// Check the HTTP Header Content-Type
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/x-www-form-urlencoded":
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			s := buf.String()
			m, err := url.ParseQuery(s)
			content := make(map[string]interface{}, 0)
			if err != nil {
				rt.logger.Errorf("Error while parsing query string: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			for key, val := range m {
				if len(val) == 1 {
					content[key] = val[0]
				} else {
					content[key] = val[0]
				}
			}

			out.Content = content
		case "application/json":
			var content interface{}
			err := json.NewDecoder(r.Body).Decode(&content)
			if err != nil {
				switch {
				case err == io.EOF:
					// empty body
					//todo should handler say if content is expected?
				case err != nil:
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			out.Content = content
		default:
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			out.Content = string(b)
		}

		results, err := handler.Handle(context.Background(), out)

		reply := &Reply{}
		reply.FromMap(results)

		if err != nil {
			rt.logger.Debugf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if reply.Data != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if reply.Code == 0 {
				reply.Code = 200
			}
			w.WriteHeader(reply.Code)
			if err := json.NewEncoder(w).Encode(reply.Data); err != nil {
				log.Error(err)
			}
			return
		}

		if reply.Code > 0 {
			w.WriteHeader(reply.Code)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

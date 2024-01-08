package swaggerapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/ml444/swaggerapi/generator"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Server is api meta server
type Server struct {
	lock     sync.Mutex
	services map[string]*descriptorpb.FileDescriptorSet
	methods  map[string][]string

	opts []generator.Option
}

// NewServer create server instance
func NewServer(opts ...generator.Option) *Server {
	return &Server{
		services: make(map[string]*descriptorpb.FileDescriptorSet),
		methods:  make(map[string][]string),
		opts:     opts,
	}
}

type ServicesReply struct {
	Services []string `json:"services,omitempty"`
	Methods  []string `json:"methods,omitempty"`
}

func (s *Server) load() error {
	if len(s.services) > 0 {
		return nil
	}
	var err error
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if fd.Services() == nil {
			return true
		}
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			fdp, e := fileDescriptorProto(fd.Path())
			if e != nil {
				err = e
				return false
			}
			fdps, e := allDependency(fdp)
			if e != nil {
				if errors.Is(e, protoregistry.NotFound) {
					// Skip this service if one of its dependencies is not found.
					continue
				}
				err = e
				return false
			}
			s.services[string(svc.FullName())] = &descriptorpb.FileDescriptorSet{File: fdps}
			if svc.Methods() == nil {
				continue
			}
			for j := 0; j < svc.Methods().Len(); j++ {
				method := svc.Methods().Get(j)
				s.methods[string(svc.FullName())] = append(s.methods[string(svc.FullName())], string(method.Name()))
			}
		}
		return true
	})
	return err
}

// ListServices return all services
func (s *Server) ListServices(w http.ResponseWriter, _ *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.load(); err != nil {
		response(w, http.StatusInternalServerError, err)
		return
	}
	reply := &ServicesReply{
		Services: make([]string, 0, len(s.services)),
		Methods:  make([]string, 0, len(s.methods)),
	}
	for name := range s.services {
		reply.Services = append(reply.Services, name)
	}
	for name, methods := range s.methods {
		for _, method := range methods {
			reply.Methods = append(reply.Methods, fmt.Sprintf("/%s/%s", name, method))
		}
	}
	sort.Strings(reply.Services)
	sort.Strings(reply.Methods)
	buf, err := json.Marshal(reply)
	if err != nil {
		response(w, http.StatusInternalServerError, err)
		return
	}
	response(w, http.StatusOK, buf)
	return
}

// GetServiceDesc return service meta by name
func (s *Server) GetServiceDesc(w http.ResponseWriter, r *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.load(); err != nil {
		if _, e := w.Write([]byte(err.Error())); e != nil {
			println(e.Error())
		}
		return
	}
	name := ""
	raws := strings.Split(r.URL.Path, "/")
	if rawsLen := len(raws); rawsLen > 0 {
		name = raws[rawsLen-1]
	}
	fds, ok := s.services[name]
	if !ok {
		response(w, http.StatusNotFound, fmt.Errorf("service %s not found", name))
		return
	}
	files := fds.File
	var target string
	if len(files) == 0 {
		response(w, http.StatusNotFound, fmt.Errorf("proto file is empty"))
		return
	}
	if files[len(files)-1].Name == nil {
		response(w, http.StatusNotFound, fmt.Errorf("proto file's name is null"))
		return
	}
	target = *files[len(files)-1].Name
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	req := new(pluginpb.CodeGeneratorRequest)
	req.FileToGenerate = []string{target}
	var para = ""
	req.Parameter = &para
	req.ProtoFile = files
	g := generator.NewGenerator(s.opts...)
	resp, err := g.Gen(req)
	if err != nil {
		response(w, http.StatusInternalServerError, err)
		return
	}
	if len(resp.File) == 0 || resp.File[0].Content == nil {
		response(w, http.StatusOK, "{}")
		return
	}

	response(w, http.StatusOK, *resp.File[0].Content)
	return
}

func response(w http.ResponseWriter, status int, v interface{}) {
	if status < 0 {
		status = http.StatusInternalServerError
	} else if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	switch v.(type) {
	case error:
		if _, e := w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, v))); e != nil {
			println(e.Error())
		}
	case string:
		if _, e := w.Write([]byte(v.(string))); e != nil {
			println(e.Error())
		}
	case []byte:
		if _, e := w.Write(v.([]byte)); e != nil {
			println(e.Error())
		}
	}
}
func allDependency(fd *descriptorpb.FileDescriptorProto) ([]*descriptorpb.FileDescriptorProto, error) {
	var files []*descriptorpb.FileDescriptorProto
	for _, dep := range fd.Dependency {
		fdDep, err := fileDescriptorProto(dep)
		if err != nil {
			continue
		}
		temp, err := allDependency(fdDep)
		if err != nil {
			return nil, err
		}
		files = append(files, temp...)
	}
	files = append(files, fd)
	return files, nil
}

func fileDescriptorProto(path string) (*descriptorpb.FileDescriptorProto, error) {
	fd, err := protoregistry.GlobalFiles.FindFileByPath(path)
	if err != nil {
		return nil, fmt.Errorf("find proto by path failed, path: %s, err: %s", path, err)
	}
	fdpb := protodesc.ToFileDescriptorProto(fd)
	return fdpb, nil
}

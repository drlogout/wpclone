package store

import (
	"croox/wpclone/config/global"
	"croox/wpclone/pkg/util"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	log "github.com/sirupsen/logrus"
)

//go:embed *.tmpl
var tpls embed.FS

var mutex sync.Mutex

type Store struct {
	Remotes Remotes `yaml:"remotes"`
}

type Remotes map[string]Remote

type Remote struct {
	DB DB `yaml:"db"`
}

type DB struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

func Load() (*Store, error) {
	if err := ensureStore(); err != nil {
		return nil, err
	}

	var s Store
	if err := util.LoadYAML(storeFilePath(), &s); err != nil {
		return nil, err
	}

	if s.Remotes == nil {
		s.Remotes = make(Remotes)
	}

	return &s, nil
}

func (s *Store) GetRemote(name string) Remote {
	remote, ok := s.Remotes[name]
	if !ok {
		remote = Remote{}
	}

	return remote
}

func (s *Store) GetRemoteDB(name string) DB {
	remote := s.GetRemote(name)
	return remote.DB
}

func (s *Store) SetRemote(name string, remote Remote) error {
	s.Remotes[name] = remote

	return saveStore(s)
}

func (s *Store) SetRemoteDB(remoteName string, db DB) error {
	if remoteName == "" {
		return fmt.Errorf("remote name is required")
	}

	remote := s.GetRemote(remoteName)

	remote.DB = db

	return s.SetRemote(remoteName, remote)
}

func NewDB(values map[string]string) DB {
	port, _ := strconv.Atoi(values["DB_PORT"])
	if port == 0 {
		port = 3306
	}

	return DB{
		Host:     values["DB_HOST"],
		Name:     values["DB_NAME"],
		User:     values["DB_USER"],
		Password: values["DB_PASSWORD"],
		Port:     port,
	}
}

func ParseDB(db DB) (DB, error) {
	if strings.Contains(db.Host, ":") {
		split := strings.Split(db.Host, ":")
		db.Host = split[0]
		p := split[1]

		port, err := strconv.Atoi(p)
		if err != nil {
			return DB{}, err
		}

		db.Port = port

	}

	if db.Port == 0 {
		db.Port = 3306
	}

	return db, nil
}

func ensureStore() error {
	if util.FileExists(storeFilePath()) {
		return nil
	}

	log.Debug("Initializing store")
	return saveStore(&Store{})
}

func saveStore(s *Store) error {
	mutex.Lock()
	defer mutex.Unlock()

	t, err := template.ParseFS(tpls, "*")
	if err != nil {
		return err
	}

	file, err := os.Create(storeFilePath())
	if err != nil {
		return err
	}
	defer file.Close()

	return t.ExecuteTemplate(file, "store.yml.tmpl", s)
}

func storeFilePath() string {
	return filepath.Join(global.ConfigDir(), "store.yml")
}

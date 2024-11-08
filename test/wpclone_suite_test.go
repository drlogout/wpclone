package main_test

import (
	"bytes"
	"croox/wpclone/config"
	"croox/wpclone/pkg/exec"
	"croox/wpclone/pkg/util"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"text/template"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

//go:embed *.tmpl
var tpls embed.FS

func TestWpclone(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wpclone Suite")
}

type Paths struct {
	LocalFolder     string
	LocalWPFolder   string
	LocalConfigFile string
}

var paths Paths

func setupLocalFiles(opts ...string) (Paths, error) {
	var testName string
	if len(opts) > 0 {
		testName = opts[0]
	}

	// tmpFolder, err := os.MkdirTemp("", "wpclone")
	// if err != nil {
	// 	return "", err
	// }
	tmpFolder := filepath.Join(util.Getwd(), "test_wpclone")
	localFolder := filepath.Join(tmpFolder, "local")
	configFile := filepath.Join(tmpFolder, "wpclone.yml")
	sshKeyFile := filepath.Join(util.Getwd(), "remote_docker/ssh/tux@croox.com")

	// ensire local folder and config file are removed
	os.RemoveAll(localFolder)
	os.RemoveAll(configFile)

	os.MkdirAll(localFolder, 0755)

	cfg := config.NewConfig()
	cfg.Local.Path = localFolder
	cfg.Local.URL = "http://wpclone.test"
	cfg.Remote.SSH.Key = sshKeyFile

	if testName != "" {
		if err := saveInit(configFile, fmt.Sprintf("%s.yml.tmpl", testName), cfg); err != nil {
			return paths, err
		}
	}

	paths = Paths{
		LocalFolder:     tmpFolder,
		LocalWPFolder:   localFolder,
		LocalConfigFile: configFile,
	}

	return paths, nil
}

func startRemoteDocker() error {
	opts := exec.RunOpts{
		Verbose: false,
		Dir:     "./remote_docker",
	}

	if err := exec.RunWithOpts(opts, "docker", "compose", "up", "-d"); err != nil {
		return err
	}

	return waitForRemote()

}

func wpcloneCLI(arg ...string) (string, error) {
	var output bytes.Buffer
	var writer io.Writer = &output

	opts := exec.RunOpts{
		Verbose: false,
		Stdout:  writer,
		Dir:     paths.LocalFolder,
	}

	arg = append([]string{"--quiet"}, arg...)

	if err := exec.RunWithOpts(opts, "wpclone", arg...); err != nil {
		return "", err
	}

	return output.String(), nil
}

func parseConfig(path string) (*config.Config, error) {
	cfg := config.NewConfig()

	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(body, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func saveInit(filePath string, tpl string, cfg *config.Config) error {
	t, err := template.ParseFS(tpls, "*")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.ExecuteTemplate(file, tpl, cfg)
}

func waitForRemote() error {
	timeout := time.Now().Add(60 * time.Second)

	for {
		if time.Now().After(timeout) {
			return fmt.Errorf("timeout")
		}

		resp, err := http.Get("http://localhost:8080")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func remoteHasContent(content string) (bool, error) {
	return hasContent("http://localhost:8080", content)
}

func localHasContent(content string) (bool, error) {
	return hasContent("http://wpclone.test", content)
}

func hasContent(url, content string) (bool, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if bytes.Contains(body, []byte(content)) {
		return true, nil
	}

	return false, nil
}

func modifyLocalWP() error {
	t, err := template.ParseFS(tpls, "footer.php.tmpl")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(paths.LocalWPFolder, "/wp-content/themes/blankslate/footer.php"))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, struct{}{})
}

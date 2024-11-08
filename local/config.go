package local

import (
	"bufio"
	"croox/wpclone/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func Configure(cfg *config.Config) error {
	dbSettings := map[string]string{
		"DB_HOST":     fmt.Sprintf("%s:%d", cfg.LocalDBHost(), cfg.LocalDBPort()),
		"DB_NAME":     cfg.LocalDBName(),
		"DB_USER":     cfg.LocalDBUser(),
		"DB_PASSWORD": cfg.LocalDBPassword(),
	}

	for key, value := range dbSettings {
		if err := WPCli(cfg, "config", "set", key, value); err != nil {
			return err
		}
	}

	for key, value := range cfg.LocalConfigSettings() {
		if err := WPCli(cfg, "config", "set", key, value); err != nil {
			return err
		}
	}

	// raw means: "Place the value into the wp-config.php file as is, instead of as a quoted string."
	for key, value := range cfg.LocalRAWConfigSettings() {
		if err := WPCli(cfg, "config", "set", key, value, "--raw"); err != nil {
			return err
		}
	}

	if HasUserINI(cfg.LocalPath()) {
		if err := UpdateUserINI(cfg.LocalPath(), cfg.LocalPath(), cfg.LocalPath()); err != nil {
			return err
		}
	}

	return nil
}

func UpdateUserINI(iniPath, srcDir, destDir string) error {
	userINI := filepath.Join(srcDir, ".user.ini")
	userINIDest := filepath.Join(destDir, ".user.ini")

	userINIBytes, err := os.ReadFile(userINI)
	if err != nil {
		return err
	}

	userINIString := string(userINIBytes)

	autoPrependEntry := fmt.Sprintf("auto_prepend_file = '%s'", filepath.Join(iniPath, "wordfence-waf.php"))

	re := regexp.MustCompile(`auto_prepend_file.*`)
	s := re.ReplaceAllString(userINIString, autoPrependEntry)

	return os.WriteFile(userINIDest, []byte(s), 0644)
}

func HasUserINI(path string) bool {
	userINI := filepath.Join(path, ".user.ini")
	_, err := os.Stat(userINI)
	return !errors.Is(err, os.ErrNotExist)
}

func WPConfigValue(constantName string, phpFile string) (string, error) {
	v, err := ReadWPConfigConstants([]string{constantName}, phpFile)
	if err != nil {
		return "", err
	}

	return v[constantName], nil
}

func ReadWPConfigConstants(constantNames []string, phpFile string) (map[string]string, error) {
	values := make(map[string]string)
	for _, name := range constantNames {
		value, err := scanWPConfig(name, phpFile)
		if err != nil {
			return nil, err
		}
		values[name] = value
	}

	return values, nil
}

func scanWPConfig(constantName, phpFile string) (string, error) {
	file, err := os.Open(phpFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Regular expression to match the define statement
	re := regexp.MustCompile(fmt.Sprintf(`^.*define\(.*'%s'.*,.* '([^']+)'.*\).*$`, constantName))

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			return matches[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

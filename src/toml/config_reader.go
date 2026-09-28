package toml

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Path    string
	Name    string   `toml:"name"`
	Scripts []Script `toml:"scripts"`
}

type Script struct {
	File string `toml:"file"`
	Desc string `toml:"desc"`
}

func GetProjects() []Config {
	cur_path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	base_paths := []string{
		filepath.Dir(cur_path),
		"C:/xampp/www",
	}

	projects := []Config{}
	fmt.Print("\033[38;2;255;50;0m")

	for _, base_path := range base_paths {
		dir, err := os.ReadDir(base_path)
		if err != nil {
			fmt.Println(base_path + " not read: " + err.Error())
			continue
		}

		for _, item := range dir {
			if !item.IsDir() {
				continue
			}

			project_dir := base_path + "/" + item.Name()
			project_dir, _ = filepath.Abs(project_dir)
			project_file := project_dir + "/devtoolkit.toml"
			_, err := os.Stat(project_file)
			if err != nil {
				continue
			}

			data, err := os.ReadFile(project_file)
			if err != nil {
				fmt.Println(project_file + " not read: " + err.Error())
				continue
			}

			var cfg Config
			_, err = toml.Decode(string(data), &cfg)

			if err != nil {
				fmt.Println("Config not read from [" + project_file + "]" + ": " + err.Error())
				continue
			}

			if len(cfg.Scripts) == 0 {
				fmt.Println("No scripts in the project [" + project_file + "]\033[0m")
				continue
			}

			cfg.Path = project_dir
			proj_name := cfg.Name
			if proj_name == "" {
				proj_name = item.Name()
			}
			projects = append(projects, cfg)
		}

	}

	fmt.Print("\033[0m")

	slices.SortFunc(projects, func(a, b Config) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})

	return projects
}

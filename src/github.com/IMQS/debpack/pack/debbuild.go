package pack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"
)

type DebBuild struct {
	Group        string                        `json:"group"`
	Name         string                        `json:"name"`
	Version      string                        `json:"version"`
	Author       string                        `json:"author"`
	Description  string                        `json:"description"`
	Repository   string                        `json:"repository"`
	Package      string                        `json:"package"`      // the package name for the go build command
	Binary       string                        `json:"binary"`       // the name of the application
	Distribution string                        `json:"distribution"` //required for debian package
	WorkDir      string                        `json:"workdir"`
	Control      *Control                      `json:"control"`
	RepoDir      string                        `json:"-"`
	DebDir       string                        `json:"-"`
	Templates    map[string]*template.Template `json:"-"`
	ChangeLog    []string                      `json:"-"`
}

func NewDebBuild(filename string) (*DebBuild, error) {
	d := &DebBuild{}
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer, d)
	if err != nil {
		return nil, err
	}
	d.RepoDir = path.Dir(d.WorkDir) + "/GIT"
	d.DebDir = path.Dir(d.WorkDir) + "/DEBBASE"
	d.Templates = make(map[string]*template.Template, 0)
	for name, raw := range templates {
		var err error
		d.Templates[name], err = template.New(name).Parse(raw)
		if err != nil {
			return nil, err
		}
	}
	d.ChangeLog = append(d.ChangeLog, fmt.Sprintf("%s (%s)", d.Binary, d.Version), "")
	err = d.runCmd("mkdir", []string{"-p", d.WorkDir})
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DebBuild) Build() error {

	err := d.cloneOrPull()
	if err != nil {
		return err
	}
	fmt.Println("Done Git")

	err = d.compile()
	if err != nil {
		return err
	}
	fmt.Println("Done compile")

	os.Chdir(d.DebDir)

	err = d.systemd()
	if err != nil {
		return err
	}
	fmt.Println("Done systemd")

	err = d.doc()
	if err != nil {
		return err
	}
	fmt.Println("Done doc")

	err = d.debian()
	if err != nil {
		return err
	}
	fmt.Println("Done debian")

	os.Chdir(d.WorkDir)

	err = d.dpkg()
	if err != nil {
		return err
	}
	fmt.Println("Done dpkg")

	return nil

}

func (d *DebBuild) cloneOrPull() error {
	if _, err := os.Stat(d.RepoDir); os.IsNotExist(err) {
		// Repo does not exist
		err := d.runCmd("git", []string{"clone", d.Repository, d.RepoDir})
		if err != nil {
			return err
		}
	}

	os.Chdir(d.RepoDir)

	err := d.runCmd("git", []string{"checkout", "master"})
	if err != nil {
		return err
	}
	err = d.runCmd("git", []string{"pull"})
	if err != nil {
		return err
	}
	err = d.runCmd("git", []string{"submodule", "update", "--init", "--recursive"})
	if err != nil {
		return err
	}

	out, err := exec.Command("git", "log", "--pretty=format:%s", "--reverse").Output()
	if err != nil {
		return err
	}
	logentries := strings.Split(string(out), "\n")
	for _, entry := range logentries {
		d.ChangeLog = append(d.ChangeLog, fmt.Sprintf("  * %s", entry))
	}
	d.ChangeLog = append(d.ChangeLog, "")
	d.ChangeLog = append(d.ChangeLog, fmt.Sprintf(" -- IMQS <imqs@imqs.co.za>  %s", time.Now().Format(time.RFC1123Z)))

	os.Chdir(d.WorkDir)
	return nil
}

func (d *DebBuild) compile() error {
	os.Chdir(d.RepoDir)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Setenv("GOPATH", cwd)
	if err != nil {
		return err
	}

	// Here we just ignore errors, they will be caught on compile
	d.runCmd("go", []string{"get", "./..."})
	
	binFile := path.Join(d.DebDir, "/usr/local/bin/", d.Binary)

	err = d.runCmd("go", []string{"build", "-o", binFile, d.Package})
	if err != nil {
		return err
	}
	os.Chdir(d.WorkDir)
	return nil
}

func (d *DebBuild) systemd() error {
	err := d.runCmd("mkdir", []string{"-p", "lib/systemd/system"})
	if err != nil {
		return err
	}
	buffer := bytes.NewBufferString("")
	err = d.Templates["systemd"].Execute(buffer, d)
	if err != nil {
		return err
	}

	sysdFile := path.Join(d.DebDir, "/lib/systemd/system", d.Binary+".service")
	err = ioutil.WriteFile(sysdFile, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (d *DebBuild) doc() error {
	err := d.runCmd("mkdir", []string{"-p", "usr/share/doc/" + d.Binary})
	if err != nil {
		return err
	}

	changelog := path.Join(d.DebDir, "usr/share/doc/", d.Binary, "changelog")
	err = ioutil.WriteFile(changelog, []byte(strings.Join(d.ChangeLog, "\n")), 0644)
	if err != nil {
		return err
	}

	err = d.runCmd("cp", []string{changelog, changelog + ".Debian"})
	if err != nil {
		return err
	}

	err = d.runCmd("gzip", []string{"--best", "-f", changelog})
	if err != nil {
		return err
	}
	err = d.runCmd("gzip", []string{"--best", "-f", changelog + ".Debian"})
	if err != nil {
		return err
	}

	err = d.runCmd("mkdir", []string{"-p", "usr/share/man/man1"})
	if err != nil {
		return err
	}
	return nil
}

func (d *DebBuild) debian() error {
	err := d.runCmd("mkdir", []string{"-p", "DEBIAN"})
	if err != nil {
		return err
	}
	for _, name := range []string{"postinst", "prerm", "control"} {
		buffer := bytes.NewBufferString("")
		err = d.Templates[name].Execute(buffer, d)
		if err != nil {
			return err
		}

		debFile := path.Join(d.DebDir, "DEBIAN", name)

		err = ioutil.WriteFile(debFile, buffer.Bytes(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DebBuild) dpkg() error {
	filename := fmt.Sprintf("%s-%s.deb", d.Binary, d.Version)
	err := d.runCmd("fakeroot", []string{"dpkg-deb", "--build", "DEBBASE", filename})
	if err != nil {
		return err
	}
	return nil
}

func (d *DebBuild) runCmd(cmd string, args []string) error {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Printf("%s %v %s\n", cmd, args, out)
		return err
	}
	return nil
}

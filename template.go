package main

import (
  "flag"
	"fmt"
	"text/template"
	"os"
  "io/ioutil"
  "bufio"
  "gopkg.in/yaml.v2"
)

type File struct {
  Path string
  Content string
  Owner string
  Permissions string
}

type CloudConfig struct {
  DiscoveryUrl string
  Files []*File
}

type Units struct {
  Services []*Service
}

type Service struct {
  Name string
  Image string
  Command string
  Type string
  Filename string
  Restart string
  Envs []string
  Ports []string
  Xfleet []string
}

var template_dir string
var root string

func main() {
  // Figure out which directory we are going to prep files for
  template_def := os.ExpandEnv("$GOPATH/src/github.com/winkapp/clc/templates")
  flag.StringVar(&template_dir, "templates", template_def, "Your templates.")
  flag.StringVar(&root, "root", "./", "Directory for configs and output.")
  flag.Parse()
  command := flag.Arg(0)

  switch command {
  case "cc": cc(root)
  case "vagrant": vagrant(root)
  case "units": units(root)
  case "new": newCluster(root)
  }

}

func newCluster(path string) {
  cc(root)
  vagrant(root)
  units(root)
  return
}

func vagrant(path string) {
  ud(path)
  copyFile(path + "/Vagrantfile", template_dir + "/Vagrantfile.template")
  copyFile(path + "/config.rb", template_dir + "/config.rb.template")
  return
}

func units(path string) {
  unitsYaml := getFileBytes(path + "/units.yaml")
  var units Units
  err := yaml.Unmarshal(unitsYaml, &units)
  checkError(err)
  for _, service := range units.Services {
    t := unitsTemplate(service)
    fileName := unitFileName(service)
    f, err := os.Create(path + "/units/" + fileName)
    checkError(err)
    defer f.Close()
    w := bufio.NewWriter(f)
    err = t.Execute(w, service)
    checkError(err)
    w.Flush()
  }
}

func ud(path string) {
  // TODO: Open the files document for file includes
  file1 := File{
    Content: getFile("./etcd_env.sh"),
    Path: "/home/core/etcd_env.sh",
    Owner: "core:core",
    Permissions: "744",
  }

  data := CloudConfig{
    Files: []*File{&file1},
  }

  t := udTemplate(data)
  err := t.Execute(os.Stdout, data)
  checkError(err)
}

func cc(path string) {
  data := ccData(path)
  t := ccTemplate(data)
  err := t.Execute(os.Stdout, data)
	checkError(err)
}

func ccData(path string) (data CloudConfig) {
  // Open the files document for file includes
  file1 := File{
    Content: getFile("./etcd_env.sh"),
    Path: "/home/core/etcd_env.sh",
    Owner: "core:core",
    Permissions: "744",
  }

  // Open the document for the discovery url
  discovery_url := getFile(path + "/discovery_url")

  // Include files on cloud config
  data = CloudConfig{
    DiscoveryUrl: discovery_url,
    Files: []*File{&file1},
  }
  return
}

func unitsTemplate(data *Service) (t *template.Template) {
  t = getTemplate("Service Template", "service.template")
  return
}

func ccTemplate(data CloudConfig) (t *template.Template) {
  t = getTemplate("Cloud Config Template", "cloud-config.template")
  return
}

func udTemplate(data CloudConfig) (t *template.Template) {
  t = getTemplate("User Data Template", "user-data.template")
  return
}

func unitFileName(service *Service) string {
  switch service.Type {
  case "single":
    service.Filename = (service.Name + ".service")
  case "multi":
    service.Filename = (service.Name + "@.service")
  }

  return service.Filename
}

func addFile(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s\n", path)
  return nil
}

func getFileBytes(path string) []byte {
  dat, err := ioutil.ReadFile(path)
  checkError(err)
  return dat
}

func getFile(path string) string {
  return string(getFileBytes(path))
}

func copyFile(dstPath string, srcPath string) {
  src, err := ioutil.ReadFile(srcPath)
  checkError(err)
  mode := int(0644)
  err = ioutil.WriteFile(dstPath, src, os.FileMode(mode))
  checkError(err)
  return
}

func getTemplate(name string, filename string) (t *template.Template) {
  templ := getFile(template_dir + "/" + filename)

  t = template.New(name)
  t, err := t.Parse(templ)
  checkError(err)
  return
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}

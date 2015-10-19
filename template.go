package main

import (
  "flag"
	"fmt"
	"text/template"
	"os"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "github.com/winkapp/libclc"
  "path"
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

var template_dir string
var root string
var config_file string

func main() {
  // Figure out which directory we are going to prep files for
  flag.StringVar(&template_dir, "templates", "", "Your templates.")
  flag.StringVar(&root, "root", "./", "Directory for configs and output.")
  flag.StringVar(&config_file, "config", "clc.yaml", "Optional config file.")
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
  var units libclc.UnitConfig
  err := yaml.Unmarshal(unitsYaml, &units)
  checkError(err)
  t := unitsTemplate()
  err = libclc.WriteUnits(&units, t, path + "/units")
  checkError(err)
}

func ud(path string) {
  // TODO: Open the files document for file includes
  file1 := libclc.File{
    Content: getFile("./etcd_env.sh"),
    Path: "/home/core/etcd_env.sh",
    Owner: "core:core",
    Permissions: "744",
  }

  data := libclc.CloudConfig{
    Files: []*libclc.File{&file1},
  }

  t := udTemplate(data)
  err := libclc.WriteUserData(&data, t, path + "/user-data")
  checkError(err)
}

func cc(path string) {
  data := ccData(path)
  t := ccTemplate(data)
  err := libclc.WriteCloudConfig(&data, t, path + "/cloud-config")
	checkError(err)
}

func ccData(path string) (data libclc.CloudConfig) {
  // Open the files document for file includes
  file1 := libclc.File{
    Content: getFile("./etcd_env.sh"),
    Path: "/home/core/etcd_env.sh",
    Owner: "core:core",
    Permissions: "744",
  }

  // Open the document for the discovery url
  discovery_url := getFile(path + "/discovery_url")

  // Include files on cloud config
  data = libclc.CloudConfig{
    DiscoveryUrl: discovery_url,
    Files: []*libclc.File{&file1},
  }
  return
}

func unitsTemplate() (t *template.Template) {
  t = getTemplate("Service Template", "service.template")
  return
}

func ccTemplate(data libclc.CloudConfig) (t *template.Template) {
  t = getTemplate("Cloud Config Template", "cloud-config.template")
  return
}

func udTemplate(data libclc.CloudConfig) (t *template.Template) {
  t = getTemplate("User Data Template", "user-data.template")
  return
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
  // If a template directory not set, use default template.
  if template_dir == "" {
    return nil
  }
  // If a template file does not exist, use default template.
  _, err := os.Stat(path.Join(template_dir, filename))
  if err != nil {
    return nil
  }
  templ := getFile(path.Join(template_dir, filename))

  t = template.New(name)
  t, err = t.Parse(templ)
  checkError(err)
  return
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}

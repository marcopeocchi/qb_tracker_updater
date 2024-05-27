package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	confPath string
	profile  int

	profiles = []string{
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best_ip.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_udp.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_http.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_https.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_ws.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_ip.txt",
	}
)

func currentUser() (string, error) {
	user, err := exec.Command("whoami").Output()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimRight(user, "\n")), nil
}

func getGithubResponse() ([]byte, error) {
	res, err := http.Get(profiles[profile-1])
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parseGithubResponse(b []byte) []byte {
	var buf bytes.Buffer

	b = bytes.ReplaceAll(b, []byte("\n\n"), []byte("\n"))
	bb := bytes.Split(b, []byte("\n"))

	for _, line := range bb[:len(bb)-1] {
		buf.Write(line)
		buf.WriteRune('\n')
	}

	buf.Truncate(buf.Len() - 2)
	return buf.Bytes()
}

func main() {
	user, err := currentUser()
	if err != nil {
		log.Fatalln("cannot infer current user", err.Error())
	}

	flag.StringVar(
		&confPath,
		"conf",
		"/home/"+user[0:len(user)-1]+"/.config/qBittorrent/qBittorrent.conf",
		"qBittorrent config path",
	)
	flag.IntVar(
		&profile,
		"profile",
		2,
		"select profile level:\n 1: best\n 2: best_ip\n 3: all_udp\n 4: all_http\n 5: all_https\n 6: all_ws\n 7: all_ip\n",
	)
	flag.Parse()

	ghResponse, err := getGithubResponse()
	if err != nil {
		log.Fatalln(err)
	}

	fd, err := os.Open(confPath)
	if err != nil {
		log.Fatalln(err)
	}

	configContent, err := io.ReadAll(fd)
	if err != nil {
		log.Fatalln(err)
	}

	sep := []byte("Bittorrent\\TrackersList=")

	if !bytes.Contains(configContent, sep) {
		log.Fatalln("failed to update trackers")
	}

	var (
		lineSep   = []byte("/announce")
		fromIndex = bytes.Index(configContent, sep)
		toIndex   = bytes.LastIndex(configContent, lineSep)
		cut       = configContent[fromIndex : toIndex+len(lineSep)]
	)

	modified := bytes.Replace(
		configContent,
		cut,
		append(sep, parseGithubResponse(ghResponse)...),
		1,
	)

	if err := os.WriteFile(confPath, modified, os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	log.Println("updated trackers")
}

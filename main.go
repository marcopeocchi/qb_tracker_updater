package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	_user, err := exec.Command("whoami").Output()

	if err != nil {
		panic("Error: Cannot infer current user")
	}

	user := string(_user)

	configPath := flag.String("conf", "/home/"+user[0:len(user)-1]+"/.config/qBittorrent/qBittorrent.conf", "qBittorrent config path")
	profile := flag.Int("profile", 2, "select profile level:\n 1: best\n 2: best_ip\n 3: all_udp\n 4: all_http\n 5: all_https\n 6: all_ws\n 7: all_ip\n")

	flag.Parse()

	profiles := [7]string{
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best_ip.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_udp.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_http.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_https.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_ws.txt",
		"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_ip.txt",
	}

	res, err := http.Get(profiles[*profile-1])
	if err != nil {
		panic("Network error: Cannot retrieve trackers from Github")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("Network error: Cannot parse response")
	}

	trackers := strings.Split(string(body), "\n")

	acc := 0
	for _, line := range trackers {
		if len(line) != 0 {
			trackers[acc] = line
			acc++
		}
	}

	trackers = trackers[:acc]

	list := reduce(trackers, func(prev, curr string) string {
		return prev + "\\n" + curr
	}, "")

	list = list[2:]

	content, err := ioutil.ReadFile(*configPath)

	if err != nil {
		panic("IO Error: Config file not found")
	}

	buff := ""

	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "Bittorrent\\TrackersList") {
			config := strings.Split(line, "=")
			config[1] = list
			line = config[0] + "=" + config[1]
		}
		buff += (line + "\n")
	}

	ioutil.WriteFile(*configPath, []byte(buff), 0)
}

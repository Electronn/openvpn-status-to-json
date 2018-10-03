package main

import (
	"fmt"
	"strings"
	"regexp"
	"io/ioutil"
	"flag"
	"log"
	"os"
)

var (
        ovpnlog = flag.String("ovpn.log", "", "Absolute path for OpenVPN server log")
	resultstring,resclient,resroute,resline string
)

func main() {
	flag.Parse()	
        if *ovpnlog == "" {
                log.Fatal("OpenVPN status log absolute path must be set with '-ovpn.log' flag")
        }
        if _, err := os.Stat(*ovpnlog); os.IsNotExist(err) { log.Fatal("File: ",*ovpnlog," does not exists")}
	b, _:= ioutil.ReadFile(*ovpnlog)
        lines := strings.Split(string(b), "\n")
        for _, line := range lines { resline = resline + " " + line }
	reg := regexp.MustCompile("OpenVPN CLIENT LIST Updated,(.*?) Common Name,Real Address,Bytes Received,Bytes Sent,Connected Since (.*?) ROUTING TABLE Virtual Address,Common Name,Real Address,Last Ref (.*?) GLOBAL STATS Max bcast/mcast queue length,(.*?) END")
	regclients := regexp.MustCompile("(.*),(.*),(.*),(.*),")
	regroutes := regexp.MustCompile("(.*),(.*),(.*),")
	regipport := regexp.MustCompile("(.*):(.*)")
	s := regexp.MustCompile("((Mon|Tue|Thu|Wed|Fri|Sat|Sun) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \\d{1,2} \\d{2}:\\d{2}:\\d{2} \\d{4})").Split(reg.FindStringSubmatch(resline)[2], -1)
	s2 := regexp.MustCompile("((Mon|Tue|Thu|Wed|Fri|Sat|Sun) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \\d{1,2} \\d{2}:\\d{2}:\\d{2} \\d{4})").Split(reg.FindStringSubmatch(resline)[3], -1)
	resultstring = resultstring + "{" + "\"updated\": " + "\"" + reg.FindStringSubmatch(resline)[1] + "\"" + "," + "\"clients\": ["
	for i := 0; i < len(s) - 1; i++ {	
		resclient = resclient + "{" + "\"client\": " + "\"" + regclients.FindStringSubmatch(s[i])[1] + "\"," + "\"remote\": " + "\"" + regipport.FindStringSubmatch(regclients.FindStringSubmatch(s[i])[2])[1] + "\"," + "\"bytes_received\": " + "\"" + regclients.FindStringSubmatch(s[i])[3] + "\"," + "\"bytes_sent\": " + "\"" + regclients.FindStringSubmatch(s[i])[4] + "\"}"
		if i == len(s) - 2 { resclient = resclient + "]," } else { resclient = resclient + "," }
	}
	resultstring = resultstring + resclient
	resultstring = resultstring + "\"routing\": ["
	for i := 0; i < len(s2) - 1; i++ {
		resroute = resroute + "{" + "\"local_ip\": " + "\"" + regroutes.FindStringSubmatch(s2[i])[1] + "\"," + "\"client\": " + "\"" + regroutes.FindStringSubmatch(s2[i])[2] + "\"," + "\"real_ip\" :" + "\"" + regipport.FindStringSubmatch(regroutes.FindStringSubmatch(s2[i])[3])[1] + "\"}"
		if i == len(s2) - 2 { resroute = resroute + "]," } else { resroute = resroute + "," }
	}
	resultstring = resultstring + resroute
	resultstring = resultstring + "\"max_bcast_mcast_queue\": " + "\"" + reg.FindStringSubmatch(resline)[4] + "\"}"
	fmt.Println(resultstring)
}

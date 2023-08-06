package util

import (
	"Server/graph/model"
	"bufio"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"os"
	"strings"
	"time"
)

type netStat struct {
	name string
	sent uint64
	recv uint64
}

type prettyNetStat struct {
	name string
	sent string
	recv string
}

func readNetstat() map[string]string {
	netFile, err := os.Open("/proc/net/netstat")
	if err != nil {
		log.Fatal("error occured when reading /proc/net/netstat: ", err)
	}
	scanner := bufio.NewScanner(netFile)
	assembly := make([]map[string]string, 0)
	titles := make([]string, 0)
	datas := make([]string, 0)
	counter := 1
	for scanner.Scan() {
		if counter%2 != 0 {
			titles = strings.Split(scanner.Text(), " ")
		} else {
			datas = strings.Split(scanner.Text(), " ")
			assembly = append(assembly, make(map[string]string))
			for i, data := range datas {
				assembly[len(assembly)-1][titles[i]] = data
			}
		}
		counter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error occured when scanning /proc/net/netstat: ", err)
	}
	return map[string]string{"InBytes": assembly[1]["InOctets"], "OutBytes": assembly[1]["OutOctets"]}
}

func CurrentStat(eachNic bool) []netStat {
	stat, err := net.IOCounters(eachNic)
	if err != nil {
		log.Fatal("error occured when running net.IOCounters: ", err)
	}
	rt := make([]netStat, 0)
	for _, value := range stat {
		rt = append(rt, netStat{
			name: value.Name,
			sent: value.BytesSent,
			recv: value.BytesRecv,
		})
	}
	return rt
}

func NetRate(interval int, eachNic bool) []netStat {
	startTraffic := CurrentStat(eachNic)
	time.Sleep(time.Duration(interval) * time.Millisecond)
	endTraffic := CurrentStat(eachNic)

	stat := make([]netStat, 0)

	for i, val := range startTraffic {
		stat = append(stat, netStat{
			val.name,
			endTraffic[i].recv - val.recv,
			endTraffic[i].sent - val.sent,
		})
	}
	return stat
}

func goPretty(stat []netStat, postfix string) []prettyNetStat {
	pretty := make([]prettyNetStat, 0)
	for _, val := range stat {
		pretty = append(pretty, prettyNetStat{
			name: val.name,
			sent: ProperUnit(val.sent, 1) + postfix,
			recv: ProperUnit(val.recv, 1) + postfix,
		})
	}
	return pretty
}

func NetRatePretty(stat []netStat) []prettyNetStat {
	return goPretty(stat, "/s")
}

func NetPretty(stat []netStat) []prettyNetStat {
	return goPretty(stat, "")
}

func Netstat2gql(stat []netStat, pretty []prettyNetStat) []*model.Net {
	good := make([]*model.Net, 0)
	for i, val := range stat {
		good = append(good, &model.Net{
			Name: val.name,
			Sent: int(val.sent),
			Recv: int(val.recv),
			PrettyRecv: pretty[i].recv,
			PrettySent: pretty[i].sent,
		})
	}
	return good
}

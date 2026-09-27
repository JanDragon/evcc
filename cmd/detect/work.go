package detect

import (
	"sort"
	"sync"

	"github.com/evcc-io/evcc/cmd/detect/tasks"
	"github.com/evcc-io/evcc/util"
)

func workers(log *util.Logger, num int, ips <-chan string, hits chan<- []tasks.Result) *sync.WaitGroup {
	var wg sync.WaitGroup

	for range num {
		wg.Go(func() {
			for ip := range ips {
				res := taskList.Test(log, "", tasks.ResultDetails{IP: ip})
				hits <- res
			}
		})
	}

	return &wg
}

func Work(log *util.Logger, num int, hosts []string) []tasks.Result {
	ip := make(chan string)
	hits := make(chan []tasks.Result)
	done := make(chan struct{})

	// log.INFO.Println(
	// 	"\n" +
	// 		strings.Join(
	// 			lo.Map(taskList.tasks, func(t tasks.Task) string {
	// 				return fmt.Sprintf("task: %s\ttype: %s\tdepends: %s\n", t.ID, t.Type, t.Depends)
	// 			}).([]string),
	// 			"",
	// 		),
	// )

	wg := workers(log, num, ip, hits)

	var res []tasks.Result
	go func() {
		for hits := range hits {
			res = append(res, hits...)
		}
		done <- struct{}{}
	}()

	for _, host := range hosts {
		ip <- host
	}

	close(ip)
	wg.Wait()

	close(hits)
	<-done

	return postProcess(res)
}

func postProcess(res []tasks.Result) []tasks.Result {
	// sort by host
	sort.Slice(res, func(i, j int) bool {
		if res[i].ResultDetails.IP == res[j].ResultDetails.IP {
			return res[i].Type < res[j].Type
		}
		return res[i].ResultDetails.IP < res[j].ResultDetails.IP
	})

	return res
}

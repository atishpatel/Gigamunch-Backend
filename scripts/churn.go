package main

import "fmt"

type sub struct {
	StartWeek int
	EndWeek   int
}

var subs = []sub{
	sub{StartWeek: 3, EndWeek: 5},
}

func main() {
	printWeeklyChrunRate()
	printMontlyChurnRate()
}
func printWeeklyChrunRate() {
	var churnRate [60]float32
	for i := range churnRate {
		subsThisWeek := 0
		cancelsThisWeek := 0
		for j := range subs {
			if i+1 >= subs[j].StartWeek && i+1 < subs[j].EndWeek {
				subsThisWeek++
			}
			if subs[j].EndWeek == i+1 {
				cancelsThisWeek++
			}
		}
		churnRate[i] = float32(cancelsThisWeek) / float32(subsThisWeek)
	}
	for i := range churnRate {
		fmt.Printf("week: %d \tchurn rate: %3.2f %% \n", i, churnRate[i]*100)
	}
}

func printMontlyChurnRate() {
	var churnRate [13]float32
	for i := range churnRate {
		subsThisMonth := 0
		cancelsThisMonth := 0
		for j := range subs {
			if (i*4)+1 >= subs[j].StartWeek && (i+1)*4 < subs[j].EndWeek {
				subsThisMonth++
			}
			if (i*4)+1 >= subs[j].EndWeek && (i+1)*4 < subs[j].EndWeek {
				cancelsThisMonth++
			}
		}
		churnRate[i] = float32(cancelsThisMonth) / float32(subsThisMonth)
	}
	for i := range churnRate {
		fmt.Printf("month: %d \tchurn rate: %3.2f %% \n", i, churnRate[i]*100)
	}
}

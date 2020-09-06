package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"REGSPLT/rs"
	ars "regsplt/rs"
)

var errStream = os.Stderr
var errFile string
var err error

func main() {

	errFile = "data.txt"
	f := "DATA"
	fmt.Println("Welcome to the playground!")
	errStream, err = os.Create(errFile)
	defer errStream.Close()
	fmt.Println("The time is", time.Now())

	file, err := os.Open("cigar.txt")
	//file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//	for scanner.Scan() {				// loop until EOF
	//		fmt.Println(scanner.Text())
	//	}

	for j := 1; j <= 28; j++ {
		if scanner.Scan() { //  test for EOF
			f = scanner.Text()
			//			fmt.Println("Fields in f:", len(strings.Fields(f)))
			f = strings.Fields(f)[5]
			// fmt.Println(f)
			//fmt.Println(j)
			s := regexp.MustCompile("[0-9]+[S,D,I,H]").Split(f, -1)
			// fmt.Println(s)
			mTotal := 0
			rS := new(ars.RStats)
			for i := 0; i < len(s); i++ {
				//fmt.Println(">",s[i])
				aStr := strings.Split(s[i], "M")[0]
				//fmt.Println(aStr)
				n, _ := strconv.Atoi(aStr)
				ars.RollingStat(float64(n), rS)
				mTotal = mTotal + n
			}
			//			fmt.Println(" M count: ", mTotal)
			//			fmt.Println("Count: ", rS.N, "Mean: ", rS.M1, "Std Dev: ", (math.Sqrt(rS.M2 / ((float64(rS.N)) - 1.0))), " Skew :", (math.Pow(float64(rS.N)-1.0, 1.5) / float64(rS.N) * rS.M3 / (math.Pow(rS.M2, 1.5))))

			// method 2
			rSM := new(ars.RStats)
			rSI := new(ars.RStats)
			rSD := new(ars.RStats)
			rSO := new(ars.RStats)
			oStr := ""
			rQI := new(ars.RQuant)
			ars.Reinit(rQI, 0.5)
			rQM := new(ars.RQuant)
			ars.Reinit(rQM, 0.5)
			rQD := new(ars.RQuant)
			ars.Reinit(rQD, 0.5)
			mStr := ""
			iStr := ""
			dStr := ""
			sCig1 := regexp.MustCompile(`[A-Z]`).Split(f, -1)
			sCig2 := regexp.MustCompile(`[0-9]+`).Split(f, -1)
			for i := 1; i < len(sCig1); i++ {
				// fmt.Print(sCig1[i-1], sCig2[i])
				switch sCig2[i] {
				case "M":
					if cgI, err := strconv.Atoi(sCig1[i-1]); err == nil {
						ars.QuantRoller(float64(cgI), rQM)
						mStr = mStr + sCig1[i-1] + ","
						ars.RollingStat(float64(cgI), rSM)
					}
				case "I":
					if cgI, err := strconv.Atoi(sCig1[i-1]); err == nil {
						ars.QuantRoller(float64(cgI), rQI)
						iStr = iStr + sCig1[i-1] + ","
						ars.RollingStat(float64(cgI), rSI)
					}
				case "D":
					if cgI, err := strconv.Atoi(sCig1[i-1]); err == nil {
						ars.QuantRoller(float64(cgI), rQD)
						dStr = dStr + sCig1[i-1] + ","
						ars.RollingStat(float64(cgI), rSD)
					}
				default:
					if cgI, err := strconv.Atoi(sCig1[i-1]); err == nil {
						ars.RollingStat(float64(cgI), rSO)
						oStr = oStr + sCig1[i-1] + sCig2[i] + ","
					}
				}

			}
			fmt.Printf(" M %% of MDI: %.3f", (rSM.Sum / (rSM.Sum + rSI.Sum + rSD.Sum)))
			fmt.Print(" Cigar:", len(sCig1))
			fmt.Printf(" M quant 0.5: %.2f", ars.RQuantResult(rQM))
			fmt.Print(" Count: ", rSM.N, " Max: ", rSM.Max, " Sum: ", rSM.Sum)
			fmt.Printf(" I quant: %.2f", ars.RQuantResult(rQI))
			fmt.Print(" Count: ", rSI.N, " Max: ", rSI.Max, " Sum: ", rSI.Sum)
			fmt.Printf(" D quant: %.2f", ars.RQuantResult(rQD))
			fmt.Print(" Count: ", rSD.N, " Max: ", rSD.Max, " Sum: ", rSD.Sum) //"\n", dStr
			fmt.Println(" Other: ", oStr, " Count: ", rSO.N, " Max: ", rSO.Max, " Sum: ", rSO.Sum)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

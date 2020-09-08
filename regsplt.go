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
	fmt.Println("Welcome! NOTE: regsplt will stop after 500,000 lines")
	errStream, err = os.Create(errFile)
	defer errStream.Close()
	fmt.Println("The time is", time.Now())

	// /hpcfs/groups/phoenix-hpc-rc003/joe/correction/chr6HG002
	// AZ-POL-chr6alignAshRef.sam
	// AZNOPOLalignchr6-Azref.sam

	file, err := os.Open("c6polAref.sam")
	//file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	/* 	file, err := os.Open("cigar.txt")
	   	//file, err := os.Open("file.txt")
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	defer file.Close() */

	scanner := bufio.NewScanner(file)
	//	for scanner.Scan() {				// loop until EOF
	//		fmt.Println(scanner.Text())
	//	}

	watchdog := 0
	for { // j := 1; j <= ; j++                ENDLESS LOOP
		if scanner.Scan() { //  test for EOF
			fdata := scanner.Text()
			//			fmt.Println("Fields in f:", len(strings.Fields(f)))
			if fdata[:1] == "@" {
				continue
			} // header of sam file. continue go to next iteration of loop.
			watchdog = watchdog + 1
			if watchdog > 500000 {
				fmt.Println("\n Break on watchdog counter, ", watchdog)
				break
			}
			f = strings.Fields(fdata)[5]
			/* 			// fmt.Println(f)
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
			*/
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
			fmt.Printf(" M %% of MDI:, %.3f,", (rSM.Sum / (rSM.Sum + rSI.Sum + rSD.Sum)))
			fmt.Printf(" M,I/D,S/H:, %.0f, %.0f, %.0f,", rSM.Sum, (rSI.Sum + rSD.Sum), rSO.Sum)
			fmt.Print(" c-elements:,", len(sCig1))
			/*   			//			fmt.Printf(" M quant 0.5: %.2f", ars.RQuantResult(rQM))
			fmt.Print(", M Count: ", rSM.N, " Max: ", rSM.Max, " Sum: ", rSM.Sum)
			//			fmt.Printf(" I quant: %.2f", ars.RQuantResult(rQI))
			fmt.Print(", I Count: ", rSI.N, " Max: ", rSI.Max, " Sum: ", rSI.Sum)
			//			fmt.Printf(" D quant: %.2f", ars.RQuantResult(rQD))
			fmt.Print(", D Count: ", rSD.N, " Max: ", rSD.Max, " Sum: ", rSD.Sum) //"\n", dStr
			fmt.Print(", Other: ", oStr, " Count: ", rSO.N, " Max: ", rSO.Max, " Sum: ", rSO.Sum) */
			fmt.Println(", Read:, ", strings.Fields(fdata)[0], ", chr:, ", strings.Fields(fdata)[2], ", pos:, ", strings.Fields(fdata)[3])
		} else {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

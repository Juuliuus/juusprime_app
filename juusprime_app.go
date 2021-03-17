package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	jm "github.com/Juuliuus/juusmenu"
	jup "github.com/Juuliuus/juusprime"
)

var (
	menuMain      *jm.Menu
	menu29        *jm.Menu
	menu29Details *jm.Menu
	menu31        *jm.Menu
	menuCalcs     *jm.Menu
	menuGen       *jm.Menu
	input         string
	wasCanceled   bool
)

const (
	muMyMainID = iota
	muMymenu29
	muMymenu29Details
	muMymenu31
	muMymenuCalcs
	muMymenuGen
)

var askExit bool

func main() {

	if x := handleFlags(); x > -1 {
		os.Exit(x)
	}

	if !jup.IsConfigured() {
		jup.Configure()
	}
	initMenus()

	switch askExit {
	case false:
		menuMain.Start() //kill phrase exits immediately
	case true:
		runMenus() //has confirmation on Kill Phrase use
	}
}

func buildMenu(menu *jm.Menu) {
	switch menu.GetID() {

	case muMyMainID:
		menu.SetMenuBreakItem("q", "QUIT", func() { fmt.Println("Ok, quitting. Bye-Bye!") })
		menu.AddSubMenu(menu29, "29", "LTE 29 primes")
		menu.AddSubMenu(menu31, "31", "GTE 31 primes")
		menu.AddSubMenu(menuCalcs, "calc", "Various calculations to help with testing")
		menu.AddSubMenu(menuGen, "1", "Generate juusprime Tuplets")

		menu.AddMenuEntry("s", "Show Symbol help", func() {
			jup.HelpSymbols()
		})
		menu.AddMenuEntry("sm", "Show Symbol Math", func() {
			jup.HelpSymbolsMath()
		})
		menu.AddMenuEntry("config", "Set/edit your configuration settings/paths", func() {
			jup.Configure()
		})
		menu.AddMenuEntry("f", "Show Output File help", func() {
			jup.HelpOutputFiles()
		})
		menu.AddMenuEntry("c", "Show current config settings", func() {
			cfg := jup.ConfigFilename()
			if jup.FileExists(cfg) {
				fmt.Println("config file:", cfg)
			}
			fmt.Println(jup.Basis29PathStr, ":", jup.Basis29Path)
			fmt.Println(jup.DataPathStr, ":", jup.DataPath)
		})
	case muMymenuGen:
		menu.SetMenuBreakItem("b", "Back", func() {})
		menu.AddMenuEntry("basis", "Generate 29Basis file (required file, only needs to be done once)", func() {
			jup.GenerateBasis()
		})
		menu.AddMenuEntry("1", "Generate juusPrime Tuplets (29basis file required)", func() {
			jup.GeneratePrimeTupletsInteractive()
		})
	case muMymenuCalcs:
		menu.SetMenuBreakItem("b", "Back", func() {})
		menu.AddMenuEntry("tb", "TNumber to basis number (Tnum 1 based, basis 0 based)", func() {
			basis := big.NewInt(0)
			tNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInput("Enter TNumber(s) comma separated:", "28,1000000,1000000000000", "x"); wasCanceled {
				return
			}
			sl := strings.Split(input, ",")

			for i := range sl {
				fmt.Sscan(sl[i], tNum)
				jup.TNumToBasisNum(tNum, basis)
				fmt.Println(fmt.Sprintf("TNumber %v is in basis-%v", tNum, basis))
				fmt.Println("")
			}
		})
		menu.AddMenuEntry("bt", "Basis number to from/to TNumber (Tnum 1 based, basis 0 based)", func() {
			basis := big.NewInt(0)
			fromNum := big.NewInt(0)
			toNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInput("Enter Basis Numbers comma separated:", "0,1,1000000", "x"); wasCanceled {
				return
			}
			sl := strings.Split(input, ",")

			for i := range sl {
				fmt.Sscan(sl[i], basis)
				jup.BasisToTNumRange(basis, fromNum, toNum)
				fmt.Println("basis", basis, ":", fromNum, "-", toNum)
				fmt.Println("")
			}
		})
		menu.AddMenuEntry("it", "Integer to TNumber", func() {
			I := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInput("Enter Integers comma separated:", "25,101,1000000", "x"); wasCanceled {
				return
			}
			sl := strings.Split(input, ",")

			for i := range sl {
				fmt.Sscan(sl[i], I)
				tNum := jup.IntToTNum(I)
				fmt.Println(fmt.Sprintf("Integer %v is in TNumber %v", I, tNum))
				fmt.Println("")
			}
		})
		menu.AddMenuEntry("ti", "TNumber to Integer, and its range", func() {
			end := big.NewInt(0)
			tNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInput("Enter TNumbers comma separated:", "94090, 946644", "x"); wasCanceled {
				return
			}
			sl := strings.Split(input, ",")

			big29 := big.NewInt(29)

			for i := range sl {
				fmt.Sscan(sl[i], tNum)
				I := jup.TNumToInt(tNum)
				end.Add(I, big29)
				fmt.Println(fmt.Sprintf("TNumber %v starts at Integer %v, and ends at %v", tNum, I, end))
				fmt.Println("")
			}
		})
		menu.AddMenuEntry("N", "Get N from a TNumber", func() {
			n := big.NewInt(0)
			tNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInput("Enter TNumbers comma separated:", "94090, 946644", "x"); wasCanceled {
				return
			}
			sl := strings.Split(input, ",")

			p31 := jup.NewPrimeGTE31(big.NewInt(31))

			for i := range sl {
				fmt.Sscan(sl[i], tNum)
				jup.GetNfromTNum(tNum, p31, n)
				fmt.Println(fmt.Sprintf("TNum: %v, n: %v", tNum, n))
				fmt.Println("")
			}

		})
		menu.AddMenuEntry("H", "Human Readable from raw data (Give a TNum and Effect ID [0, 1, 2, or 3])", func() {
			tNum := big.NewInt(0)
			effect := 0
			str1, str2 := "n/a", "n/a"

			if input, wasCanceled = jup.GetUserInput("Enter TNumbers comma separated:", "94090, 94664", "x"); wasCanceled {
				return
			}
			sl := strings.Split(input, ",")

			if input, wasCanceled = jup.GetUserInputInteger("Enter Effect ID [0, 1, 2, or 3]:", "0", "x"); wasCanceled {
				return
			}
			effect, _ = strconv.Atoi(input)

			for i := range sl {
				fmt.Sscan(sl[i], tNum)
				jup.HumanReadable(tNum, &effect, &str1, &str2, os.Stdout)
				fmt.Println("")
			}

		})
	case muMymenu31:
		menu.SetMenuBreakItem("b", "Back", func() {})
		menu.AddMenuEntry("31", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(31)))
		})
		menu.AddMenuEntry("37", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(37)))
		})
		menu.AddMenuEntry("41", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(41)))
		})
		menu.AddMenuEntry("43", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(43)))
		})
		menu.AddMenuEntry("47", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(47)))
		})
		menu.AddMenuEntry("49", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(49)))
		})
		menu.AddMenuEntry("53", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(53)))
		})
		menu.AddMenuEntry("59", "Details", func() {
			printGTE31s(jup.NewPrimeGTE31(big.NewInt(59)))
		})
	case muMymenu29:
		menu.SetMenuBreakItem("b", "Back", func() {})
		menu.AddSubMenu(menu29Details, "d", "The LTE 29 primes details")
		menu.AddMenuEntry("gen29", "Generate alternative/custom 29 Basis files (for experienced users)", func() {
			jup.GenerateBasisInteractive()
		})
		menu.AddMenuEntry("gen23", "Generate the LTE 23 Tuplets (TNumbers 1 to 27) ", func() {
			jup.GeneratePrimes7to23()
		})
	case muMymenu29Details:
		menu.SetMenuBreakItem("b", "Back", func() {})
		menu.AddMenuEntry("M", "back to Main Menu", func() {
			menuMain.Start()
		})
		menu.AddMenuEntry("7", "7 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(7)))
		})
		menu.AddMenuEntry("11", "11 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(11)))
		})
		menu.AddMenuEntry("13", "13 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(13)))
		})
		menu.AddMenuEntry("17", "17 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(17)))
		})
		menu.AddMenuEntry("19", "19 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(19)))
		})
		menu.AddMenuEntry("23", "23 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(23)))
		})
		menu.AddMenuEntry("29", "29 details and raw data", func() {
			printLTE29s(jup.NewPrimeLTE29(big.NewInt(29)))
		})
	}
}

func initMenus() {
	jm.MenuOptions.SetKillPhrase("bye")
	jm.MenuOptions.AlignRight()

	menuMain = jm.NewMenu("Main")
	menu29 = jm.NewMenu("Primes LTE 29")
	menu29Details = jm.NewMenu("Details LTE 29")
	menu31 = jm.NewMenu("Primes GTE 31")
	menuCalcs = jm.NewMenu("Calcs")
	menuGen = jm.NewMenu("Generation")

	menuMain.SetID(muMyMainID)
	menu29.SetID(muMymenu29)
	menu29Details.SetID(muMymenu29Details)
	menu31.SetID(muMymenu31)
	menuCalcs.SetID(muMymenuCalcs)
	menuGen.SetID(muMymenuGen)

	buildMenu(menuMain)
	buildMenu(menu29)
	buildMenu(menu29Details)
	buildMenu(menu31)
	buildMenu(menuCalcs)
	buildMenu(menuGen)
}

func runMenus() {
	exitFlag := false
	menuExit := jm.NewMenu("Really Exit?")
	menuExit.SetMenuBreakItem("c", "cancel", func() {})
	menuExit.AddMenuEntry("e", "Exit", func() { exitFlag = true })
	menuExit.SetChooseOne(true)

	for {
		menuMain.Start()
		if jm.MenuSystem.WasKilled() {
			jm.MenuSystem.UnKill()
			menuExit.Start()
			//they might have used the killPhrase on the ChooseOne menu too!
			exitFlag = exitFlag || jm.MenuSystem.WasKilled()
			if exitFlag {
				break
			} else {
				continue
			}
		} else {
			break
		}
	}
}

func handleFlags() (exitcode int) {

	if len(os.Args) > 1 {

		const overwrites = "Be aware automation overwrites output files, if they exist, automatically."

		switch strings.ToUpper(os.Args[1]) {
		case "--HELP", "-H", "HELP":
			fmt.Println("")
			fmt.Println("For automation flags help use `juusprime_app automate --help`")
			fmt.Println(overwrites)
			fmt.Println("")
			//don't return from here this is the only way, it seems, to intercept help messages.
		case "--VERSION", "-V", "VERSION":
			fmt.Println("")
			fmt.Println("juusprime-app, Version 1.0.1, March 2021")
			fmt.Println("https://github.com/Juuliuus/juusprime")
			fmt.Println("https://github.com/Juuliuus/juusprime_app")
			fmt.Println("")
			return 0
		case "AUTOMATE":
			if len(os.Args) > 2 {
				switch strings.ToUpper(os.Args[2]) {
				case "--HELP", "-H", "HELP":
					fmt.Println("")
					fmt.Println("automate subcommand")
					fmt.Println(overwrites)
					fmt.Println("")
				}
			}
		}

		flag.BoolVar(&askExit, "x", false, "if false (default) menu kill phrase exits menu system without confirmation")

		automateCmd := flag.NewFlagSet("automate", flag.ExitOnError)

		basisFilePtr := automateCmd.String("bf", "", "Path to 29basis file")
		outputPathPtr := automateCmd.String("out", "", "Path to output folder")
		fromBasisNumPtr := automateCmd.String("bfrom", "0", "from Basis Number (0 based)")
		toBasisNumPtr := automateCmd.String("bto", "0", "to Basis Number")
		//overwritePtr := automateCmd.Bool("ask", false, "if false [default] asks whether to overwrite output files")
		filterPtr := automateCmd.Int("filter", 1, "Filter by 0-5 (none,6,L5,R5,L5R5,Q)\ndefault is 1 (sextuplets)")

		flag.Parse()

		switch os.Args[1] {

		//For every subcommand, we parse its own flags and have access to trailing positional arguments.

		case "automate":
			automateCmd.Parse(os.Args[2:])
			fmt.Println("subcommand 'automate'")

			auto := jup.GetNewAutomationStructure()
			auto.BasisFile = *basisFilePtr
			auto.OutputPath = *outputPathPtr
			auto.FromBasisNum = *fromBasisNumPtr
			auto.ToBasisNum = *toBasisNumPtr
			//auto.Overwrite = *overwritePtr
			auto.Filter = *filterPtr
			exitCode := jup.GeneratePrimeTupletsAutomated(auto)

			fmt.Println("Finished.")
			return exitCode
		}

	}

	return -1

}

func printLTE29s(p *jup.PrimeLTE29) {
	p.ShowDetails(true)
}

func printGTE31s(p *jup.PrimeGTE31) {
	p.ShowDetails(true)
}

package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"

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
)

const (
	muMyMainID = iota
	muMymenu29
	muMymenu29Details
	muMymenu31
	muMymenuCalcs
	muMymenuGen
)

var exitImmediately *bool

func main() {

	setFlags()
	if !jup.IsConfigured() {
		jup.Configure()
	}
	initMenus()

	switch *exitImmediately {
	case true:
		menuMain.Start() //kill phrase exits immediately
	case false:
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
			var (
				input       string
				wasCanceled bool
			)
			basis := big.NewInt(0)
			tNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInputInteger("Enter From TNumber:", "28", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, tNum)
			jup.TNumToBasisNum(tNum, basis)
			fmt.Println(basis)
		})
		menu.AddMenuEntry("bt", "Basis number to from/to TNumber (Tnum 1 based, basis 0 based)", func() {
			var (
				input       string
				wasCanceled bool
			)
			basis := big.NewInt(0)
			fromNum := big.NewInt(0)
			toNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInputInteger("Enter Basis Number:", "0", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, basis)
			jup.BasisToTNumRange(basis, fromNum, toNum)
			fmt.Println(fromNum, " : ", toNum)
		})
		menu.AddMenuEntry("it", "Integer to TNumber", func() {
			var (
				input       string
				wasCanceled bool
			)
			i := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInputInteger("Enter Integer:", "25", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, i)
			tNum := jup.IntToTNum(i)
			fmt.Println(fmt.Sprintf("Integer %v is in TNumber %v", i, tNum))
		})
		menu.AddMenuEntry("ti", "TNumber to Integer, and its range", func() {
			var (
				input       string
				wasCanceled bool
			)
			i := big.NewInt(0)
			end := big.NewInt(0)
			tNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInputInteger("Enter TNumber:", "1", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, tNum)
			i = jup.TNumToInt(tNum)
			big29 := big.NewInt(29)
			end.Add(i, big29)

			fmt.Println(fmt.Sprintf("TNumber %v starts at Integer %v, and ends at %v", tNum, i, end))
		})
		menu.AddMenuEntry("N", "Get N from a TNumber", func() {
			var (
				input       string
				wasCanceled bool
			)
			n := big.NewInt(0)
			//fromNum := big.NewInt(0)
			tNum := big.NewInt(0)
			if input, wasCanceled = jup.GetUserInputInteger("Enter TNumber Number:", "28", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, tNum)
			p31 := jup.NewPrimeGTE31(big.NewInt(31))
			_ = p31

			jup.GetNfromTNum(tNum, p31, n)
			fmt.Println(fmt.Sprintf("TNum: %v, n: %v", tNum, n))
		})
		menu.AddMenuEntry("H", "Human Readable from raw data (Give a TNum and Effect ID [0, 1, 2, or 3])", func() {
			var (
				input       string
				wasCanceled bool
			)
			tNum := big.NewInt(0)
			effect := 0
			str1, str2 := "n/a", "n/a"

			if input, wasCanceled = jup.GetUserInputInteger("Enter TNumber Number:", "28", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, tNum)

			if input, wasCanceled = jup.GetUserInputInteger("Enter Effect ID [0, 1, 2, or 3]:", "0", "x"); wasCanceled {
				return
			}
			fmt.Sscan(input, effect)
			jup.HumanReadable(tNum, &effect, &str1, &str2, os.Stdout)
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

func setFlags() {
	exitImmediately = flag.Bool("x", true, "if true (default) menu kill phrase exits menu system without confirmation")
	flag.Parse()
}

func printLTE29s(p *jup.PrimeLTE29) {
	p.ShowDetails(true)
}

func printGTE31s(p *jup.PrimeGTE31) {
	p.ShowDetails(true)
}

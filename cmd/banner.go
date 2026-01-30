package cmd

import "fmt"

func PrintBanner() {
	fmt.Println(
		"   ____ ___  ____ ___  ______ _\n" +
			"  / __ `__ \\/ __ `/ / / / __ `/\n" +
			" / / / / / / /_/ / /_/ / /_/ / \n" +
			"/_/ /_/ /_/\\__,_/\\__, /\\__,_/  \n" +
			"                /____/         \n\n" +
			"Maya CLI - Modular Command Line Tool (" + Version + ")\n")
}

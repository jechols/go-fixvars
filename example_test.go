package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Example_goFix() {
	var tmpDir, err = os.MkdirTemp("", "gofix-vars-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	var toolPath = filepath.Join(tmpDir, "fixvars")

	var buildCmd = exec.Command("go", "build", "-o", toolPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to build tool: %v\nOutput: %s\n", err, string(out))
		return
	}

	var projDir = filepath.Join(tmpDir, "testproj")
	if err := os.Mkdir(projDir, 0755); err != nil {
		panic(err)
	}

	var modCmd = exec.Command("go", "mod", "init", "testproj")
	modCmd.Dir = projDir
	if out, err := modCmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to init module: %v\nOutput: %s\n", err, string(out))
		return
	}

	var srcFile = filepath.Join(projDir, "main.go")
	srcCode, err := os.ReadFile("testdata/example_main.go")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(srcFile, srcCode, 0644); err != nil {
		panic(err)
	}

	var fixCmd = exec.Command("go", "fix", "-diff", "-fixtool="+toolPath, "./...")
	fixCmd.Dir = projDir
	var out, _ = fixCmd.CombinedOutput()

	var outStr = strings.TrimSpace(string(out))
	var lines = strings.Split(outStr, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(line, "---") {
			fmt.Println("--- main.go (old)")
		} else if strings.HasPrefix(line, "+++") {
			fmt.Println("+++ main.go (new)")
		} else {
			fmt.Println(line)
		}
	}

	// Output:
	// --- main.go (old)
	// +++ main.go (new)
	// @@ -6,10 +6,11 @@
	//  )
	//
	//  func main() {
	// -	x := 1
	// +	var x = 1
	//  	var y = x + 1
	//
	// -	if info, err := os.Stat("no"); err == nil {
	// +	var info, err = os.Stat("no")
	// +	if err == nil {
	//  		fmt.Printf("Info: %#v\n", info)
	//  	} else {
	//  		fmt.Printf("y is %d\n", y)
}

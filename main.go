package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	dirFmt := flag.String("dir-fmt", "2006/01", "formating directory structure inside the backup root with go time format.\nmultiple depth of dirs are possible.")
	fileFmt := flag.String("file-fmt", "20060102_150405.db", "formatting backup filename with go time format.")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %v [-flags...] [db] [backup_root]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		return
	}
	db := args[0]
	_, err := os.Stat(db)
	if err != nil {
		log.Fatalf("db not found: %v", db)
	}
	backupRoot := args[1]
	now := time.Now()
	dirs := now.Format(*dirFmt)
	file := now.Format(*fileFmt)
	dest := filepath.Clean(filepath.Join(backupRoot, dirs, file))
	destDir := filepath.Dir(dest)
	// It is sane that only a user who performs the backup can see the content.
	err = os.MkdirAll(destDir, 0700)
	if err != nil {
		log.Fatalf("cannot create backup directory: %v\n", destDir)
	}
	cmd := exec.Command("sqlite3", db, fmt.Sprintf(".backup '%v'", dest))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("backup failed: %v", err)
	}
}

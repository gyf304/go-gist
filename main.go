package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

func init() {
	args.GithubURL = "https://github.com"
	args.ClientID = "0c01d146f93dfe836694"
	usr, err := user.Current()
	if err == nil {
		args.TokenPath = filepath.Join(usr.HomeDir, ".config", "go-gist", "token.txt")
	}
}

var args struct {
	ClientID    string   `arg:"--client-id,env:GITHUB_CLIENT_ID" help:"specify a custom GitHub OAuth client ID"`
	GithubURL   string   `arg:"--github-url,env:GITHUB_URL" help:"specify a custom GitHub URL"`
	AccessToken string   `arg:"--access-token,env:GITHUB_ACCESS_TOKEN" help:"specify a custom GitHub access token"`
	Description string   `arg:"-d" help:"specify a description for gist"`
	FileNames   []string `arg:"-f,separate" help:"specify filenames for gist"`
	Files       []string `arg:"positional"`
	Base64      bool     `arg:"--base64" help:"use base64, useful for binary files"`
	Private     bool     `arg:"-p" help:"make gist private"`
	Read        bool     `arg:"-r" help:"read a gist"`
	Output      string   `arg:"-o" help:"specify a output file, only valid if [-r,--read] is specified"`
	Login       bool     `arg:"-l" help:"force a login"`
	Timeout     int      `arg:"-t" help:"specify a timeout for login, 0 or negative values disables timeout"`
	TokenPath   string   `arg:"--token-path" help:"specify a path for storing / reading GitHub access token"`
}

func doPost() {
	files := make(map[string][]byte)
	if len(args.Files) == 0 {
		fmt.Fprintf(os.Stderr, "No input files\n")
		os.Exit(1)
		return
	}
	if len(args.FileNames) == 0 {
		for _, v := range args.Files {
			_, filename := filepath.Split(v)
			args.FileNames = append(args.FileNames, filename)
		}
	}
	if len(args.FileNames) != len(args.Files) {
		fmt.Fprintf(os.Stderr, "%d filenames provided but there are %d files\n", len(args.FileNames), len(args.Files))
		os.Exit(1)
		return
	}
	for i, filepath := range args.Files {
		filename := args.FileNames[i]
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
			return
		}
		files[filename] = file
	}
	err := authenticate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot authenticate: %s\n", err.Error())
		os.Exit(1)
		return
	}
	postResp, err := post(args.Description, files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot upload: %s\n", err.Error())
		os.Exit(1)
		return
	}
	fmt.Println(postResp.ID)
}

func doRead() {
	var err error
	var outFile io.WriteCloser
	outFile = os.Stdout
	if args.Output != "" {
		outFile, err = os.OpenFile(args.Output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file: %s\n", err.Error())
			os.Exit(1)
			return
		}
	}
	if len(args.Files) == 0 {
		fmt.Fprintf(os.Stderr, "Need to provide gist\n")
		os.Exit(1)
		return
	}
	if len(args.Files) != 1 {
		fmt.Fprintf(os.Stderr, "Only one gist needed\n")
		os.Exit(1)
		return
	}
	if len(args.FileNames) > 1 {
		fmt.Fprintf(os.Stderr, "Can only provide one filename\n")
		os.Exit(1)
		return
	}
	files, err := get(args.Files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read gist %s: %s\n", args.Files[0], err.Error())
		os.Exit(1)
		return
	}
	if len(files) == 0 {
		return
	}
	if len(args.FileNames) == 0 {
		var file []byte
		for _, f := range files {
			file = f
			break
		}
		outFile.Write(file)
	} else {
		filename := args.FileNames[0]
		file := files[filename]
		if file == nil {
			fmt.Fprintf(os.Stderr, "File %s not found in gist %s\n", filename, args.Files[0])
			os.Exit(1)
			return
		}
		outFile.Write(file)
	}
}

func main() {
	arg.MustParse(&args)
	loadToken()
	// compile data
	if args.Login {
		args.AccessToken = ""
		err := authenticate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot login: %s\n", err.Error())
			os.Exit(1)
			return
		}
	}
	if args.Read {
		doRead()
	} else {
		doPost()
	}
}

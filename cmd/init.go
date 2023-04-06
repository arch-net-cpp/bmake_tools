/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"errors"
	"github.com/arch-net-cpp/bmake_tools/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

var useGitee bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new project with the given project name",
	Long:  `Create a new project with the given project name`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a project name")
		}
		return utils.ValidateDirectoryName(args[0])
	},
	Run: workFlow,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolVar(&useGitee, "use_gitee", false, "Use Gitee instead of GitHub as remote repository")
}

const (
	githubBaseURL = "https://github.com/arch-net-cpp/"
	giteeBaseURL  = "https://github.com/arch-net-cpp/"
)

func workFlow(cmd *cobra.Command, args []string) {
	// create new directory
	dirName := args[0]
	archnetDir := dirName + "/" + "arch_net"
	bmakeDir := dirName + "/" + "bmake"
	cpp3rdlibDir := dirName + "/" + "cpp3rdlib"

	err := os.Mkdir(dirName, 0777)
	if err != nil {
		utils.ErrorFmtPrintf("create directory error: %v", err)
	}
	err = os.MkdirAll(archnetDir, 0777)
	if err != nil {
		utils.ErrorFmtPrintf("create directory error: %v", err)
	}
	err = os.Mkdir(bmakeDir, 0777)
	if err != nil {
		utils.ErrorFmtPrintf("create directory error: %v", err)
	}
	err = os.Mkdir(cpp3rdlibDir, 0777)
	if err != nil {
		utils.ErrorFmtPrintf("create directory error: %v", err)
	}
	// git clone option

	token := os.Getenv("GITHUB_TOKEN")
	auth := &http.BasicAuth{
		Username: "abc123", // yes, this can be anything except an empty string
		Password: token,
	}

	// git clone arch_net
	_, err = git.PlainClone(archnetDir, false, &git.CloneOptions{
		URL:   genGitURL("arch_net"),
		Depth: 1,
		Auth:  auth,
	})
	if err != nil {
		utils.ErrorFmtPrintf("git clone arch_net from: %v with error: %v ", genGitURL("arch_net"), err)
	}
	utils.DefaultFmtPrintf("clone arch_net into %v successfully", dirName)

	// git clone bmake
	_, err = git.PlainClone(bmakeDir, false, &git.CloneOptions{
		URL:   genGitURL("bmake"),
		Depth: 1,
		Auth:  auth,
	})
	if err != nil {
		utils.ErrorFmtPrintf("git clone bmake from: %v with error: %v ", genGitURL("bmake"), err)
	}
	utils.DefaultFmtPrintf("clone bmake into %v successfully", dirName)

	// git clone cpp3rdlib
	_, err = git.PlainClone(cpp3rdlibDir, false, &git.CloneOptions{
		URL:   genGitURL("cpp3rdlib"),
		Depth: 1,
		Auth:  auth,
	})
	if err != nil {
		utils.ErrorFmtPrintf("git clone cpp3rdlib from: %v with error: %v ", genGitURL("cpp3rdlib"), err)
	}
	utils.DefaultFmtPrintf("clone cpp3rdlib into %v successfully", dirName)

	// creat template CMakeLists.txt
	file, err := os.Create(dirName + "/" + "CMakeLists.txt")
	if err != nil {
		utils.ErrorFmtPrintf("create CMakeLists.txt with error: %v ", err)
	}
	defer file.Close()

	// 写入内容到文件
	// 创建写入器对象
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	projectVariable := map[string]interface{}{
		"projectName": dirName,
	}
	err = cmakeTemplate.Execute(writer, projectVariable)
	if err != nil {
		utils.ErrorFmtPrintf("generate CMakeLists.txt with error: %v ", err)
	}

	utils.DefaultFmtPrintf("bmake init new project [%v] successfully ! ! !", dirName)
}

var cmakeTemplate = template.Must(template.New("CMakeLists.txt").
	Parse(`cmake_minimum_required(VERSION 3.21)
project({{.projectName}})
set(CMAKE_CXX_STANDARD 14)

# for debug uasge
include(CMakePrintHelpers)
cmake_print_variables(PROJECT_SOURCE_DIR)
cmake_print_variables(CMAKE_MODULE_PATH)

# include bmake
set(BMAKE_ROOT_DIR ${CMAKE_SOURCE_DIR})
list(APPEND CMAKE_MODULE_PATH ${BMAKE_ROOT_DIR}/bmake)
cmake_print_variables(CMAKE_MODULE_PATH)
include(bmake)
# bmake will automatically load all dependent libraries under [cpp3rdlib]

# include arch_net
add_subdirectory(arch_net)

# custom Customized cmake code


`))

func genGitURL(name string) string {
	if useGitee {
		return giteeBaseURL + name
	}
	return githubBaseURL + name
}

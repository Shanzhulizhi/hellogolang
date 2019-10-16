package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"errors"
	flag "github.com/spf13/pflag" 
)

type Args struct {
	start int   //start page
	end int		//end page
	length int  //lines
	readType string  //-l or -f
	dst string  //destination
	input string  //input file
}
var err error

func errorHandler(err error) { //func to handle error
	if err != nil {
		fmt.Println("\nError!\n",err)
		os.Exit(1)
	}
}

func execute(inargs Args) {
	fin := os.Stdin
	fout := os.Stdout
	curLine := 0  //current line
	curPage := 1  //current page
	var inpipe io.WriteCloser//pipe

	//judge the input type
	if inargs.input != "" {
		fin,err = os.Open(inargs.input)
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error:Can't find input file \"%s\"!\n",inargs.input)
			os.Exit(1)
		}
		defer fin.Close()  //延迟结束
	}

	//输出
	if inargs.dst != "" {
		ins := exec.Command("grep","-nf","keyword")
		inpipe,err = ins.StdinPipe()
		errorHandler(err)
		defer inpipe.Close()  //最后执行
		ins.Stdout = fout
		ins.Start() 
	}
	
	if inargs.readType == "l" { //按行读取
		line := bufio.NewScanner(fin)
		for line.Scan() {
			if curPage >= inargs.start && curPage <= inargs.end {
				fout.Write([]byte(line.Text() + "\n"))
				if inargs.dst != "" {
					inpipe.Write([]byte(line.Text() + "\n"))
				}
			}
			curLine++ //下一行
			if curLine == inargs.length {
				curPage++
				curLine = 0  //新的一页开始
			}
		}
	}else {  //用换行符'\f'分页
		readText := bufio.NewReader(fin)
		for {
			page, ferr := readText.ReadString('\f')
			if ferr != nil || ferr == io.EOF {
				if ferr == io.EOF {
					if curPage >= inargs.start && curPage <= inargs.end {
						fmt.Fprintf(fout,"%s",page)
					}
				}
				break
			}
			page = strings.Replace(page,"\f","",-1)
			curPage++
			if curPage >= inargs.start && curPage <= inargs.end {
				fmt.Fprintf(fout,"%s",page)
			}
		}
	}

	if curPage < inargs.end {
		fmt.Fprintf(os.Stderr, "./selpg: end (%d) greater than total pages (%d), less output than expected\n", inargs.end, curPage)
	}
}
func main() {
	inargs := new(Args)
	flag.IntVar(&inargs.start,"s",0,"the start page")    //开始页码默认为0
	flag.IntVar(&inargs.end,"e",0,"the end page")       //结束页码默认为0
	flag.IntVar(&inargs.length,"l",72,"the length of the page") //每页默认72行
	flag.StringVar(&inargs.dst,"d","","the destiny of printing")//输出位置默认为空

	//读取方式
	isf := flag.Bool("f",false,"")  //默认不是-f
	flag.Parse() 
	
	//如果是-f,则行数取-1，否则为-l
	if *isf {
		inargs.readType = "f"
		inargs.length = -1
	}else {
		inargs.readType = "l"
	}

	//如果是文件输入，则设置文件名
	inargs.input = ""
	if flag.NArg() == 1 {
		inargs.input = flag.Arg(0)
	}

	//检查参数合法性
	if nargs := flag.NArg();nargs != 1 && flag.NArg() !=0 {
		err = errors.New("No enough arguments\n")
	}
	if inargs.start > inargs.end || inargs.start < 1 {
		err = errors.New("Wrong start or end page\n")
	}
	if inargs.readType == "f" && inargs.length != -1 {
		err = errors.New("-l and -f both come up\n")
	}
	errorHandler(err) //处理错误
	execute(*inargs) //执行指令
}
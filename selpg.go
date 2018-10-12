package main


import (
	"fmt"
	flag "github.com/spf13/pflag"
	"io"
	//"bytes"
	"os/exec"
	"os"
	//"log"
	"bufio"
	//"errors"
)





/*================================= globals =======================*/

/* program name, for error messages */
var progname string


const (
	BUFSIZ = 1024
	/* INBUFSIZ is size of array inbuf */
	INBUFSIZ = 16 * 1024
)

var inbuf []byte = make([]byte,INBUFSIZ)

/*================================= types =========================*/
type selpg_args struct {
	Start_page int
	End_page int
	In_filename string
	Page_len int
	Page_type int

	Print_dest string
}

var Start_page int
var End_page int
var In_filename string
var Page_len int = 72
var Page_type bool
var Print_dest string

func main() {
	/* save name by which program is invoked, for error messages */
	progname = os.Args[0]

	processArgs(os.Args)
	processInput()

}



func processArgs(args []string) {

	flag.IntVarP(&Start_page,"startpage","s",1,"start page")
	flag.IntVarP(&End_page,"endpage","e",1,"end page")

	flag.IntVarP(&Page_len,"pagelen","l",72,"page length (default 72), only valid if '-f' not set")
	flag.BoolVarP(&Page_type,"pagetype","f",false,"specify the page type that pages are separated by '\f'")
	flag.StringVarP(&Print_dest,"printdest","d","","printer destination, the specified range of pages will be printed to the destination")

	if len(args) < 3 {
		fmt.Fprintf(os.Stderr,"%s: not enough arguments\n", progname)
		usage()
		os.Exit(1)
	}

	//todo debug

	

	flag.Parse()
    pArgs := flag.Args()
	if len(pArgs) >= 1 {
		In_filename = pArgs[0]
	} else {
		In_filename = ""
	}
}


func processInput() {
	//var file *os.File
	var scanner *bufio.Scanner
	if In_filename != "" {
		file, err := os.OpenFile(In_filename,os.O_RDONLY, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr,"%s: failed to open file\n",err)
			os.Exit(1)
		}
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}
	scanner.Buffer(inbuf,INBUFSIZ)

	var  out io.Writer
	var in io.Reader
	ch := make(chan int)
	if Print_dest != "" {
		in, out = io.Pipe()
		go print(ch,in,Print_dest)
	} else {
		out = os.Stdout
	}
	writer := bufio.NewWriter(out)
	//	read the selected pages
	var lineCtr int
	var pageCtr int
	//page type is default(by lengths)
	if Page_type == false {
		lineCtr = 0
		pageCtr = 1
		for scanner.Scan() {
			lineCtr++
			line := scanner.Text()
			if lineCtr > Page_len {
				pageCtr++
				lineCtr = 1
			}

			if pageCtr >= Start_page && pageCtr <= End_page {
				writer.WriteString(line)
				writer.WriteString("\n")
			}
			
		}

		
	} else {
		pageCtr = 0
		split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF {
				advance = 0
				token = nil
				err = nil
				return
			}
			for index,value := range data {
				if value == byte('\f') {
					advance = index + 1
					token = data[:index]
					err = nil
					return 
				}
			}
			
			advance = len(data)
			token = data
			err = nil 
			return
		}
		scanner.Split(split)
		for scanner.Scan() {
			pageCtr++
			if pageCtr >= Start_page && pageCtr <= End_page {
				writer.WriteString(scanner.Text())
			}		
		}
	}

	
	writer.Flush()
	if Print_dest != "" {
		<-ch
	}
	doneMsg := progname + ": done\n" 
	writer.WriteString(doneMsg)

	if Start_page > pageCtr {
		fmt.Fprintf(os.Stderr, "%s: start page (%d) greater than total page (%d), no output written\n", progname, Start_page, pageCtr)
	} else if pageCtr < End_page {
		fmt.Fprintf(os.Stderr, 
			"%s: end page (%d) greater than total pages (%d), less output than expected\n", progname, End_page, pageCtr)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr,
		  	"%s: system error [%s] occured on scanner\n",
		  	progname, err )
		os.Exit(14)
	}
	
	writer.Flush()	

}


func print(ch chan int,in io.Reader, dest string) {
	cmd := exec.Command("lp","-d",dest)
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr,"%s\n",err)
		os.Exit(1)
	}
	ch<- 0
} 


func usage() {
	fmt.Fprintf(os.Stderr,"\nUSAGE: %s -s startpage -e endpage [OPTION...] [file]\n",progname)
	flag.PrintDefaults()
}


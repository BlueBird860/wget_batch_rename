package main
 
import (
    "flag"
    "fmt"
	"os"
	"os/exec"
	"io"
	"log"
	"bufio"
	"bytes"
    "sync"  
)

func isFlagPassed(name string) bool {
    found := false
    flag.Visit(func(f *flag.Flag) {
        if f.Name == name {
            found = true
        }
    })
    return found
}

func main() {
	
    var workfile=flag.String("i", "", "【必须指定】下载列表文件，规则：每行一个url")
	var extension=flag.String("ext", "" ,"指定下载文件后缀名")
	var start_num=flag.Int("index",0,"下载文件开始序号(默认0)")
    var prefix=flag.String("prefix","file","文件默认前缀")
    var wget_header=flag.String("header","","set http header (wget --header)")
    var threads_num=flag.Int("threads",8,"运行wget的线程数")
    var override=flag.Bool("override",false,"已存在文件是否覆盖（默认不覆盖）")
    flag.Parse()

    
	if (!(isFlagPassed("i")&&isFlagPassed("ext"))) {
		println("请输入-h查询用法")
		os.Exit(0)

	}

	file, err := os.Open(*workfile)
    if err != nil {
		log.Fatal("不存在文件："+*workfile)
	}
    

    var tasks []string  
	var buffer = bufio.NewReader(file)
	for {
        line, _, c := buffer.ReadLine()
        if c == io.EOF {
            break
        }

        tasks = append(tasks, string(line))

    }
    file.Close()

  
    url_index:=0
    var mutex sync.Mutex
    var wait sync.WaitGroup
    for i:=0;i<*threads_num;i++{
        go  exec_wget_cmd_mul(wget_header,tasks,&url_index, prefix, start_num, extension,*override,&mutex,&wait,i)
        wait.Add(1)
    }
	wait.Wait()
}



func getNextFileName(pre *string,index *int, ext *string) string {
	var filename string
    filename=fmt.Sprintf("%s%03d.%s",*pre,*index,*ext)
	*index=*index + 1
	return filename
}


func exec_wget_cmd_mul(wget_header *string, urls []string, url_index *int, prefix *string, start_num *int, extension *string ,override bool,mutex *sync.Mutex,wait *sync.WaitGroup,id int) {
    var url,dfile string
    for {
        mutex.Lock()
        url=urls[*url_index]
        
        if *url_index<len(urls)-1 {
            *url_index++
        } else{
            mutex.Unlock()
            fmt.Printf("=======threads[%d] done!=======\n",id)
            wait.Done()
            break;
        }
        dfile=getNextFileName(prefix,start_num,extension)
        fmt.Printf("thread[%d]: %s -> %s\n",id,url,dfile)
        mutex.Unlock()

        
        if !override {
            _, err := os.Stat(dfile)
            if err == nil {
               fmt.Printf("%s is existence, ignored.\n",dfile)
                continue
            }
        }
        
        cmd_line:=fmt.Sprintf("\twget --header='%s' %s --output-document='%s'  --prefer-family=IPv4",*wget_header,url,dfile)
        println(cmd_line)
        cmd := exec.Command("wget",url,"--header",*wget_header, "--output-document", dfile,"--prefer-family","IPv4")
    
        var out bytes.Buffer
        var stderr bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &stderr
    
        err := cmd.Run()
        if err != nil {
    
             fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
            
            
        }
    }
}



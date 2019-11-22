package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)
var LeftTags = []string{"<div","<p","<span","<h1","<h2","<h3","<h4","<h5","<a","<ul","<li"}
var RightTags = []string{"</div","</p","</span","</h1","</h2","</h3","</h4","</h5","</a","</ul","</li"}
var TagMap = map[string]string{
	"</div":"<div",
	"</p":"<p",
	"</span":"<span",
	"</h1":"<h1",
	"</h2":"<h2",
	"</h3":"<h3",
	"</h4":"<h4",
	"</h5":"<h5",
	"</a":"<a",
	"</ul":"<ul",
	"</li":"<li",
}

var line = 1;
var cara = 1;
var tmp_char = "";
var tmp_tag = ""
var php_start = false
var out_tag  = ""
var errinfo []string
var originStack []string
var phpStack []string
var keyname = ""
var out_php = ""
var tmp_if []string
var last_php = ""
var phpStackTag = make(map[string][]string)
var errorinfo = ""
func main(){
	args := os.Args
	var content = readFile(args[1])
	for  i:=0;i<len(content);i++ {
		tmp_char = string(content[i])
		if tmp_tag != ""{
			if( tmp_char != " " && tmp_char != ">" ){
				tmp_tag += tmp_char
			}else{
				if len(phpStack) > 0{
					phpStackTag[phpStack[len(phpStack) - 1]],errorinfo = rangeAddTag(phpStackTag[phpStack[len(phpStack) - 1]])
				}else{
					originStack,errorinfo = rangeAddTag(originStack)
				}
				if errorinfo != "" {
					break
				}
			}
		}
		if php_start == true {   //每个phpStackTag[keyname] 再重新计算标签闭合，最后弹出时，最后一个一次总结，合给前面一个结果集
			keyname = ""
			if (tmp_char=="i" && string(content[i+1]) == "f"){ //if
				keyname = "if_"+strconv.Itoa(line)+"_"+strconv.Itoa(cara);
			}
			if(tmp_char=="f" && content[i-1:i+3] == " for"){ //for /foreach
				keyname = "for_"+strconv.Itoa(line)+"_"+strconv.Itoa(cara);

			}
			if (tmp_char=="e" && content[i:i+4] == "else"){ //if
				keyname = "else_"+strconv.Itoa(line)+"_"+strconv.Itoa(cara);
			}
			if keyname != ""{
				phpStack = append(phpStack,keyname)
				phpStackTag[keyname] = []string{}
			}
			if (tmp_char == "}" || tmp_char == "endif" || tmp_char== "endforeach"){ //结尾 else 弹出时，判断tmp if 如果两个标签不同，报错。
				if len(phpStack) > 0{
					out_php = phpStack[len(phpStack) - 1] //判断要弹出的元素
					if "if" == out_php[0:2]{  //弹出的元素是if 临时存储if 剩余标签
						tmp_if = phpStackTag[out_php]
					}
					if (len(out_php) >=4 && "else" == out_php[0:4] ){
						if reflect.DeepEqual(tmp_if,phpStackTag[out_php]) {
							phpStackTag[out_php] = nil
						}else{
							errinfo = append(errinfo,"if else tag not same,line:"+strconv.Itoa(line)+",position:"+strconv.Itoa(cara))
							break
						}
					}
					phpStack = phpStack[0:len(phpStack) - 1] // 弹出标签
					if len(phpStackTag[out_php]) > 0{
						if len(phpStack) > 0 {
							last_php = phpStack[len(phpStack)-1] //获取弹出元素后最后一个元素
							phpStackTag[last_php] = append(phpStackTag[last_php], phpStackTag[out_php]...)
						}else{
							originStack = append(originStack, phpStackTag[out_php]...)
						}
						delete(phpStackTag,out_php) //释放占用内存
					}
				}else{
					errinfo = append(errinfo,"php end no start,line:"+strconv.Itoa(line)+",position:"+strconv.Itoa(cara))
					break
				}


			}

		}

		if tmp_char == "<"{
			if string(content[i+1]) == "?"{
				php_start = true
			}else{
				tmp_tag = "<"
			}

		}
		if (tmp_char == "?" && string(content[i+1]) == ">") {
				php_start = false
		}

		if tmp_char == "\n" {
			line++
			cara = 0
		}else{
			cara++
		}
	}
	if len(originStack) > 0{
		var extra_tags = ""
		for _,value := range originStack{
			extra_tags = extra_tags+value+","
		}
		errinfo = append(errinfo,"end but extra:"+extra_tags)
	}
	if len(errinfo) > 0 {
		fmt.Println(errinfo)
	}else{
		fmt.Println(0)
	}


}
func inArray(str string,arrays []string) (exist_key bool){
	exist_key = false
	for _,value := range arrays{
		if str == value{
			exist_key = true
			break
		}
	}
	return
}
func readFile(filename string) (ct_str string){
	f,_ := os.Open(filename)
	defer f.Close()
	ct ,_:= ioutil.ReadAll(f)
	ct_str = string(ct)
	return
}
func rangeAddTag(pStack [] string)  ([]string ,string){
	if inArray(strings.ToLower(tmp_tag),LeftTags){
		pStack = append(pStack, tmp_tag)
	}
	if inArray(strings.ToLower(tmp_tag),RightTags) {
		if len(pStack) > 0{
			out_tag = pStack[len(pStack) - 1]
			if TagMap[tmp_tag] == out_tag{
				pStack = pStack[0:len(pStack) - 1]
			}else{
				errinfo = append(errinfo,tmp_tag+" start tag error,line:"+strconv.Itoa(line)+",position:"+strconv.Itoa(cara))
				return errinfo,"error"
			}
		}else{
			errinfo = append(errinfo,tmp_tag+" no start,line:"+strconv.Itoa(line)+",position:"+strconv.Itoa(cara))
			return errinfo,"error"
		}

	}
	tmp_tag = ""
	return pStack,""
}

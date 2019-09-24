package udwTest

import "fmt"

func Ok(isOk bool,msgObj ...interface{}){
	if isOk==false{
		if len(msgObj)==0{
			panic("fail")
		}else{
			panic(fmt.Sprintln(msgObj...))
		}
	}
}
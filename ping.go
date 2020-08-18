package main

import (
	"./util"
	"github.com/huichen/sego"
)

var segmenter sego.Segmenter

func init()  {
	segmenter.LoadDictionary("/Users/hjzhou/Gocode/src/github.com/huichen/sego/data/dictionary.txt")
}

type Ping struct {
	Id 			string
	OnePing		string
	Comments 	[4]string
	Tags 		[3][]string
}

var tags = [3][]string{
	{"美国","特朗普","台湾","民进党","蔡英文","川普","印度","莫迪","越南","东南亚","大选","华为"},
	{"人民币","汇率","华裔","北美","非洲","做题家","日本","韩国","安倍晋三","安倍","俄罗斯","普京","美帝"},
	{"湖北","武汉","新疆","英国","欧洲","德国","法国","意大利","上海","不列颠","欧盟"},
}

func ExtractKeyWords(ping string) (keyWords [3][]string){
	text := []byte(ping)
	segments := segmenter.Segment(text)
	//fmt.Println(sego.SegmentsToString(segments, false))
	nums := sego.SegmentsToString(segments,false)
	nums = util.RemoveRepeatElement(nums)
	for i :=0;i<3 ; i++ {
		for j := 0;j< len(tags[i]); j++ {
			for k := 0; k< len(nums); k++ {
				if tags[i][j]==nums[k] {
					keyWords[i] = append(keyWords[i],nums[k])
				}
			}
		}
	}
	return
}


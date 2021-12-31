package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"os"
	"strconv"
)

var baseUrl = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2021/"
var tableName = "areas_standard"
var f *os.File
var id = 0

func GetSQL() {
	//地址
	var url string
	url = baseUrl + "index.html"
	//打开area.sql文件，准备写入sql语句
	f, _ = os.Create("area.sql")

	FindProvince(url)

	fmt.Println("数据已写入 area.sql 中，共: " + strconv.Itoa(id) + " 条数据")
}

func FindProvince(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".provincetr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		parentId := 0
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			province := td.Find("a")
			provinceName := UseNewEncoder(province.Text(), "gbk", "utf-8")
			fmt.Println("省份：" + provinceName)
			parentId = 0
			id = id + 1
			io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+provinceName+"',1,"+strconv.Itoa(parentId)+");\r\n")
			hrefUrl, res := province.Attr("href")
			if res {
				FindCity(baseUrl+hrefUrl, id)
			}
		})
	})
}

func FindCity(url string, parentId int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".citytr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		cityId := tr.Find("td").First().Find("a").Text()
		cityName := tr.Find("td").Last().Find("a").Text()
		cityName = UseNewEncoder(cityName, "gbk", "utf-8")
		fmt.Println("城市：" + cityId + "  ==> " + cityName)
		id = id + 1
		io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+cityName+"',2,"+strconv.Itoa(parentId)+");\r\n")
		hrefUrl, res := tr.Find("td").First().Find("a").Attr("href")
		if res {
			FindCounty(baseUrl+hrefUrl, id)
		}
	})
}

func FindCounty(url string, parentId int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".countytr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		countyId := tr.Find("td").First().Text()
		countyName := tr.Find("td").Last().Text()
		if countyId == "" {
			countyId = tr.Find("td").First().Find("a").Text()
			countyName = tr.Find("td").Last().Find("a").Text()
		}
		countyName = UseNewEncoder(countyName, "gbk", "utf-8")
		fmt.Println("区县：" + countyId + "  ==> " + countyName)
		id++
		io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+countyName+"',3,"+strconv.Itoa(parentId)+");\r\n")
	})
}

func UseNewEncoder(src string, oldEncoder string, newEncoder string) string {
	srcDecoder := mahonia.NewDecoder(oldEncoder)
	desDecoder := mahonia.NewDecoder(newEncoder)
	resStr := srcDecoder.ConvertString(src)
	_, resBytes, _ := desDecoder.Translate([]byte(resStr), true)
	return string(resBytes)
}

func main() {
	GetSQL()
}

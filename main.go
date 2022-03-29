package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

var baseUrl = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2021/"
var tableName = "areas_standard"
var f *os.File
var id = 43170

var appointProvince = "宁夏回族自治区"
var runProvince = false

var GarbledCode = false

var currentProvince = ""
var currentCity = ""
var currentCounty = ""
var currentTown = ""

func GetSQL() {
	//地址
	FindProvince(baseUrl)

	fmt.Println("数据已写入 area.sql 中，共: " + strconv.Itoa(id) + " 条数据")
}

// 查找省份
func FindProvince(url string) {
	url = url + "index.html"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".provincetr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		parentId := 0
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			//打开area.sql文件，准备写入sql语句
			province := td.Find("a")
			currentProvince = UseNewEncoder(province.Text(), "gbk", "utf-8")
			if appointProvince != "" {
				if appointProvince == currentProvince {
					runProvince = true
				}
			} else {
				runProvince = true
			}
			if runProvince {
				f, _ = os.Create("area_" + currentProvince + ".sql")
				fmt.Println("省份：" + currentProvince + string(id+1))
				parentId = 0
				id = id + 1
				io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+currentProvince+"',1,"+strconv.Itoa(parentId)+");\r\n")
				hrefUrl, res := province.Attr("href")
				if res {
					handleUrl := handleUrl(url)
					FindCity(handleUrl+hrefUrl, id)
					fmt.Println("=== 等待下一个省份开始 ===")
					//time.Sleep(time.Second * 20)
				}
			}
		})
	})
}

// 查找城市
func FindCity(url string, parentId int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".citytr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		cityId := tr.Find("td").First().Find("a").Text()
		currentCity = tr.Find("td").Last().Find("a").Text()
		currentCity = UseNewEncoder(currentCity, "gbk", "utf-8")
		fmt.Println("城市：" + currentProvince + cityId + "  ==> " + currentCity)
		id = id + 1
		io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+currentCity+"',2,"+strconv.Itoa(parentId)+");\r\n")
		//time.Sleep(time.Second * 5)
		hrefUrl, res := tr.Find("td").First().Find("a").Attr("href")
		if res {
			handleUrl := handleUrl(url)
			FindCounty(handleUrl+hrefUrl, id)
		}
	})
}

// 查找区县
func FindCounty(url string, parentId int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".countytr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		countyId := tr.Find("td").First().Text()
		currentCounty = tr.Find("td").Last().Text()
		if countyId == "" {
			countyId = tr.Find("td").First().Find("a").Text()
			currentCounty = tr.Find("td").Last().Find("a").Text()
		}
		currentCounty = UseNewEncoder(currentCounty, "gbk", "utf-8")
		fmt.Println("区县：" + currentProvince + countyId + "  ==> " + currentCounty)
		id++
		io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+currentCounty+"',3,"+strconv.Itoa(parentId)+");\r\n")
		hrefUrl, res := tr.Find("td").First().Find("a").Attr("href")
		if res {
			handleUrl := handleUrl(url)
			FindTown(handleUrl+hrefUrl, id)
		}
	})

	if doc.Find(".countytr").Text() == "" {
		FindTown(url, parentId)
	}
}

// 查找镇/街道
func FindTown(url string, parentId int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".towntr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		townId := tr.Find("td").First().Text()
		currentTown = tr.Find("td").Last().Text()
		if townId == "" {
			townId = tr.Find("td").First().Find("a").Text()
			currentTown = tr.Find("td").Last().Find("a").Text()
		}
		currentTown = UseNewEncoder(currentTown, "gbk", "utf-8")
		fmt.Println("镇/街道：" + currentProvince + townId + "  ==> " + currentTown)
		id++
		io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+currentTown+"',4,"+strconv.Itoa(parentId)+");\r\n")
		//hrefUrl, res := tr.Find("td").First().Find("a").Attr("href")
		//if res {
		//	handleUrl := handleUrl(url)
		//	FindVillage(handleUrl+hrefUrl, id)
		//}
	})
}

// 查找社区/村
func FindVillage(url string, parentId int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".villagetr").Each(func(i int, tr *goquery.Selection) {
		//根据页面特点，有加粗<b>标签的是省级数据
		villageId := tr.Find("td").First().Text()
		villageName := tr.Find("td").Last().Text()
		if villageId == "" {
			villageId = tr.Find("td").First().Find("a").Text()
			villageName = tr.Find("td").Last().Find("a").Text()
		}
		villageName = UseNewEncoder(villageName, "gbk", "utf-8")
		fmt.Println("社区/村：" + currentProvince + villageId + "  ==> " + villageName)
		id++
		io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`level`,`parent_id`) values("+strconv.Itoa(id)+",'"+villageName+"',5,"+strconv.Itoa(parentId)+");\r\n")
	})
}

func handleUrl(url string) string {
	reg := regexp.MustCompile("[a-z0-9]+.html")
	url = reg.ReplaceAllString(url, "")
	fmt.Println("地址：" + url)
	return url
}

func UseNewEncoder(src string, oldEncoder string, newEncoder string) string {
	if !GarbledCode {
		return src
	}
	srcDecoder := mahonia.NewDecoder(oldEncoder)
	desDecoder := mahonia.NewDecoder(newEncoder)
	resStr := srcDecoder.ConvertString(src)
	fmt.Println("==" + resStr)
	_, resBytes, _ := desDecoder.Translate([]byte(resStr), true)
	return string(resBytes)
}

func main() {
	GarbledCode = false
	GetSQL()
}

package main

import (
	"carscraper/pkg/adsdb"
	"log"
)

type Test struct {
	tt string
}

type Res struct {
	result map[int]map[int]map[int]*Test
}

func newRes() *Res {
	r := make(map[int]map[int]map[int]*Test)
	return &Res{result: r}
}

func (r Res) addToMap(a int, b int, c int, t Test) {
	if r.result[a] == nil {
		m1 := make(map[int]*Test)
		m1[c] = &t

		m2 := make(map[int]map[int]*Test)
		m2[b] = m1

		r.result[a] = m2
		//m[a][b][c] = &Test{tt: "xxx"}
	} else {
		if r.result[a][b] == nil {
			m1 := make(map[int]*Test)
			m1[c] = &t

			m2 := make(map[int]map[int]*Test)
			m2[b] = m1

			r.result[a] = m2
		}
		r.result[a][b][c] = &t
	}
}

func main() {

	rr := newRes()

	rr.addToMap(1, 2, 3, Test{tt: "123"})
	rr.addToMap(1, 2, 4, Test{tt: "124"})
	rr.addToMap(1, 3, 1, Test{tt: "131"})
	rr.addToMap(1, 3, 2, Test{tt: "131"})

	rr.addToMap(2, 2, 3, Test{tt: "123"})
	rr.addToMap(2, 2, 4, Test{tt: "124"})
	rr.addToMap(2, 3, 1, Test{tt: "131"})
	rr.addToMap(2, 3, 2, Test{tt: "131"})

	log.Printf("%+v", rr)
	log.Println("DONE")
	return
	//m := make(map[int]map[int]map[int]*Test)
	//
	//if m[1] == nil {
	//	m1 := make(map[int]*Test)
	//	m1[3] = &Test{tt: "xxsds"}
	//
	//	m2 := make(map[int]map[int]*Test)
	//	m2[2] = m1
	//
	//	m[1] = m2
	//	m[1][2][3] = &Test{tt: "xxx"}
	//} else {
	//
	//}
	//
	//m[1][2][4] = &Test{tt: "yyyy"}
	//
	//log.Printf("%+v", m)

	//
	//market1 := adsdb.Market{
	//	Name: "autovit",
	//	PageURL:  "www.autovit.ro",
	//}
	//
	//market2 := adsdb.Market{
	//	Name: "mobile.de",
	//	PageURL:  "www.mobile.de",
	//}
	//
	//var mkts []adsdb.Market
	//mkts = append(mkts, market1)
	//mkts = append(mkts, market2)
	//
	criteria := adsdb.Criteria{
		Brand:        "mazda",
		CarModel:     "cx-5",
		YearFrom:     pOf(2019),
		YearTo:       pOf(2023),
		Fuel:         "diesel",
		KmFrom:       pOf(0),
		KmTo:         pOf(125000),
		AllowProcess: true,
		Markets:      nil,
		ScrapeLogs:   nil,
	}
	//
	//log.Println(criteria)
	//
	repo := adsdb.NewRepository("carsfinder")

	repo.AddCriteria(criteria)
	//err := repo.Migrate()
	//
	//errorshandler.HandleErr(err)
	//
	//repo.AddCriteria(criteria)

	//databaseName := "carsfinder"
	//dsn := fmt.Sprintf("root:rootpass@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseName)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//
	//errorshandler.HandleErr(err)
	//
	//jc := jobs.NewJobsCoordinatorService(db)
	//jc.Start()

}

func pOf(n int) *int {
	return &n
}

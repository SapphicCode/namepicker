package namepicker

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Name represents a name item
type Name struct {
	Name   string
	Gender string
	Rank   int
}

// Name list sorting
type nameList []*Name

func (l nameList) Len() int {
	return len(l)
}

func (l nameList) Less(i, j int) bool {
	return l[i].Rank > l[j].Rank
}

func (l nameList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Template arguments
type templateArgs struct {
	Title        string
	BodyTemplate string
	Content      string
	Years        []string
}

// NewEngine returns a populated *gin.Engine
func NewEngine(years map[string][]*Name) *gin.Engine {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*.html")

	yearList := make([]string, 0, len(years))
	for year := range years {
		yearList = append(yearList, year)
	}
	sort.Strings(yearList)

	rand.Seed(time.Now().UnixNano())

	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "root.html", templateArgs{
			Title:        "Generator",
			BodyTemplate: "generator",
			Years:        yearList,
		})
	})

	engine.GET("/names", func(ctx *gin.Context) {
		// find year
		year := ctx.Query("year")
		_, exists := years[year]
		if !exists {
			ctx.AbortWithError(401, errors.New("data for year does not exist"))
			return
		}

		// sort the list
		names := make(nameList, len(years[year]))
		copy(names, years[year])
		sort.Sort(names)
		fmt.Println(len(names))

		// filter by gender
		switch gender := ctx.Query("gender"); gender {
		case "M", "F":
			newNames := make(nameList, 0, len(names))
			for _, item := range names {
				if item.Gender == gender {
					newNames = append(newNames, item)
				}
			}
			names = newNames
		}

		// filter to user-specified length
		limit := ctx.Query("limit")
		if limit != "" {
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				ctx.AbortWithError(500, err)
			}
			if limitInt < 0 {
				ctx.AbortWithError(500, errors.New("Invalid top limit"))
			}
			names = names[:limitInt]
		}

		// pick x random names
		number := ctx.Query("n")
		if number != "" {
			newNames := make(nameList, 0, len(names))
			numberInt, _ := strconv.Atoi(number)
			// TODO: errors

			for i := 0; i < numberInt; i++ {
				n := rand.Intn(len(names))
				newNames = append(newNames, names[n])
			}
			names = newNames
		}

		nameString := "# Note: Refresh to generate another set!\n"
		for _, name := range names {
			nameString += fmt.Sprintf("%s,%d\t%s\n", name.Gender, name.Rank, name.Name)
		}

		ctx.String(200, nameString)
	})

	return engine
}

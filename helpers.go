package boop

import (
	"fmt"
	"junk/boop-server/pgdb"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func filterBoops(filter func(pgdb.Boop) bool, boops []pgdb.Boop) []pgdb.Boop {
	filtered_boops := []pgdb.Boop{}
	for _, boop := range boops {
		if filter(boop) {
			filtered_boops = append(filtered_boops, boop)
		}
	}
	return filtered_boops
}

func transformBoop(boop pgdb.Boop) string {
	return fmt.Sprintf("%v %v %v", boop.ID, boop.Created.Unix(), boop.Text)
}

func transformBoops(boops []pgdb.Boop) []string {
	transformedBoops := make([]string, len(boops))
	for i, boop := range boops {
		transformedBoops[i] = transformBoop(boop)
	}
	return transformedBoops
}

func containsEvent(boop pgdb.Boop) bool {
	matcher, err := regexp.MatchString(`\!(\w+)(\s|$)`, boop.Text)
	if err != nil {
		log.Print(err)
		return false
	}
	return matcher
}

const BEAN_MATCHER = `(?:(?:^)|(?:\s+))([\+-])(\$|\w+)(?::(\d+(\.\d+)?))?`

func containsBean(boop pgdb.Boop) bool {
	matcher, err := regexp.MatchString(BEAN_MATCHER, boop.Text)
	if err != nil {
		log.Print(err)
		return false
	}
	return matcher
}

func extractBeans(boop pgdb.Boop) []Bean {
	beans := []Bean{}

	re := regexp.MustCompile(BEAN_MATCHER)
	groups := re.FindAllStringSubmatch(strings.ToLower(boop.Text), 32)

	for _, group := range groups {
		bean := Bean{
			Name:  group[2],
			Value: 1,
		}

		if group[3] != "" {
			val, err := strconv.ParseFloat(group[3], 32)
			if err != nil {
				log.Print(err)
			}

			bean.Value = float32(val)
		}

		if group[1] != "+" {
			bean.Value = -1 * bean.Value
		}

		beans = append(beans, bean)
	}

	return beans
}

func simplifyBeans(beans []Bean) []Bean {
	total := map[string]float32{}
	for _, bean := range beans {
		val, ok := total[bean.Name]
		if !ok {
			total[bean.Name] = bean.Value
		} else {
			total[bean.Name] = val + bean.Value
		}
	}
	simplifiedBeans := []Bean{}
	for name, value := range total {
		simplifiedBeans = append(simplifiedBeans, Bean{
			Name:  name,
			Value: value,
		})
	}
	sort.Slice(simplifiedBeans, func(i, j int) bool {
		return simplifiedBeans[i].Name < simplifiedBeans[j].Name
	})
	return simplifiedBeans
}

package another

import "fmt"

var mine1 []string = []string{"ore", "ore", "rock", "mud", "ore", "rock"}

func findOre() func([]string) (int, []string) {
	slotInfo := -1
	return func(mine []string) (int, []string) {
		slotInfo++
		material := (mine)[0]
		mine[0] = mine[len(mine)-1]
		mine[len(mine)-1] = ""
		mine = mine[:len(mine)-1]
		// fmt.Println(mine)
		if material == "ore" {
			return slotInfo, mine
		}
		return -1, mine
	}
}

func main() {
	miner1 := findOre()

	for i := 0; i < 3; i++ {
		fmt.Println(miner1(mine1))
	}
}

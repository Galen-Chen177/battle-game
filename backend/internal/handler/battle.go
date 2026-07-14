package handler

import (
	"net/http"

	"battle-game/internal/battle"
	"battle-game/internal/engine"

	"github.com/gin-gonic/gin"
)

func StartBattle(c *gin.Context) {
	var units []battle.Unit

	if err := c.ShouldBindJSON(&units); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var leftUnits, rightUnits []*battle.Unit

	for i := range units {
		u := &units[i]

		// 校验 camp
		if u.Camp != 0 && u.Camp != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid camp value, must be 0 (Left) or 1 (Right)"})
			return
		}

		// 归一化状态
		u.Alive = true
		if u.MaxHP == 0 {
			u.MaxHP = u.HP
		}

		if u.Camp == 0 {
			leftUnits = append(leftUnits, u)
		} else {
			rightUnits = append(rightUnits, u)
		}
	}

	if len(leftUnits) == 0 || len(rightUnits) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both camps must have at least one unit"})
		return
	}

	left := &battle.Team{Name: "Left", Units: leftUnits}
	right := &battle.Team{Name: "Right", Units: rightUnits}

	b := battle.NewBattle(left, right)
	e := engine.New(b)
	e.Run()

	c.JSON(http.StatusOK, b.Events)
}

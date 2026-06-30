package flags

import (
	"danzmen/db"
	"fmt"
	"os"
	"strconv"
)

func (p *ProgramOpts) FlagToggle(db *db.SqliteDB, dbTasks []*db.DBJoin_DateRecord_Tasks) {
	f := p.Args[0]
	id, err := strconv.Atoi(f)
	if err != nil {
		fmt.Println("Invalid {id}. Only accepting positive numbers.")
		os.Exit(1)
	}

	for _, v := range dbTasks {
		if v.Id == id {
			var c int = 0
			if v.Completed == 0 {
				c = 1
			}

			if err := db.UpdateCompletedTask(id, c); err != nil {
				fmt.Println("Error: ", err.Error())
				os.Exit(1)
				return
			}
			fmt.Println("Success")
			return
		}
	}
	fmt.Println("Id not found. Does it exists in today's date?")
	os.Exit(1)

}

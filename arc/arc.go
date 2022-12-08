/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId:
Description:
*****************************************************************************/

package arc

import (
	"fmt"
)

type ARC struct{}

func initARC() (*ARC, error) {
	fmt.Println("ARC created.")
	return new(ARC), nil
}

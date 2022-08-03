package artifacts

import "fmt"

// HandlerError trata menssagem de erro
func HandlerError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

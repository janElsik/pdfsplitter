package functions

import (
	"fmt"
	"os/exec"
)

func ConvertToPdf(fileInput string) {
	cmd := exec.Command("unoconv", "-f", "pdf", fileInput)

	if err := cmd.Run(); err != nil {
		fmt.Printf("error with converting to pdf: %v \n", err)
		fmt.Println(fileInput)
	}

}

package mercedes_benz_de

import (
	"carscraper/pkg/jobs"
	"fmt"
)

type MercedesBenzRoURLBuilder struct {
}

func NewMercedesBenzRoURLBuilder() *MercedesBenzRoURLBuilder {
	return &MercedesBenzRoURLBuilder{}
}

func (b MercedesBenzRoURLBuilder) GetURL(job jobs.SessionJob) string {
	return fmt.Sprintf("https://cdn.sip.mercedes-benz.com/api/vs/v3/UCui/DE/overview")
}

func (b MercedesBenzRoURLBuilder) GetCountURL(job jobs.SessionJob) string {
	return fmt.Sprintf("https://cdn.sip.mercedes-benz.com/api/vs/v3/UCui/DE/count")
}
package mobile

import (
	"carscraper/pkg/jobs"
	"io"
	"net/http"
)

type Request struct {
	urlBuilder *URLBuilder
}

func NewRequest(criteria jobs.Criteria) *Request {
	builder := NewURLBuilder(criteria)
	return &Request{
		urlBuilder: builder,
	}
}

func (r Request) _GetPage(pageNumber int) []byte {

	url := r.urlBuilder.GetPageURL(pageNumber)

	httpMethod := "GET"
	httpClient := &http.Client{}
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Add("Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "show_qs_e=vhc%3Acar%2Cms1%3A17200_-58_%2Cfrn%3A2019%2Cful%3Adiesel%2Cmlx%3A125000; sorting_e=PRICE_ASC; show_qs_e=vhc%3Acar%2Cms1%3A17200_-58_%2Cfrn%3A2019%2Cful%3Adiesel%2Cmlx%3A125000; sorting_e=PRICE_ASC; _abck=C11B91BF8F400E4A6B3812393C9F69AB~-1~YAAQh15swacqESyMAQAAp6/5RQuqy9WlSi65V2iHNOUY7E76BK1Jl/MeLlXQDSH7mVdien9Qxd1PjFU30bcjTLcVKzvPhPK8vtu9pTM3SVVLxn98eq0PDpCXf2yRPudlycehtD/Ilkr3p+6arlynn5qKqKnXtnSY1eYGTm0bMDzS6Kn7JglAyFaX0qobF7W8KlrXYALD0Kweq8WNYCGqOQm6ZV2K3y1zdngr91XxM0BrNVqQVUMr/6xUaJtgQDSe77T42vSBFq1jvIRVtA0VhGosXffI4C+fDZao++XISSzud9twWqUAKhb+Utz/eQ+mqpGAatiaW5FjNnopikTnM607SLLXqTovy9MHb2LEKZTwW0Btry3pPkbAa0mSmx/h4yO6kQIZ~-1~-1~-1; ak_bmsc=425589E5170E187FCD3A4E8FE60FD452~000000000000000000000000000000~YAAQh15swagqESyMAQAAp6/5RRYmRPVXDb0NNOGkHh1uFaYZ/N+gv5xwVbxEiwSeX0ZvhGoXXqujaZvmBbx1aYVCERrqz71kND47WkOxA614ZdEAnPeZVqlUiHpsaqrx78BFyS+ATINnGBdESk0ZNm63+mJyGFB+cvAN3biucYXnYiTTg0u8hM3YloIcU3kPPAp179ZUQMKEsLPgeGAp8mhbRFze70i4bCObrNQPQPUPmxGAeyEXowLzSR0pyXGHLLL3XnWB7/Xb85XOU4GuNSpoYsw28kP90GJn6Ytb0e1boEU5T0dEalApHTs6vlGeeNhnInhhtKoqyGfHaUYRlWvhZ2OKGtzxTszorh11xHDdvpxfKhvplg8V; bm_mi=A296828ABDB3901E7F871C1516195623~YAAQtnp7XL0Vl9aLAQAAhEViRhaAoqPwAS0lzcDbiDabwCZ2DNufJ84FGd1GEgeb/iZgPhRTDKgJYW459AXJ6iotj3aHochUlAJxfaLYL2Tv+N7kdrla0uYdnMZ5U3Ji7OmlrQhANVSMOmfBYxi+BAEPCM1Te8bduOdOgijCdLIkJeiq87vPwfnIIfxdFm5XOLnGyIa+I3kw8Ox1odQnSjxPT7Hvt+hpfxKzKGrWd5ePBNq/rvBiGEeEEPyG3YgHqPDg2beINndLb0YgYxl3x47FwMLBEy3axM9fU1r1TKkBRdc/Utqtvs4pMkSJEsDrKW13FisyFIKWn98NqZOXe033kQqfpsrzOiakFl/jzJ184NHy5NwDqR638gsYxA9yyYoosLRH7rs7gjAW1mg3LOqq/itJtV1idsFWtMsl4khqRtNuHU/fMKk4BY8Y3x+c2Vuvtf3PIzdf+UR9rKdkXPyIXBJ2wz5GdCL9hWGD4w==~1; bm_sv=5E4B61EAB72A54817F83F71C4A62AC6E~YAAQtnp7XL4Vl9aLAQAAhEViRhYHRhiNeF0qO4NCwxeQrjiK9maJ4ilx70Bfb3OLeizCKT49VRDW1iZtvXbwodlwGIIrihQW3yPX2ViLjtzQFcnUVYkpAWzWB9qmjJtpxmZs28I+6BR/TJCG1IZacFwKj4EAt1sbj2QpX9DnENjCLu077dL8r6ZClhwo2P0ulrXbqdKfVT//bNtp0F8VXbuQj0wZJt1e4+CPZakYRfmf4ZeI3MXFdspVNCD1HLE=~1; bm_sz=119943C7EC81472A503D92AD61FC0AA2~YAAQh15swaoqESyMAQAAp6/5RRbpctURtsbsEBYaD0UnUX/OqFrl4LwMe8F8fYytx36p3fmYFKK7woVuJ3MvJH85VeVy8KkamKgkl/ioIAe7TU4UzJEhYFZzWjaPKC0X8bPZBYm79Vb2DxSlvZSk1EiCnM4UDBcWZR8FSAwZX7RZKlNOpR4tuk7cztxtEYZudkIG/io3vUWtd/riRTK+cCByP7t/i22hMhaPxr9zbUIHU4tAierGUsIWjydixaGn5O0jWs3e0WYeuLjbCrKkdcFFSKX0P5eCKWTEZLugJzWqhg==~3159619~3551284; vi=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaWQiOiJmN2FkNmM3YS1lYmNkLTRhMjctYWFlOC0yYTcwNGQwZDg3NWIiLCJkbnQiOnRydWUsImlhdCI6MTY5NjAyMTA1NywiYXVkIjpbXX0.fHtDyJqqroj1vbj7Enyst-XdLkJPQeQ_8K4FsEI2j6M; _abck=695DE7246B311E2008C5952DB9FC4D90~-1~YAAQol5swRtFjjSMAQAA11u5QAtYJpZc01ghfnMT8n0CN7p6bLoxuGHUuT8BQ/0G2ZqROsML2mHizNDN3EXh+BJMFP50hhT6CiazdE25t84quMCoOsho8cNkLZm5Q+IQXo1uqCXdNQyDVOviA3yB8rMaN6VZtCvLDMLzthGDD+QklV19cLF8laEu5BR/16JIQI5m9qoNzgs0eW01VZ5RC4hVTX0ugWRcArD8OKcpP5msDxAp2juHE67UuQBZie/JDeLofa/J/jAw1hU7Cuyc2+w5DXeEAWcbAPT7mZFqY9YowVA6rCQbhTEoSz+cZoemh/fYyDpevgdppGq3d7DvlxtuhPwoG5EWWD3UEOC2o/4SJsi2NwmKqi+3EGS96aQM3IaV01bp0H3GRPPqGwVctuBX1cHrD2c=~0~-1~-1")

	response, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return bodyBytes
}

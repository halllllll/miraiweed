package ready

import "fmt"

type URLs struct {
	Num                   string
	Base                  string
	Login                 string
	Service               string
	EducationBoard        string
	StudentsSearch        string
	StudentsSearchReflesh string
	TeacherSearch         string
	Search                string
}

func NewUrls() *URLs {
	return &URLs{}
}

func (u *URLs) PrepareUrl(num string) {
	u.Num = num
	u.Base = fmt.Sprintf("https://miraiseed%s.benesse.ne.jp", u.Num)
	u.Login = fmt.Sprintf("%s/seed/vw020101/displayLogin/1", u.Base)
	u.Service = fmt.Sprintf("%s/seed/vw030101/displaySchoolAdminMenu", u.Base)
	u.EducationBoard = fmt.Sprintf("%s/seed/vw030101/displayEducationBoardMenu", u.Base)
	u.StudentsSearch = fmt.Sprintf("%s/seed/vw030501/displaySearchChildInfo", u.Base)
	u.StudentsSearchReflesh = fmt.Sprintf("%s/seed/vw030501/refresh", u.Base)
	u.Search = fmt.Sprintf("%s/seed/vw030501/search", u.Base)
	u.TeacherSearch = fmt.Sprintf("%s/seed/vw030401/", u.Base)
}

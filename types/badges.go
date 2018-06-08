package types

type Badge struct {
	Name string
	ImageUrl string
}

var Badges []Badge
var NamedBadges map[string]Badge

func init() {
	NamedBadges = map[string]Badge{}
	NamedBadges["Complexity"] = Badge{Name: "Complexity", ImageUrl: "https://image.flaticon.com/icons/svg/887/887904.svg"}
	NamedBadges["Social"] = Badge{Name: "Social", ImageUrl: "https://image.flaticon.com/icons/svg/437/437553.svg"}
	NamedBadges["Profit"] = Badge{Name: "Profit", ImageUrl: "https://image.flaticon.com/icons/svg/138/138255.svg"}
	NamedBadges["Impact"] = Badge{Name: "Impact", ImageUrl: "https://image.flaticon.com/icons/svg/815/815376.svg"}
	NamedBadges["Coolness"] = Badge{Name: "Coolness", ImageUrl: "https://image.flaticon.com/icons/svg/576/576804.svg"}
	NamedBadges["Open Source"] = Badge{Name: "Open Source", ImageUrl: "/public/open-source.png"}

	Badges = []Badge{
		NamedBadges["Complexity"],
		NamedBadges["Social"],
		NamedBadges["Profit"],
		NamedBadges["Impact"],
		NamedBadges["Coolness"],
		NamedBadges["Open Source"],
	}
}

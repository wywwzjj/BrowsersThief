package browsers

type BrowserData interface {
	GetLoginData()
	GetCookies()
	GetHistory()
	GetBookmarks()
	GetWebData()
}

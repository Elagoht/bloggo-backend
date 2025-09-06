package models

type ViewStatistics struct {
	ViewsToday     int64 `json:"viewsToday"`
	ViewsThisWeek  int64 `json:"viewsThisWeek"`
	ViewsThisMonth int64 `json:"viewsThisMonth"`
	ViewsThisYear  int64 `json:"viewsThisYear"`
	TotalViews     int64 `json:"totalViews"`
}

type HourlyViewCount struct {
	Hour      int   `json:"hour"`
	ViewCount int64 `json:"viewCount"`
}

type Last24HoursViews struct {
	Hours []HourlyViewCount `json:"hours"`
}

type DailyViewCount struct {
	Day       int   `json:"day"`
	ViewCount int64 `json:"viewCount"`
}

type LastMonthViews struct {
	Days []DailyViewCount `json:"days"`
}

type MonthlyViewCount struct {
	Year      int   `json:"year"`
	Month     int   `json:"month"`
	ViewCount int64 `json:"viewCount"`
}

type LastYearViews struct {
	Months []MonthlyViewCount `json:"months"`
}

type CategoryViewDistribution struct {
	CategoryId   int64   `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	ViewCount    int64   `json:"viewCount"`
	Percentage   float64 `json:"percentage"`
}

type MostViewedBlog struct {
	PostId       int64  `json:"postId"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	ViewCount    int64  `json:"viewCount"`
	Author       string `json:"author"`
	CategoryName string `json:"categoryName"`
}

type BlogStatistics struct {
	TotalPublishedBlogs int64   `json:"totalPublishedBlogs"`
	TotalDraftedBlogs   int64   `json:"totalDraftedBlogs"`
	TotalPendingBlogs   int64   `json:"totalPendingBlogs"`
	TotalReadTime       int64   `json:"totalReadTime"`
	AverageReadTime     float64 `json:"averageReadTime"`
	AverageViews        float64 `json:"averageViews"`
}

type CategoryBlogDistribution struct {
	CategoryId   int64   `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	BlogCount    int64   `json:"blogCount"`
	Percentage   float64 `json:"percentage"`
}

type CategoryReadTimeDistribution struct {
	CategoryId      int64   `json:"categoryId"`
	CategoryName    string  `json:"categoryName"`
	TotalReadTime   int64   `json:"totalReadTime"`
	AverageReadTime float64 `json:"averageReadTime"`
	Percentage      float64 `json:"percentage"`
}

type CategoryLengthDistribution struct {
	CategoryId    int64   `json:"categoryId"`
	CategoryName  string  `json:"categoryName"`
	TotalLength   int64   `json:"totalLength"`
	AverageLength float64 `json:"averageLength"`
	Percentage    float64 `json:"percentage"`
}

type LongestBlog struct {
	PostId       int64  `json:"postId"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	ReadTime     int    `json:"readTime"`
	Author       string `json:"author"`
	CategoryName string `json:"categoryName"`
}

type UserAgentStat struct {
	UserAgent  string  `json:"userAgent"`
	ViewCount  int64   `json:"viewCount"`
	Percentage float64 `json:"percentage"`
}

type DeviceTypeStat struct {
	DeviceType string  `json:"deviceType"`
	ViewCount  int64   `json:"viewCount"`
	Percentage float64 `json:"percentage"`
}

type OSStatistic struct {
	OS         string  `json:"operatingSystem"`
	ViewCount  int64   `json:"viewCount"`
	Percentage float64 `json:"percentage"`
}

type BrowserStat struct {
	Browser    string  `json:"browser"`
	ViewCount  int64   `json:"viewCount"`
	Percentage float64 `json:"percentage"`
}

type AuthorStatistics struct {
	AuthorId        int64   `json:"authorId"`
	AuthorName      string  `json:"authorName"`
	TotalBlogs      int64   `json:"totalBlogs"`
	TotalViews      int64   `json:"totalViews"`
	AverageViews    float64 `json:"averageViews"`
	TotalReadTime   int64   `json:"totalReadTime"`
	AverageReadTime float64 `json:"averageReadTime"`
}

type ResponseAllStatistics struct {
	ViewStats                    *ViewStatistics                `json:"viewStatistics"`
	Last24Hours                  *Last24HoursViews              `json:"last24HoursViews"`
	LastMonth                    *LastMonthViews                `json:"lastMonthViews"`
	LastYear                     *LastYearViews                 `json:"lastYearViews"`
	CategoryViewsDistribution    []CategoryViewDistribution     `json:"categoryViewsDistribution"`
	MostViewedBlogs              []MostViewedBlog               `json:"mostViewedBlogs"`
	BlogStats                    *BlogStatistics                `json:"blogStatistics"`
	LongestBlogs                 []LongestBlog                  `json:"longestBlogs"`
	CategoryBlogsDistribution    []CategoryBlogDistribution     `json:"categoryBlogsDistribution"`
	CategoryReadTimeDistribution []CategoryReadTimeDistribution `json:"categoryReadTimeDistribution"`
	CategoryLengthDistribution   []CategoryLengthDistribution   `json:"categoryLengthDistribution"`
	TopUserAgents                []UserAgentStat                `json:"topUserAgents"`
	DeviceTypeDistribution       []DeviceTypeStat               `json:"deviceTypeDistribution"`
	OSDistribution               []OSStatistic                  `json:"operatingSystemDistribution"`
	BrowserDistribution          []BrowserStat                  `json:"browserDistribution"`
}

type ResponseAuthorStatistics struct {
	AuthorStats                  *AuthorStatistics              `json:"authorStatistics"`
	ViewStats                    *ViewStatistics                `json:"viewStatistics"`
	Last24Hours                  *Last24HoursViews              `json:"last24HoursViews"`
	LastMonth                    *LastMonthViews                `json:"lastMonthViews"`
	LastYear                     *LastYearViews                 `json:"lastYearViews"`
	CategoryViewsDistribution    []CategoryViewDistribution     `json:"categoryViewsDistribution"`
	MostViewedBlogs              []MostViewedBlog               `json:"mostViewedBlogs"`
	BlogStats                    *BlogStatistics                `json:"blogStatistics"`
	LongestBlogs                 []LongestBlog                  `json:"longestBlogs"`
	CategoryBlogsDistribution    []CategoryBlogDistribution     `json:"categoryBlogsDistribution"`
	CategoryReadTimeDistribution []CategoryReadTimeDistribution `json:"categoryReadTimeDistribution"`
	CategoryLengthDistribution   []CategoryLengthDistribution   `json:"categoryLengthDistribution"`
	TopUserAgents                []UserAgentStat                `json:"topUserAgents"`
	DeviceTypeDistribution       []DeviceTypeStat               `json:"deviceTypeDistribution"`
	OSDistribution               []OSStatistic                  `json:"operatingSystemDistribution"`
	BrowserDistribution          []BrowserStat                  `json:"browserDistribution"`
}

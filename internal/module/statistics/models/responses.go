package models

type ViewStatistics struct {
	ViewsToday     int64 `json:"views_today"`
	ViewsThisWeek  int64 `json:"views_this_week"`
	ViewsThisMonth int64 `json:"views_this_month"`
	ViewsThisYear  int64 `json:"views_this_year"`
	TotalViews     int64 `json:"total_views"`
}

type HourlyViewCount struct {
	Hour       int   `json:"hour"`
	ViewCount  int64 `json:"view_count"`
}

type Last24HoursViews struct {
	Hours []HourlyViewCount `json:"hours"`
}

type DailyViewCount struct {
	Day       int   `json:"day"`
	ViewCount int64 `json:"view_count"`
}

type LastMonthViews struct {
	Days []DailyViewCount `json:"days"`
}

type MonthlyViewCount struct {
	Month     int   `json:"month"`
	ViewCount int64 `json:"view_count"`
}

type LastYearViews struct {
	Months []MonthlyViewCount `json:"months"`
}

type CategoryViewDistribution struct {
	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	ViewCount    int64  `json:"view_count"`
	Percentage   float64 `json:"percentage"`
}

type MostViewedBlog struct {
	PostId      int64  `json:"post_id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	ViewCount   int64  `json:"view_count"`
	Author      string `json:"author"`
	CategoryName string `json:"category_name"`
}

type BlogStatistics struct {
	TotalPublishedBlogs int64   `json:"total_published_blogs"`
	TotalDraftedBlogs   int64   `json:"total_drafted_blogs"`
	TotalPendingBlogs   int64   `json:"total_pending_blogs"`
	TotalReadTime       int64   `json:"total_read_time"`
	AverageReadTime     float64 `json:"average_read_time"`
	AverageViews        float64 `json:"average_views"`
}

type CategoryBlogDistribution struct {
	CategoryId    int64   `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	BlogCount     int64   `json:"blog_count"`
	Percentage    float64 `json:"percentage"`
}

type CategoryReadTimeDistribution struct {
	CategoryId       int64   `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	TotalReadTime    int64   `json:"total_read_time"`
	AverageReadTime  float64 `json:"average_read_time"`
	Percentage       float64 `json:"percentage"`
}

type CategoryLengthDistribution struct {
	CategoryId       int64   `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	TotalLength      int64   `json:"total_length"`
	AverageLength    float64 `json:"average_length"`
	Percentage       float64 `json:"percentage"`
}

type LongestBlog struct {
	PostId       int64  `json:"post_id"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	ReadTime     int    `json:"read_time"`
	Author       string `json:"author"`
	CategoryName string `json:"category_name"`
}

type UserAgentStat struct {
	UserAgent string `json:"user_agent"`
	ViewCount int64  `json:"view_count"`
	Percentage float64 `json:"percentage"`
}

type DeviceTypeStat struct {
	DeviceType string  `json:"device_type"`
	ViewCount  int64   `json:"view_count"`
	Percentage float64 `json:"percentage"`
}

type OSStatistic struct {
	OS         string  `json:"operating_system"`
	ViewCount  int64   `json:"view_count"`
	Percentage float64 `json:"percentage"`
}

type BrowserStat struct {
	Browser    string  `json:"browser"`
	ViewCount  int64   `json:"view_count"`
	Percentage float64 `json:"percentage"`
}

type AuthorStatistics struct {
	AuthorId            int64   `json:"author_id"`
	AuthorName          string  `json:"author_name"`
	TotalBlogs          int64   `json:"total_blogs"`
	TotalViews          int64   `json:"total_views"`
	AverageViews        float64 `json:"average_views"`
	TotalReadTime       int64   `json:"total_read_time"`
	AverageReadTime     float64 `json:"average_read_time"`
}

type ResponseAllStatistics struct {
	ViewStats                     *ViewStatistics                    `json:"view_statistics"`
	Last24Hours                   *Last24HoursViews                  `json:"last_24_hours_views"`
	LastMonth                     *LastMonthViews                    `json:"last_month_views"`
	LastYear                      *LastYearViews                     `json:"last_year_views"`
	CategoryViewsDistribution     []CategoryViewDistribution         `json:"category_views_distribution"`
	MostViewedBlogs               []MostViewedBlog                   `json:"most_viewed_blogs"`
	BlogStats                     *BlogStatistics                    `json:"blog_statistics"`
	LongestBlogs                  []LongestBlog                      `json:"longest_blogs"`
	CategoryBlogsDistribution     []CategoryBlogDistribution         `json:"category_blogs_distribution"`
	CategoryReadTimeDistribution  []CategoryReadTimeDistribution     `json:"category_read_time_distribution"`
	CategoryLengthDistribution    []CategoryLengthDistribution       `json:"category_length_distribution"`
	TopUserAgents                 []UserAgentStat                    `json:"top_user_agents"`
	DeviceTypeDistribution        []DeviceTypeStat                   `json:"device_type_distribution"`
	OSDistribution                []OSStatistic                      `json:"operating_system_distribution"`
	BrowserDistribution           []BrowserStat                      `json:"browser_distribution"`
}

type ResponseAuthorStatistics struct {
	AuthorStats                   *AuthorStatistics                  `json:"author_statistics"`
	ViewStats                     *ViewStatistics                    `json:"view_statistics"`
	Last24Hours                   *Last24HoursViews                  `json:"last_24_hours_views"`
	LastMonth                     *LastMonthViews                    `json:"last_month_views"`
	LastYear                      *LastYearViews                     `json:"last_year_views"`
	CategoryViewsDistribution     []CategoryViewDistribution         `json:"category_views_distribution"`
	MostViewedBlogs               []MostViewedBlog                   `json:"most_viewed_blogs"`
	BlogStats                     *BlogStatistics                    `json:"blog_statistics"`
	LongestBlogs                  []LongestBlog                      `json:"longest_blogs"`
	CategoryBlogsDistribution     []CategoryBlogDistribution         `json:"category_blogs_distribution"`
	CategoryReadTimeDistribution  []CategoryReadTimeDistribution     `json:"category_read_time_distribution"`
	CategoryLengthDistribution    []CategoryLengthDistribution       `json:"category_length_distribution"`
	TopUserAgents                 []UserAgentStat                    `json:"top_user_agents"`
	DeviceTypeDistribution        []DeviceTypeStat                   `json:"device_type_distribution"`
	OSDistribution                []OSStatistic                      `json:"operating_system_distribution"`
	BrowserDistribution           []BrowserStat                      `json:"browser_distribution"`
}
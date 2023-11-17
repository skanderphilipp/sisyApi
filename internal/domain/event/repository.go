package event

type Repository interface {
	FindUpcomingByVenueID()
	FindPastEventsByVenueID()
	FindAllUpcoming()
	FindToday()
	FindTomorrow()
	FindCurrent()
	Save()
	Update()
	Delete()
}

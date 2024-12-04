package session

type PomoSession struct {
	status         PomoSessionStatus
	WorkSessionsBreak uint // CONFIG VARIABLE...
	workedSessions uint
}

func (s *PomoSession) Status() PomoSessionStatus {
	return s.status
}

func (s *PomoSession) CompletedWorkSessions() uint {
	return s.workedSessions
}


func (s *PomoSession) GetNextStatus() PomoSessionStatus {
	// EITHER LONG OR SHORT BREAK
	if s.status != PomoSessionWork {
		return PomoSessionWork
	}

	nWorkedSessions := s.workedSessions
	if s.status == PomoSessionWork{
		nWorkedSessions = nWorkedSessions + 1
	}

	if (nWorkedSessions % s.WorkSessionsBreak) == 0 {
		return PomoSessionLongBreak
	}

	return PomoSessionShortBreak
}


func (s *PomoSession) SetNextStatus(status PomoSessionStatus) {
	if s.status != PomoSessionWork && status == PomoSessionWork {
		s.workedSessions++
	}
	s.status = status
}

func (s *PomoSession) Reset() {
	s.status = PomoSessionWork
	s.workedSessions = 0
}

package http


type httpController struct {
	*httpService
}
func (self *httpController) New(hService *httpService) *httpController {
	self.httpService = hService
	return self
}

package logging

type TestReporter struct {
}

func (tr *TestReporter) ReportException(error) {

}
func (tr *TestReporter) ReportFatal(error) {

}
func NewTestReporter() Reporter {
	return &TestReporter{}
}

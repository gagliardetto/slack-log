package slacklog

// FakeSlackLogClient is a slack client that does not send any slack hook (used for testing only)
type FakeSlackLogClient struct{}

func (c FakeSlackLogClient) Log15(msg string, ctx ...interface{}) MessageInterface {
	return c
}
func (c FakeSlackLogClient) WithField(key string, value interface{}) MessageInterface {
	return c
}
func (c FakeSlackLogClient) WithFields(fields Fields) MessageInterface {
	return c
}

func (c FakeSlackLogClient) SetTitle(title string) MessageInterface {
	return c
}

func (c FakeSlackLogClient) Send() error {
	return nil
}
func (c FakeSlackLogClient) Say(text string) error {
	return nil
}
func (c FakeSlackLogClient) Shout(text string) error {
	return nil
}

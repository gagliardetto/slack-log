package slacklog

// FakeSlackLogClient is a slack client that does not send any slack hook (used for testing only)
type FakeSlackLogClient struct{}

func (c FakeSlackLogClient) Log15(msg string, ctx ...interface{}) ClientInterface {
	return c
}
func (c FakeSlackLogClient) WithField(key string, value interface{}) ClientInterface {
	return c
}
func (c FakeSlackLogClient) WithFields(fields ContextFields) ClientInterface {
	return c
}

func (c FakeSlackLogClient) SetMessage(msg string) ClientInterface {
	return c
}
func (c FakeSlackLogClient) SetColor(hex string) ClientInterface {
	return c
}
func (c FakeSlackLogClient) SetIcon(iconEmoji string) ClientInterface {
	return c
}
func (c FakeSlackLogClient) OverrideChannel(channel string) ClientInterface {
	return c
}

func (c FakeSlackLogClient) Send() error {
	return nil
}
func (c FakeSlackLogClient) SimpleSend(text, iconEmoji string) error {
	return nil
}

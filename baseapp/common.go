package baseapp

func DereferenceBaseApp(app interface{}) *BaseApp {
	ba, ok := app.(*BaseApp)
	if !ok {
		return nil
	}
	return ba
}

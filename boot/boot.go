package boot

func Boot() error {
	if err := initConfig(); err != nil {
		return err
	}

	if err := initLog(); err != nil {
		return err
	}

	if err := initDB(); err != nil {
		return err
	}

	if err := initHTTP(); err != nil {
		return err
	}

	return nil
}
